package mempool

import (
	"context"
	"fmt"
	"math"

	"github.com/huandu/skiplist"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

var (
	_ Mempool  = (*PriorityNonceMempool[int])(nil)
	_ Iterator = (*PriorityNonceIterator[int])(nil)
)

type Comparable interface {
	comparable
}

type TxPriority[P comparable] struct {
	// GetTxPriority returns the priority of the transaction. A priority must be
	// comparable via Compare.
	GetTxPriority func(ctx context.Context, tx sdk.Tx) P

	// Compare compares two transaction priorities. The result should be
	// 0 if a == b, -1 if a < b, and +1 if a > b.
	Compare func(a, b P) int

	MinValue P
}

// NewDefaultTxPriority returns a TxPriority comparator using ctx.Priority as
// the defining transaction priority.
func NewDefaultTxPriority() TxPriority[int64] {
	return TxPriority[int64]{
		GetTxPriority: func(goCtx context.Context, tx sdk.Tx) int64 {
			return sdk.UnwrapSDKContext(goCtx).Priority()
		},
		Compare: func(a, b int64) int {
			return skiplist.Int64.Compare(a, b)
		},
		MinValue: math.MinInt64,
	}
}

// PriorityNonceMempool is a mempool implementation that stores txs
// in a partially ordered set by 2 dimensions: priority, and sender-nonce
// (sequence number). Internally it uses one priority ordered skip list and one
// skip list per sender ordered by sender-nonce (sequence number). When there
// are multiple txs from the same sender, they are not always comparable by
// priority to other sender txs and must be partially ordered by both sender-nonce
// and priority.
type PriorityNonceMempool[P comparable] struct {
	priorityIndex  *skiplist.SkipList
	priorityCounts map[P]int
	senderIndices  map[string]*skiplist.SkipList
	scores         map[txMeta[P]]txMeta[P]
	config         PriorityNonceMempoolConfig[P]
}

type PriorityNonceIterator[P comparable] struct {
	mempool       *PriorityNonceMempool[P]
	priorityNode  *skiplist.Element
	senderCursors map[string]*skiplist.Element
	sender        string
	nextPriority  *skiplist.Element
}

// txMeta stores transaction metadata used in indices
type txMeta[P comparable] struct {
	// nonce is the sender's sequence number
	nonce uint64
	// priority is the transaction's priority
	priority P
	// sender is the transaction's sender
	sender string
	// weight is the transaction's weight, used as a tiebreaker for transactions
	// with the same priority
	weight P
	// senderElement is a pointer to the transaction's element in the sender index
	senderElement *skiplist.Element
}

// skiplistComparable is a comparator for txKeys that first compares priority,
// then weight, then sender, then nonce, uniquely identifying a transaction.
//
// Note, skiplistComparable is used as the comparator in the priority index.
func skiplistComparable[P comparable](txPriority TxPriority[P]) skiplist.Comparable {
	return skiplist.LessThanFunc(func(a, b any) int {
		keyA := a.(txMeta[P])
		keyB := b.(txMeta[P])

		res := txPriority.Compare(keyA.priority, keyB.priority)
		if res != 0 {
			return res
		}

		// Weight is used as a tiebreaker for transactions with the same priority.
		// Weight is calculated in a single pass in .Select(...) and so will be 0
		// on .Insert(...).
		res = txPriority.Compare(keyA.weight, keyB.weight)
		if res != 0 {
			return res
		}

		// Because weight will be 0 on .Insert(...), we must also compare sender and
		// nonce to resolve priority collisions. If we didn't then transactions with
		// the same priority would overwrite each other in the priority index.
		res = skiplist.String.Compare(keyA.sender, keyB.sender)
		if res != 0 {
			return res
		}

		return skiplist.Uint64.Compare(keyA.nonce, keyB.nonce)
	})
}

type PriorityNonceMempoolConfig[P comparable] struct {
	// TxPriority defines the transaction priority and comparator.
	TxPriority TxPriority[P]

	// OnRead is a callback to be called when a tx is read from the mempool.
	OnRead func(tx sdk.Tx)

	// TxReplacement is a callback to be called when duplicated transaction nonce
	// detected during mempool insert. An application can define a transaction
	// replacement rule based on tx priority or certain transaction fields.
	TxReplacement func(op, np P, oTx, nTx sdk.Tx) bool

	// MaxTx sets the maximum number of transactions allowed in the mempool with
	// the semantics:
	// - if MaxTx == 0, there is no cap on the number of transactions in the mempool
	// - if MaxTx > 0, the mempool will cap the number of transactions it stores,
	//   and will prioritize transactions by their priority and sender-nonce
	//   (sequence number) when evicting transactions.
	// - if MaxTx < 0, `Insert` is a no-op.
	MaxTx int
}

func DefaultPriorityNonceMempoolConfig() PriorityNonceMempoolConfig[int64] {
	return PriorityNonceMempoolConfig[int64]{
		TxPriority: NewDefaultTxPriority(),
	}
}

// NewPriorityMempool returns the SDK's default mempool implementation which
// returns txs in a partial order by 2 dimensions; priority, and sender-nonce.
func NewPriorityMempool[P comparable](config PriorityNonceMempoolConfig[P]) *PriorityNonceMempool[P] {
	mp := &PriorityNonceMempool[P]{
		priorityIndex:  skiplist.New(skiplistComparable(config.TxPriority)),
		priorityCounts: make(map[P]int),
		senderIndices:  make(map[string]*skiplist.SkipList),
		scores:         make(map[txMeta[P]]txMeta[P]),
		config:         config,
	}

	return mp
}

// NextSenderTx returns the next transaction for a given sender by nonce order,
// i.e. the next valid transaction for the sender. If no such transaction exists,
// nil will be returned.
func (mp *PriorityNonceMempool[P]) NextSenderTx(sender string) sdk.Tx {
	senderIndex, ok := mp.senderIndices[sender]
	if !ok {
		return nil
	}

	cursor := senderIndex.Front()
	return cursor.Value.(sdk.Tx)
}

// Insert attempts to insert a Tx into the app-side mempool in O(log n) time,
// returning an error if unsuccessful. Sender and nonce are derived from the
// transaction's first signature.
//
// Transactions are unique by sender and nonce. Inserting a duplicate tx is an
// O(log n) no-op.
//
// Inserting a duplicate tx with a different priority overwrites the existing tx,
// changing the total order of the mempool.
func (mp *PriorityNonceMempool[P]) Insert(ctx context.Context, tx sdk.Tx) error {
	maxTx := mp.config.MaxTx
	if maxTx > 0 && mp.CountTx() >= maxTx {
		return ErrMempoolTxMaxCapacity
	} else if maxTx < 0 {
		return nil
	}

	sigs, err := tx.(signing.SigVerifiableTx).GetSignaturesV2()
	if err != nil {
		return err
	}
	if len(sigs) == 0 {
		return fmt.Errorf("tx must have at least one signer")
	}

	sig := sigs[0]
	sender := sdk.AccAddress(sig.PubKey.Address()).String()
	priority := mp.config.TxPriority.GetTxPriority(ctx, tx)
	nonce := sig.Sequence
	key := txMeta[P]{nonce: nonce, priority: priority, sender: sender}

	senderIndex, ok := mp.senderIndices[sender]
	if !ok {
		senderIndex = skiplist.New(skiplist.LessThanFunc(func(a, b any) int {
			return skiplist.Uint64.Compare(b.(txMeta[P]).nonce, a.(txMeta[P]).nonce)
		}))

		// initialize sender index if not found
		mp.senderIndices[sender] = senderIndex
	}

	// Since mp.priorityIndex is scored by priority, then sender, then nonce, a
	// changed priority will create a new key, so we must remove the old key and
	// re-insert it to avoid having the same tx with different priorityIndex indexed
	// twice in the mempool.
	//
	// This O(log n) remove operation is rare and only happens when a tx's priority
	// changes.
	sk := txMeta[P]{nonce: nonce, sender: sender}
	if oldScore, txExists := mp.scores[sk]; txExists {
		if mp.config.TxReplacement != nil && !mp.config.TxReplacement(oldScore.priority, priority, senderIndex.Get(key).Value.(sdk.Tx), tx) {
			return fmt.Errorf(
				"tx doesn't fit the replacement rule, oldPriority: %v, newPriority: %v, oldTx: %v, newTx: %v",
				oldScore.priority,
				priority,
				senderIndex.Get(key).Value.(sdk.Tx),
				tx,
			)
		}

		mp.priorityIndex.Remove(txMeta[P]{
			nonce:    nonce,
			sender:   sender,
			priority: oldScore.priority,
			weight:   oldScore.weight,
		})
		mp.priorityCounts[oldScore.priority]--
	}

	mp.priorityCounts[priority]++

	// Since senderIndex is scored by nonce, a changed priority will overwrite the
	// existing key.
	key.senderElement = senderIndex.Set(key, tx)

	mp.scores[sk] = txMeta[P]{priority: priority}
	mp.priorityIndex.Set(key, tx)

	return nil
}

func (i *PriorityNonceIterator[P]) iteratePriority() Iterator {
	// beginning of priority iteration
	if i.priorityNode == nil {
		i.priorityNode = i.mempool.priorityIndex.Front()
	} else {
		i.priorityNode = i.priorityNode.Next()
	}

	// end of priority iteration
	if i.priorityNode == nil {
		return nil
	}

	i.sender = i.priorityNode.Key().(txMeta[P]).sender

	i.nextPriority = i.priorityNode.Next()
	return i.Next()
}

func (i *PriorityNonceIterator[P]) Next() Iterator {
	if i.priorityNode == nil {
		return nil
	}

	cursor, ok := i.senderCursors[i.sender]
	if !ok {
		// beginning of sender iteration
		cursor = i.mempool.senderIndices[i.sender].Front()
	} else {
		// middle of sender iteration
		cursor = cursor.Next()
	}

	// end of sender iteration
	if cursor == nil {
		return i.iteratePriority()
	}

	key := cursor.Key().(txMeta[P])

	// We've reached a transaction with a priority lower than the next highest
	// priority in the pool.
	if i.nextPriority == nil {
		i.senderCursors[i.sender] = cursor
		return i
	}
	nextPriorityKey := i.nextPriority.Key().(txMeta[P])
	nextPriority := nextPriorityKey.priority
	if i.mempool.config.TxPriority.Compare(key.priority, nextPriority) < 0 {
		return i.iteratePriority()
	} else if i.mempool.config.TxPriority.Compare(key.priority, nextPriority) == 0 {
		// Weight is incorporated into the priority index key only (not sender index)
		// so we must fetch it here from the scores map.
		weight := i.mempool.scores[txMeta[P]{nonce: key.nonce, sender: key.sender}].weight
		if i.mempool.config.TxPriority.Compare(weight, nextPriorityKey.weight) < 0 {
			return i.iteratePriority()
		}
	}

	i.senderCursors[i.sender] = cursor
	return i
}

func (i *PriorityNonceIterator[P]) Tx() sdk.Tx {
	return i.senderCursors[i.sender].Value.(sdk.Tx)
}

// Select returns a set of transactions from the mempool, ordered by priority
// and sender-nonce in O(n) time. The passed in list of transactions are ignored.
// This is a readonly operation, the mempool is not modified.
//
// The maxBytes parameter defines the maximum number of bytes of transactions to
// return.
func (mp *PriorityNonceMempool[P]) Select(_ context.Context, _ [][]byte) Iterator {
	if mp.priorityIndex.Len() == 0 {
		return nil
	}

	mp.reorderPriorityTies()

	iterator := &PriorityNonceIterator[P]{
		mempool:       mp,
		senderCursors: make(map[string]*skiplist.Element),
	}

	return iterator.iteratePriority()
}

type reorderKey[P comparable] struct {
	deleteKey txMeta[P]
	insertKey txMeta[P]
	tx        sdk.Tx
}

func (mp *PriorityNonceMempool[P]) reorderPriorityTies() {
	node := mp.priorityIndex.Front()

	var reordering []reorderKey[P]
	for node != nil {
		key := node.Key().(txMeta[P])
		if mp.priorityCounts[key.priority] > 1 {
			newKey := key
			newKey.weight = senderWeight[P](mp.config.TxPriority, key.senderElement)
			reordering = append(reordering, reorderKey[P]{deleteKey: key, insertKey: newKey, tx: node.Value.(sdk.Tx)})
		}

		node = node.Next()
	}

	for _, k := range reordering {
		mp.priorityIndex.Remove(k.deleteKey)
		delete(mp.scores, txMeta[P]{nonce: k.deleteKey.nonce, sender: k.deleteKey.sender})
		mp.priorityIndex.Set(k.insertKey, k.tx)
		mp.scores[txMeta[P]{nonce: k.insertKey.nonce, sender: k.insertKey.sender}] = k.insertKey
	}
}

// senderWeight returns the weight of a given tx (t) at senderCursor. Weight is
// defined as the first (nonce-wise) same sender tx with a priority not equal to
// t. It is used to resolve priority collisions, that is when 2 or more txs from
// different senders have the same priority.
func senderWeight[P comparable](txPriority TxPriority[P], senderCursor *skiplist.Element) P {
	if senderCursor == nil {
		return txPriority.MinValue
	}

	weight := senderCursor.Key().(txMeta[P]).priority
	senderCursor = senderCursor.Next()
	for senderCursor != nil {
		p := senderCursor.Key().(txMeta[P]).priority
		if txPriority.Compare(p, weight) != 0 {
			weight = p
		}

		senderCursor = senderCursor.Next()
	}

	return weight
}

// CountTx returns the number of transactions in the mempool.
func (mp *PriorityNonceMempool[P]) CountTx() int {
	return mp.priorityIndex.Len()
}

// Remove removes a transaction from the mempool in O(log n) time, returning an
// error if unsuccessful.
func (mp *PriorityNonceMempool[P]) Remove(tx sdk.Tx) error {
	sigs, err := tx.(signing.SigVerifiableTx).GetSignaturesV2()
	if err != nil {
		return err
	}
	if len(sigs) == 0 {
		return fmt.Errorf("attempted to remove a tx with no signatures")
	}

	sig := sigs[0]
	sender := sdk.AccAddress(sig.PubKey.Address()).String()
	nonce := sig.Sequence

	scoreKey := txMeta[P]{nonce: nonce, sender: sender}
	score, ok := mp.scores[scoreKey]
	if !ok {
		return ErrTxNotFound
	}
	tk := txMeta[P]{nonce: nonce, priority: score.priority, sender: sender, weight: score.weight}

	senderTxs, ok := mp.senderIndices[sender]
	if !ok {
		return fmt.Errorf("sender %s not found", sender)
	}

	mp.priorityIndex.Remove(tk)
	senderTxs.Remove(tk)
	delete(mp.scores, scoreKey)
	mp.priorityCounts[score.priority]--

	return nil
}

func (mp *PriorityNonceMempool[P]) IsEmpty() error {
	if mp.priorityIndex.Len() != 0 {
		return fmt.Errorf("priorityIndex not empty")
	}

	var countKeys []P
	for k := range mp.priorityCounts {
		countKeys = append(countKeys, k)
	}

	for _, k := range countKeys {
		if mp.priorityCounts[k] != 0 {
			return fmt.Errorf("priorityCounts not zero at %v, got %v", k, mp.priorityCounts[k])
		}
	}

	var senderKeys []string
	for k := range mp.senderIndices {
		senderKeys = append(senderKeys, k)
	}

	for _, k := range senderKeys {
		if mp.senderIndices[k].Len() != 0 {
			return fmt.Errorf("senderIndex not empty for sender %v", k)
		}
	}

	return nil
}
