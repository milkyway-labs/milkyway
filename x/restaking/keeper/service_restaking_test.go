package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_GetAllServiceWhitelistedOperators() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		serviceID uint32
		expected  []uint32
	}{
		{
			name: "whitelisted pools are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistPool(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 1,
		},
		{
			name: "whitelisted operators for different service are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistOperator(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 2,
		},
		{
			name: "whitelisted operators are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistOperator(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 1,
			expected:  []uint32{1, 2},
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

			whitelisted, err := suite.k.GetAllServiceWhitelistedOperators(ctx, tc.serviceID)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expected, whitelisted)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ServiceIsOperatorWhitelisted() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		serviceID uint32
		expected  []uint32
	}{
		{
			name: "whitelisted pools are not considered whitelisted",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistPool(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 1,
		},
		{
			name: "whitelisted operators for different service are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistOperator(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 2,
		},
		{
			name: "whitelisted operators are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistOperator(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 1,
			expected:  []uint32{1, 2},
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

			for _, expected := range tc.expected {
				whitelisted, err := suite.k.ServiceIsOperatorWhitelisted(ctx, tc.serviceID, expected)
				suite.Require().NoError(err)
				suite.Require().True(whitelisted)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ServiceIsOpertorsWhitelistConfigured() {
	testCases := []struct {
		name          string
		store         func(ctx sdk.Context)
		serviceID     uint32
		isInitialized bool
	}{
		{
			name:          "no whitelisted operators means not initialized",
			serviceID:     1,
			isInitialized: false,
		},
		{
			name: "no whitelisted operators for service means not initialized",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistOperator(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID:     2,
			isInitialized: false,
		},
		{
			name: "with whitelisted operator is initialized",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			serviceID:     1,
			isInitialized: true,
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

			isInitialized, err := suite.k.ServiceIsOpertorsWhitelistConfigured(ctx, tc.serviceID)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.isInitialized, isInitialized)
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_GetAllServiceWhitelistedPools() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		serviceID uint32
		expected  []uint32
	}{
		{
			name: "whitelisted operators are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistOperator(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 1,
		},
		{
			name: "whitelisted pools for different service are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistPool(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 2,
		},
		{
			name: "whitelisted pools are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistPool(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 1,
			expected:  []uint32{1, 2},
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

			whitelisted, err := suite.k.GetAllServiceWhitelistedPools(ctx, tc.serviceID)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.expected, whitelisted)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ServiceIsPoolWhitelisted() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		serviceID uint32
		expected  []uint32
	}{
		{
			name: "whitelisted operators are not considered whitelisted",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistOperator(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistOperator(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 1,
		},
		{
			name: "whitelisted pools for different service are not returned",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistPool(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 2,
		},
		{
			name: "pools are whitelisted properly",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistPool(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID: 1,
			expected:  []uint32{1, 2},
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

			for _, expected := range tc.expected {
				whitelisted, err := suite.k.ServiceIsPoolWhitelisted(ctx, tc.serviceID, expected)
				suite.Require().NoError(err)
				suite.Require().True(whitelisted)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_ServiceIsPoolWhitelistConfigured() {
	testCases := []struct {
		name          string
		store         func(ctx sdk.Context)
		serviceID     uint32
		isInitialized bool
	}{
		{
			name:          "no whitelisted pools means not initialized",
			serviceID:     1,
			isInitialized: false,
		},
		{
			name: "no whitelisted pools for service means not initialized",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)

				err = suite.k.ServiceWhitelistPool(ctx, 1, 2)
				suite.Require().NoError(err)
			},
			serviceID:     2,
			isInitialized: false,
		},
		{
			name: "with whitelisted pools is initialized",
			store: func(ctx sdk.Context) {
				err := suite.k.ServiceWhitelistPool(ctx, 1, 1)
				suite.Require().NoError(err)
			},
			serviceID:     1,
			isInitialized: true,
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

			isInitialized, err := suite.k.ServiceIsPoolsWhitelistConfigured(ctx, tc.serviceID)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.isInitialized, isInitialized)
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_SaveServiceDelegation() {
	testCases := []struct {
		name       string
		setup      func()
		store      func(ctx sdk.Context)
		delegation types.Delegation
		check      func(ctx sdk.Context)
	}{
		{
			name: "service delegation is stored properly",
			delegation: types.NewServiceDelegation(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			check: func(ctx sdk.Context) {
				store := ctx.KVStore(suite.storeKey)

				// Make sure the user-service delegation key exists and contains the delegation
				delegationBz := store.Get(types.UserServiceDelegationStoreKey("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", 1))
				suite.Require().NotNil(delegationBz)

				delegation, err := types.UnmarshalDelegation(suite.cdc, delegationBz)
				suite.Require().NoError(err)

				suite.Require().Equal(types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
				), delegation)

				// Make sure the service-user delegation key exists
				hasDelegationsByServiceKey := store.Has(types.DelegationByServiceIDStoreKey(1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"))
				suite.Require().True(hasDelegationsByServiceKey)
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

			err := suite.k.SetDelegation(ctx, tc.delegation)
			suite.Require().NoError(err)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetServiceDelegation() {
	testCases := []struct {
		name          string
		setup         func()
		store         func(ctx sdk.Context)
		serviceID     uint32
		userAddress   string
		expFound      bool
		expDelegation types.Delegation
		check         func(ctx sdk.Context)
	}{
		{
			name:        "not found delegation returns false",
			serviceID:   1,
			userAddress: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			expFound:    false,
		},
		{
			name: "found delegation is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			serviceID:   1,
			userAddress: "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			expFound:    true,
			expDelegation: types.NewServiceDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
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

			delegation, found := suite.k.GetServiceDelegation(ctx, tc.serviceID, tc.userAddress)
			if !tc.expFound {
				suite.Require().False(found)
			} else {
				suite.Require().True(found)
				suite.Require().Equal(tc.expDelegation, delegation)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_AddServiceTokensAndShares() {
	testCases := []struct {
		name           string
		setup          func()
		store          func(ctx sdk.Context)
		service        servicestypes.Service
		tokensToAdd    sdk.Coins
		shouldErr      bool
		expService     servicestypes.Service
		expAddedShares sdk.DecCoins
		check          func(ctx sdk.Context)
	}{
		{
			name: "adding tokens to an empty service works properly",
			service: servicestypes.Service{
				ID:      1,
				Address: servicestypes.GetServiceAddress(1).String(),
			},
			tokensToAdd: sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			shouldErr:   false,
			expService: servicestypes.Service{
				ID:      1,
				Address: servicestypes.GetServiceAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(100)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
				),
			},
			expAddedShares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
			),
		},
		{
			name: "adding tokens to a non-empty service works properly",
			service: servicestypes.Service{
				ID:      1,
				Address: servicestypes.GetServiceAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(50)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
				),
			},
			tokensToAdd: sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(20))),
			shouldErr:   false,
			expService: servicestypes.Service{
				ID:      1,
				Address: servicestypes.GetServiceAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(70)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(140)),
				),
			},
			expAddedShares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(40)),
			),
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

			service, addedShares, err := suite.k.AddServiceTokensAndShares(ctx, tc.service, tc.tokensToAdd)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expService, service)
				suite.Require().Equal(tc.expAddedShares, addedShares)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_DelegateToService() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		serviceID uint32
		amount    sdk.Coins
		delegator string
		shouldErr bool
		expShares sdk.DecCoins
		check     func(ctx sdk.Context)
	}{
		{
			name:      "service not found returns error",
			serviceID: 1,
			amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: true,
		},
		{
			name: "inactive service returns error",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:              1,
					Status:          servicestypes.SERVICE_STATUS_INACTIVE,
					Address:         servicestypes.GetServiceAddress(1).String(),
					Tokens:          sdk.NewCoins(),
					DelegatorShares: sdk.NewDecCoins(),
				})
			},
			serviceID: 1,
			amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: true,
		},
		{
			name: "invalid exchange rate service returns error",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:              1,
					Status:          servicestypes.SERVICE_STATUS_ACTIVE,
					Address:         servicestypes.GetServiceAddress(1).String(),
					Tokens:          sdk.NewCoins(),
					DelegatorShares: sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
				})
			},
			serviceID: 1,
			amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			store: func(ctx sdk.Context) {
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
				})
			},
			serviceID: 1,
			amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator: "invalid",
			shouldErr: true,
		},
		{
			name: "insufficient funds return error",
			store: func(ctx sdk.Context) {
				// Create the service
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Address: servicestypes.GetServiceAddress(1).String(),
				})

				// Set the next service id
				suite.sk.SetNextServiceID(ctx, 2)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
			},
			serviceID: 1,
			amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: true,
		},
		{
			name: "delegating to an existing service works properly",
			store: func(ctx sdk.Context) {
				// Create the service
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})

				// Set the correct service tokens amount
				suite.fundAccount(
					ctx,
					servicestypes.GetServiceAddress(1).String(),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(20))),
				)

				// Set the next service id
				suite.sk.SetNextServiceID(ctx, 2)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			serviceID: 1,
			amount:    sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: false,
			expShares: sdk.NewDecCoins(sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(500))),
			check: func(ctx sdk.Context) {
				// Make sure the service now exists
				service, found := suite.sk.GetService(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(120)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(600)),
					),
				}, service)

				// Make sure the delegation exists
				delegation, found := suite.k.GetServiceDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(found)
				suite.Require().Equal(types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(500)),
					),
				), delegation)

				// Make sure the user balance has been reduced properly
				userBalance := suite.bk.GetBalance(ctx, sdk.AccAddress("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(0)), userBalance)

				// Make sure the service account balance has increased properly
				serviceBalance := suite.bk.GetBalance(ctx, servicestypes.GetServiceAddress(1), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(120)), serviceBalance)
			},
		},
		{
			name: "delegating another token denom works properly",
			store: func(ctx sdk.Context) {
				// Create the service
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(80)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(125)),
					),
				})

				// Set the correct service tokens amount
				suite.fundAccount(
					ctx,
					servicestypes.GetServiceAddress(1).String(),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(80))),
				)

				// Set the next service id
				suite.sk.SetNextServiceID(ctx, 2)

				// Save the existing delegation
				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(125)),
					),
				))

				suite.Require().NoError(err)
				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("uinit", sdkmath.NewInt(100))),
				)
			},
			serviceID: 1,
			amount:    sdk.NewCoins(sdk.NewCoin("uinit", sdkmath.NewInt(100))),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: false,
			expShares: sdk.NewDecCoins(sdk.NewDecCoinFromDec("service/1/uinit", sdkmath.LegacyNewDec(100))),
			check: func(ctx sdk.Context) {
				// Make sure the service now exists
				service, found := suite.sk.GetService(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(80)),
						sdk.NewCoin("uinit", sdkmath.NewInt(100)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(125)),
						sdk.NewDecCoinFromDec("service/1/uinit", sdkmath.LegacyNewDec(100)),
					),
				}, service)

				// Make sure the delegation has been updated properly
				delegation, found := suite.k.GetServiceDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(found)
				suite.Require().Equal(types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(125)),
						sdk.NewDecCoinFromDec("service/1/uinit", sdkmath.LegacyNewDec(100)),
					),
				), delegation)

				// Make sure the user balance has been reduced properly
				userBalance := suite.bk.GetBalance(ctx, sdk.AccAddress("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(0)), userBalance)

				// Make sure the service account balance has increased properly
				serviceBalance := suite.bk.GetAllBalances(ctx, servicestypes.GetServiceAddress(1))
				suite.Require().Equal(sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(80)),
					sdk.NewCoin("uinit", sdkmath.NewInt(100)),
				), serviceBalance)
			},
		},
		{
			name: "delegating more tokens works properly",
			store: func(ctx sdk.Context) {
				// Create the service
				suite.sk.SaveService(ctx, servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(80)),
						sdk.NewCoin("uinit", sdkmath.NewInt(75)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(125)),
						sdk.NewDecCoinFromDec("service/1/uinit", sdkmath.LegacyNewDec(200)),
					),
				})

				// Set the correct service tokens amount
				suite.fundAccount(
					ctx,
					servicestypes.GetServiceAddress(1).String(),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(80)),
						sdk.NewCoin("uinit", sdkmath.NewInt(75)),
					),
				)

				// Set the next service id
				suite.sk.SetNextServiceID(ctx, 2)

				// Save the existing delegation
				err := suite.k.SetDelegation(ctx, types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDec(100)),
						sdk.NewDecCoinFromDec("service/1/uinit", sdkmath.LegacyNewDec(60)),
					),
				))
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(100)),
						sdk.NewCoin("uinit", sdkmath.NewInt(225)),
					),
				)
			},
			serviceID: 1,
			amount: sdk.NewCoins(
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
				sdk.NewCoin("uinit", sdkmath.NewInt(225)),
			),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: false,
			expShares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDecWithPrec(15625, 2)),
				sdk.NewDecCoinFromDec("service/1/uinit", sdkmath.LegacyNewDec(600)),
			),
			check: func(ctx sdk.Context) {
				// Make sure the service now exists
				service, found := suite.sk.GetService(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(servicestypes.Service{
					ID:      1,
					Status:  servicestypes.SERVICE_STATUS_ACTIVE,
					Address: servicestypes.GetServiceAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(180)),
						sdk.NewCoin("uinit", sdkmath.NewInt(300)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDecWithPrec(28125, 2)),
						sdk.NewDecCoinFromDec("service/1/uinit", sdkmath.LegacyNewDec(800)),
					),
				}, service)

				// Make sure the delegation has been updated properly
				delegation, found := suite.k.GetServiceDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(found)
				suite.Require().Equal(types.NewServiceDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("service/1/umilk", sdkmath.LegacyNewDecWithPrec(25625, 2)), // 100 (existing) + 156.25 (new)
						sdk.NewDecCoinFromDec("service/1/uinit", sdkmath.LegacyNewDec(660)),              // 60 (existing) + 600 (new)
					),
				), delegation)

				// Make sure the user balance has been reduced properly
				userBalance := suite.bk.GetBalance(ctx, sdk.AccAddress("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(0)), userBalance)

				// Make sure the service account balance has increased properly
				serviceBalance := suite.bk.GetAllBalances(ctx, servicestypes.GetServiceAddress(1))
				suite.Require().Equal(sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(180)),
					sdk.NewCoin("uinit", sdkmath.NewInt(300)),
				), serviceBalance)
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

			shares, err := suite.k.DelegateToService(ctx, tc.serviceID, tc.amount, tc.delegator)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expShares, shares)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
