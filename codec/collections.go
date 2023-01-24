package codec

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/gogoproto/proto"
)

type collInterfaceValue[T proto.Message] struct {
	cdc Codec
}

func (c collInterfaceValue[T]) Encode(value T) ([]byte, error) {
	return c.cdc.MarshalInterface(value)
}

func (c collInterfaceValue[T]) Decode(b []byte) (T, error) {
	var iface T
	err := c.cdc.UnmarshalInterface(b, &iface)
	return iface, err
}

func (c collInterfaceValue[T]) EncodeJSON(value T) ([]byte, error) {
	return c.cdc.MarshalInterfaceJSON(value)
}

func (c collInterfaceValue[T]) DecodeJSON(b []byte) (T, error) {
	var iface T
	err := c.cdc.UnmarshalInterfaceJSON(b, &iface)
	return iface, err
}

func (c collInterfaceValue[T]) Stringify(value T) string {
	return value.String()
}

func (c collInterfaceValue[T]) ValueType() string {
	return "protointerface"
}

func CollInterfaceValue[T proto.Message](bCdc BinaryCodec) collections.ValueCodec[T] {
	cdc := bCdc.(Codec)
	return collInterfaceValue[T]{cdc}
}

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
