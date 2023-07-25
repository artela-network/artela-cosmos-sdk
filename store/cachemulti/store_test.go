package cachemulti

import (
	"fmt"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store/cachekv"
	"cosmossdk.io/store/iavl"
	"cosmossdk.io/store/types"
	dbm "github.com/cosmos/cosmos-db"
	iavltree "github.com/cosmos/iavl"
	"github.com/stretchr/testify/require"
)

func TestStoreGetKVStore(t *testing.T) {
	require := require.New(t)

	s := Store{stores: map[types.StoreKey]types.CacheWrap{}}
	key := types.NewKVStoreKey("abc")
	errMsg := fmt.Sprintf("kv store with key %v has not been registered in stores", key)

	require.PanicsWithValue(errMsg,
		func() { s.GetStore(key) })

	require.PanicsWithValue(errMsg,
		func() { s.GetKVStore(key) })
}

func createCacheMultiTree(b *testing.B) types.CacheMultiStore {
	db := dbm.NewMemDB()
	defer db.Close()

	storeKeys := make(map[string]types.StoreKey, 3)
	storeKeys["store1"] = types.NewKVStoreKey("store1")
	storeKeys["store2"] = types.NewKVStoreKey("store2")
	storeKeys["store3"] = types.NewKVStoreKey("store3")

	stores := make(map[types.StoreKey]types.CacheWrapper)
	for _, key := range storeKeys {
		dbStore := dbm.NewPrefixDB(db, []byte(key.Name()))
		tree := iavltree.NewMutableTreeWithOpts(dbStore, 1000, &iavltree.Options{InitialVersion: 1}, false, log.NewNopLogger())
		store := iavl.UnsafeNewStore(tree)
		stores[key] = cachekv.NewStore(store)
	}

	dbCms := dbm.NewPrefixDB(db, []byte("cms"))
	cms := NewStore(dbCms, stores, storeKeys, nil, nil)

	// set new key
	for _, key := range storeKeys {
		store := cms.GetKVStore(key)
		for i := 0; i < 10000; i++ {
			store.Set([]byte(fmt.Sprintf("new_key%d", i)), []byte(fmt.Sprintf("new_value%d", i)))
		}
	}

	return cms
}

func BenchmarkWrite(b *testing.B) {
	b.Run("single-3-stores", func(sub *testing.B) {
		sub.ReportAllocs()
		for i := 0; i < sub.N; i++ {
			cms := createCacheMultiTree(sub)
			sub.StartTimer()
			cms.Write()
			sub.StopTimer()
		}
	})
	b.Run("parallel-3-stores", func(sub *testing.B) {
		sub.ReportAllocs()
		for i := 0; i < sub.N; i++ {
			cms := createCacheMultiTree(sub)
			sub.StartTimer()
			cms.WriteParallel()
			sub.StopTimer()
		}
	})
}
