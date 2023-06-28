package baseapp

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	baseapptestutil "github.com/cosmos/cosmos-sdk/baseapp/testutil"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CounterServerImpl struct {
	t          *testing.T
	capKey     storetypes.StoreKey
	deliverKey []byte
}

func (m CounterServerImpl) IncrementCounter(ctx context.Context, msg *baseapptestutil.MsgCounter) (*baseapptestutil.MsgCreateCounterResponse, error) {
	return incrementCounter(ctx, m.t, m.capKey, m.deliverKey, msg)
}

func incrementCounter(ctx context.Context,
	t *testing.T,
	capKey storetypes.StoreKey,
	deliverKey []byte,
	msg sdk.Msg,
) (*baseapptestutil.MsgCreateCounterResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(capKey)

	sdkCtx.GasMeter().ConsumeGas(5, "test")

	var msgCount int64

	switch m := msg.(type) {
	case *baseapptestutil.MsgCounter:
		if m.FailOnHandler {
			return nil, errors.New("message handler failure")
			//return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "message handler failure")
		}
		msgCount = m.Counter
	case *baseapptestutil.MsgCounter2:
		if m.FailOnHandler {
			return nil, errors.New("message handler failure")
			//return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "message handler failure")
		}
		msgCount = m.Counter
	}

	sdkCtx.EventManager().EmitEvents(
		counterEvent(sdk.EventTypeMessage, msgCount),
	)

	_, err := incrementingCounter(t, store, deliverKey, msgCount)
	if err != nil {
		return nil, err
	}

	return &baseapptestutil.MsgCreateCounterResponse{}, nil
}

func parseTxMemo(t *testing.T, tx sdk.Tx) (counter int64, failOnAnte bool) {
	txWithMemo, ok := tx.(sdk.TxWithMemo)
	require.True(t, ok)

	memo := txWithMemo.GetMemo()
	vals, err := url.ParseQuery(memo)
	require.NoError(t, err)

	counter, err = strconv.ParseInt(vals.Get("counter"), 10, 64)
	require.NoError(t, err)

	failOnAnte = vals.Get("failOnAnte") == "true"
	return counter, failOnAnte
}
