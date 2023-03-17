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

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	"cosmossdk.io/x/evidence/testutil"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
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

	app, err := simtestutil.Setup(testutil.AppConfig, &evidenceKeeper)
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
	app := integration.SetupTestApp(t, &types.MsgSubmitEvidence{})

	evidence := &types.Equivocation{
		Height:           0,
		Time:             time.Unix(0, 0),
		Power:            100,
		ConsensusAddress: sdk.ConsAddress(ed25519.GenPrivKey().PubKey().Address()).String(),
	}

	any, err := codectypes.NewAnyWithValue(evidence)
	assert.NilError(t, err)

	msgs, err := app.RunMsgs(app.Ctx, &types.MsgSubmitEvidence{
		Submitter: sdk.AccAddress("test").String(),
		Evidence:  any,
	})
	assert.NilError(t, err)

	fmt.Println(msgs)
}

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
