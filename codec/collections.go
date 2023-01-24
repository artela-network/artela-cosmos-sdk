package codec

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/gogoproto/proto"
)

type collValue[T any, PT interface {
	*T
	proto.Message
}] struct {
	cdc Codec
}

func (c collValue[T, PT]) Encode(value T) ([]byte, error) {
	return c.cdc.Marshal(PT(&value))
}

func (c collValue[T, PT]) Decode(b []byte) (T, error) {
	x := new(T)
	err := c.cdc.Unmarshal(b, PT(x))
	return *x, err
}

func (c collValue[T, PT]) EncodeJSON(value T) ([]byte, error) {
	return c.cdc.MarshalJSON(PT(&value))
}

func (c collValue[T, PT]) DecodeJSON(b []byte) (T, error) {
	x := new(T)
	err := c.cdc.UnmarshalJSON(b, PT(x))
	return *x, err
}

func (c collValue[T, PT]) Stringify(value T) string {
	return PT(&value).String()
}

func (c collValue[T, PT]) ValueType() string {
	return "proto"
}

func CollValue[T any, PT interface {
	*T
	proto.Message
}](bCdc BinaryCodec) collections.ValueCodec[T] {
	cdc := bCdc.(Codec)
	return collValue[T, PT]{
		cdc,
	}
}
