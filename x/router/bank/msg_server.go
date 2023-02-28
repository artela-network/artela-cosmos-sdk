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

// This could all be generated with a hook for the overrides to handle past module versions

func (m msgServer) Send(ctx context.Context, send *router.MsgSend) (*router.MsgSendResponse, error) {
	c := send.ToPrimitive()

	// we've advanced the API but a user has decided to hold bank back.  The user receives a compile time error
	// and must correct it. In this case the author decides it's OK to throw away the content of the memo field
	// if its empty, but error if it's not empty.

	// c2 := (*bank.PrimitiveMsgSend)(c)

	if c.Memo != "" {
		return nil, fmt.Errorf("memo field is not empty and bank is held back")
	}

	c2 := bank.PrimitiveMsgSend{FromAddress: c.FromAddress,
		ToAddress: c.ToAddress,
		Amount:    c.Amount,
	}

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
