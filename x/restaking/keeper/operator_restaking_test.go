package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_SaveOperatorDelegation() {
	testCases := []struct {
		name       string
		setup      func()
		store      func(ctx sdk.Context)
		delegation types.OperatorDelegation
		check      func(ctx sdk.Context)
	}{
		{
			name: "operator delegation is stored properly",
			delegation: types.NewOperatorDelegation(
				1,
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
			),
			check: func(ctx sdk.Context) {
				store := ctx.KVStore(suite.storeKey)

				// Make sure the user-operator delegation key exists and contains the delegation
				delegationBz := store.Get(types.UserOperatorDelegationStoreKey("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", 1))
				suite.Require().NotNil(delegationBz)

				delegation, err := types.UnmarshalOperatorDelegation(suite.cdc, delegationBz)
				suite.Require().NoError(err)

				suite.Require().Equal(types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
				), delegation)

				// Make sure the operator-user delegation key exists
				hasDelegationsByOperatorKey := store.Has(types.DelegationByOperatorIDStoreKey(1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"))
				suite.Require().True(hasDelegationsByOperatorKey)
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

			suite.k.SaveOperatorDelegation(ctx, tc.delegation)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetOperatorDelegation() {
	testCases := []struct {
		name          string
		setup         func()
		store         func(ctx sdk.Context)
		operatorID    uint32
		userAddress   string
		expFound      bool
		expDelegation types.OperatorDelegation
		check         func(ctx sdk.Context)
	}{
		{
			name:        "not found delegation returns false",
			operatorID:  1,
			userAddress: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			expFound:    false,
		},
		{
			name: "found delegation is returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
					sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
				))
			},
			operatorID:  1,
			userAddress: "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			expFound:    true,
			expDelegation: types.NewOperatorDelegation(
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

			delegation, found := suite.k.GetOperatorDelegation(ctx, tc.operatorID, tc.userAddress)
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

func (suite *KeeperTestSuite) TestKeeper_AddOperatorTokensAndShares() {
	testCases := []struct {
		name           string
		setup          func()
		store          func(ctx sdk.Context)
		operator       operatorstypes.Operator
		tokensToAdd    sdk.Coins
		shouldErr      bool
		expOperator    operatorstypes.Operator
		expAddedShares sdk.DecCoins
		check          func(ctx sdk.Context)
	}{
		{
			name: "adding tokens to an empty operator works properly",
			operator: operatorstypes.Operator{
				ID:      1,
				Address: operatorstypes.GetOperatorAddress(1).String(),
			},
			tokensToAdd: sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			shouldErr:   false,
			expOperator: operatorstypes.Operator{
				ID:      1,
				Address: operatorstypes.GetOperatorAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(100)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
				),
			},
			expAddedShares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
			),
		},
		{
			name: "adding tokens to a non-empty operator works properly",
			operator: operatorstypes.Operator{
				ID:      1,
				Address: operatorstypes.GetOperatorAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(50)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
				),
			},
			tokensToAdd: sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(20))),
			shouldErr:   false,
			expOperator: operatorstypes.Operator{
				ID:      1,
				Address: operatorstypes.GetOperatorAddress(1).String(),
				Tokens: sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(70)),
				),
				DelegatorShares: sdk.NewDecCoins(
					sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(140)),
				),
			},
			expAddedShares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(40)),
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

			operator, addedShares, err := suite.k.AddOperatorTokensAndShares(ctx, tc.operator, tc.tokensToAdd)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expOperator, operator)
				suite.Require().Equal(tc.expAddedShares, addedShares)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_DelegateToOperator() {
	testCases := []struct {
		name       string
		setup      func()
		store      func(ctx sdk.Context)
		operatorID uint32
		amount     sdk.Coins
		delegator  string
		shouldErr  bool
		expShares  sdk.DecCoins
		check      func(ctx sdk.Context)
	}{
		{
			name:       "operator not found returns error",
			operatorID: 1,
			amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr:  true,
		},
		{
			name: "inactive operator returns error",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:              1,
					Status:          operatorstypes.OPERATOR_STATUS_INACTIVE,
					Address:         operatorstypes.GetOperatorAddress(1).String(),
					Tokens:          sdk.NewCoins(),
					DelegatorShares: sdk.NewDecCoins(),
				})
			},
			operatorID: 1,
			amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr:  true,
		},
		{
			name: "invalid exchange rate operator returns error",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:              1,
					Status:          operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address:         operatorstypes.GetOperatorAddress(1).String(),
					Tokens:          sdk.NewCoins(),
					DelegatorShares: sdk.NewDecCoins(sdk.NewDecCoinFromDec("umilk", sdkmath.LegacyNewDec(100))),
				})
			},
			operatorID: 1,
			amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr:  true,
		},
		{
			name: "invalid delegator address returns error",
			store: func(ctx sdk.Context) {
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
				})
			},
			operatorID: 1,
			amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator:  "invalid",
			shouldErr:  true,
		},
		{
			name: "insufficient funds return error",
			store: func(ctx sdk.Context) {
				// Create the operator
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Address: operatorstypes.GetOperatorAddress(1).String(),
				})

				// Set the next operator id
				suite.ok.SetNextOperatorID(ctx, 2)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(50))),
				)
			},
			operatorID: 1,
			amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr:  true,
		},
		{
			name: "delegating to an existing operator works properly",
			store: func(ctx sdk.Context) {
				// Create the operator
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(20)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
					),
				})

				// Set the correct operator tokens amount
				suite.fundAccount(
					ctx,
					operatorstypes.GetOperatorAddress(1).String(),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(20))),
				)

				// Set the next operator id
				suite.ok.SetNextOperatorID(ctx, 2)

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
				)
			},
			operatorID: 1,
			amount:     sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(100))),
			delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr:  false,
			expShares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(500)),
			),
			check: func(ctx sdk.Context) {
				// Make sure the operator now exists
				operator, found := suite.ok.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(120)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(600)),
					),
				}, operator)

				// Make sure the delegation exists
				delegation, found := suite.k.GetOperatorDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(500)),
					),
				), delegation)

				// Make sure the user balance has been reduced properly
				userBalance := suite.bk.GetBalance(ctx, sdk.AccAddress("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(0)), userBalance)

				// Make sure the operator account balance has increased properly
				operatorBalance := suite.bk.GetBalance(ctx, operatorstypes.GetOperatorAddress(1), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(120)), operatorBalance)
			},
		},
		{
			name: "delegating another token denom works properly",
			store: func(ctx sdk.Context) {
				// Create the operator
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(80)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(125)),
					),
				})

				// Set the correct operator tokens amount
				suite.fundAccount(
					ctx,
					operatorstypes.GetOperatorAddress(1).String(),
					sdk.NewCoins(sdk.NewCoin("umilk", sdkmath.NewInt(80))),
				)

				// Set the next operator id
				suite.ok.SetNextOperatorID(ctx, 2)

				// Save the existing delegation
				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(125)),
					),
				))

				// Send some funds to the user
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("uinit", sdkmath.NewInt(100))),
				)
			},
			operatorID: 1,
			amount:     sdk.NewCoins(sdk.NewCoin("uinit", sdkmath.NewInt(100))),
			delegator:  "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr:  false,
			expShares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("operator/1/uinit", sdkmath.LegacyNewDec(100)),
			),
			check: func(ctx sdk.Context) {
				// Make sure the operator now exists
				operator, found := suite.ok.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(80)),
						sdk.NewCoin("uinit", sdkmath.NewInt(100)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(125)),
						sdk.NewDecCoinFromDec("operator/1/uinit", sdkmath.LegacyNewDec(100)),
					),
				}, operator)

				// Make sure the delegation has been updated properly
				delegation, found := suite.k.GetOperatorDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(125)),
						sdk.NewDecCoinFromDec("operator/1/uinit", sdkmath.LegacyNewDec(100)),
					),
				), delegation)

				// Make sure the user balance has been reduced properly
				userBalance := suite.bk.GetBalance(ctx, sdk.AccAddress("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(0)), userBalance)

				// Make sure the operator account balance has increased properly
				operatorBalance := suite.bk.GetAllBalances(ctx, operatorstypes.GetOperatorAddress(1))
				suite.Require().Equal(sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(80)),
					sdk.NewCoin("uinit", sdkmath.NewInt(100)),
				), operatorBalance)
			},
		},
		{
			name: "delegating more tokens works properly",
			store: func(ctx sdk.Context) {
				// Create the operator
				suite.ok.SaveOperator(ctx, operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(80)),
						sdk.NewCoin("uinit", sdkmath.NewInt(75)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(125)),
						sdk.NewDecCoinFromDec("operator/1/uinit", sdkmath.LegacyNewDec(200)),
					),
				})

				// Set the correct operator tokens amount
				suite.fundAccount(
					ctx,
					operatorstypes.GetOperatorAddress(1).String(),
					sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(80)),
						sdk.NewCoin("uinit", sdkmath.NewInt(75)),
					),
				)

				// Set the next operator id
				suite.ok.SetNextOperatorID(ctx, 2)

				// Save the existing delegation
				suite.k.SaveOperatorDelegation(ctx, types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDec(100)),
						sdk.NewDecCoinFromDec("operator/1/uinit", sdkmath.LegacyNewDec(60)),
					),
				))

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
			operatorID: 1,
			amount: sdk.NewCoins(
				sdk.NewCoin("umilk", sdkmath.NewInt(100)),
				sdk.NewCoin("uinit", sdkmath.NewInt(225)),
			),
			delegator: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			shouldErr: false,
			expShares: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDecWithPrec(15625, 2)),
				sdk.NewDecCoinFromDec("operator/1/uinit", sdkmath.LegacyNewDec(600)),
			),
			check: func(ctx sdk.Context) {
				// Make sure the operator now exists
				operator, found := suite.ok.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(operatorstypes.Operator{
					ID:      1,
					Status:  operatorstypes.OPERATOR_STATUS_ACTIVE,
					Address: operatorstypes.GetOperatorAddress(1).String(),
					Tokens: sdk.NewCoins(
						sdk.NewCoin("umilk", sdkmath.NewInt(180)),
						sdk.NewCoin("uinit", sdkmath.NewInt(300)),
					),
					DelegatorShares: sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDecWithPrec(28125, 2)),
						sdk.NewDecCoinFromDec("operator/1/uinit", sdkmath.LegacyNewDec(800)),
					),
				}, operator)

				// Make sure the delegation has been updated properly
				delegation, found := suite.k.GetOperatorDelegation(ctx, 1, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperatorDelegation(
					1,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewDecCoins(
						sdk.NewDecCoinFromDec("operator/1/umilk", sdkmath.LegacyNewDecWithPrec(25625, 2)), // 100 (existing) + 156.25 (new)
						sdk.NewDecCoinFromDec("operator/1/uinit", sdkmath.LegacyNewDec(660)),              // 60 (existing) + 600 (new)
					),
				), delegation)

				// Make sure the user balance has been reduced properly
				userBalance := suite.bk.GetBalance(ctx, sdk.AccAddress("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"), "umilk")
				suite.Require().Equal(sdk.NewCoin("umilk", sdkmath.NewInt(0)), userBalance)

				// Make sure the operator account balance has increased properly
				operatorBalance := suite.bk.GetAllBalances(ctx, operatorstypes.GetOperatorAddress(1))
				suite.Require().Equal(sdk.NewCoins(
					sdk.NewCoin("umilk", sdkmath.NewInt(180)),
					sdk.NewCoin("uinit", sdkmath.NewInt(300)),
				), operatorBalance)
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

			shares, err := suite.k.DelegateToOperator(ctx, tc.operatorID, tc.amount, tc.delegator)
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