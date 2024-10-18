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
			name: "operator joined services are exported correctly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				err := suite.k.SaveOperatorJoinedServices(ctx, 1,
					types.NewOperatorJoinedServices([]uint32{1, 2}))
				suite.Require().NoError(err)

				err = suite.k.SaveOperatorJoinedServices(ctx, 2,
					types.NewOperatorJoinedServices([]uint32{3, 4}))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				OperatorsJoinedServices: []types.OperatorJoinedServicesRecord{
					types.NewOperatorJoinedServicesRecord(1,
						types.NewOperatorJoinedServices(
							[]uint32{1, 2},
						)),
					types.NewOperatorJoinedServicesRecord(2,
						types.NewOperatorJoinedServices(
							[]uint32{3, 4},
						)),
				},
			},
		},
		{
			name: "service params are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				suite.k.SaveServiceParams(ctx, 1, types.NewServiceParams(
					[]uint32{1, 2},
					nil,
				))

				suite.k.SaveServiceParams(ctx, 2, types.NewServiceParams(
					[]uint32{3, 4},
					[]uint32{5, 6},
				))
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				ServicesParams: []types.ServiceParamsRecord{
					{
						ServiceID: 1,
						Params: types.NewServiceParams(
							[]uint32{1, 2},
							nil,
						),
					},
					{
						ServiceID: 2,
						Params: types.NewServiceParams(
							[]uint32{3, 4},
							[]uint32{5, 6},
						),
					},
				},
			},
		},
		{
			name: "pool delegations are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				err := suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(200))),
				))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				Delegations: []types.Delegation{
					types.NewPoolDelegation(
						1,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
					),
					types.NewPoolDelegation(
						2,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(200))),
					),
				},
			},
		},
		{
			name: "service delegations are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					2,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("services/2/utia", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				Delegations: []types.Delegation{
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

				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					2,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operators/2/utia", sdkmath.LegacyNewDec(50)),
					),
				))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				Delegations: []types.Delegation{
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
			name: "pool unbonding delegations are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				_, err := suite.k.SetUnbondingDelegation(ctx, types.NewPoolUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("pool/1/umilk", sdkmath.NewInt(100))),
					1,
				))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDelegations: []types.UnbondingDelegation{
					types.NewPoolUnbondingDelegation(
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						1,
						10,
						time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
						sdk.NewCoins(sdk.NewCoin("pool/1/umilk", sdkmath.NewInt(100))),
						1,
					),
				},
			},
		},
		{
			name: "operator unbonding delegations are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				_, err := suite.k.SetUnbondingDelegation(ctx, types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("operator/1/umilk", sdkmath.NewInt(100))),
					1,
				))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDelegations: []types.UnbondingDelegation{
					types.NewOperatorUnbondingDelegation(
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						1,
						10,
						time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
						sdk.NewCoins(sdk.NewCoin("operator/1/umilk", sdkmath.NewInt(100))),
						1,
					),
				},
			},
		},
		{
			name: "service unbonding delegations are exported properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.DefaultParams())

				_, err := suite.k.SetUnbondingDelegation(ctx, types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("service/1/umilk", sdkmath.NewInt(100))),
					1,
				))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDelegations: []types.UnbondingDelegation{
					types.NewServiceUnbondingDelegation(
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						1,
						10,
						time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
						sdk.NewCoins(sdk.NewCoin("service/1/umilk", sdkmath.NewInt(100))),
						1,
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
			name: "operator params are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				OperatorsJoinedServices: []types.OperatorJoinedServicesRecord{
					types.NewOperatorJoinedServicesRecord(1,
						types.NewOperatorJoinedServices(
							[]uint32{1, 2},
						),
					),
					types.NewOperatorJoinedServicesRecord(2,
						types.NewOperatorJoinedServices(
							[]uint32{3, 4},
						),
					),
				},
			},
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetOperatorJoinedServices(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewOperatorJoinedServices(
					[]uint32{1, 2},
				), stored)

				stored, err = suite.k.GetOperatorJoinedServices(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewOperatorJoinedServices(
					[]uint32{3, 4},
				), stored)
			},
		},
		{
			name: "service params are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				ServicesParams: []types.ServiceParamsRecord{
					{
						ServiceID: 1,
						Params: types.NewServiceParams(
							[]uint32{1, 2},
							nil,
						),
					},
					{
						ServiceID: 2,
						Params: types.NewServiceParams(
							[]uint32{3, 4},
							[]uint32{5, 6},
						),
					},
				},
			},
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetServiceParams(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewServiceParams(
					[]uint32{1, 2},
					nil,
				), stored)

				stored, err = suite.k.GetServiceParams(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewServiceParams(
					[]uint32{3, 4},
					[]uint32{5, 6},
				), stored)
			},
		},
		{
			name: "pool delegations are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				Delegations: []types.Delegation{
					types.NewPoolDelegation(
						1,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
					),
					types.NewPoolDelegation(
						2,
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						sdk.NewDecCoins(sdk.NewDecCoinFromDec("utia", sdkmath.LegacyNewDec(200))),
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
				Delegations: []types.Delegation{
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
				Delegations: []types.Delegation{
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
			name: "pool unbonding delegations are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDelegations: []types.UnbondingDelegation{
					types.NewPoolUnbondingDelegation(
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						1,
						10,
						time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
						sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
						1,
					),
				},
			},
			check: func(ctx sdk.Context) {
				stored, poolUnbondingDelegationFound := suite.k.GetPoolUnbondingDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(poolUnbondingDelegationFound)

				suite.Require().Equal(types.NewPoolUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				), stored)
			},
		},
		{
			name: "operator unbonding delegations are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDelegations: []types.UnbondingDelegation{
					types.NewOperatorUnbondingDelegation(
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						1,
						10,
						time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
						sdk.NewCoins(sdk.NewCoin("operators/1/umilk", sdkmath.NewInt(100))),
						1,
					),
				},
			},
			check: func(ctx sdk.Context) {
				stored, operatorUnbondingDelegationFound := suite.k.GetOperatorUnbondingDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(operatorUnbondingDelegationFound)

				suite.Require().Equal(types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("operators/1/umilk", sdkmath.NewInt(100))),
					1,
				), stored)
			},
		},
		{
			name: "service unbonding delegations are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDelegations: []types.UnbondingDelegation{
					types.NewServiceUnbondingDelegation(
						"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
						1,
						10,
						time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
						sdk.NewCoins(sdk.NewCoin("services/1/umilk", sdkmath.NewInt(100))),
						1,
					),
				},
			},
			check: func(ctx sdk.Context) {
				stored, serviceUnbondingDelegationFound := suite.k.GetServiceUnbondingDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(serviceUnbondingDelegationFound)

				suite.Require().Equal(types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("services/1/umilk", sdkmath.NewInt(100))),
					1,
				), stored)
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
