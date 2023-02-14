package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func cloneAppend(bz []byte, tail []byte) (res []byte) {
	res = make([]byte, len(bz)+len(tail))
	copy(res, bz)
	copy(res[len(bz):], tail)
	return
}

func TestCreateDenomAddressPrefix(t *testing.T) {
	require := require.New(t)

	key := types.CreateDenomAddressPrefix("")
	require.Len(key, len(types.DenomAddressPrefix)+1)
	require.Equal(append(types.DenomAddressPrefix, 0), key)

	key = types.CreateDenomAddressPrefix("abc")
	require.Len(key, len(types.DenomAddressPrefix)+4)
	require.Equal(append(types.DenomAddressPrefix, 'a', 'b', 'c', 0), key)
}
