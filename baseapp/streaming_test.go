package baseapp

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	baseapptestutil "github.com/cosmos/cosmos-sdk/baseapp/testutil"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

var _ storetypes.ABCIListener = (*MockABCIListener)(nil)

type MockABCIListener struct {
	name      string
	ChangeSet []*storetypes.StoreKVPair
}

func NewMockABCIListener(name string) MockABCIListener {
	return MockABCIListener{
		name:      name,
		ChangeSet: make([]*storetypes.StoreKVPair, 0),
	}
}

func (m *MockABCIListener) ListenBeginBlock(context.Context, abci.RequestBeginBlock, abci.ResponseBeginBlock) error {
	return nil
}

func (m *MockABCIListener) ListenEndBlock(context.Context, abci.RequestEndBlock, abci.ResponseEndBlock) error {
	return nil
}

func (m *MockABCIListener) ListenDeliverTx(context.Context, abci.RequestDeliverTx, abci.ResponseDeliverTx) error {
	return nil
}

func (m *MockABCIListener) ListenCommit(_ context.Context, _ abci.ResponseCommit, changeSet []*storetypes.StoreKVPair) error {
	m.ChangeSet = changeSet
	return nil
}

var distKey1 = storetypes.NewKVStoreKey("distKey1")

type BaseAppSuite struct {
	baseApp  *BaseApp
	cdc      *codec.ProtoCodec
	txConfig client.TxConfig
}

func newBaseAppSuite(t *testing.T, opts ...func(*BaseApp)) *BaseAppSuite {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	baseapptestutil.RegisterInterfaces(cdc.InterfaceRegistry())

	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	db := dbm.NewMemDB()

	app := NewBaseApp(t.Name(), defaultLogger(), db, txConfig.TxDecoder(), opts...)
	require.Equal(t, t.Name(), app.Name())

	app.SetInterfaceRegistry(cdc.InterfaceRegistry())
	app.MsgServiceRouter().SetInterfaceRegistry(cdc.InterfaceRegistry())
	app.MountStores(capKey1, capKey2)
	app.SetParamStore(&paramStore{db: dbm.NewMemDB()})

	// mount stores and seal
	require.Nil(t, app.LoadLatestVersion())

	return &BaseAppSuite{
		baseApp:  app,
		cdc:      cdc,
		txConfig: txConfig,
	}
}

func newWrappedTxCounter(cfg client.TxConfig, counter int64, msgCounters ...int64) signing.Tx {
	msgs := make([]sdk.Msg, 0, len(msgCounters))
	for _, c := range msgCounters {
		msg := &baseapptestutil.MsgCounter{Counter: c, FailOnHandler: false}
		msgs = append(msgs, msg)
	}

	builder := cfg.NewTxBuilder()
	builder.SetMsgs(msgs...)
	builder.SetMemo("counter=" + strconv.FormatInt(counter, 10) + "&failOnAnte=false")

	return builder.GetTx()
}

