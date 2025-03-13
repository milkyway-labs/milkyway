package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/v10/app/testutil"
	"github.com/milkyway-labs/milkyway/v10/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/v10/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v10/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/v10/x/rewards/types"
)

func (suite *KeeperTestSuite) TestQuerier_RewardsPlans() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		req       *types.QueryRewardsPlansRequest
		shouldErr bool
		expPlans  []types.RewardsPlan
	}{
		{
			name: "query without pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.keeper.RewardsPlans.Set(ctx, 1, types.NewRewardsPlan(
					1,
					"Plan 1",
					1,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				))
				suite.Require().NoError(err)

				err = suite.keeper.RewardsPlans.Set(ctx, 2, types.NewRewardsPlan(
					2,
					"Plan 2",
					2,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				))
				suite.Require().NoError(err)
			},
			req:       types.NewQueryRewardsPlansRequest(nil),
			shouldErr: false,
			expPlans: []types.RewardsPlan{
				types.NewRewardsPlan(
					1,
					"Plan 1",
					1,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				),
				types.NewRewardsPlan(
					2,
					"Plan 2",
					2,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				),
			},
		},
		{
			name: "query with pagination returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.keeper.RewardsPlans.Set(ctx, 1, types.NewRewardsPlan(
					1,
					"Plan 1",
					1,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				))
				suite.Require().NoError(err)

				err = suite.keeper.RewardsPlans.Set(ctx, 2, types.NewRewardsPlan(
					2,
					"Plan 2",
					2,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				))
				suite.Require().NoError(err)
			},
			req: types.NewQueryRewardsPlansRequest(&query.PageRequest{
				Limit:  1,
				Offset: 1,
			}),
			shouldErr: false,
			expPlans: []types.RewardsPlan{
				types.NewRewardsPlan(
					2,
					"Plan 2",
					2,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				),
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.RewardsPlans(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				for i, plan := range res.RewardsPlans {
					suite.Require().True(tc.expPlans[i].Equal(plan))
				}
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_RewardsPlan() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		req       *types.QueryRewardsPlanRequest
		shouldErr bool
		expPlan   types.RewardsPlan
	}{
		{
			name:      "invalid plan ID returns error",
			req:       types.NewQueryRewardsPlanRequest(0),
			shouldErr: true,
		},
		{
			name:      "not found plan returns error",
			req:       types.NewQueryRewardsPlanRequest(1),
			shouldErr: true,
		},
		{
			name: "found plan returns data properly",
			store: func(ctx sdk.Context) {
				err := suite.keeper.RewardsPlans.Set(ctx, 1, types.NewRewardsPlan(
					1,
					"Plan 1",
					1,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					types.NewEgalitarianPoolsDistribution(1),
					types.NewEgalitarianOperatorsDistribution(1),
					types.NewBasicUsersDistribution(1),
				))
				suite.Require().NoError(err)
			},
			req:       types.NewQueryRewardsPlanRequest(1),
			shouldErr: false,
			expPlan: types.NewRewardsPlan(
				1,
				"Plan 1",
				1,
				utils.MustParseCoin("100_000000service"),
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				types.NewEgalitarianPoolsDistribution(1),
				types.NewEgalitarianOperatorsDistribution(1),
				types.NewBasicUsersDistribution(1),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.RewardsPlan(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().True(tc.expPlan.Equal(res.RewardsPlan))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_PoolOutstandingRewards() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		updateCtx  func(ctx sdk.Context) sdk.Context
		req        *types.QueryPoolOutstandingRewardsRequest
		shouldErr  bool
		expRewards types.OutstandingRewards
	}{
		{
			name:      "invalid pool ID returns error",
			req:       types.NewQueryPoolOutstandingRewardsRequest(0),
			shouldErr: true,
		},
		{
			name:      "pool not found",
			req:       types.NewQueryPoolOutstandingRewardsRequest(1),
			shouldErr: true,
		},
		{
			name: "existing outstanding rewards are returned properly",
			store: func(ctx sdk.Context) {
				// Create a service and a rewards plan
				service, _ := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to the pool
				delAddr := testutil.TestAddress(1)
				suite.DelegatePool(ctx, utils.MustParseCoin("100_000000umilk"), delAddr.String(), true)
				suite.SetUserPreferences(ctx, delAddr.String(), []restakingtypes.TrustedServiceEntry{
					restakingtypes.NewTrustedServiceEntry(service.ID, nil),
				})
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			req:       types.NewQueryPoolOutstandingRewardsRequest(1),
			shouldErr: false,
			expRewards: types.OutstandingRewards{
				Rewards: types.DecPools{
					{
						Denom:    "umilk",
						DecCoins: utils.MustParseDecCoins("11574service"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.PoolOutstandingRewards(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expRewards, res.Rewards)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_OperatorOutstandingRewards() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		updateCtx  func(ctx sdk.Context) sdk.Context
		req        *types.QueryOperatorOutstandingRewardsRequest
		shouldErr  bool
		expRewards types.OutstandingRewards
	}{
		{
			name:      "invalid operator ID returns error",
			req:       types.NewQueryOperatorOutstandingRewardsRequest(0),
			shouldErr: true,
		},
		{
			name:      "operator not found returns error",
			req:       types.NewQueryOperatorOutstandingRewardsRequest(2),
			shouldErr: true,
		},
		{
			name: "existing rewards are returned properly",
			store: func(ctx sdk.Context) {
				// Create a service and a rewards plan
				service, operator := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to the operator
				delAddr := testutil.TestAddress(1)
				suite.DelegateOperator(ctx, operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			req:       types.NewQueryOperatorOutstandingRewardsRequest(1),
			shouldErr: false,
			expRewards: types.OutstandingRewards{
				Rewards: types.DecPools{
					{
						Denom:    "umilk",
						DecCoins: utils.MustParseDecCoins("11574service"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.OperatorOutstandingRewards(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expRewards, res.Rewards)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_ServiceOutstandingRewards() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		updateCtx  func(ctx sdk.Context) sdk.Context
		req        *types.QueryServiceOutstandingRewardsRequest
		shouldErr  bool
		expRewards types.OutstandingRewards
	}{
		{
			name:      "invalid service ID returns error",
			req:       types.NewQueryServiceOutstandingRewardsRequest(0),
			shouldErr: true,
		},
		{
			name:      "service not found",
			req:       types.NewQueryServiceOutstandingRewardsRequest(1),
			shouldErr: true,
		},
		{
			name: "exiting rewards are returned properly",
			store: func(ctx sdk.Context) {
				// Create a service and a rewards plan
				service, _ := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to the service
				delAddr := testutil.TestAddress(1)
				suite.DelegateService(ctx, service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			req:       types.NewQueryServiceOutstandingRewardsRequest(1),
			shouldErr: false,
			expRewards: types.OutstandingRewards{
				Rewards: types.DecPools{
					{
						Denom:    "umilk",
						DecCoins: utils.MustParseDecCoins("11574service"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.ServiceOutstandingRewards(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expRewards, res.Rewards)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_OperatorCommission() {
	testCases := []struct {
		name          string
		store         func(ctx sdk.Context)
		updateCtx     func(ctx sdk.Context) sdk.Context
		req           *types.QueryOperatorCommissionRequest
		shouldErr     bool
		expCommission types.AccumulatedCommission
	}{
		{
			name:      "invalid operator ID returns error",
			req:       types.NewQueryOperatorCommissionRequest(0),
			shouldErr: true,
		},
		{
			name:          "operator not found returns zero",
			req:           types.NewQueryOperatorCommissionRequest(1),
			shouldErr:     false,
			expCommission: types.AccumulatedCommission{},
		},
		{
			name: "existing commission is returned properly",
			store: func(ctx sdk.Context) {
				// Create a service and a rewards plan
				service, operator := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to the operator
				delAddr := testutil.TestAddress(1)
				suite.DelegateOperator(ctx, operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			req:       types.NewQueryOperatorCommissionRequest(1),
			shouldErr: false,
			expCommission: types.AccumulatedCommission{
				Commissions: types.DecPools{
					{
						Denom:    "umilk",
						DecCoins: utils.MustParseDecCoins("1157.4service"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.OperatorCommission(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expCommission, res.Commission)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_PoolDelegationRewards() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		updateCtx  func(ctx sdk.Context) sdk.Context
		req        *types.QueryPoolDelegationRewardsRequest
		shouldErr  bool
		expRewards types.DecPools
	}{
		{
			name:      "invalid delegator address returns error",
			req:       types.NewQueryPoolDelegationRewardsRequest(1, "invalid"),
			shouldErr: true,
		},
		{
			name:      "invalid pool ID returns error",
			req:       types.NewQueryPoolDelegationRewardsRequest(0, testutil.TestAddress(1).String()),
			shouldErr: true,
		},
		{
			name:      "pool not found returns error",
			req:       types.NewQueryPoolDelegationRewardsRequest(1, testutil.TestAddress(1).String()),
			shouldErr: true,
		},
		{
			name: "existing rewards are returned properly",
			store: func(ctx sdk.Context) {
				// Create a service and a rewards plan
				service, _ := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to the pool
				delAddr := testutil.TestAddress(1)
				suite.DelegatePool(ctx, utils.MustParseCoin("100_000000umilk"), delAddr.String(), true)
				suite.SetUserPreferences(ctx, delAddr.String(), []restakingtypes.TrustedServiceEntry{
					restakingtypes.NewTrustedServiceEntry(service.ID, nil),
				})
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			req: types.NewQueryPoolDelegationRewardsRequest(
				1, testutil.TestAddress(1).String(),
			),
			shouldErr: false,
			expRewards: types.DecPools{
				{
					Denom:    "umilk",
					DecCoins: utils.MustParseDecCoins("11574service"),
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.PoolDelegationRewards(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expRewards, res.Rewards)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_OperatorDelegationRewards() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		updateCtx  func(ctx sdk.Context) sdk.Context
		req        *types.QueryOperatorDelegationRewardsRequest
		shouldErr  bool
		expRewards types.DecPools
	}{
		{
			name:      "invalid delegator address returns error",
			req:       types.NewQueryOperatorDelegationRewardsRequest(1, "invalid"),
			shouldErr: true,
		},
		{
			name:      "invalid operator ID returns error",
			req:       types.NewQueryOperatorDelegationRewardsRequest(0, testutil.TestAddress(1).String()),
			shouldErr: true,
		},
		{
			name:      "operator not found returns error",
			req:       types.NewQueryOperatorDelegationRewardsRequest(1, testutil.TestAddress(1).String()),
			shouldErr: true,
		},
		{
			name: "existing rewards are returned properly",
			store: func(ctx sdk.Context) {
				// Create a service and a rewards plan
				service, operator := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to the operator
				delAddr := testutil.TestAddress(1)
				suite.DelegateOperator(ctx, operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			req:       types.NewQueryOperatorDelegationRewardsRequest(1, testutil.TestAddress(1).String()),
			shouldErr: false,
			expRewards: types.DecPools{
				{
					Denom:    "umilk",
					DecCoins: utils.MustParseDecCoins("10416.6service"), // After deducting commission
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.OperatorDelegationRewards(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expRewards, res.Rewards)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_ServiceDelegationRewards() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		updateCtx  func(ctx sdk.Context) sdk.Context
		req        *types.QueryServiceDelegationRewardsRequest
		shouldErr  bool
		expRewards types.DecPools
	}{
		{
			name:      "invalid delegator address returns error",
			req:       types.NewQueryServiceDelegationRewardsRequest(1, "invalid"),
			shouldErr: true,
		},
		{
			name:      "invalid service ID returns error",
			req:       types.NewQueryServiceDelegationRewardsRequest(0, testutil.TestAddress(1).String()),
			shouldErr: true,
		},
		{
			name:      "service not found returns error",
			req:       types.NewQueryServiceDelegationRewardsRequest(1, testutil.TestAddress(1).String()),
			shouldErr: true,
		},
		{
			name: "existing rewards are returned properly",
			store: func(ctx sdk.Context) {
				// Create a service and a rewards plan
				service, _ := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to the service
				delAddr := testutil.TestAddress(1)
				suite.DelegateService(ctx, service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			req:       types.NewQueryServiceDelegationRewardsRequest(1, testutil.TestAddress(1).String()),
			shouldErr: false,
			expRewards: types.DecPools{
				{
					Denom:    "umilk",
					DecCoins: utils.MustParseDecCoins("11574service"),
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.ServiceDelegationRewards(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expRewards, res.Rewards)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorTotalRewards() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		updateCtx func(ctx sdk.Context) sdk.Context
		req       *types.QueryDelegatorTotalRewardsRequest
		shouldErr bool
		expRes    *types.QueryDelegatorTotalRewardsResponse
	}{
		{
			name:      "invalid delegator address returns error",
			req:       types.NewQueryDelegatorTotalRewardsRequest("invalid"),
			shouldErr: true,
		},
		{
			name:      "no delegations found return empty response",
			req:       types.NewQueryDelegatorTotalRewardsRequest("cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"),
			shouldErr: false,
			expRes: &types.QueryDelegatorTotalRewardsResponse{
				Rewards: nil,
				Total:   types.DecPools{},
			},
		},
		{
			name: "existing rewards are returned properly",
			store: func(ctx sdk.Context) {
				// Create a service and operator, and a rewards plan
				service, operator := suite.setupSampleServiceAndOperator(ctx)
				suite.CreateBasicRewardsPlan(
					ctx,
					service.ID,
					utils.MustParseCoin("100_000000service"),
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					utils.MustParseCoins("10000_000000service"),
				)

				// Delegate to pools, operators and services
				delAddr := testutil.TestAddress(1)
				suite.DelegatePool(ctx, utils.MustParseCoin("100_000000umilk"), delAddr.String(), true)
				suite.SetUserPreferences(ctx, delAddr.String(), []restakingtypes.TrustedServiceEntry{
					restakingtypes.NewTrustedServiceEntry(service.ID, nil),
				})
				suite.DelegateOperator(ctx, operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)
				suite.DelegateService(ctx, service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)
			},
			updateCtx: func(ctx sdk.Context) sdk.Context {
				// Allocate rewards
				return suite.allocateRewards(ctx, 10*time.Second)
			},
			req:       types.NewQueryDelegatorTotalRewardsRequest(testutil.TestAddress(1).String()),
			shouldErr: false,
			expRes: &types.QueryDelegatorTotalRewardsResponse{
				Rewards: []types.DelegationDelegatorReward{
					types.NewDelegationDelegatorReward(
						restakingtypes.DELEGATION_TYPE_POOL, 1,
						types.DecPools{
							{
								Denom:    "umilk",
								DecCoins: utils.MustParseDecCoins("3858service"),
							},
						},
					),
					types.NewDelegationDelegatorReward(
						restakingtypes.DELEGATION_TYPE_OPERATOR, 1,
						types.DecPools{
							{
								Denom:    "umilk",
								DecCoins: utils.MustParseDecCoins("3472.2service"),
							},
						},
					),
					types.NewDelegationDelegatorReward(
						restakingtypes.DELEGATION_TYPE_SERVICE, 1,
						types.DecPools{
							{
								Denom:    "umilk",
								DecCoins: utils.MustParseDecCoins("3858service"),
							},
						},
					),
				},
				Total: types.DecPools{
					{
						Denom:    "umilk",
						DecCoins: utils.MustParseDecCoins("11188.2service"),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			if tc.updateCtx != nil {
				ctx = tc.updateCtx(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.DelegatorTotalRewards(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expRes, res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_DelegatorWithdrawAddress() {
	testCases := []struct {
		name               string
		store              func(ctx sdk.Context)
		req                *types.QueryDelegatorWithdrawAddressRequest
		shouldErr          bool
		expWithdrawAddress string
	}{
		{
			name:      "invalid delegator address returns error",
			req:       types.NewQueryDelegatorWithdrawAddressRequest("invalid"),
			shouldErr: true,
		},
		{
			name: "delegator without custom address returns default address",
			req: types.NewQueryDelegatorWithdrawAddressRequest(
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr:          false,
			expWithdrawAddress: "cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
		},
		{
			name: "delegator with different withdraw address set returns proper value",
			store: func(ctx sdk.Context) {
				msgServer := keeper.NewMsgServer(suite.keeper)
				_, err := msgServer.SetWithdrawAddress(ctx, types.NewMsgSetWithdrawAddress(
					"cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
					"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
				))
				suite.Require().NoError(err)
			},
			req: types.NewQueryDelegatorWithdrawAddressRequest(
				"cosmos13t6y2nnugtshwuy0zkrq287a95lyy8vzleaxmd",
			),
			shouldErr:          false,
			expWithdrawAddress: "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.DelegatorWithdrawAddress(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expWithdrawAddress, res.WithdrawAddress)
			}
		})
	}
}
