package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
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
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 3)
				suite.Require().NoError(err)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 4)
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				OperatorsJoinedServices: []types.OperatorJoinedServices{
					types.NewOperatorJoinedServices(1, []uint32{1, 2}),
					types.NewOperatorJoinedServices(2, []uint32{3, 4}),
				},
			},
		},
		{
			name: "services allowed operators are exported properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 2, 3)
				suite.Require().NoError(err)
				err = suite.k.AddOperatorToServiceAllowList(ctx, 2, 4)
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				ServicesAllowedOperators: []types.ServiceAllowedOperators{
					types.NewServiceAllowedOperators(1, []uint32{1, 2}),
					types.NewServiceAllowedOperators(2, []uint32{3, 4}),
				},
			},
		},
		{
			name: "services securing pools are exported properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				err = suite.k.AddPoolToServiceSecuringPools(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.AddPoolToServiceSecuringPools(ctx, 2, 3)
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 2, 4)
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				ServicesSecuringPools: []types.ServiceSecuringPools{
					types.NewServiceSecuringPools(1, []uint32{1, 2}),
					types.NewServiceSecuringPools(2, []uint32{3, 4}),
				},
			},
		},
		{
			name: "pool delegations are exported properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
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
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
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
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
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
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewPoolUnbondingDelegation(
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
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewOperatorUnbondingDelegation(
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
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewServiceUnbondingDelegation(
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
			name: "user preferences are exported properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.DefaultParams())
				suite.Require().NoError(err)

				err = suite.k.SetUserPreferences(ctx, "cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw", types.NewUserPreferences(
					true,
					false,
					[]uint32{1, 2, 3},
				))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
				UsersPreferences: []types.UserPreferencesEntry{
					types.NewUserPreferencesEntry(
						"cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw",
						types.NewUserPreferences(
							true,
							false,
							[]uint32{1, 2, 3},
						),
					),
				},
			},
		},
		{
			name: "params are exported properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					30*24*time.Hour,
					nil,
					sdkmath.LegacyNewDec(100000),
				))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				Params: types.NewParams(
					30*24*time.Hour,
					nil,
					sdkmath.LegacyNewDec(100000),
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
		name      string
		genesis   *types.GenesisState
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "operator params are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				OperatorsJoinedServices: []types.OperatorJoinedServices{
					types.NewOperatorJoinedServices(1, []uint32{1, 2}),
					types.NewOperatorJoinedServices(2, []uint32{3, 4}),
				},
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				joined, err := suite.k.HasOperatorJoinedService(ctx, 1, 1)
				suite.Require().NoError(err)
				suite.Require().True(joined)

				joined, err = suite.k.HasOperatorJoinedService(ctx, 1, 2)
				suite.Require().NoError(err)
				suite.Require().True(joined)

				joined, err = suite.k.HasOperatorJoinedService(ctx, 2, 3)
				suite.Require().NoError(err)
				suite.Require().True(joined)

				joined, err = suite.k.HasOperatorJoinedService(ctx, 2, 4)
				suite.Require().NoError(err)
				suite.Require().True(joined)
			},
		},
		{
			name: "services allowed operators are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				ServicesAllowedOperators: []types.ServiceAllowedOperators{
					types.NewServiceAllowedOperators(1, []uint32{1, 2}),
					types.NewServiceAllowedOperators(2, []uint32{3, 4}),
				},
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetAllServiceAllowedOperators(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal([]uint32{1, 2}, stored)

				stored, err = suite.k.GetAllServiceAllowedOperators(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().Equal([]uint32{3, 4}, stored)
			},
		},
		{
			name: "services securing pools are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				ServicesSecuringPools: []types.ServiceSecuringPools{
					types.NewServiceSecuringPools(1, []uint32{1, 2}),
					types.NewServiceSecuringPools(2, []uint32{3, 4}),
				},
			},
			shouldErr: false,
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetAllServiceSecuringPools(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal([]uint32{1, 2}, stored)

				stored, err = suite.k.GetAllServiceSecuringPools(ctx, 2)
				suite.Require().NoError(err)
				suite.Require().Equal([]uint32{3, 4}, stored)
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
			shouldErr: false,
			check: func(ctx sdk.Context) {
				_, found, err := suite.k.GetPoolDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)

				_, found, err = suite.k.GetPoolDelegation(ctx, 2, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)
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
				_, found, err := suite.k.GetServiceDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)

				_, found, err = suite.k.GetServiceDelegation(ctx, 2, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)
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
				_, found, err := suite.k.GetOperatorDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)

				_, found, err = suite.k.GetOperatorDelegation(ctx, 2, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)
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
				stored, found, err := suite.k.GetPoolUnbondingDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)

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
				stored, found, err := suite.k.GetOperatorUnbondingDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)

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
				stored, found, err := suite.k.GetServiceUnbondingDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)

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
			name: "user preferences are stored properly",
			genesis: &types.GenesisState{
				Params: types.DefaultParams(),
				UsersPreferences: []types.UserPreferencesEntry{
					types.NewUserPreferencesEntry(
						"cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw",
						types.NewUserPreferences(
							true,
							false,
							[]uint32{1, 2, 3},
						),
					),
				},
			},
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetUserPreferencesEntries(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal([]types.UserPreferencesEntry{
					types.NewUserPreferencesEntry(
						"cosmos1jseuux3pktht0kkhlcsv4kqff3mql65udqs4jw",
						types.NewUserPreferences(
							true,
							false,
							[]uint32{1, 2, 3},
						),
					),
				}, stored)
			},
		},
		{
			name: "params are stored properly",
			genesis: &types.GenesisState{
				Params: types.NewParams(
					30*24*time.Hour,
					nil,
					sdkmath.LegacyNewDec(100000),
				),
			},
			check: func(ctx sdk.Context) {
				params, err := suite.k.GetParams(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewParams(
					30*24*time.Hour,
					nil,
					sdkmath.LegacyNewDec(100000),
				), params)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()

			err := suite.k.InitGenesis(ctx, tc.genesis)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
