package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// InitGenesis - Init store state from genesis data
//
// CONTRACT: old coins from the FeeCollectionKeeper need to be transferred through
// a genesis port script to the new fee collector account
func (ak AccountKeeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	if err := ak.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}

	accounts, err := types.UnpackAccounts(data.Accounts)
	if err != nil {
		panic(err)
	}
	accounts = types.SanitizeGenesisAccounts(accounts)

	// Set the accounts and make sure the global account number matches the largest account number (even if zero).
	var lastAccNum *uint64
	for _, acc := range accounts {
		accNum := acc.GetAccountNumber()
		for lastAccNum == nil || *lastAccNum < accNum {
			n := ak.NextAccountNumber(ctx)
			lastAccNum = &n
		}
		ak.SetAccount(ctx, acc)
	}

	ak.GetModuleAccount(ctx, types.FeeCollectorName)
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (ak AccountKeeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := ak.GetParams(ctx)

	var genAccounts types.GenesisAccounts
	iter, err := ak.AccountsState.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		v, err := iter.Value()
		if err != nil {
			panic(err)
		}
		genAccounts = append(genAccounts, v.(types.GenesisAccount))
	}

	return types.NewGenesisState(params, genAccounts)
}
