package integration

import (
	"fmt"
	"testing"

	"cosmossdk.io/log"
	"github.com/cosmos/gogoproto/proto"
	"gotest.tools/v3/assert"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	servermock "github.com/cosmos/cosmos-sdk/server/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

type TestMsg struct {
	Name    string
	Message proto.Message
}

type IntegrationTestApp struct {
	t *testing.T

	Baseapp  *baseapp.BaseApp
	Registry codectypes.InterfaceRegistry
	Logger   log.Logger
	Ctx      sdk.Context
}

func SetupTestApp(t *testing.T, msgs ...proto.Message) *IntegrationTestApp {
	logger := log.NewTestLogger(t)
	db := dbm.NewMemDB()

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	logger.Info("Registering messages", "msgs", msgs)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), msgs...)

	txConfig := authtx.NewTxConfig(codec.NewProtoCodec(interfaceRegistry), authtx.DefaultSignModes)
	testStore := storetypes.NewKVStoreKey("integration")

	bApp := baseapp.NewBaseApp(t.Name(), logger, db, txConfig.TxDecoder())
	bApp.MountStores(testStore)
	bApp.SetInitChainer(servermock.InitChainer(testStore))

	router := baseapp.NewMsgServiceRouter()
	router.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetMsgServiceRouter(router)

	assert.NilError(t, bApp.LoadLatestVersion())

	ctx := bApp.NewContext(true, cmtproto.Header{})

	return &IntegrationTestApp{
		t: t,

		Baseapp:  bApp,
		Registry: interfaceRegistry,
		Logger:   logger,
		Ctx:      ctx,
	}
}

func (app *IntegrationTestApp) RunMsgs(ctx sdk.Context, msgs ...sdk.Msg) ([]*codectypes.Any, error) {
	app.Logger.Info("Running msg", "msgs", msgs)

	var responses []*codectypes.Any
	for i, msg := range msgs {
		handler := app.Baseapp.MsgServiceRouter().Handler(msg)
		if handler == nil {
			return nil, fmt.Errorf("can't route message %+v", msg)
		}

		msgResult, err := handler(ctx, msg)
		if err != nil {
			return nil, fmt.Errorf("failed to execute message; ")
		}

		if len(msgResult.MsgResponses) > 0 {
			msgResponse := msgResult.MsgResponses[0]
			if msgResponse == nil {
				return nil, fmt.Errorf("go nil Msg response at index %d for msg %s", i, msg)
			}

			responses = append(responses, msgResponse)
		}
	}
	return responses, nil
}
