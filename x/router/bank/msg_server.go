package bank

import (
	"context"
	router "cosmossdk.io/x/router/api/cosmos/bank/v1beta1"
	"fmt"
	bank "github.com/cosmos/cosmos-sdk/x/bank/v2/api/cosmos/bank/v1beta1"
)

type msgServer struct {
	router.UnsafeMsgServer
	bankServer bank.MsgServer
}

// boiler plate glue code, default implementation and happy path.
//
// this implementation surfaces state machine breaking (non proto-breaking) changes as compile time errors.
//
// an alternative approach could be a deep copy from proto.Message -> proto.Message which errors on unknown fields
// with a test suite exercising every message type in the SDK.

func (m msgServer) Send(ctx context.Context, send *router.MsgSend) (*router.MsgSendResponse, error) {
	c := send.ToPrimitive()

	// c2 := (*bank.PrimitiveMsgSend)(c)
	// bank ahead of router, I've decided that a default empty value is OK.

	if send.Memo != "" {
		return nil, fmt.Errorf("memo not supported, bank is held back")
	}

	c2 := &bank.PrimitiveMsgSend{
		FromAddress: c.FromAddress,
		ToAddress:   c.ToAddress,
		Amount:      c.Amount,
		Memo:        "",
	}

	send2 := c2.FromPrimitive()
	res, err := m.bankServer.Send(ctx, send2)
	if err != nil {
		return nil, err
	}
	cRes := res.ToPrimitive()
	res2 := (*router.PrimitiveMsgSendResponse)(cRes)

	return res2.FromPrimitive(), nil
}

func (m msgServer) MultiSend(ctx context.Context, send *router.MsgMultiSend) (*router.MsgMultiSendResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m msgServer) UpdateParams(ctx context.Context, params *router.MsgUpdateParams) (*router.MsgUpdateParamsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m msgServer) SetSendEnabled(ctx context.Context, enabled *router.MsgSetSendEnabled) (*router.MsgSetSendEnabledResponse, error) {
	//TODO implement me
	panic("implement me")
}

var _ router.MsgServer = msgServer{}
