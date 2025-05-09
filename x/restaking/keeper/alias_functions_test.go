package keeper_test

import (
	"fmt"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v12/x/operators/types"
	servicestypes "github.com/milkyway-labs/milkyway/v12/x/services/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetDelegation() {
	testCases := []struct {
		name       string
		delegation types.Delegation
		shouldErr  bool
		check      func(ctx sdk.Context)
	}{
		{
			name: "pool delegation is set properly",
			delegation: types.NewPoolDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(50))),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				store := suite.storeService.OpenKVStore(ctx)

				// Make sure the delegation key exists
				found, err := store.Has(types.UserPoolDelegationStoreKey("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", 1))
				suite.Require().NoError(err)
				suite.Require().True(found)

				// Make sure the delegation-by-pool-id exists
				found, err = store.Has(types.DelegationByPoolIDStoreKey(1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"))
				suite.Require().NoError(err)
				suite.Require().True(found)
			},
		},
		{
			name: "operator delegation is set properly",
			delegation: types.NewOperatorDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				store := suite.storeService.OpenKVStore(ctx)

				// Make sure the delegation key exists
				found, err := store.Has(types.UserOperatorDelegationStoreKey("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", 1))
				suite.Require().NoError(err)
				suite.Require().True(found)

				// Make sure the delegation-by-operator-id exists
				found, err = store.Has(types.DelegationByOperatorIDStoreKey(1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"))
				suite.Require().NoError(err)
				suite.Require().True(found)
			},
		},
		{
			name: "service delegation is set properly",
			delegation: types.NewServiceDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				store := suite.storeService.OpenKVStore(ctx)

				// Make sure the delegation key exists
				stored, err := store.Get(types.UserServiceDelegationStoreKey("cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd", 1))
				suite.Require().NoError(err)
				suite.Require().NotEmpty(stored)

				// Make sure the delegation-by-service-id exists
				found, err := store.Has(types.DelegationByServiceIDStoreKey(1, "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd"))
				suite.Require().NoError(err)
				suite.Require().True(found)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()

			err := suite.k.SetDelegation(ctx, tc.delegation)
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

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_GetAllOperatorsJoinedServicesRecord() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		shouldErr  bool
		expRecords []types.OperatorJoinedServices
	}{
		{
			name:       "operator without joined services returns nil",
			shouldErr:  false,
			expRecords: nil,
		},
		{
			name: "operators joined services are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 3)
				suite.Require().NoError(err)
				err = suite.k.AddServiceToOperatorJoinedServices(ctx, 2, 4)
				suite.Require().NoError(err)
			},
			expRecords: []types.OperatorJoinedServices{
				types.NewOperatorJoinedServices(1, []uint32{1, 2}),
				types.NewOperatorJoinedServices(2, []uint32{3, 4}),
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

			records, err := suite.k.GetAllOperatorsJoinedServices(ctx)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expRecords, records)
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_GetAllServicesAllowedOperators() {
	testCases := []struct {
		name                string
		store               func(ctx sdk.Context)
		shouldErr           bool
		expectedWhitelisted []types.ServiceAllowedOperators
	}{
		{
			name:      "no allowed operators returns nil",
			shouldErr: false,
		},
		{
			name: "securing pools are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.AddPoolToServiceSecuringPools(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			shouldErr: false,
		},
		{
			name: "allowed operators are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.AddOperatorToServiceAllowList(ctx, 2, 4)
				suite.Require().NoError(err)
				err = suite.k.AddOperatorToServiceAllowList(ctx, 2, 5)
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expectedWhitelisted: []types.ServiceAllowedOperators{
				types.NewServiceAllowedOperators(1, []uint32{1, 2}),
				types.NewServiceAllowedOperators(2, []uint32{4, 5}),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			whitelistedOperators, err := suite.k.GetAllServicesAllowedOperators(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expectedWhitelisted, whitelistedOperators)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetAllServicesSecuringPools() {
	testCases := []struct {
		name                string
		store               func(ctx sdk.Context)
		shouldErr           bool
		expectedWhitelisted []types.ServiceSecuringPools
	}{
		{
			name:      "no securing pools returns nil",
			shouldErr: false,
		},
		{
			name: "allowed operators are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.AddOperatorToServiceAllowList(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddOperatorToServiceAllowList(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			shouldErr: false,
		},
		{
			name: "securing pools are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.AddPoolToServiceSecuringPools(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.AddPoolToServiceSecuringPools(ctx, 2, 4)
				suite.Require().NoError(err)
				err = suite.k.AddPoolToServiceSecuringPools(ctx, 2, 5)
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expectedWhitelisted: []types.ServiceSecuringPools{
				types.NewServiceSecuringPools(1, []uint32{1, 2}),
				types.NewServiceSecuringPools(2, []uint32{4, 5}),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx := suite.ctx
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			whitelistedOperators, err := suite.k.GetAllServicesSecuringPools(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expectedWhitelisted, whitelistedOperators)
			}
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
			shouldErr: false,
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

			delegations, err := suite.k.GetAllPoolDelegations(ctx)
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
			shouldErr: false,
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

			delegations, err := suite.k.GetAllOperatorDelegations(ctx)
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
			shouldErr: false,
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

			delegations, err := suite.k.GetAllServiceDelegations(ctx)
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
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					2,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					3,
				))
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expUnbonding: []types.UnbondingDelegation{
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

			unbonding, err := suite.k.GetAllPoolUnbondingDelegations(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbonding, unbonding)
			}

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
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					2,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					3,
				))
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expUnbonding: []types.UnbondingDelegation{
				types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
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

			unbonding, err := suite.k.GetAllOperatorUnbondingDelegations(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbonding, unbonding)
			}

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
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					1,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewOperatorUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					2,
				))
				suite.Require().NoError(err)

				_, err = suite.k.SetUnbondingDelegation(ctx, types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
					3,
				))
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expUnbonding: []types.UnbondingDelegation{
				types.NewServiceUnbondingDelegation(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					2,
					10,
					time.Date(2024, 1, 8, 12, 0, 0, 0, time.UTC),
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

			unbonding, err := suite.k.GetAllServiceUnbondingDelegations(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expUnbonding, unbonding)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_UnbondRestakedAssets() {
	delegator := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		account   string
		amount    sdk.Coins
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "undelegate without delegations fails",
			account:   delegator,
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
			shouldErr: true,
		},
		{
			name: "undelegate more then delegated fails",
			store: func(ctx sdk.Context) {
				// Fund the account
				suite.fundAccount(ctx, delegator, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

				err := suite.pk.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)

				_, err = suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", 300), delegator)
				suite.Assert().NoError(err)

				// Delegate to service
				err = suite.sk.CreateService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"",
					"",
					"",
					"",
					"",
					false,
				))
				suite.Assert().NoError(err)
				_, err = suite.k.DelegateToService(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin("stake", 350)), delegator)
				suite.Assert().NoError(err)

				// Delegate to operator
				err = suite.ok.CreateOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE, "", "", "", "",
				))
				suite.Assert().NoError(err)
				_, err = suite.k.DelegateToOperator(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin("stake", 350)), delegator)
				suite.Assert().NoError(err)
			},
			account:   delegator,
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 1100)),
			shouldErr: true,
		},
		{
			name:    "undelegate with multiple delegations torward a pool",
			account: delegator,
			amount:  sdk.NewCoins(sdk.NewInt64Coin("stake", 300)),
			store: func(ctx sdk.Context) {
				// Fund the delegator account
				suite.fundAccount(ctx, delegator, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

				// Set the first pool id
				err := suite.pk.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)

				// Create delegators delegations
				for i := 0; i < 100; i++ {
					amount := int64(i*5 + 300)
					d := authtypes.NewModuleAddress(fmt.Sprintf("delegator-%d", i)).String()
					suite.fundAccount(ctx, d, sdk.NewCoins(sdk.NewInt64Coin("stake", amount)))
					_, err = suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", amount), d)
					suite.Assert().NoError(err)
				}

				// Delegate to pool
				_, err = suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", 300), delegator)
				suite.Assert().NoError(err)
			},
		},
		{
			name: "partial undelegate leaves balances to operator",
			store: func(ctx sdk.Context) {
				// Fund the account
				suite.fundAccount(ctx, delegator, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

				err := suite.pk.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)

				_, err = suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", 300), delegator)
				suite.Assert().NoError(err)

				// Delegate to service
				err = suite.sk.CreateService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"",
					"",
					"",
					"",
					"",
					false,
				))
				suite.Assert().NoError(err)
				_, err = suite.k.DelegateToService(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin("stake", 350)), delegator)
				suite.Assert().NoError(err)

				// Delegate to operator
				err = suite.ok.CreateOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE, "", "", "", "",
				))
				suite.Assert().NoError(err)
				_, err = suite.k.DelegateToOperator(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin("stake", 350)), delegator)
				suite.Assert().NoError(err)
			},
			account: delegator,
			amount:  sdk.NewCoins(sdk.NewInt64Coin("stake", 950)),
			check: func(ctx sdk.Context) {
				del, found, err := suite.k.GetServiceDelegation(ctx, 1, delegator)
				suite.Require().NoError(err)
				suite.Assert().True(found)
				suite.Assert().Equal(types.DELEGATION_TYPE_OPERATOR, del.Type)
				operator, err := suite.ok.GetOperator(ctx, 1)
				suite.Require().NoError(err)
				suite.Assert().Equal(
					sdk.NewDecCoins(sdk.NewInt64DecCoin("stake", 50)),
					operator.TokensFromSharesTruncated(del.Shares))
			},
		},
		{
			name: "undelegate  with amounts shared between pool, service and operator",
			store: func(ctx sdk.Context) {
				// Prepare pool, service and operator
				err := suite.pk.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.sk.CreateService(ctx, servicestypes.NewService(
					1,
					servicestypes.SERVICE_STATUS_ACTIVE,
					"",
					"",
					"",
					"",
					"",
					false,
				))
				err = suite.ok.CreateOperator(ctx, operatorstypes.NewOperator(
					1, operatorstypes.OPERATOR_STATUS_ACTIVE, "", "", "", "",
				))

				// Create delegators delegations
				for i := 0; i < 100; i++ {
					amount := int64(i*5 + 300)
					totalAmount := amount * 3
					d := authtypes.NewModuleAddress(fmt.Sprintf("delegator-%d", i)).String()
					suite.fundAccount(ctx, d, sdk.NewCoins(sdk.NewInt64Coin("stake", totalAmount)))

					// Perform delegations
					_, err = suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", amount), d)
					suite.Assert().NoError(err)
					_, err = suite.k.DelegateToService(ctx, 1,
						sdk.NewCoins(sdk.NewInt64Coin("stake", amount)), d)
					suite.Assert().NoError(err)
					_, err = suite.k.DelegateToOperator(ctx, 1,
						sdk.NewCoins(sdk.NewInt64Coin("stake", amount)), d)
					suite.Assert().NoError(err)
				}

				// Fund the account
				suite.fundAccount(ctx, delegator, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

				// Perform delegations
				_, err = suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", 300), delegator)
				suite.Assert().NoError(err)
				_, err = suite.k.DelegateToService(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin("stake", 350)), delegator)
				suite.Assert().NoError(err)
				_, err = suite.k.DelegateToOperator(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin("stake", 350)), delegator)
				suite.Assert().NoError(err)
			},
			account:   "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			amount:    sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)),
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			accAddr := sdk.MustAccAddressFromBech32(tc.account)
			completionTime, err := suite.k.UnbondRestakedAssets(ctx, accAddr, tc.amount)

			if !tc.shouldErr {
				suite.Assert().NoError(err)

				unbondingTime, err := suite.k.UnbondingTime(ctx)
				suite.Require().NoError(err)

				expectedCompletion := ctx.BlockHeader().Time.Add(unbondingTime)
				suite.Assert().Equal(expectedCompletion, completionTime)
			} else {
				suite.Assert().Error(err)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_GetUserPreferencesEntries() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		shouldErr  bool
		expEntries []types.UserPreferencesEntry
	}{
		{
			name:       "empty store returns empty entries",
			store:      func(ctx sdk.Context) {},
			expEntries: []types.UserPreferencesEntry(nil),
		},
		{
			name: "entries are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetUserPreferences(ctx, "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47", types.NewUserPreferences([]types.TrustedServiceEntry{
					types.NewTrustedServiceEntry(1, nil),
					types.NewTrustedServiceEntry(2, nil),
					types.NewTrustedServiceEntry(3, nil),
				}))
				suite.Require().NoError(err)

				err = suite.k.SetUserPreferences(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", types.NewUserPreferences([]types.TrustedServiceEntry{
					types.NewTrustedServiceEntry(4, nil),
					types.NewTrustedServiceEntry(5, nil),
				}))
				suite.Require().NoError(err)
			},
			expEntries: []types.UserPreferencesEntry{
				types.NewUserPreferencesEntry(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					types.NewUserPreferences([]types.TrustedServiceEntry{
						types.NewTrustedServiceEntry(4, nil),
						types.NewTrustedServiceEntry(5, nil),
					}),
				),
				types.NewUserPreferencesEntry(
					"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
					types.NewUserPreferences([]types.TrustedServiceEntry{
						types.NewTrustedServiceEntry(1, nil),
						types.NewTrustedServiceEntry(2, nil),
						types.NewTrustedServiceEntry(3, nil),
					}),
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

			entries, err := suite.k.GetUserPreferencesEntries(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expEntries, entries)
			}
		})
	}
}
