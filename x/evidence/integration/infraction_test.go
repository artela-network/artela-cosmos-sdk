package integration_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/x/evidence/exported"
	"cosmossdk.io/x/evidence/keeper"
	"cosmossdk.io/x/evidence/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"gotest.tools/v3/assert"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	evidencetestutil "cosmossdk.io/x/evidence/testutil"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	pubkeys = []cryptotypes.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52"),
	}

	valAddresses = []sdk.ValAddress{
		sdk.ValAddress(pubkeys[0].Address()),
		sdk.ValAddress(pubkeys[1].Address()),
		sdk.ValAddress(pubkeys[2].Address()),
	}

	// The default power validators are initialized to have within tests
	initAmt   = sdk.TokensFromConsensusPower(200, sdk.DefaultPowerReduction)
	initCoins = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt))
)

type fixture struct {
	ctx sdk.Context
	app *runtime.App

	evidenceKeeper keeper.Keeper
}

func initFixture(t assert.TestingT) *fixture {
	f := &fixture{}
	var evidenceKeeper keeper.Keeper

	app, err := simtestutil.Setup(evidencetestutil.AppConfig, &evidenceKeeper)
	assert.NilError(t, err)

	router := types.NewRouter()
	router = router.AddRoute(types.RouteEquivocation, testEquivocationHandler(evidenceKeeper))
	evidenceKeeper.SetRouter(router)

	f.ctx = app.BaseApp.NewContext(false, cmtproto.Header{Height: 1})
	f.app = app
	f.evidenceKeeper = evidenceKeeper

	return f
}

func TestHandleDoubleSign(t *testing.T) {
	t.Parallel()
	// f := initFixture(t)
	pk := ed25519.GenPrivKey()
	app := testutil.SetupTestApp(t, &types.Equivocation{}, &types.MsgSubmitEvidence{})

	evidence := &types.Equivocation{
		Height:           1,
		Time:             time.Unix(1, 0),
		Power:            100,
		ConsensusAddress: sdk.ConsAddress(pk.PubKey().Address()).String(),
	}
	evidenceAny, err := codectypes.NewAnyWithValue(evidence)
	assert.NilError(t, err)

	msgs, err := app.ExecMsgs(app.Ctx, &types.MsgSubmitEvidence{
		Submitter: sdk.AccAddress("test").String(),
		Evidence:  evidenceAny,
	})
	assert.NilError(t, err)

	fmt.Println(msgs)

	// ctx := app.Ctx
	// populateValidators(t, f)

	// power := int64(100)
	// stakingParams := f.stakingKeeper.GetParams(ctx)
	// operatorAddr, val := valAddresses[0], pubkeys[0]
	// tstaking := stakingtestutil.NewHelper(t, ctx, f.stakingKeeper)

	// selfDelegation := tstaking.CreateValidatorWithValPower(operatorAddr, val, power, true)
	// // execute end-blocker and verify validator attributes
	// staking.EndBlocker(ctx, f.stakingKeeper)
	// assert.DeepEqual(t,
	// 	f.bankKeeper.GetAllBalances(ctx, sdk.AccAddress(operatorAddr)).String(),
	// 	sdk.NewCoins(sdk.NewCoin(stakingParams.BondDenom, initAmt.Sub(selfDelegation))).String(),
	// )
	// assert.DeepEqual(t, selfDelegation, f.stakingKeeper.Validator(ctx, operatorAddr).GetBondedTokens())

	// // handle a signature to set signing info
	// f.slashingKeeper.HandleValidatorSignature(ctx, val.Address(), selfDelegation.Int64(), true)

	// // double sign less than max age
	// oldTokens := f.stakingKeeper.Validator(ctx, operatorAddr).GetTokens()
	// evidence := abci.RequestBeginBlock{
	// 	ByzantineValidators: []abci.Misbehavior{{
	// 		Validator: abci.Validator{Address: val.Address(), Power: power},
	// 		Type:      abci.MisbehaviorType_DUPLICATE_VOTE,
	// 		Time:      time.Unix(0, 0),
	// 		Height:    0,
	// 	}},
	// }

	// f.evidenceKeeper.BeginBlocker(ctx, evidence)

	// // should be jailed and tombstoned
	// assert.Assert(t, f.stakingKeeper.Validator(ctx, operatorAddr).IsJailed())
	// assert.Assert(t, f.slashingKeeper.IsTombstoned(ctx, sdk.ConsAddress(val.Address())))

	// // tokens should be decreased
	// newTokens := f.stakingKeeper.Validator(ctx, operatorAddr).GetTokens()
	// assert.Assert(t, newTokens.LT(oldTokens))

	// // submit duplicate evidence
	// f.evidenceKeeper.BeginBlocker(ctx, evidence)

	// // tokens should be the same (capped slash)
	// assert.Assert(t, f.stakingKeeper.Validator(ctx, operatorAddr).GetTokens().Equal(newTokens))

	// // jump to past the unbonding period
	// ctx = ctx.WithBlockTime(time.Unix(1, 0).Add(stakingParams.UnbondingTime))

	// // require we cannot unjail
	// assert.Error(t, f.slashingKeeper.Unjail(ctx, operatorAddr), slashingtypes.ErrValidatorJailed.Error())

	// // require we be able to unbond now
	// ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	// del, _ := f.stakingKeeper.GetDelegation(ctx, sdk.AccAddress(operatorAddr), operatorAddr)
	// validator, _ := f.stakingKeeper.GetValidator(ctx, operatorAddr)
	// totalBond := validator.TokensFromShares(del.GetShares()).TruncateInt()
	// tstaking.Ctx = ctx
	// tstaking.Denom = stakingParams.BondDenom
	// tstaking.Undelegate(sdk.AccAddress(operatorAddr), operatorAddr, totalBond, true)

	// // query evidence from store
	// evidences := f.evidenceKeeper.GetAllEvidence(ctx)
	// assert.Assert(t, len(evidences) == 1)
}

// func populateValidators(t assert.TestingT, f *fixture) {
// 	// add accounts and set total supply
// 	totalSupplyAmt := initAmt.MulRaw(int64(len(valAddresses)))
// 	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, totalSupplyAmt))
// 	assert.NilError(t, f.bankKeeper.MintCoins(f.ctx, minttypes.ModuleName, totalSupply))

// 	for _, addr := range valAddresses {
// 		assert.NilError(t, f.bankKeeper.SendCoinsFromModuleToAccount(f.ctx, minttypes.ModuleName, (sdk.AccAddress)(addr), initCoins))
// 	}
// }

func newPubKey(pk string) (res cryptotypes.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}

	pubkey := &ed25519.PubKey{Key: pkBytes}

	return pubkey
}

func testEquivocationHandler(_ interface{}) types.Handler {
	return func(ctx sdk.Context, e exported.Evidence) error {
		if err := e.ValidateBasic(); err != nil {
			return err
		}

		ee, ok := e.(*types.Equivocation)
		if !ok {
			return fmt.Errorf("unexpected evidence type: %T", e)
		}
		if ee.Height%2 == 0 {
			return fmt.Errorf("unexpected even evidence height: %d", ee.Height)
		}

		return nil
	}
}
