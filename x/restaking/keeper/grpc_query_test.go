package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (suite *KeeperTestSuite) TestQuerier_OperatorParams() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		request   *types.QueryOperatorParamsRequest
		shouldErr bool
		expParams types.OperatorParams
	}{
		{
			name:      "invalid request returns error",
			request:   nil,
			shouldErr: true,
		},
		{
			name:      "invalid operator id returns error",
			request:   types.NewQueryOperatorParamsRequest(0),
			shouldErr: true,
		},
		{
			name:      "not found operator params returns default value",
			request:   types.NewQueryOperatorParamsRequest(1),
			shouldErr: false,
			expParams: types.DefaultOperatorParams(),
		},
		{
			name: "found operator params are returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperatorParams(ctx, 1, types.NewOperatorParams(
					sdkmath.LegacyNewDecWithPrec(1, 2),
					[]uint32{1, 2},
				))
			},
			request:   types.NewQueryOperatorParamsRequest(1),
			shouldErr: false,
			expParams: types.NewOperatorParams(
				sdkmath.LegacyNewDecWithPrec(1, 2),
				[]uint32{1, 2},
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.OperatorParams(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expParams, res.OperatorParams)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_ServiceParams() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		request   *types.QueryServiceParamsRequest
		shouldErr bool
		expParams types.ServiceParams
	}{
		{
			name:      "invalid request returns error",
			request:   nil,
			shouldErr: true,
		},
		{
			name:      "invalid service id returns error",
			request:   types.NewQueryServiceParamsRequest(0),
			shouldErr: true,
		},
		{
			name:      "not found service params returns default value",
			request:   types.NewQueryServiceParamsRequest(1),
			shouldErr: false,
			expParams: types.DefaultServiceParams(),
		},
		{
			name: "found service params are returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveServiceParams(ctx, 1, types.NewServiceParams(
					sdkmath.LegacyNewDecWithPrec(1, 2),
					[]uint32{1, 2},
					nil,
				))
			},
			request:   types.NewQueryServiceParamsRequest(1),
			shouldErr: false,
			expParams: types.NewServiceParams(
				sdkmath.LegacyNewDecWithPrec(1, 2),
				[]uint32{1, 2},
				nil,
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.ServiceParams(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expParams, res.ServiceParams)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_PoolDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryPoolDelegationsRequest
		shouldErr      bool
		expDelegations []types.DelegationResponse
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				})
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryPoolDelegationsRequest(1, nil),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewPoolDelegation(
						1,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
					),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				),
				types.NewDelegationResponse(
					types.NewPoolDelegation(
						1,
						"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
					),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				})
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryPoolDelegationsRequest(1, &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewPoolDelegation(
						1,
						"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
					),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.PoolDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegations, res.Delegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_PoolDelegation() {
	testCases := []struct {
		name          string
		store         func(ctx sdk.Context)
		request       *types.QueryPoolDelegationRequest
		shouldErr     bool
		expDelegation types.DelegationResponse
	}{
		{
			name: "not found delegation returns error",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				})
				suite.Require().NoError(err)
			},
			request: types.NewQueryPoolDelegationRequest(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "found delegation is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				})
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryPoolDelegationRequest(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			expDelegation: types.NewDelegationResponse(
				types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				),
				sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.PoolDelegation(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegation, res.Delegation)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_PoolUnbondingDelegations() {
	testCases := []struct {
		name                    string
		store                   func(ctx sdk.Context)
		request                 *types.QueryPoolUnbondingDelegationsRequest
		shouldErr               bool
		expUnbondingDelegations []types.UnbondingDelegation
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				pool := poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				}
				err := suite.pk.SavePool(ctx, pool)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &pool,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &pool,
						Delegator:                "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request:   types.NewQueryPoolUnbondingDelegationsRequest(1, nil),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewPoolUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					1,
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				),
				types.NewPoolUnbondingDelegation(
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					1,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
					2,
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				pool := poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				}
				err := suite.pk.SavePool(ctx, pool)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &pool,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &pool,
						Delegator:                "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request: types.NewQueryPoolUnbondingDelegationsRequest(1, &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewPoolUnbondingDelegation(
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					1,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
					2,
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.PoolUnbondingDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbondingDelegations, res.UnbondingDelegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_PoolUnbondingDelegation() {
	testCases := []struct {
		name              string
		store             func(ctx sdk.Context)
		request           *types.QueryPoolUnbondingDelegationRequest
		shouldErr         bool
		expUnbondingEntry types.UnbondingDelegation
	}{
		{
			name: "not found unbonding delegation returns error",
			store: func(ctx sdk.Context) {
				pool := poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				}
				err := suite.pk.SavePool(ctx, pool)
				suite.Require().NoError(err)
			},
			request: types.NewQueryPoolUnbondingDelegationRequest(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "found unbonding delegation is returned properly",
			store: func(ctx sdk.Context) {
				pool := poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				}
				err := suite.pk.SavePool(ctx, pool)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &pool,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)
			},
			request: types.NewQueryPoolUnbondingDelegationRequest(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			expUnbondingEntry: types.NewPoolUnbondingDelegation(
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				1,
				10,
				time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
				sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				1,
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.PoolUnbondingDelegation(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbondingEntry, res.UnbondingDelegation)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_OperatorDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryOperatorDelegationsRequest
		shouldErr      bool
		expDelegations []types.DelegationResponse
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryOperatorDelegationsRequest(1, nil),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewOperatorDelegation(
						1,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(50)),
					),
				),
				types.NewDelegationResponse(
					types.NewOperatorDelegation(
						1,
						"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryOperatorDelegationsRequest(1, &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewOperatorDelegation(
						1,
						"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.OperatorDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegations, res.Delegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_OperatorDelegation() {
	testCases := []struct {
		name          string
		store         func(ctx sdk.Context)
		request       *types.QueryOperatorDelegationRequest
		shouldErr     bool
		expDelegation types.DelegationResponse
	}{
		{
			name: "not found delegation returns error",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})
			},
			request: types.NewQueryOperatorDelegationRequest(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "found delegation is returned properly",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryOperatorDelegationRequest(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: false,
			expDelegation: types.NewDelegationResponse(
				types.NewOperatorDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				),
				sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(100)),
				),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.OperatorDelegation(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegation, res.Delegation)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_OperatorUnbondingDelegations() {
	testCases := []struct {
		name                    string
		store                   func(ctx sdk.Context)
		request                 *types.QueryOperatorUnbondingDelegationsRequest
		shouldErr               bool
		expUnbondingDelegations []types.UnbondingDelegation
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				operator := operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.ok.SaveOperator(ctx, operator)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &operator,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &operator,
						Delegator:                "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request:   types.NewQueryOperatorUnbondingDelegationsRequest(1, nil),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewOperatorUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					1,
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				),
				types.NewOperatorUnbondingDelegation(
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					1,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
					2,
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				operator := operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.ok.SaveOperator(ctx, operator)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &operator,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &operator,
						Delegator:                "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request: types.NewQueryOperatorUnbondingDelegationsRequest(1, &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewOperatorUnbondingDelegation(
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					1,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
					2,
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.OperatorUnbondingDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbondingDelegations, res.UnbondingDelegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_OperatorUnbondingDelegation() {
	testCases := []struct {
		name              string
		store             func(ctx sdk.Context)
		request           *types.QueryOperatorUnbondingDelegationRequest
		shouldErr         bool
		expUnbondingEntry types.UnbondingDelegation
	}{
		{
			name: "not found unbonding delegation returns error",
			store: func(ctx sdk.Context) {
				operator := operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.ok.SaveOperator(ctx, operator)
			},
			request: types.NewQueryOperatorUnbondingDelegationRequest(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "found unbonding delegation is returned properly",
			store: func(ctx sdk.Context) {
				operator := operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.ok.SaveOperator(ctx, operator)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &operator,
						Delegator:                "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request: types.NewQueryOperatorUnbondingDelegationRequest(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: false,
			expUnbondingEntry: types.NewOperatorUnbondingDelegation(
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				1,
				20,
				time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
				sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				1,
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.OperatorUnbondingDelegation(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbondingEntry, res.UnbondingDelegation)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_ServiceDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryServiceDelegationsRequest
		shouldErr      bool
		expDelegations []types.DelegationResponse
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryServiceDelegationsRequest(1, nil),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewServiceDelegation(
						1,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(50)),
					),
				),
				types.NewDelegationResponse(
					types.NewServiceDelegation(
						1,
						"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryServiceDelegationsRequest(1, &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewServiceDelegation(
						1,
						"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
					),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.ServiceDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegations, res.Delegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_ServiceDelegation() {
	testCases := []struct {
		name          string
		store         func(ctx sdk.Context)
		request       *types.QueryServiceDelegationRequest
		shouldErr     bool
		expDelegation types.DelegationResponse
	}{
		{
			name: "not found delegation returns error",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})
			},
			request: types.NewQueryServiceDelegationRequest(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "found delegation is returned properly",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryServiceDelegationRequest(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: false,
			expDelegation: types.NewDelegationResponse(
				types.NewServiceDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				),
				sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(100)),
				),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.ServiceDelegation(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegation, res.Delegation)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_ServiceUnbondingDelegations() {
	testCases := []struct {
		name                    string
		store                   func(ctx sdk.Context)
		request                 *types.QueryServiceUnbondingDelegationsRequest
		shouldErr               bool
		expUnbondingDelegations []types.UnbondingDelegation
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				service := servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.sk.SaveService(ctx, service)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &service,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &service,
						Delegator:                "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request:   types.NewQueryServiceUnbondingDelegationsRequest(1, nil),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewServiceUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					1,
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				),
				types.NewServiceUnbondingDelegation(
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					1,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
					2,
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				service := servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.sk.SaveService(ctx, service)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &service,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &service,
						Delegator:                "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request: types.NewQueryServiceUnbondingDelegationsRequest(1, &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewServiceUnbondingDelegation(
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					1,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
					2,
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.ServiceUnbondingDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbondingDelegations, res.UnbondingDelegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_ServiceUnbondingDelegation() {
	testCases := []struct {
		name              string
		store             func(ctx sdk.Context)
		request           *types.QueryServiceUnbondingDelegationRequest
		shouldErr         bool
		expUnbondingEntry types.UnbondingDelegation
	}{
		{
			name: "not found unbonding delegation returns error",
			store: func(ctx sdk.Context) {
				service := servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.sk.SaveService(ctx, service)
			},
			request: types.NewQueryServiceUnbondingDelegationRequest(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: true,
		},
		{
			name: "found unbonding delegation is returned properly",
			store: func(ctx sdk.Context) {
				service := servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.sk.SaveService(ctx, service)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &service,
						Delegator:                "cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
						BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request: types.NewQueryServiceUnbondingDelegationRequest(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: false,
			expUnbondingEntry: types.NewServiceUnbondingDelegation(
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				1,
				20,
				time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
				sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				1,
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.ServiceUnbondingDelegation(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbondingEntry, res.UnbondingDelegation)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorPoolDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryDelegatorPoolDelegationsRequest
		shouldErr      bool
		expDelegations []types.DelegationResponse
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				})
				suite.Require().NoError(err)

				err = suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              2,
					Denom:           "utia",
					Address:         poolstypes.GetPoolAddress(2).String(),
					Tokens:          sdkmath.NewInt(100),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorPoolDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewPoolDelegation(
						1,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
					),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				),
				types.NewDelegationResponse(
					types.NewPoolDelegation(
						2,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(100))),
					),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(100))),
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				})
				suite.Require().NoError(err)

				err = suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              2,
					Denom:           "utia",
					Address:         poolstypes.GetPoolAddress(2).String(),
					Tokens:          sdkmath.NewInt(100),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryDelegatorPoolDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewPoolDelegation(
						2,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(100))),
					),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(100))),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorPoolDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegations, res.Delegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorPoolUnbondingDelegations() {
	testCases := []struct {
		name                    string
		store                   func(ctx sdk.Context)
		request                 *types.QueryDelegatorPoolUnbondingDelegationsRequest
		shouldErr               bool
		expUnbondingDelegations []types.UnbondingDelegation
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				firstPool := poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				}
				err := suite.pk.SavePool(ctx, firstPool)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &firstPool,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				secondPool := poolstypes.Pool{
					ID:              2,
					Denom:           "utia",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				}
				err = suite.pk.SavePool(ctx, secondPool)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &secondPool,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorPoolUnbondingDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewPoolUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					1,
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				),
				types.NewPoolUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					2,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
					2,
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				firstPool := poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				}
				err := suite.pk.SavePool(ctx, firstPool)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &firstPool,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				secondPool := poolstypes.Pool{
					ID:              2,
					Denom:           "utia",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				}
				err = suite.pk.SavePool(ctx, secondPool)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &secondPool,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewPoolUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request: types.NewQueryDelegatorPoolUnbondingDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Limit:  1,
				Offset: 1,
			}),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewPoolUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					2,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
					2,
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorPoolUnbondingDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbondingDelegations, res.UnbondingDelegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorOperatorDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryDelegatorOperatorDelegationsRequest
		shouldErr      bool
		expDelegations []types.DelegationResponse
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      2,
					Address: operatorstypes.GetOperatorAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorOperatorDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewOperatorDelegation(
						1,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(50)),
					),
				),
				types.NewDelegationResponse(
					types.NewOperatorDelegation(
						2,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(100)),
					),
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      2,
					Address: operatorstypes.GetOperatorAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryDelegatorOperatorDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewOperatorDelegation(
						2,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(100)),
					),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorOperatorDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegations, res.Delegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorOperatorUnbondingDelegations() {
	testCases := []struct {
		name                    string
		store                   func(ctx sdk.Context)
		request                 *types.QueryDelegatorOperatorUnbondingDelegationsRequest
		shouldErr               bool
		expUnbondingDelegations []types.UnbondingDelegation
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				firstOperator := operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.ok.SaveOperator(ctx, firstOperator)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &firstOperator,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				secondOperator := operatorstypes.Operator{
					ID:      2,
					Address: operatorstypes.GetOperatorAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.ok.SaveOperator(ctx, secondOperator)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &secondOperator,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorOperatorUnbondingDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewOperatorUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					1,
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				),
				types.NewOperatorUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					2,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
					2,
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				firstOperator := operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.ok.SaveOperator(ctx, firstOperator)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &firstOperator,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				secondOperator := operatorstypes.Operator{
					ID:      2,
					Address: operatorstypes.GetOperatorAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.ok.SaveOperator(ctx, secondOperator)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &secondOperator,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewOperatorUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request: types.NewQueryDelegatorOperatorUnbondingDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Limit:  1,
				Offset: 1,
			}),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewOperatorUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					2,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
					2,
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorOperatorUnbondingDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbondingDelegations, res.UnbondingDelegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorServiceDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryDelegatorServiceDelegationsRequest
		shouldErr      bool
		expDelegations []types.DelegationResponse
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorServiceDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewServiceDelegation(
						1,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(50)),
					),
				),
				types.NewDelegationResponse(
					types.NewServiceDelegation(
						2,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(100)),
					),
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryDelegatorServiceDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.DelegationResponse{
				types.NewDelegationResponse(
					types.NewServiceDelegation(
						2,
						"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
						),
					),
					sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(100)),
					),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorServiceDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegations, res.Delegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorServiceUnbondingDelegations() {
	testCases := []struct {
		name                    string
		store                   func(ctx sdk.Context)
		request                 *types.QueryDelegatorServiceUnbondingDelegationsRequest
		shouldErr               bool
		expUnbondingDelegations []types.UnbondingDelegation
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				firstService := servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.sk.SaveService(ctx, firstService)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &firstService,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				secondService := servicestypes.Service{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.sk.SaveService(ctx, secondService)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &secondService,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorServiceUnbondingDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewServiceUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					1,
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				),
				types.NewServiceUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					2,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
					2,
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				firstService := servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.sk.SaveService(ctx, firstService)

				_, err := suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &firstService,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
					},
					10,
					time.Date(2024, 1, 1, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
				suite.Require().NoError(err)

				secondService := servicestypes.Service{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				}
				suite.sk.SaveService(ctx, secondService)
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegationEntry(ctx,
					types.UndelegationData{
						Target:                   &secondService,
						Delegator:                "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
						BuildUnbondingDelegation: types.NewServiceUnbondingDelegation,
					},
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
				)
				suite.Require().NoError(err)
			},
			request: types.NewQueryDelegatorServiceUnbondingDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Limit:  1,
				Offset: 1,
			}),
			shouldErr: false,
			expUnbondingDelegations: []types.UnbondingDelegation{
				types.NewServiceUnbondingDelegation(
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					2,
					20,
					time.Date(2024, 1, 2, 12, 00, 00, 00, time.UTC),
					sdk.NewCoins(sdk.NewCoin("utia", sdkmath.NewInt(50))),
					2,
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorServiceUnbondingDelegations(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbondingDelegations, res.UnbondingDelegations)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorPools() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		request   *types.QueryDelegatorPoolsRequest
		shouldErr bool
		expPools  []poolstypes.Pool
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				})
				suite.Require().NoError(err)

				err = suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              2,
					Denom:           "utia",
					Address:         poolstypes.GetPoolAddress(2).String(),
					Tokens:          sdkmath.NewInt(100),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorPoolsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expPools: []poolstypes.Pool{
				{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				},
				{
					ID:              2,
					Denom:           "utia",
					Address:         poolstypes.GetPoolAddress(2).String(),
					Tokens:          sdkmath.NewInt(100),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				},
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				})
				suite.Require().NoError(err)

				err = suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              2,
					Denom:           "utia",
					Address:         poolstypes.GetPoolAddress(2).String(),
					Tokens:          sdkmath.NewInt(100),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryDelegatorPoolsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expPools: []poolstypes.Pool{
				{
					ID:              2,
					Denom:           "utia",
					Address:         poolstypes.GetPoolAddress(2).String(),
					Tokens:          sdkmath.NewInt(100),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorPools(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expPools, res.Pools)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorPool() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		request   *types.QueryDelegatorPoolRequest
		shouldErr bool
		expPool   poolstypes.Pool
	}{
		{
			name:      "non existing pool returns error",
			request:   types.NewQueryDelegatorPoolRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", 1),
			shouldErr: true,
		},
		{
			name: "existing pool is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(150),
					DelegatorShares: sdkmath.LegacyNewDec(150),
				})
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorPoolRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", 1),
			shouldErr: false,
			expPool: poolstypes.Pool{
				ID:              1,
				Denom:           "umilk",
				Address:         poolstypes.GetPoolAddress(1).String(),
				Tokens:          sdkmath.NewInt(150),
				DelegatorShares: sdkmath.LegacyNewDec(150),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorPool(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expPool, res.Pool)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorOperators() {
	testCases := []struct {
		name         string
		store        func(ctx sdk.Context)
		request      *types.QueryDelegatorOperatorsRequest
		shouldErr    bool
		expOperators []operatorstypes.Operator
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      2,
					Address: operatorstypes.GetOperatorAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorOperatorsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expOperators: []operatorstypes.Operator{
				{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				},
				{
					ID:      2,
					Address: operatorstypes.GetOperatorAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				},
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      2,
					Address: operatorstypes.GetOperatorAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryDelegatorOperatorsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expOperators: []operatorstypes.Operator{
				{
					ID:      2,
					Address: operatorstypes.GetOperatorAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(150)),
					),
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorOperators(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expOperators, res.Operators)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorOperator() {
	testCases := []struct {
		name        string
		store       func(ctx sdk.Context)
		request     *types.QueryDelegatorOperatorRequest
		shouldErr   bool
		expOperator operatorstypes.Operator
	}{
		{
			name:      "non existing operator returns error",
			request:   types.NewQueryDelegatorOperatorRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", 1),
			shouldErr: true,
		},
		{
			name: "existing operator is returned properly",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorOperatorRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", 1),
			shouldErr: false,
			expOperator: operatorstypes.Operator{

				ID:      1,
				Address: operatorstypes.GetOperatorAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(150)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(150)),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorOperator(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expOperator, res.Operator)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorServices() {
	testCases := []struct {
		name        string
		store       func(ctx sdk.Context)
		request     *types.QueryDelegatorServicesRequest
		shouldErr   bool
		expServices []servicestypes.Service
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorServicesRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expServices: []servicestypes.Service{
				{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				},
				{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(150)),
					),
				},
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)
			},
			request: types.NewQueryDelegatorServicesRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expServices: []servicestypes.Service{
				{
					ID:      2,
					Address: servicestypes.GetServiceAddress(2).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("utia", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(150)),
					),
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorServices(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expServices, res.Services)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorService() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		request    *types.QueryDelegatorServiceRequest
		shouldErr  bool
		expService servicestypes.Service
	}{
		{
			name:      "non existing service returns error",
			request:   types.NewQueryDelegatorServiceRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", 1),
			shouldErr: true,
		},
		{
			name: "existing service is returned properly",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(150)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
					),
				})

				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)
			},
			request:   types.NewQueryDelegatorServiceRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", 1),
			shouldErr: false,
			expService: servicestypes.Service{
				ID:      1,
				Address: servicestypes.GetServiceAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(150)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(150)),
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.DelegatorService(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expService, res.Service)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_Params() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		request   *types.QueryParamsRequest
		shouldErr bool
		expParams types.Params
	}{
		{
			name: "params are returned properly",
			store: func(ctx sdk.Context) {
				params := types.NewParams(30 * 24 * time.Hour)
				suite.k.SetParams(ctx, params)
			},
			request:   types.NewQueryParamsRequest(),
			shouldErr: false,
			expParams: types.NewParams(30 * 24 * time.Hour),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			querier := keeper.NewQuerier(suite.k)
			res, err := querier.Params(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expParams, res.Params)
			}
		})
	}
}
