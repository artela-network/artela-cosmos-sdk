package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	version "github.com/cometbft/cometbft/proto/tendermint/version"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// IsValSetSorted reports whether valset is sorted.
func IsValSetSorted(data []stakingtypes.Validator, powerReduction math.Int) bool {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		if stakingtypes.ValidatorsByVotingPower(data).Less(i, i-1, powerReduction) {
			return false
		}
	}
	return true
}

func (s *KeeperTestSuite) TestHistoricalInfo() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	_, addrVals := createValAddrs(50)

	validators := make([]stakingtypes.Validator, len(addrVals))

	for i, valAddr := range addrVals {
		validators[i] = testutil.NewValidator(s.T(), valAddr, PKs[i])
	}

	hi := stakingtypes.NewHistoricalInfo(ctx.BlockHeader(), validators, keeper.PowerReduction(ctx))
	require.NoError(keeper.SetHistoricalInfo(ctx, 2, &hi))

	recv, err := keeper.GetHistoricalInfo(ctx, 2)
	require.NoError(err, "HistoricalInfo not found after set")
	require.Equal(hi, recv, "HistoricalInfo not equal")
	require.True(IsValSetSorted(recv.Valset, keeper.PowerReduction(ctx)), "HistoricalInfo validators is not sorted")

	require.NoError(keeper.DeleteHistoricalInfo(ctx, 2))

	recv, err = keeper.GetHistoricalInfo(ctx, 2)
	require.ErrorIs(err, stakingtypes.ErrNoHistoricalInfo, "HistoricalInfo found after delete")
	require.Equal(stakingtypes.HistoricalInfo{}, recv, "HistoricalInfo is not empty")
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func CollsMigration(
	ctx sdk.Context,
	storeKey *storetypes.KVStoreKey,
	iterations int,
	writeElem func(int64),
	targetHash string,
) error {
	for i := int64(0); i < int64(iterations); i++ {
		writeElem(i)
	}

	allkvs := []byte{}
	it := ctx.KVStore(storeKey).Iterator(nil, nil)
	defer it.Close()
	for ; it.Valid(); it.Next() {
		kv := append(it.Key(), it.Value()...)
		allkvs = append(allkvs, kv...)
	}

	hash := sha256.Sum256(allkvs)
	if hex.EncodeToString(hash[:]) != targetHash {
		return fmt.Errorf("hashes don't match: %s != %s\n", hex.EncodeToString(hash[:]), targetHash)
	}

	return nil
}

func (s *KeeperTestSuite) TestHistoricalInfoCollMigration() {
	// reset suite
	s.SetupTest()
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	_, addrVals := createValAddrs(50)

	validators := make([]stakingtypes.Validator, len(addrVals))

	for i, valAddr := range addrVals {
		validators[i] = testutil.NewValidator(s.T(), valAddr, PKs[i])
	}

	err := CollsMigration(
		ctx,
		s.key,
		10000,
		func(i int64) {
			header := cmtproto.Header{
				Version: version.Consensus{},
				ChainID: "HelloChain" + strconv.Itoa(int(i)),
				Height:  i,
				Time:    time.Unix(123456789+i, 12345),
				LastBlockId: cmtproto.BlockID{
					Hash: []byte("LastBlockHash"),
					PartSetHeader: cmtproto.PartSetHeader{
						Total: 100 + uint32(i),
						Hash:  []byte("LastBlockPartHash"),
					},
				},
				LastCommitHash:     getTestHash("LastCommitHash", i),
				DataHash:           getTestHash("DataHash", i),
				ValidatorsHash:     getTestHash("ValidatorsHash", i),
				NextValidatorsHash: getTestHash("NextValidatorsHash", i),
				ConsensusHash:      getTestHash("ConsensusHash", i),
				AppHash:            getTestHash("AppHash", i),
				LastResultsHash:    getTestHash("LastResultsHash", i),
				EvidenceHash:       getTestHash("EvidenceHash", i),
				ProposerAddress:    getTestHash("ProposerAddress", i),
			}
			hi := stakingtypes.NewHistoricalInfo(header, validators, keeper.PowerReduction(ctx))
			require.NoError(keeper.SetHistoricalInfo(ctx, i, &hi))
		},
		"a6a907d7c465d0cd2d437f7b8b28f573a4d98452aca0a21c5e7c8997cd9866b8",
	)

	require.NoError(err)
}

func getTestHash(field string, i int64) []byte {
	hash := sha256.Sum256([]byte(field + strconv.Itoa(int(i))))
	return hash[:]
}

