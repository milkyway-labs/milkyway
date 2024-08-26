package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_GetAllOperatorsParams() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		shouldErr  bool
		expRecords []types.OperatorParamsRecord
	}{
		{
			name: "operators params are returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperatorParams(ctx, 1, types.NewOperatorParams(
					sdkmath.LegacyNewDec(10),
					[]uint32{1, 2},
				))

				suite.k.SaveOperatorParams(ctx, 2, types.NewOperatorParams(
					sdkmath.LegacyNewDec(3),
					[]uint32{3, 4},
				))
			},
			expRecords: []types.OperatorParamsRecord{
				{
					OperatorID: 1,
					Params: types.NewOperatorParams(
						sdkmath.LegacyNewDec(10),
						[]uint32{1, 2},
					),
				},
				{
					OperatorID: 2,
					Params: types.NewOperatorParams(
						sdkmath.LegacyNewDec(3),
						[]uint32{3, 4},
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

			operators := suite.k.GetAllOperatorsParams(ctx)
			suite.Require().Equal(tc.expRecords, operators)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetAllServicesParams() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		shouldErr  bool
		expRecords []types.ServiceParamsRecord
	}{
		{
			name: "services params are returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveServiceParams(ctx, 1, types.NewServiceParams(
					sdkmath.LegacyNewDec(10),
					[]uint32{1, 2},
					nil,
				))

				suite.k.SaveServiceParams(ctx, 2, types.NewServiceParams(
					sdkmath.LegacyNewDec(3),
					[]uint32{3, 4},
					[]uint32{1, 2},
				))
			},
			expRecords: []types.ServiceParamsRecord{
				{
					ServiceID: 1,
					Params: types.NewServiceParams(
						sdkmath.LegacyNewDec(10),
						[]uint32{1, 2},
						nil,
					),
				},
				{
					ServiceID: 2,
					Params: types.NewServiceParams(
						sdkmath.LegacyNewDec(3),
						[]uint32{3, 4},
						[]uint32{1, 2},
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

			services := suite.k.GetAllServicesParams(ctx)
			suite.Require().Equal(tc.expRecords, services)
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_GetAllPoolDelegations() {
	testCases := []struct {
		name           string
		setup          func()
		store          func(ctx sdk.Context)
		shouldErr      bool
		expDelegations []types.Delegation
		check          func(ctx sdk.Context)
	}{
		{
			name: "delegations are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			expDelegations: []types.Delegation{
				types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
				),
				types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/2/utia", sdkmath.LegacyNewDec(50))),
				),
				types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
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

			delegations := suite.k.GetAllPoolDelegations(ctx)
			suite.Require().Equal(tc.expDelegations, delegations)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetAllOperatorDelegations() {
	testCases := []struct {
		name           string
		setup          func()
		store          func(ctx sdk.Context)
		shouldErr      bool
		expDelegations []types.Delegation
		check          func(ctx sdk.Context)
	}{
		{
			name: "delegations are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("operator/2/utia", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			expDelegations: []types.Delegation{
				types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100))),
				),
				types.NewOperatorDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("operator/2/utia", sdkmath.LegacyNewDec(50))),
				),
				types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100))),
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

			delegations := suite.k.GetAllOperatorDelegations(ctx)
			suite.Require().Equal(tc.expDelegations, delegations)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetAllServiceDelegations() {
	testCases := []struct {
		name           string
		setup          func()
		store          func(ctx sdk.Context)
		shouldErr      bool
		expDelegations []types.Delegation
		check          func(ctx sdk.Context)
	}{
		{
			name: "delegations are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("service/2/utia", sdkmath.LegacyNewDec(50))),
				))
				suite.Require().NoError(err)

				err = suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			expDelegations: []types.Delegation{
				types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100))),
				),
				types.NewServiceDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("service/2/utia", sdkmath.LegacyNewDec(50))),
				),
				types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100))),
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

			delegations := suite.k.GetAllServiceDelegations(ctx)
			suite.Require().Equal(tc.expDelegations, delegations)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_GetAllPoolUnbondingDelegations() {
	testCases := []struct {
		name         string
		setup        func()
		store        func(ctx sdk.Context)
		shouldErr    bool
		expUnbonding []types.UnbondingDelegation
		check        func(ctx sdk.Context)
	}{
		{
			name: "unbonding delegations are returned properly",
			store: func(ctx sdk.Context) {
				_, err := suite.k.SetUnbondingDelegation(ctx, types.NewPoolUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					2,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					3,
				))
				suite.Require().NoError(err)
			},
			expUnbonding: []types.UnbondingDelegation{
				types.NewPoolUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
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

			unbonding := suite.k.GetAllPoolUnbondingDelegations(ctx)
			suite.Require().Equal(tc.expUnbonding, unbonding)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetAllOperatorUnbondingDelegations() {
	testCases := []struct {
		name         string
		setup        func()
		store        func(ctx sdk.Context)
		shouldErr    bool
		expUnbonding []types.UnbondingDelegation
		check        func(ctx sdk.Context)
	}{
		{
			name: "unbonding delegations are returned properly",
			store: func(ctx sdk.Context) {
				_, err := suite.k.SetUnbondingDelegation(ctx, types.NewPoolUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					2,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					3,
				))
				suite.Require().NoError(err)
			},
			expUnbonding: []types.UnbondingDelegation{
				types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					2,
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

			unbonding := suite.k.GetAllOperatorUnbondingDelegations(ctx)
			suite.Require().Equal(tc.expUnbonding, unbonding)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetAllServiceUnbondingDelegations() {
	testCases := []struct {
		name         string
		setup        func()
		store        func(ctx sdk.Context)
		shouldErr    bool
		expUnbonding []types.UnbondingDelegation
		check        func(ctx sdk.Context)
	}{
		{
			name: "unbonding delegations are returned properly",
			store: func(ctx sdk.Context) {
				_, err := suite.k.SetUnbondingDelegation(ctx, types.NewPoolUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					1,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					2,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					3,
				))
				suite.Require().NoError(err)
			},
			expUnbonding: []types.UnbondingDelegation{
				types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 00, 00, 000, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					3,
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

			unbonding := suite.k.GetAllServiceUnbondingDelegations(ctx)
			suite.Require().Equal(tc.expUnbonding, unbonding)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
