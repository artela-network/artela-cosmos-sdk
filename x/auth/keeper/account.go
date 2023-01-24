package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAccountWithAddress implements AccountKeeperI.
func (ak AccountKeeper) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) sdk.AccountI {
	acc := ak.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		panic(err)
	}

	return ak.NewAccount(ctx, acc)
}

// NewAccount sets the next account number to a given account interface
func (ak AccountKeeper) NewAccount(ctx sdk.Context, acc sdk.AccountI) sdk.AccountI {
	if err := acc.SetAccountNumber(ak.NextAccountNumber(ctx)); err != nil {
		panic(err)
	}

	return acc
}

// HasAccount implements AccountKeeperI.
func (ak AccountKeeper) HasAccount(ctx sdk.Context, addr sdk.AccAddress) bool {
	has, _ := ak.AccountsState.Has(ctx, addr)
	return has
}

// HasAccountAddressByID checks account address exists by id.
func (ak AccountKeeper) HasAccountAddressByID(ctx sdk.Context, id uint64) bool {
	_, err := ak.AccountsState.Indexes.Number.MatchExact(ctx, id)
	return err == nil
}

// GetAccount implements AccountKeeperI.
func (ak AccountKeeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) sdk.AccountI {
	acc, _ := ak.AccountsState.Get(ctx, addr)
	return acc
}

// GetAccountAddressById returns account address by id.
func (ak AccountKeeper) GetAccountAddressByID(ctx sdk.Context, id uint64) string {
	addr, err := ak.AccountsState.Indexes.Number.MatchExact(ctx, id)
	if err != nil {
		return ""
	}
	return addr.String()
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (ak AccountKeeper) GetAllAccounts(ctx sdk.Context) (accounts []sdk.AccountI) {
	iter, err := ak.AccountsState.Iterate(ctx, nil)
	if err != nil {
		return
	}
	values, err := iter.Values()
	if err != nil {
		return nil
	}
	return values
}

// SetAccount implements AccountKeeperI.
func (ak AccountKeeper) SetAccount(ctx sdk.Context, acc sdk.AccountI) {
	err := ak.AccountsState.Set(ctx, acc.GetAddress(), acc)
	if err != nil {
		panic(err)
	}
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (ak AccountKeeper) RemoveAccount(ctx sdk.Context, acc sdk.AccountI) {
	_ = ak.AccountsState.Remove(ctx, acc.GetAddress())
}