func (s *KeeperTestSuite) TestTrackHistoricalInfo() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	_, addrVals := createValAddrs(50)

	// set historical entries in params to 5
	params := stakingtypes.DefaultParams()
	params.HistoricalEntries = 5
	require.NoError(keeper.SetParams(ctx, params))

	// set historical info at 5, 4 which should be pruned
	// and check that it has been stored
	h4 := cmtproto.Header{
		ChainID: "HelloChain",
		Height:  4,
	}
	h5 := cmtproto.Header{
		ChainID: "HelloChain",
		Height:  5,
	}
	valSet := []stakingtypes.Validator{
		testutil.NewValidator(s.T(), addrVals[0], PKs[0]),
		testutil.NewValidator(s.T(), addrVals[1], PKs[1]),
	}
	hi4 := stakingtypes.NewHistoricalInfo(h4, valSet, keeper.PowerReduction(ctx))
	hi5 := stakingtypes.NewHistoricalInfo(h5, valSet, keeper.PowerReduction(ctx))
	require.NoError(keeper.SetHistoricalInfo(ctx, 4, &hi4))
	require.NoError(keeper.SetHistoricalInfo(ctx, 5, &hi5))
	recv, err := keeper.GetHistoricalInfo(ctx, 4)
	require.NoError(err)
	require.Equal(hi4, recv)
	recv, err = keeper.GetHistoricalInfo(ctx, 5)
	require.NoError(err)
	require.Equal(hi5, recv)

	// Set bonded validators in keeper
	val1 := testutil.NewValidator(s.T(), addrVals[2], PKs[2])
	val1.Status = stakingtypes.Bonded // when not bonded, consensus power is Zero
	val1.Tokens = keeper.TokensFromConsensusPower(ctx, 10)
	require.NoError(keeper.SetValidator(ctx, val1))
	require.NoError(keeper.SetLastValidatorPower(ctx, val1.GetOperator(), 10))
	val2 := testutil.NewValidator(s.T(), addrVals[3], PKs[3])
	val1.Status = stakingtypes.Bonded
	val2.Tokens = keeper.TokensFromConsensusPower(ctx, 80)
	require.NoError(keeper.SetValidator(ctx, val2))
	require.NoError(keeper.SetLastValidatorPower(ctx, val2.GetOperator(), 80))

	vals := []stakingtypes.Validator{val1, val2}
	require.True(IsValSetSorted(vals, keeper.PowerReduction(ctx)))

	// Set Header for BeginBlock context
	header := cmtproto.Header{
		ChainID: "HelloChain",
		Height:  10,
	}
	ctx = ctx.WithBlockHeader(header)

	require.NoError(keeper.TrackHistoricalInfo(ctx))

	// Check HistoricalInfo at height 10 is persisted
	expected := stakingtypes.HistoricalInfo{
		Header: header,
		Valset: vals,
	}
	recv, err = keeper.GetHistoricalInfo(ctx, 10)
	require.NoError(err, "GetHistoricalInfo failed after BeginBlock")
	require.Equal(expected, recv, "GetHistoricalInfo returned unexpected result")

	// Check HistoricalInfo at height 5, 4 is pruned
	recv, err = keeper.GetHistoricalInfo(ctx, 4)
	require.ErrorIs(err, stakingtypes.ErrNoHistoricalInfo, "GetHistoricalInfo did not prune earlier height")
	require.Equal(stakingtypes.HistoricalInfo{}, recv, "GetHistoricalInfo at height 4 is not empty after prune")
	recv, err = keeper.GetHistoricalInfo(ctx, 5)
	require.ErrorIs(err, stakingtypes.ErrNoHistoricalInfo, "GetHistoricalInfo did not prune first prune height")
	require.Equal(stakingtypes.HistoricalInfo{}, recv, "GetHistoricalInfo at height 5 is not empty after prune")
}

func (s *KeeperTestSuite) TestGetAllHistoricalInfo() {
	ctx, keeper := s.ctx, s.stakingKeeper
	require := s.Require()

	_, addrVals := createValAddrs(50)

	valSet := []stakingtypes.Validator{
		testutil.NewValidator(s.T(), addrVals[0], PKs[0]),
		testutil.NewValidator(s.T(), addrVals[1], PKs[1]),
	}

	header1 := cmtproto.Header{ChainID: "HelloChain", Height: 9}
	header2 := cmtproto.Header{ChainID: "HelloChain", Height: 10}
	header3 := cmtproto.Header{ChainID: "HelloChain", Height: 11}

	hist1 := stakingtypes.HistoricalInfo{Header: header1, Valset: valSet}
	hist2 := stakingtypes.HistoricalInfo{Header: header2, Valset: valSet}
	hist3 := stakingtypes.HistoricalInfo{Header: header3, Valset: valSet}

	expHistInfos := []stakingtypes.HistoricalInfo{hist1, hist2, hist3}

	for i, hi := range expHistInfos {
		require.NoError(keeper.SetHistoricalInfo(ctx, int64(9+i), &hi)) //nolint:gosec // G601: Implicit memory aliasing in for loop.
	}

	infos, err := keeper.GetAllHistoricalInfo(ctx)
	require.NoError(err)
	require.Equal(expHistInfos, infos)
}
