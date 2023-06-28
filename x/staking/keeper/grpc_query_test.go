package keeper_test

import (
	gocontext "context"
	"fmt"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/testutil"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (s *KeeperTestSuite) TestGRPCQueryValidator() {
	ctx, keeper, queryClient := s.ctx, s.stakingKeeper, s.queryClient
	require := s.Require()

	validator := testutil.NewValidator(s.T(), sdk.ValAddress(PKs[0].Address().Bytes()), PKs[0])
	require.NoError(keeper.SetValidator(ctx, validator))
	var req *types.QueryValidatorRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryValidatorRequest{}
			},
			false,
		},
		{
			"with valid and not existing address",
			func() {
				req = &types.QueryValidatorRequest{
					ValidatorAddr: "cosmosvaloper15jkng8hytwt22lllv6mw4k89qkqehtahd84ptu",
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryValidatorRequest{ValidatorAddr: validator.OperatorAddress}
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Validator(gocontext.Background(), req)
			if tc.expPass {
				require.NoError(err)
				require.True(validator.Equal(&res.Validator))
			} else {
				require.Error(err)
				require.Nil(res)
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryValidators() {
	ctx, keeper, queryClient := s.ctx, s.stakingKeeper, s.queryClient
	require := s.Require()
	validator := testutil.NewValidator(s.T(), PKs[0].Address().Bytes(), PKs[0])
	require.NoError(keeper.SetValidator(ctx, validator))
	var req *types.QueryValidatorsRequest
	res, err := queryClient.Validators(gocontext.Background(), req)
	require.NoError(err)
	require.Equal(1, len(res.Validators))
}

func (s *KeeperTestSuite) TestGRPCQueryDelegation() {
	ctx, keeper, queryClient, msgServer := s.ctx, s.stakingKeeper, s.queryClient, s.msgServer
	require := s.Require()
	s.execExpectCalls()

	pk := ed25519.GenPrivKey().PubKey()
	require.NotNil(pk)

	comm := types.NewCommissionRates(math.LegacyNewDec(0), math.LegacyNewDec(0), math.LegacyNewDec(0))

	msg, err := types.NewMsgCreateValidator(ValAddr, pk, sdk.NewCoin("stake", sdk.NewInt(10)), types.Description{Moniker: "NewVal"}, comm, math.OneInt())
	require.NoError(err)

	res, err := msgServer.CreateValidator(ctx, msg)
	require.NoError(err)
	require.NotNil(res)

	testCases := []struct {
		name     string
		preRun   func()
		req      types.QueryDelegationRequest
		expecErr bool
		errorMsg string
	}{
		{
			name: "happy path",
			preRun: func() {
				msg := &types.MsgDelegate{
					DelegatorAddress: Addr.String(),
					ValidatorAddress: ValAddr.String(),
					Amount:           sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: keeper.TokensFromConsensusPower(s.ctx, int64(100))},
				}
				_, err := msgServer.Delegate(ctx, msg)
				require.NoError(err)
			},
			req: types.QueryDelegationRequest{
				DelegatorAddr: Addr.String(),
				ValidatorAddr: ValAddr.String(),
			},
			expecErr: false,
			errorMsg: "",
		},
		{
			name:   "invalid addr",
			preRun: func() {},
			req: types.QueryDelegationRequest{
				DelegatorAddr: "invalid-addr",
				ValidatorAddr: ValAddr.String(),
			},
			expecErr: true,
			errorMsg: "decoding bech32 failed",
		},
		{
			name:   "invalid validator addr",
			preRun: func() {},
			req: types.QueryDelegationRequest{
				DelegatorAddr: Addr.String(),
				ValidatorAddr: "invalid-val-addr",
			},
			expecErr: true,
			errorMsg: "decoding bech32 failed",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.preRun()
			res, err := queryClient.Delegation(gocontext.Background(), &tc.req)
			if tc.expecErr {
				require.Error(err, tc.errorMsg)
			} else {
				require.NoError(err)
				require.Equal(res.DelegationResponse.Delegation.DelegatorAddress, Addr.String())
				require.Equal(res.DelegationResponse.Delegation.ValidatorAddress, ValAddr.String())
			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryValidatorDelegation() {
	ctx, keeper, queryClient, msgServer := s.ctx, s.stakingKeeper, s.queryClient, s.msgServer
	require := s.Require()
	s.execExpectCalls()

	pk := ed25519.GenPrivKey().PubKey()
	require.NotNil(pk)

	comm := types.NewCommissionRates(math.LegacyNewDec(0), math.LegacyNewDec(0), math.LegacyNewDec(0))

	msg, err := types.NewMsgCreateValidator(ValAddr, pk, sdk.NewCoin("stake", sdk.NewInt(10)), types.Description{Moniker: "NewVal"}, comm, math.OneInt())
	require.NoError(err)

	res, err := msgServer.CreateValidator(ctx, msg)
	require.NoError(err)
	require.NotNil(res)

	testCases := []struct {
		name     string
		preRun   func()
		req      types.QueryValidatorDelegationsRequest
		expecErr bool
		errorMsg string
	}{
		{
			name: "happy path",
			preRun: func() {
				msg := &types.MsgDelegate{
					DelegatorAddress: Addr.String(),
					ValidatorAddress: ValAddr.String(),
					Amount:           sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: keeper.TokensFromConsensusPower(s.ctx, int64(100))},
				}
				_, err := msgServer.Delegate(ctx, msg)
				require.NoError(err)
			},
			req: types.QueryValidatorDelegationsRequest{
				ValidatorAddr: ValAddr.String(),
			},
			expecErr: false,
			errorMsg: "",
		},
		{
			name:   "invalid validator addr",
			preRun: func() {},
			req: types.QueryValidatorDelegationsRequest{
				ValidatorAddr: "invalid-val-addr",
			},
			expecErr: true,
			errorMsg: "decoding bech32 failed",
		},
		{
			name:   "empty validator addr",
			preRun: func() {},
			req: types.QueryValidatorDelegationsRequest{
				ValidatorAddr: "",
			},
			expecErr: true,
			errorMsg: "decoding bech32 failed",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.preRun()
			res, err := queryClient.ValidatorDelegations(gocontext.Background(), &tc.req)
			if tc.expecErr {
				require.Error(err, tc.errorMsg)
			} else {
				require.NoError(err)
				require.Equal(len(res.DelegationResponses), 1)

			}
		})
	}
}

func (s *KeeperTestSuite) TestGRPCQueryUnbondingDelegations() {
	ctx, keeper, queryClient, msgServer := s.ctx, s.stakingKeeper, s.queryClient, s.msgServer
	require := s.Require()
	s.execExpectCalls()

	pk := ed25519.GenPrivKey().PubKey()
	require.NotNil(pk)

	comm := types.NewCommissionRates(math.LegacyNewDec(0), math.LegacyNewDec(0), math.LegacyNewDec(0))

	msg, err := types.NewMsgCreateValidator(ValAddr, pk, sdk.NewCoin("stake", sdk.NewInt(10)), types.Description{Moniker: "NewVal"}, comm, math.OneInt())
	require.NoError(err)

	res, err := msgServer.CreateValidator(ctx, msg)
	require.NoError(err)
	require.NotNil(res)

	testCases := []struct {
		name     string
		preRun   func()
		req      types.QueryUnbondingDelegationRequest
		expecErr bool
		errorMsg string
	}{
		{
			name: "happy path",
			preRun: func() {
				msg := &types.MsgDelegate{
					DelegatorAddress: Addr.String(),
					ValidatorAddress: ValAddr.String(),
					Amount:           sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: keeper.TokensFromConsensusPower(s.ctx, int64(100))},
				}
				_, err := msgServer.Delegate(ctx, msg)
				require.NoError(err)

				unMsg := &types.MsgUndelegate{
					DelegatorAddress: Addr.String(),
					ValidatorAddress: ValAddr.String(),
					Amount:           sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(50)),
				}
				_, err = msgServer.Undelegate(ctx, unMsg)
				require.NoError(err)
			},
			req: types.QueryUnbondingDelegationRequest{
				DelegatorAddr: Addr.String(),
				ValidatorAddr: ValAddr.String(),
			},
			expecErr: false,
			errorMsg: "",
		},
		{
			name:   "invalid delegator addr",
			preRun: func() {},
			req: types.QueryUnbondingDelegationRequest{
				DelegatorAddr: "invalid delegator addr",
				ValidatorAddr: ValAddr.String(),
			},
			expecErr: true,
			errorMsg: "decoding bech32 failed",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.preRun()
			res, err := queryClient.UnbondingDelegation(gocontext.Background(), &tc.req)
			if tc.expecErr {
				require.Error(err, tc.errorMsg)
			} else {
				require.NoError(err)
				require.Equal(res.Unbond.DelegatorAddress, Addr.String())
				require.Equal(res.Unbond.ValidatorAddress, ValAddr.String())

			}
		})
	}
}
