package testutil

import (
	"fmt"
	"testing"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"github.com/cosmos/gogoproto/proto"
	"gotest.tools/v3/assert"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	servermock "github.com/cosmos/cosmos-sdk/server/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

type IntegrationTestApp struct {
	t                 *testing.T
	Baseapp           *baseapp.BaseApp
	InterfaceRegistry codectypes.InterfaceRegistry
	Ctx               sdk.Context
}

func createIntegrationTestRegistry(msgs ...proto.Message) types.InterfaceRegistry {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterInterface("sdk.Msg",
		(*sdk.Msg)(nil),
		msgs...,
	)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), msgs...)
	fmt.Println("msgs: ", msgs)
	fmt.Println("interface registry: ", interfaceRegistry.ListAllInterfaces())

	return interfaceRegistry
}

func SetupTestApp(t *testing.T, msgs ...proto.Message) *IntegrationTestApp {
	logger := log.NewTestLogger(t)
	db := dbm.NewMemDB()
	interfaceRegistry := createIntegrationTestRegistry(msgs...)

	txConfig := authtx.NewTxConfig(codec.NewProtoCodec(interfaceRegistry), authtx.DefaultSignModes)
	testStore := storetypes.NewKVStoreKey("test")

	bApp := baseapp.NewBaseApp(t.Name(), logger, db, txConfig.TxDecoder())
	bApp.MountStores(testStore)
	bApp.SetInitChainer(servermock.InitChainer(testStore))

	router := baseapp.NewMsgServiceRouter()
	router.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetMsgServiceRouter(router)

	assert.NilError(t, bApp.LoadLatestVersion())
	// testdata.RegisterMsgServer(bApp.MsgServiceRouter(), )

	ctx := bApp.NewContext(true, cmtproto.Header{})
	return &IntegrationTestApp{
		t:                 t,
		Baseapp:           bApp,
		InterfaceRegistry: interfaceRegistry,
		Ctx:               ctx,
	}
}

func (app *IntegrationTestApp) ExecMsgs(ctx sdk.Context, msgs ...sdk.Msg) ([]sdk.Result, error) {
	results := make([]sdk.Result, len(msgs))
	for i, msg := range msgs {
		fmt.Println("msg in loop : ", msg)
		fmt.Printf("sdk.MsgTypeURL(msg): %v\n", sdk.MsgTypeURL(msg))
		msgServiceHandler := app.Baseapp.MsgServiceRouter().HandlerByTypeURL(sdk.MsgTypeURL(msg))
		// app.Baseapp.MsgServiceRouter()
		if msgServiceHandler == nil {
			return nil, fmt.Errorf("handler not found can't route message %q", msg)
		}
		msgResult, err := msgServiceHandler(ctx, msg)
		if err != nil {
			return nil, errorsmod.Wrapf(err, "message %s at position %d", sdk.MsgTypeURL(msg), i)
		}
		// Handler should always return non-nil sdk.Result.
		if msgResult == nil {
			return nil, fmt.Errorf("got nil sdk.Result for message %q at position %d", msg, i)
		}

		if len(msgResult.MsgResponses) != 0 {
			msgResponse := msgResult.MsgResponses[0]
			if msgResponse == nil {
				return nil, fmt.Errorf("got nil Msg response at index %d for msg %s", i, msg)
			}

			results[i] = *msgResult
		}
		
	}
	return results, nil
}
