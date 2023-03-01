// this would be generated

package bankv1beta1

type PrimitiveMsgSend struct {
	FromAddress string
	ToAddress   string
	Amount      []struct {
		Denom  string
		Amount string
	}
}

func (msg *MsgSend) ToPrimitive() *PrimitiveMsgSend {
	up := &PrimitiveMsgSend{
		FromAddress: msg.FromAddress,
		ToAddress:   msg.ToAddress,
	}
	var amounts []struct{ Denom, Amount string }
	for _, a := range msg.Amount {
		amounts = append(amounts, struct{ Denom, Amount string }{a.Denom, a.Amount})
	}
	return up
}

type PrimitiveMsgSendResponse struct{}

func (res *MsgSendResponse) ToPrimitive() *PrimitiveMsgSendResponse {
	up := &PrimitiveMsgSendResponse{}
	return up
}

func (c *PrimitiveMsgSendResponse) FromPrimitive() (res *MsgSendResponse) {
	return res
}
