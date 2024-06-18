package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_GetAllPoolDelegations() {
	testCases := []struct {
		name           string
		setup          func()
		store          func(ctx sdk.Context)
		shouldErr      bool
		expDelegations []types.PoolDelegation
		check          func(ctx sdk.Context)
	}{
		{
			name: "delegations are returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(100),
				))
			},
			expDelegations: []types.PoolDelegation{
				types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				),
				types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				),
				types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(100),
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

func (suite *KeeperTestSuite) TestKeeper_GetAllDelegatorPoolDelegations() {
	testCases := []struct {
		name           string
		setup          func()
		store          func(ctx sdk.Context)
		delegator      string
		expDelegations []types.PoolDelegation
		check          func(ctx sdk.Context)
	}{
		{
			name:           "user without delegations returns empty list",
			delegator:      "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			expDelegations: nil,
		},
		{
			name: "user with single delegation returns it properly",
			store: func(ctx sdk.Context) {
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(100),
				))
			},
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			expDelegations: []types.PoolDelegation{
				types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(100),
				),
			},
		},
		{
			name: "user with multiple delegations returns them properly",
			store: func(ctx sdk.Context) {
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(100),
				))
			},
			delegator: "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			expDelegations: []types.PoolDelegation{
				types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				),
				types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
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

			delegations := suite.k.GetAllDelegatorPoolDelegations(ctx, tc.delegator)
			suite.Require().Equal(tc.expDelegations, delegations)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetPoolDelegations() {
	testCases := []struct {
		name           string
		setup          func()
		store          func(ctx sdk.Context)
		poolID         uint32
		shouldErr      bool
		expDelegations []types.PoolDelegation
		check          func(ctx sdk.Context)
	}{
		{
			name:           "pool without delegations returns empty list",
			poolID:         1,
			shouldErr:      false,
			expDelegations: nil,
		},
		{
			name: "pool with single delegation returns it properly",
			store: func(ctx sdk.Context) {
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(100),
				))
			},
			poolID:    2,
			shouldErr: false,
			expDelegations: []types.PoolDelegation{
				types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				),
			},
		},
		{
			name: "pool with multiple delegations returns them properly",
			store: func(ctx sdk.Context) {
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					2,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(50),
				))
				suite.k.SavePoolDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(100),
				))
			},
			poolID:    1,
			shouldErr: false,
			expDelegations: []types.PoolDelegation{
				types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdkmath.LegacyNewDec(100),
				),
				types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdkmath.LegacyNewDec(100),
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

			delegations, err := suite.k.GetPoolDelegations(ctx, tc.poolID)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expDelegations, delegations)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
