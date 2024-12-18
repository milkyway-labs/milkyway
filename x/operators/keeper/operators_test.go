package keeper_test

import (
	"time"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v5/x/operators/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetNextOperatorID() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		id        uint32
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "next operator id is saved correctly",
			id:   1,
			check: func(ctx sdk.Context) {
				nextOperatorID, err := suite.k.GetNextOperatorID(ctx)
				suite.Require().NoError(err)
				suite.Require().EqualValues(1, nextOperatorID)
			},
		},
		{
			name: "next operator id is overridden properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextOperatorID(ctx, 1)
				suite.Require().NoError(err)
			},
			id: 2,
			check: func(ctx sdk.Context) {
				nextOperatorID, err := suite.k.GetNextOperatorID(ctx)
				suite.Require().NoError(err)
				suite.Require().EqualValues(2, nextOperatorID)
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

			err := suite.k.SetNextOperatorID(ctx, tc.id)
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

func (suite *KeeperTestSuite) TestKeeper_GetNextOperatorID() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		expNext   uint32
	}{
		{
			name:      "non existing next service returns 1",
			shouldErr: false,
			expNext:   1,
		},
		{
			name: "exiting next operator id is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextOperatorID(ctx, 2)
				suite.Require().NoError(err)
			},
			expNext: 2,
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

func (suite *KeeperTestSuite) TestKeeper_CreateOperator() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		operator  types.Operator
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "operator is registered correctly",
			store: func(ctx sdk.Context) {
				// Set the registration fee
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					24*time.Hour,
				))
				suite.Require().NoError(err)

				// Fund the user account
				suite.fundAccount(
					ctx,
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(200_000_000))),
				)
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
				stored, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)

				// Make sure the hook has been called
				suite.Require().True(suite.hooks.CalledMap["AfterOperatorRegistered"])
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			err := suite.k.CreateOperator(ctx, tc.operator)
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
				err := suite.k.CreateOperator(ctx, types.NewOperator(
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
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			operator, err := suite.k.GetOperator(ctx, tc.operatorID)
			if tc.expFound {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expOperator, operator)
			} else {
				suite.Require().ErrorIs(err, collections.ErrNotFound)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_SaveOperator() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		operator  types.Operator
		shouldErr bool
		check     func(ctx sdk.Context)
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
			shouldErr: false,
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
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
				err := suite.k.CreateOperator(ctx, types.NewOperator(
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
			shouldErr: false,
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
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
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			err := suite.k.SaveOperator(ctx, tc.operator)
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

func (suite *KeeperTestSuite) TestKeeper_StartOperatorInactivation() {
	testCases := []struct {
		name      string
		setup     func()
		setupCtx  func(ctx sdk.Context) sdk.Context
		store     func(ctx sdk.Context)
		operator  types.Operator
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "inactivating operator returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					12*time.Hour,
				))
				suite.Require().NoError(err)
			},
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_INACTIVATING,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "inactive operator returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					12*time.Hour,
				))
				suite.Require().NoError(err)
			},
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_INACTIVE,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			shouldErr: true,
		},
		{
			name: "operator inactivation is started properly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					12*time.Hour,
				))
				suite.Require().NoError(err)
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
				stored, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVATING,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)

				// Make sure the operator has been inserted into the inactivating queue
				inactivatingQueue, _ := suite.k.GetInactivatingOperators(ctx)
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
			suite.SetupTest()

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

			err := suite.k.StartOperatorInactivation(ctx, tc.operator)
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

func (suite *KeeperTestSuite) TestKeeper_CompleteOperatorInactivation() {
	testCases := []struct {
		name      string
		setup     func()
		setupCtx  func(ctx sdk.Context) sdk.Context
		store     func(ctx sdk.Context)
		operator  types.Operator
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "operator inactivation is completed properly",
			setupCtx: func(ctx sdk.Context) sdk.Context {
				return ctx.WithBlockTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
			},
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(100_000_000))),
					12*time.Hour,
				))
				suite.Require().NoError(err)

				err = suite.k.SaveOperatorParams(ctx, 1, types.NewOperatorParams(
					sdkmath.LegacyNewDec(100),
				))
				suite.Require().NoError(err)
			},
			operator: types.NewOperator(
				1,
				types.OPERATOR_STATUS_INACTIVATING,
				"MilkyWay Operator",
				"https://milkyway.com",
				"https://milkyway.com/picture",
				"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
			),
			check: func(ctx sdk.Context) {
				// Make sure the operator status has been updated
				stored, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), stored)

				// Make sure the operator has been removed from the inactivating queue
				inactivatingQueue, _ := suite.k.GetInactivatingOperators(ctx)
				suite.Require().Len(inactivatingQueue, 0)

				// Make sure the params are no longer there
				operatorParams, err := suite.k.GetOperatorParams(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.DefaultOperatorParams(), operatorParams)

				// Make sure the hook has been called
				suite.Require().True(suite.hooks.CalledMap["AfterOperatorInactivatingCompleted"])
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

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

			err := suite.k.CompleteOperatorInactivation(ctx, tc.operator)
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

func (suite *KeeperTestSuite) TestKeeper_ReactivateInactiveOperator() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		operatorID uint32
		shouldErr  bool
		check      func(ctx sdk.Context)
	}{
		{
			name: "reactivate active operator fails",
			store: func(ctx sdk.Context) {
				err := suite.k.CreateOperator(ctx, types.NewOperator(
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
			shouldErr:  true,
		},
		{
			name: "reactivate inactivating operator fails",
			store: func(ctx sdk.Context) {
				err := suite.k.CreateOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVATING,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			operatorID: 1,
			shouldErr:  true,
		},
		{
			name: "reactivate inactive operator works properly",
			store: func(ctx sdk.Context) {
				err := suite.k.CreateOperator(ctx, types.NewOperator(
					1,
					types.OPERATOR_STATUS_INACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				))
				suite.Require().NoError(err)
			},
			operatorID: 1,
			shouldErr:  false,
			check: func(ctx sdk.Context) {
				operator, err := suite.k.GetOperator(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewOperator(
					1,
					types.OPERATOR_STATUS_ACTIVE,
					"MilkyWay Operator",
					"https://milkyway.com",
					"https://milkyway.com/picture",
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
				), operator)

				// Check hook called
				called := suite.hooks.CalledMap["AfterOperatorReactivated"]
				suite.Require().True(called)
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
			operator, err := suite.k.GetOperator(ctx, tc.operatorID)
			suite.Require().NoError(err)

			err = suite.k.ReactivateInactiveOperator(ctx, operator)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check(ctx)
				}
			}
		})
	}
}
