package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func (suite *KeeperTestSuite) TestQuerier_PoolDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryPoolDelegationsRequest
		shouldErr      bool
		expDelegations []types.PoolDelegationResponse
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

				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdkmath.LegacyNewDec(100),
				))
			},
			request:   types.NewQueryPoolDelegationsRequest(1, nil),
			shouldErr: false,
			expDelegations: []types.PoolDelegationResponse{
				types.NewPoolDelegationResponse(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
					sdk.NewCoin("umilk", sdkmath.NewInt(50)),
				),
				types.NewPoolDelegationResponse(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdkmath.LegacyNewDec(100),
					sdk.NewCoin("umilk", sdkmath.NewInt(100)),
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

				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdkmath.LegacyNewDec(100),
				))
			},
			request: types.NewQueryPoolDelegationsRequest(1, &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.PoolDelegationResponse{
				types.NewPoolDelegationResponse(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdkmath.LegacyNewDec(100),
					sdk.NewCoin("umilk", sdkmath.NewInt(100)),
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
		expDelegation types.PoolDelegationResponse
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

				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
			},
			request: types.NewQueryPoolDelegationRequest(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr: false,
			expDelegation: types.NewPoolDelegationResponse(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdkmath.LegacyNewDec(50),
				sdk.NewCoin("umilk", sdkmath.NewInt(50)),
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

func (suite *KeeperTestSuite) TestQuerier_OperatorDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryOperatorDelegationsRequest
		shouldErr      bool
		expDelegations []types.OperatorDelegationResponse
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

				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request:   types.NewQueryOperatorDelegationsRequest(1, nil),
			shouldErr: false,
			expDelegations: []types.OperatorDelegationResponse{
				types.NewOperatorDelegationResponse(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(50)),
					),
				),
				types.NewOperatorDelegationResponse(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
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

				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request: types.NewQueryOperatorDelegationsRequest(1, &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.OperatorDelegationResponse{
				types.NewOperatorDelegationResponse(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
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
		expDelegation types.OperatorDelegationResponse
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

				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request: types.NewQueryOperatorDelegationRequest(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: false,
			expDelegation: types.NewOperatorDelegationResponse(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
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

func (suite *KeeperTestSuite) TestQuerier_ServiceDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryServiceDelegationsRequest
		shouldErr      bool
		expDelegations []types.ServiceDelegationResponse
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

				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request:   types.NewQueryServiceDelegationsRequest(1, nil),
			shouldErr: false,
			expDelegations: []types.ServiceDelegationResponse{
				types.NewServiceDelegationResponse(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(50)),
					),
				),
				types.NewServiceDelegationResponse(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
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

				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request: types.NewQueryServiceDelegationsRequest(1, &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.ServiceDelegationResponse{
				types.NewServiceDelegationResponse(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
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
		expDelegation types.ServiceDelegationResponse
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

				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request: types.NewQueryServiceDelegationRequest(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
			),
			shouldErr: false,
			expDelegation: types.NewServiceDelegationResponse(
				1,
				"cosmos1d03wa9qd8flfjtvldndw5csv94tvg5hzfcmcgn",
				sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
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

func (suite *KeeperTestSuite) TestQuerier_DelegatorPoolDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryDelegatorPoolDelegationsRequest
		shouldErr      bool
		expDelegations []types.PoolDelegationResponse
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

				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				))
			},
			request:   types.NewQueryDelegatorPoolDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expDelegations: []types.PoolDelegationResponse{
				types.NewPoolDelegationResponse(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
					sdk.NewCoin("umilk", sdkmath.NewInt(50)),
				),
				types.NewPoolDelegationResponse(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
					sdk.NewCoin("utia", sdkmath.NewInt(100)),
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

				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				))
			},
			request: types.NewQueryDelegatorPoolDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.PoolDelegationResponse{
				types.NewPoolDelegationResponse(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
					sdk.NewCoin("utia", sdkmath.NewInt(100)),
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

func (suite *KeeperTestSuite) TestQuerier_DelegatorOperatorDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryDelegatorOperatorDelegationsRequest
		shouldErr      bool
		expDelegations []types.OperatorDelegationResponse
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

				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request:   types.NewQueryDelegatorOperatorDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expDelegations: []types.OperatorDelegationResponse{
				types.NewOperatorDelegationResponse(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(50)),
					),
				),
				types.NewOperatorDelegationResponse(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
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

				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request: types.NewQueryDelegatorOperatorDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.OperatorDelegationResponse{
				types.NewOperatorDelegationResponse(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(100)),
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

func (suite *KeeperTestSuite) TestQuerier_DelegatorServiceDelegations() {
	testCases := []struct {
		name           string
		store          func(ctx sdk.Context)
		request        *types.QueryDelegatorServiceDelegationsRequest
		shouldErr      bool
		expDelegations []types.ServiceDelegationResponse
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

				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request:   types.NewQueryDelegatorServiceDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", nil),
			shouldErr: false,
			expDelegations: []types.ServiceDelegationResponse{
				types.NewServiceDelegationResponse(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(50)),
					),
				),
				types.NewServiceDelegationResponse(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
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

				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
					),
				))
			},
			request: types.NewQueryDelegatorServiceDelegationsRequest("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", &query.PageRequest{
				Offset: 1,
				Limit:  1,
			}),
			shouldErr: false,
			expDelegations: []types.ServiceDelegationResponse{
				types.NewServiceDelegationResponse(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(100)),
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
