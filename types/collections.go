package types

import "cosmossdk.io/collections"

var AccAddressKey collections.KeyCodec[AccAddress] = accAddressKey{}

type accAddressKey struct{}

func (a accAddressKey) Encode(buffer []byte, key AccAddress) (int, error) {
	return copy(buffer, key), nil
}

func (a accAddressKey) Decode(buffer []byte) (int, AccAddress, error) {
	return len(buffer), buffer, nil
}

func (a accAddressKey) Size(key AccAddress) int {
	return len(key)
}

func (a accAddressKey) EncodeJSON(value AccAddress) ([]byte, error) {
	return collections.StringKey.EncodeJSON(value.String())
}

func (a accAddressKey) DecodeJSON(b []byte) (AccAddress, error) {
	addrStr, err := collections.StringKey.DecodeJSON(b)
	if err != nil {
		return nil, err
	}
	return AccAddressFromBech32(addrStr)
}

func (a accAddressKey) Stringify(key AccAddress) string {
	return key.String()
}

func (a accAddressKey) KeyType() string {
	//TODO implement me
	panic("implement me")
}

func (a accAddressKey) EncodeNonTerminal(buffer []byte, key AccAddress) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (a accAddressKey) DecodeNonTerminal(buffer []byte) (int, AccAddress, error) {
	//TODO implement me
	panic("implement me")
}

func (a accAddressKey) SizeNonTerminal(key AccAddress) int {
	//TODO implement me
	panic("implement me")
}
