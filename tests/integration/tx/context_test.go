package tx

import (
	"testing"

	"cosmossdk.io/x/tx/signing"
	"github.com/stretchr/testify/require"
)

func TestDefineCustomGetSigners(t *testing.T) {
	customMsg := &testpb.CustomSignedMessage{}
	signers := [][]byte{[]byte("foo")}

	context, err := signing.NewContext(signing.Options{
		AddressCodec:          dummyAddressCodec{},
		ValidatorAddressCodec: dummyValidatorAddressCodec{},
	})
	require.NoError(t, err)

	_, err = context.GetSigners(customMsg)
	// before calling DefineCustomGetSigners, we should get an error
	require.ErrorContains(t, err, "need custom signers function")
	signing.DefineCustomGetSigners(context, func(msg *testpb.CustomSignedMessage) ([][]byte, error) {
		return signers, nil
	})
	// after calling DefineCustomGetSigners, we should get the signers
	gotSigners, err := context.GetSigners(customMsg)
	require.NoError(t, err)
	require.Equal(t, signers, gotSigners)
}
