package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_ExportGenesis() {
	testCases := []struct {
		name       string
		setup      func()
		store      func(ctx sdk.Context)
		expGenesis *types.GenesisState
	}{
		{
			name: "pool delegations are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(100),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(200),
				))
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				PoolsDelegations: []types.PoolDelegation{
					types.NewPoolDelegation(
						1,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdkmath.LegacyNewDec(100),
					),
					types.NewPoolDelegation(
						2,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdkmath.LegacyNewDec(200),
					),
				},
			},
		},
		{
			name: "service delegations are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.k.SaveServiceDelegation(ctx, types.NewServiceDelegation(
					2,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(50)),
					),
				))
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				ServicesDelegations: []types.ServiceDelegation{
					types.NewServiceDelegation(
						1,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
						),
					),
					types.NewServiceDelegation(
						2,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(50)),
						),
					),
				},
			},
		},
		{
			name: "operators delegations are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					2,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(50)),
					),
				))
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				OperatorsDelegations: []types.OperatorDelegation{
					types.NewOperatorDelegation(
						1,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
						),
					),
					types.NewOperatorDelegation(
						2,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(50)),
						),
					),
				},
			},
		},
		{
			name: "params are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.NewParams(
					30*24*time.Hour,
				))
			},
			expGenesis: &types.GenesisState{
				Params: types.NewParams(
					30 * 24 * time.Hour,
				),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			genesis := suite.k.ExportGenesis(ctx)
			suite.Require().Equal(tc.expGenesis, genesis)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_InitGenesis() {
	testCases := []struct {
		name    string
		genesis *types.GenesisState
		check   func(ctx sdk.Context)
	}{
		{
			name: "pool delegations are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				PoolsDelegations: []types.PoolDelegation{
					types.NewPoolDelegation(
						1,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdkmath.LegacyNewDec(100),
					),
					types.NewPoolDelegation(
						2,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdkmath.LegacyNewDec(200),
					),
				},
			},
			check: func(ctx sdk.Context) {
				_, pool1DelegationFound := suite.k.GetPoolDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(pool1DelegationFound)

				_, pool2DelegationFound := suite.k.GetPoolDelegation(ctx, 2, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(pool2DelegationFound)
			},
		},
		{
			name: "services delegations are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				ServicesDelegations: []types.ServiceDelegation{
					types.NewServiceDelegation(
						1,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
						),
					),
					types.NewServiceDelegation(
						2,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(50)),
						),
					),
				},
			},
			check: func(ctx sdk.Context) {
				_, service1DelegationFound := suite.k.GetServiceDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(service1DelegationFound)

				_, service2DelegationFound := suite.k.GetServiceDelegation(ctx, 2, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(service2DelegationFound)
			},
		},
		{
			name: "operators delegations are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				OperatorsDelegations: []types.OperatorDelegation{
					types.NewOperatorDelegation(
						1,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
						),
					),
					types.NewOperatorDelegation(
						2,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(
							sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(50)),
						),
					),
				},
			},
			check: func(ctx sdk.Context) {
				_, operator1DelegationFound := suite.k.GetOperatorDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(operator1DelegationFound)

				_, operator2DelegationFound := suite.k.GetOperatorDelegation(ctx, 2, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(operator2DelegationFound)
			},
		},
		{
			name: "params are stored properly",
			genesis: &types.GenesisState{
				Params: types.NewParams(
					30 * 24 * time.Hour,
				),
			},
			check: func(ctx sdk.Context) {
				params := suite.k.GetParams(ctx)
				suite.Require().Equal(types.NewParams(
					30*24*time.Hour,
				), params)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()

			suite.k.InitGenesis(ctx, tc.genesis)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}