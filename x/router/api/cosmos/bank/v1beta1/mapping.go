// this would be generated

package bankv1beta1

import basev1beta1 "github.com/cosmos/cosmos-sdk/x/router/api/cosmos/base/v1beta1"

type PrimitiveMsgSend struct {
	FromAddress string
	ToAddress   string
	Amount      []struct {
		Denom  string
		Amount string
	}
	Memo string
}

func (msg *MsgSend) ToPrimitive() *PrimitiveMsgSend {
	up := &PrimitiveMsgSend{
		FromAddress: msg.FromAddress,
		ToAddress:   msg.ToAddress,
		Memo:        msg.Memo,
	}
	var amounts []struct{ Denom, Amount string }
	for _, a := range msg.Amount {
		amounts = append(amounts, struct{ Denom, Amount string }{a.Denom, a.Amount})
	}
	return up
}

func (c *PrimitiveMsgSend) ToMsgSend() (msg *MsgSend) {
	msg.FromAddress = c.FromAddress
	msg.ToAddress = c.ToAddress
	msg.Memo = c.Memo
	var amounts []*basev1beta1.Coin
	for _, a := range c.Amount {
		amounts = append(amounts, &basev1beta1.Coin{Denom: a.Denom, Amount: a.Amount})
	}
	msg.Amount = amounts
	return msg
}

type PrimitiveMsgSendResponse struct{}

func (res *MsgSendResponse) ToPrimitive() *PrimitiveMsgSendResponse {
	up := &PrimitiveMsgSendResponse{}
	return up
}

func (c *PrimitiveMsgSendResponse) ToMsgSendResponse() (res *MsgSendResponse) {
	return res
}
