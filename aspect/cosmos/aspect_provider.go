package cosmos

import (
	"math/big"

	"github.com/artela-network/aspect-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AspectCosmosProvider interface {
	types.AspectProvider

	FilterAspectTx(tx sdk.Msg) bool
	CreateTxPointRequest(sdkCtx sdk.Context, msg sdk.Msg, txIndex int64, baseFee *big.Int, innerTx *types.EthStackTransaction) (*types.EthTxAspect, error)
	CreateBlockPointRequest(sdkCtx sdk.Context) *types.EthBlockAspect
}
