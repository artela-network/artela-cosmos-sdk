// this would be generated

package bankv1beta1

import basev1beta1 "github.com/cosmos/cosmos-sdk/x/router/api/cosmos/base/v1beta1"

type ConvertibleMsgSend struct {
	FromAddress string
	ToAddress   string
	Amount      []struct {
		Denom  string
		Amount string
	}
}

func (msg *MsgSend) ToConvertible() *ConvertibleMsgSend {
	up := &ConvertibleMsgSend{
		FromAddress: msg.FromAddress,
		ToAddress:   msg.ToAddress,
	}
	var amounts []struct{ Denom, Amount string }
	for _, a := range msg.Amount {
		amounts = append(amounts, struct{ Denom, Amount string }{a.Denom, a.Amount})
	}
	return up
}

func (c *ConvertibleMsgSend) ToMsgSend() (msg *MsgSend) {
	msg.FromAddress = c.FromAddress
	msg.ToAddress = c.ToAddress
	var amounts []*basev1beta1.Coin
	for _, a := range c.Amount {
		amounts = append(amounts, &basev1beta1.Coin{Denom: a.Denom, Amount: a.Amount})
	}
	msg.Amount = amounts
	return msg
}

type ConvertibleMsgSendResponse struct{}

func (res *MsgSendResponse) ToConvertible() *ConvertibleMsgSendResponse {
	up := &ConvertibleMsgSendResponse{}
	return up
}

func (c *ConvertibleMsgSendResponse) ToMsgSendResponse() (res *MsgSendResponse) {
	return res
}
