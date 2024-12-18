package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	poolstypes "github.com/milkyway-labs/milkyway/v6/x/pools/types"
	"github.com/milkyway-labs/milkyway/v6/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_SavePoolDelegation() {
	testCases := []struct {
		name       string
		setup      func()
		store      func(ctx sdk.Context)
		delegation types.Delegation
		check      func(ctx sdk.Context)
	}{
		{
			name: "pool delegation is stored properly",
			delegation: types.NewPoolDelegation(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			check: func(ctx sdk.Context) {
				store := suite.storeService.OpenKVStore(ctx)

				// Make sure the user-pool delegation key exists and contains the delegation
				delegationBz, err := store.Get(types.UserPoolDelegationStoreKey("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", 1))
				suite.Require().NoError(err)
				suite.Require().NotNil(delegationBz)

				delegation, err := types.UnmarshalDelegation(suite.cdc, delegationBz)
				suite.Require().NoError(err)

				suite.Require().Equal(types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
				), delegation)

				// Make sure the pool-user delegation key exists
				hasDelegationsByPoolKey, err := store.Has(types.DelegationByPoolIDStoreKey(1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"))
				suite.Require().NoError(err)
				suite.Require().True(hasDelegationsByPoolKey)
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

func (suite *KeeperTestSuite) TestKeeper_GetPoolDelegation() {
	testCases := []struct {
		name          string
		setup         func()
		store         func(ctx sdk.Context)
		poolID        uint32
		userAddress   string
		shouldErr     bool
		expFound      bool
		expDelegation types.Delegation
		check         func(ctx sdk.Context)
	}{
		{
			name:        "not found delegation returns false",
			poolID:      1,
			userAddress: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr:   false,
			expFound:    false,
		},
		{
			name: "found delegation is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)
			},
			poolID:      1,
			userAddress: "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			expFound:    true,
			shouldErr:   false,
			expDelegation: types.NewPoolDelegation(
				1,
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
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

			delegation, found, err := suite.k.GetPoolDelegation(ctx, tc.poolID, tc.userAddress)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if !tc.expFound {
					suite.Require().False(found)
				} else {
					suite.Require().True(found)
					suite.Require().Equal(tc.expDelegation, delegation)
				}
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_AddPoolTokensAndShares() {
	testCases := []struct {
		name           string
		setup          func()
		store          func(ctx sdk.Context)
		pool           poolstypes.Pool
		tokensToAdd    sdk.Coin
		shouldErr      bool
		expPool        poolstypes.Pool
		expAddedShares sdk.DecCoin
		check          func(ctx sdk.Context)
	}{
		{
			name:        "adding tokens to an empty pool works properly",
			pool:        poolstypes.NewPool(1, "umilk"),
			tokensToAdd: sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			shouldErr:   false,
			expPool: poolstypes.Pool{
				ID:              1,
				Denom:           "umilk",
				Address:         poolstypes.GetPoolAddress(1).String(),
				Tokens:          sdkmath.NewInt(100),
				DelegatorShares: sdkmath.LegacyNewDec(100),
			},
			expAddedShares: sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100)),
		},
		{
			name: "adding tokens to a non-empty pool works properly",
			pool: poolstypes.Pool{
				ID:              1,
				Denom:           "umilk",
				Address:         poolstypes.GetPoolAddress(1).String(),
				Tokens:          sdkmath.NewInt(50),
				DelegatorShares: sdkmath.LegacyNewDec(100),
			},
			tokensToAdd: sdk.NewCoin("umilk", sdkmath.NewInt(20)),
			shouldErr:   false,
			expPool: poolstypes.Pool{
				ID:              1,
				Denom:           "umilk",
				Address:         poolstypes.GetPoolAddress(1).String(),
				Tokens:          sdkmath.NewInt(70),
				DelegatorShares: sdkmath.LegacyNewDec(140),
			},
			expAddedShares: sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(40)),
		},
		{
			name: "adding tokens to a non-empty pool works properly",
			pool: poolstypes.Pool{
				ID:              1,
				Denom:           "umilk",
				Address:         poolstypes.GetPoolAddress(1).String(),
				Tokens:          sdkmath.NewInt(50),
				DelegatorShares: sdkmath.LegacyNewDec(100),
			},
			tokensToAdd: sdk.NewCoin("utia", sdkmath.NewInt(20)),
			shouldErr:   true,
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

			pool, addedShares, err := suite.k.AddPoolTokensAndShares(ctx, tc.pool, tc.tokensToAdd)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expPool, pool)
				suite.Require().Equal(tc.expAddedShares, addedShares)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_DelegateToPool() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		amount    sdk.Coin
		delegator string
		shouldErr bool
		expShares sdk.DecCoins
		check     func(ctx sdk.Context)
	}{
		{
			name: "invalid exchange rate pool returns error",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.ZeroInt(),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)
			},
			amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: true,
		},
		{
			name: "invalid delegator address returns error",
			store: func(ctx sdk.Context) {
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "umilk"))
				suite.Require().NoError(err)
			},
			amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			delegator: "invalid",
			shouldErr: true,
		},
		{
			name: "insufficient funds return error",
			store: func(ctx sdk.Context) {
				// Create the pool
				err := suite.pk.SavePool(ctx, poolstypes.NewPool(1, "umilk"))
				suite.Require().NoError(err)

				// Set the next pool id
				suite.pk.SetNextPoolID(ctx, 2)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
			},
			amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: true,
		},
		{
			name: "delegating to a non-existing pool works properly",
			store: func(ctx sdk.Context) {
				// Set the next pool id
				suite.pk.SetNextPoolID(ctx, 1)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: false,
			expShares: sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
			check: func(ctx sdk.Context) {
				// Make sure the pool now exists
				pool, err := suite.pk.GetPool(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(100),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				}, pool)

				// Make sure the delegation exists
				delegation, found, err := suite.k.GetPoolDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
				), delegation)

				// Make sure the user balance has been reduced properly
				userBalance := suite.bk.GetBalance(ctx, sdk.AccAddress("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(0)), userBalance)

				// Make sure the pool account balance has increased properly
				poolBalance := suite.bk.GetBalance(ctx, poolstypes.GetPoolAddress(1), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(100)), poolBalance)
			},
		},
		{
			name: "delegating to an existing pool works properly",
			store: func(ctx sdk.Context) {
				// Create the pool
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(20),
					DelegatorShares: sdkmath.LegacyNewDec(100),
				})
				suite.Require().NoError(err)

				// Set the correct pool tokens amount
				suite.fundAccount(
					ctx,
					poolstypes.GetPoolAddress(1).String(),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(20))),
				)

				// Set the next pool id
				suite.pk.SetNextPoolID(ctx, 2)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: false,
			expShares: sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(500))),
			check: func(ctx sdk.Context) {
				// Make sure the pool now exists
				pool, err := suite.pk.GetPool(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(120),
					DelegatorShares: sdkmath.LegacyNewDec(600),
				}, pool)

				// Make sure the delegation exists
				delegation, found, err := suite.k.GetPoolDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(500))),
				), delegation)

				// Make sure the user balance has been reduced properly
				userBalance := suite.bk.GetBalance(ctx, sdk.AccAddress("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(0)), userBalance)

				// Make sure the pool account balance has increased properly
				poolBalance := suite.bk.GetBalance(ctx, poolstypes.GetPoolAddress(1), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(120)), poolBalance)
			},
		},
		{
			name: "delegating more tokens works properly",
			store: func(ctx sdk.Context) {
				// Create the pool
				err := suite.pk.SavePool(ctx, poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(80),
					DelegatorShares: sdkmath.LegacyNewDec(125),
				})
				suite.Require().NoError(err)

				// Set the correct pool tokens amount
				suite.fundAccount(
					ctx,
					poolstypes.GetPoolAddress(1).String(),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(80))),
				)

				// Set the next pool id
				suite.pk.SetNextPoolID(ctx, 2)

				// Save the existing delegation
				err = suite.k.SetDelegation(ctx, types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDec(100))),
				))
				suite.Require().NoError(err)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			amount:    sdk.NewCoin("umilk", sdkmath.NewInt(100)),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: false,
			expShares: sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDecWithPrec(15625, 2))),
			check: func(ctx sdk.Context) {
				// Make sure the pool now exists
				pool, err := suite.pk.GetPool(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(poolstypes.Pool{
					ID:              1,
					Denom:           "umilk",
					Address:         poolstypes.GetPoolAddress(1).String(),
					Tokens:          sdkmath.NewInt(180),
					DelegatorShares: sdkmath.LegacyNewDecWithPrec(28125, 2),
				}, pool)

				// Make sure the delegation has been updated properly
				delegation, found, err := suite.k.GetPoolDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPoolDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("pool/1/umilk", sdkmath.LegacyNewDecWithPrec(25625, 2))), // 100 (existing) + 156.25 (new)
				), delegation)

				// Make sure the user balance has been reduced properly
				userBalance := suite.bk.GetBalance(ctx, sdk.AccAddress("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(0)), userBalance)

				// Make sure the pool account balance has increased properly
				poolBalance := suite.bk.GetBalance(ctx, poolstypes.GetPoolAddress(1), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(180)), poolBalance)
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

			shares, err := suite.k.DelegateToPool(ctx, tc.amount, tc.delegator)
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