func TestABCI_MultiListener_StateChanges(t *testing.T) {
	anteKey := []byte("ante-key")
	anteOpt := func(bapp *BaseApp) { bapp.SetAnteHandler(anteHandlerTxTest(t, capKey1, anteKey)) }
	distOpt := func(bapp *BaseApp) { bapp.MountStores(distKey1) }
	mockListener1 := NewMockABCIListener("lis_1")
	mockListener2 := NewMockABCIListener("lis_2")
	streamingManager := storetypes.StreamingManager{ABCIListeners: []storetypes.ABCIListener{&mockListener1, &mockListener2}}
	streamingManagerOpt := func(bapp *BaseApp) { bapp.SetStreamingManager(streamingManager) }
	addListenerOpt := func(bapp *BaseApp) { bapp.CommitMultiStore().AddListeners([]storetypes.StoreKey{distKey1}) }
	suite := newBaseAppSuite(t, anteOpt, distOpt, streamingManagerOpt, addListenerOpt)

	suite.baseApp.InitChain(abci.RequestInitChain{
		ConsensusParams: &abci.ConsensusParams{},
	})

	deliverKey := []byte("deliver-key")
	baseapptestutil.RegisterCounterServer(suite.baseApp.MsgServiceRouter(), CounterServerImpl{t, capKey1, deliverKey})

	nBlocks := 3
	txPerHeight := 5

	for blockN := 0; blockN < nBlocks; blockN++ {
		header := tmproto.Header{Height: int64(blockN) + 1}
		suite.baseApp.BeginBlock(abci.RequestBeginBlock{Header: header})
		var expectedChangeSet []*storetypes.StoreKVPair

		for i := 0; i < txPerHeight; i++ {
			counter := int64(blockN*txPerHeight + i)
			tx := newWrappedTxCounter(suite.txConfig, counter, counter)

			txBytes, err := suite.txConfig.TxEncoder()(tx)
			require.NoError(t, err)

			sKey := []byte(fmt.Sprintf("distKey%d", i))
			sVal := []byte(fmt.Sprintf("distVal%d", i))
			deliverCtx := getDeliverStateCtx(suite.baseApp)
			store := deliverCtx.KVStore(distKey1)
			store.Set(sKey, sVal)

			expectedChangeSet = append(expectedChangeSet, &storetypes.StoreKVPair{
				StoreKey: distKey1.Name(),
				Delete:   false,
				Key:      sKey,
				Value:    sVal,
			})

			res := suite.baseApp.DeliverTx(abci.RequestDeliverTx{Tx: txBytes})
			require.True(t, res.IsOK(), fmt.Sprintf("%v", res))

			events := res.GetEvents()
			require.Len(t, events, 3, "should contain ante handler, message type and counter events respectively")
			require.Equal(t, sdk.MarkEventsToIndex(counterEvent("ante_handler", counter).ToABCIEvents(), map[string]struct{}{})[0], events[0], "ante handler event")
			require.Equal(t, sdk.MarkEventsToIndex(counterEvent(sdk.EventTypeMessage, counter).ToABCIEvents(), map[string]struct{}{})[0], events[2], "msg handler update counter event")
		}

		suite.baseApp.EndBlock(abci.RequestEndBlock{})
		suite.baseApp.Commit()

		require.Equal(t, expectedChangeSet, mockListener1.ChangeSet, "should contain the same changeSet")
		require.Equal(t, expectedChangeSet, mockListener2.ChangeSet, "should contain the same changeSet")
	}
}

func Test_Ctx_with_StreamingManager(t *testing.T) {
	mockListener1 := NewMockABCIListener("lis_1")
	mockListener2 := NewMockABCIListener("lis_2")
	listeners := []storetypes.ABCIListener{&mockListener1, &mockListener2}
	streamingManager := storetypes.StreamingManager{ABCIListeners: listeners, StopNodeOnErr: true}
	streamingManagerOpt := func(bapp *BaseApp) { bapp.SetStreamingManager(streamingManager) }
	addListenerOpt := func(bapp *BaseApp) { bapp.CommitMultiStore().AddListeners([]storetypes.StoreKey{distKey1}) }
	suite := newBaseAppSuite(t, streamingManagerOpt, addListenerOpt)

	suite.baseApp.InitChain(abci.RequestInitChain{
		ConsensusParams: &abci.ConsensusParams{},
	})

	ctx := getDeliverStateCtx(suite.baseApp)
	sm := ctx.StreamingManager()
	require.NotNil(t, sm, fmt.Sprintf("nil StreamingManager: %v", sm))
	require.Equal(t, listeners, sm.ABCIListeners, fmt.Sprintf("should contain same listeners: %v", listeners))
	require.Equal(t, true, sm.StopNodeOnErr, "should contain StopNodeOnErr = true")

	nBlocks := 2

	for blockN := 0; blockN < nBlocks; blockN++ {
		header := tmproto.Header{Height: int64(blockN) + 1}
		suite.baseApp.BeginBlock(abci.RequestBeginBlock{Header: header})

		ctx := getDeliverStateCtx(suite.baseApp)
		sm := ctx.StreamingManager()
		require.NotNil(t, sm, fmt.Sprintf("nil StreamingManager: %v", sm))
		require.Equal(t, listeners, sm.ABCIListeners, fmt.Sprintf("should contain same listeners: %v", listeners))
		require.Equal(t, true, sm.StopNodeOnErr, "should contain StopNodeOnErr = true")

		suite.baseApp.EndBlock(abci.RequestEndBlock{})
		suite.baseApp.Commit()
	}
}

func getDeliverStateCtx(app *BaseApp) sdk.Context {
	v := reflect.ValueOf(app).Elem()
	f := v.FieldByName("deliverState")
	rf := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	return rf.MethodByName("Context").Call(nil)[0].Interface().(sdk.Context)
}
