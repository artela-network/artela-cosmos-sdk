package tx

import (
	"testing"

	"cosmossdk.io/x/tx/signing"
	"github.com/stretchr/testify/require"

	gogopb "github.com/cosmos/cosmos-sdk/tests/integration/aminojson/testdata/gogo/testpb"
	pulsarpb "github.com/cosmos/cosmos-sdk/tests/integration/aminojson/testdata/pulsar/testpb"
)

func TestDefineCustomGetSigners(t *testing.T) {
	pulsarMsg := &pulsarpb.Streng{}
	signers := [][]byte{[]byte("foo")}

	context, err := signing.NewContext(signing.Options{
		AddressCodec:          dummyAddressCodec{},
		ValidatorAddressCodec: dummyValidatorAddressCodec{},
	})
	require.NoError(t, err)

	_, err = context.GetSigners(pulsarMsg)
	// before calling DefineCustomGetSigners, we should get an error
	require.ErrorContains(t, err, "use DefineCustomGetSigners")
	signing.DefineCustomGetSigners(context, func(msg *pulsarpb.Streng) ([][]byte, error) {
		return signers, nil
	})
	// after calling DefineCustomGetSigners, we should get the signers
	gotSigners, err := context.GetSigners(pulsarMsg)
	require.NoError(t, err)
	require.Equal(t, signers, gotSigners)

	// Same test as above, but call DefineCustomGetSigners with a gogo proto message
	gogoMsg := &gogopb.Streng{}
	context, err = signing.NewContext(signing.Options{
		AddressCodec:          dummyAddressCodec{},
		ValidatorAddressCodec: dummyValidatorAddressCodec{},
	})
	require.NoError(t, err)
	_, err = context.GetSigners(pulsarMsg)
	// before calling DefineCustomGetSigners, we should get an error
	require.ErrorContains(t, err, "use DefineCustomGetSigners")
	signing.DefineCustomGetSigners(context, func(msg *gogopb.Streng) ([][]byte, error) {
		return signers, nil
	})
	// after calling DefineCustomGetSigners, we should get the signers
	gotSigners, err = context.GetSigners(gogoMsg)
	require.NoError(t, err)
	require.Equal(t, signers, gotSigners)
}
