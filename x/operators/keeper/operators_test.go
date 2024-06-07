package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/milkyway-labs/milkyway/x/operators/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetNextOperatorID() {
	testCases := []struct {
		name  string
		store func(ctx sdk.Context)
		id    uint32
		check func(ctx sdk.Context)
	}{
		{
			name: "next operator id is saved correctly",
			id:   1,
			check: func(ctx sdk.Context) {
				store := ctx.KVStore(suite.storeKey)
				operatorIDBz := store.Get(types.NextOperatorIDKey)
				suite.Require().Equal(uint32(1), types.GetOperatorIDFromBytes(operatorIDBz))
			},
		},
		{
			name: "next operator id is overridden properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(ctx, 1)
			},
			id: 2,
			check: func(ctx sdk.Context) {
				store := ctx.KVStore(suite.storeKey)
				operatorIDBz := store.Get(types.NextOperatorIDKey)
				suite.Require().Equal(uint32(2), types.GetOperatorIDFromBytes(operatorIDBz))
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

			suite.k.SetNextOperatorID(ctx, tc.id)
			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetNextOperatorID() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		expNext   uint32
	}{
		{
			name:      "non existing next operator returns error",
			shouldErr: true,
		},
		{
			name: "exiting next operator id is returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextOperatorID(ctx, 1)
			},
			expNext: 1,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			next, err := suite.k.GetNextOperatorID(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expNext, next)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_RegisterOperator() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		operator  types.Operator
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "user without enough funds to pay for registration fee returns erorr",
			store: func(ctx sdk.Context) {
				// Set the registration fee
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(200_000_000))),
					24*time.Hour,
				))

				// Fund the user account
				suite.fundAccount(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))))
			},
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "operator is registered correctly",
			store: func(ctx sdk.Context) {
				// Set the registration fee
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					24*time.Hour,
				))

				// Fund the user account
				suite.fundAccount(ctx, "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4", sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(200_000_000))))
			},
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the operator has been stored
				stored, found := suite.k.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)

				// Make sure the user has been charged
				userAddress, err := sdk.AccAddressFromBech32("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.Require().NoError(err)
				userBalance := suite.bk.GetBalance(ctx, userAddress, "uatom")
				suite.Require().Equal(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000)), userBalance)

				// Make sure the community pool has been funded
				poolBalance := suite.bk.GetBalance(ctx, authtypes.NewModuleAddress(authtypes.FeeCollectorName), "uatom")
				suite.Require().Equal(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000)), poolBalance)

				// Make sure the hook has been called
				suite.Require().True(suite.hooks.CalledMap["AfterOperatorRegistered"])
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

			err := suite.k.RegisterOperator(ctx, tc.operator)
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

func (suite *KeeperTestSuite) TestKeeper_GetOperator() {
	testCases := []struct {
		name        string
		setup       func()
		store       func(ctx sdk.Context)
		operatorID  uint32
		expFound    bool
		expOperator types.Operator
	}{
		{
			name:       "non existing operator returns false",
			operatorID: 1,
			expFound:   false,
		},
		{
			name: "existing operator is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			operatorID: 1,
			expFound:   true,
			expOperator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
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

			operator, found := suite.k.GetOperator(ctx, tc.operatorID)
			suite.Require().Equal(tc.expFound, found)
			suite.Require().Equal(tc.expOperator, operator)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_SaveOperator() {
	testCases := []struct {
		name     string
		setup    func()
		store    func(ctx sdk.Context)
		operator types.Operator
		check    func(ctx sdk.Context)
	}{
		{
			name: "non existing operator is stored properly",
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			check: func(ctx sdk.Context) {
				stored, found := suite.k.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)
			},
		},
		{
			name: "existing operator is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.RegisterOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_INACTIVATING,
				"MilkyWay Updated Operator",
				"https://milkyway.zone",
				"https://milkyway.zone/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			check: func(ctx sdk.Context) {
				//
				stored, found := suite.k.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVATING,
					"MilkyWay Updated Operator",
					"https://milkyway.zone",
					"https://milkyway.zone/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)
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

			suite.k.SaveOperator(ctx, tc.operator)
			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_StartOperatorInactivation() {
	testCases := []struct {
		name     string
		setup    func()
		setupCtx func(ctx sdk.Context) sdk.Context
		store    func(ctx sdk.Context)
		operator types.Operator
		check    func(ctx sdk.Context)
	}{
		{
			name: "operator inactivation is started properly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					12*time.Hour,
				))
			},
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_ACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			check: func(ctx sdk.Context) {
				// Make sure the operator status has been updated
				stored, found := suite.k.GetOperator(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVATING,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)

				// Make sure the operator has been inserted into the inactivating queue
				inactivatingQueue := suite.k.GetInactivatingOperators(ctx)
				suite.Require().Len(inactivatingQueue, 1)
				suite.Require().Equal(types.NewUnbondingOperator(
					1,
					time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				), inactivatingQueue[0])

				// Make sure the hook has been called
				suite.Require().True(suite.hooks.CalledMap["AfterOperatorInactivatingStarted"])
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
			if tc.setupCtx != nil {
				ctx = tc.setupCtx(ctx)
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			suite.k.StartOperatorInactivation(ctx, tc.operator)
			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
