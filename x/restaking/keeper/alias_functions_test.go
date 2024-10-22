package keeper_test

import (
	"fmt"
	"time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_GetAllOperatorsJoinedServicesRecord() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		shouldErr  bool
		expRecords []types.OperatorJoinedServicesRecord
	}{
		{
			name: "operators joined services are returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperatorJoinedServices(ctx, 1, types.NewOperatorJoinedServices(
					[]uint32{1, 2},
				))

				suite.k.SaveOperatorJoinedServices(ctx, 2, types.NewOperatorJoinedServices(
					[]uint32{3, 4},
				))
			},
			expRecords: []types.OperatorJoinedServicesRecord{
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

func (suite *KeeperTestSuite) TestKeeper_GetAllServicesWhitelistedOperators() {
	testCases := []struct {
		name                string
		store               func(ctx sdk.Context)
		shouldErr           bool
		expectedWhitelisted []types.ServiceWhitelistedOperators
	}{
		{
			name:      "no whitelisted operators returns nil",
			shouldErr: false,
		},
		{
			name: "whitelisted pools are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.ServiceWhitelistPool(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			shouldErr: false,
		},
		{
			name: "whitelisted operators are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.ServiceWhitelistOperator(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistOperator(ctx, 2, 4)
				suite.Require().NoError(err)
				err = suite.k.ServiceWhitelistOperator(ctx, 2, 5)
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expectedWhitelisted: []types.ServiceWhitelistedOperators{
				types.NewServiceWhitelistedOperators(1, []uint32{1, 2}),
				types.NewServiceWhitelistedOperators(2, []uint32{4, 5}),
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

			whitelistedOperators, err := suite.k.GetAllServicesWhitelistedOperators(ctx)
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

func (suite *KeeperTestSuite) TestKeeper_GetAllServicesWhitelistedPools() {
	testCases := []struct {
		name                string
		store               func(ctx sdk.Context)
		shouldErr           bool
		expectedWhitelisted []types.ServiceWhitelistedPools
	}{
		{
			name:      "no whitelisted pools returns nil",
			shouldErr: false,
		},
		{
			name: "whitelisted operators are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.ServiceWhitelistOperator(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			shouldErr: false,
		},
		{
			name: "whitelisted pools are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)
				err = suite.k.ServiceWhitelistPool(ctx, 1, 2)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistPool(ctx, 2, 4)
				suite.Require().NoError(err)
				err = suite.k.ServiceWhitelistPool(ctx, 2, 5)
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expectedWhitelisted: []types.ServiceWhitelistedPools{
				types.NewServiceWhitelistedPools(1, []uint32{1, 2}),
				types.NewServiceWhitelistedPools(2, []uint32{4, 5}),
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

			whitelistedOperators, err := suite.k.GetAllServicesWhitelistedPools(ctx)
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

			unbonding := suite.k.GetAllServiceUnbondingDelegations(ctx)
			suite.Require().Equal(tc.expUnbonding, unbonding)

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
		setup     func(ctx sdk.Context)
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
			setup: func(ctx sdk.Context) {
				// Fund the account
				suite.fundAccount(ctx, delegator, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

				suite.pk.SetNextPoolID(ctx, 1)
				_, err := suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", 300), delegator)
				suite.Assert().NoError(err)

				// Delegate to service
				err = suite.sk.CreateService(ctx, servicestypes.NewService(1, servicestypes.SERVICE_STATUS_ACTIVE,
					"", "", "", "", ""))
				suite.Assert().NoError(err)
				_, err = suite.k.DelegateToService(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin("stake", 350)), delegator)
				suite.Assert().NoError(err)

				// Delegate to operator
				err = suite.ok.RegisterOperator(ctx, operatorstypes.NewOperator(
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
			setup: func(ctx sdk.Context) {
				// Fund the delegator account
				suite.fundAccount(ctx, delegator, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

				// Set the first pool id
				suite.pk.SetNextPoolID(ctx, 1)

				// Create delegators delegations
				for i := 0; i < 100; i++ {
					amount := int64(i*5 + 300)
					d := authtypes.NewModuleAddress(fmt.Sprintf("delegator-%d", i)).String()
					suite.fundAccount(ctx, d, sdk.NewCoins(sdk.NewInt64Coin("stake", amount)))
					_, err := suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", amount), d)
					suite.Assert().NoError(err)
				}

				// Delegate to pool
				_, err := suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", 300), delegator)
				suite.Assert().NoError(err)
			},
		},
		{
			name: "partial undelegate leaves balances to operator",
			setup: func(ctx sdk.Context) {
				// Fund the account
				suite.fundAccount(ctx, delegator, sdk.NewCoins(sdk.NewInt64Coin("stake", 1000)))

				suite.pk.SetNextPoolID(ctx, 1)
				_, err := suite.k.DelegateToPool(ctx, sdk.NewInt64Coin("stake", 300), delegator)
				suite.Assert().NoError(err)

				// Delegate to service
				err = suite.sk.CreateService(ctx, servicestypes.NewService(1, servicestypes.SERVICE_STATUS_ACTIVE,
					"", "", "", "", ""))
				suite.Assert().NoError(err)
				_, err = suite.k.DelegateToService(ctx, 1,
					sdk.NewCoins(sdk.NewInt64Coin("stake", 350)), delegator)
				suite.Assert().NoError(err)

				// Delegate to operator
				err = suite.ok.RegisterOperator(ctx, operatorstypes.NewOperator(
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
				del, found := suite.k.GetServiceDelegation(ctx, 1, delegator)
				suite.Assert().True(found)
				suite.Assert().Equal(types.DELEGATION_TYPE_OPERATOR, del.Type)
				operator, _ := suite.ok.GetOperator(ctx, 1)
				suite.Assert().Equal(
					sdk.NewDecCoins(sdk.NewInt64DecCoin("stake", 50)),
					operator.TokensFromSharesTruncated(del.Shares))
			},
		},
		{
			name: "undelegate  with amounts shared between pool, service and operator",
			setup: func(ctx sdk.Context) {
				// Prepare pool, service and operator
				suite.pk.SetNextPoolID(ctx, 1)
				err := suite.sk.CreateService(ctx, servicestypes.NewService(1, servicestypes.SERVICE_STATUS_ACTIVE,
					"", "", "", "", ""))
				err = suite.ok.RegisterOperator(ctx, operatorstypes.NewOperator(
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
			if tc.setup != nil {
				tc.setup(ctx)
			}

			accAddr := sdk.MustAccAddressFromBech32(tc.account)
			completionTime, err := suite.k.UnbondRestakedAssets(ctx, accAddr, tc.amount)

			if !tc.shouldErr {
				suite.Assert().NoError(err)
				expectedCompletion := ctx.BlockHeader().Time.Add(suite.k.UnbondingTime(ctx))
				suite.Assert().Equal(expectedCompletion, completionTime)
			} else {
				suite.Assert().Error(err)
			}
		})
	}
}
