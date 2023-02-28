package bank

import (
	"context"
	router "cosmossdk.io/x/router/api/cosmos/bank/v1beta1"
	bank "github.com/cosmos/cosmos-sdk/x/bank/v2/api/cosmos/bank/v1beta1"
)

type msgServer struct {
	router.UnsafeMsgServer
	bankServer bank.MsgServer
}

// This could all be generated with a hook for the overrides to handle past module versions

func (m msgServer) Send(ctx context.Context, send *router.MsgSend) (*router.MsgSendResponse, error) {
	c := send.ToPrimitive()

	c2 := (*bank.PrimitiveMsgSend)(c)

	send2 := c2.ToMsgSend()
	res, err := m.bankServer.Send(ctx, send2)
	if err != nil {
		return nil, err
	}
	cRes := res.ToPrimitive()
	res2 := (*router.PrimitiveMsgSendResponse)(cRes)

	return res2.ToMsgSendResponse(), nil
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
