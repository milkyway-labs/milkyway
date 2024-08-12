package keeper_test

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func (s *KeeperTestSuite) TestQuerier_RewardsPlans() {
	service, _ := s.setupSampleServiceAndOperator()

	plan1 := s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	plan2 := s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	testCases := []struct {
		name        string
		req         *types.QueryRewardsPlansRequest
		expectedErr string
		check       func(resp *types.QueryRewardsPlansResponse)
	}{
		{
			name:        "query without pagination returns data properly",
			req:         &types.QueryRewardsPlansRequest{},
			expectedErr: "",
			check: func(resp *types.QueryRewardsPlansResponse) {
				s.Require().Equal([]types.RewardsPlan{plan1, plan2}, resp.RewardsPlans)
			},
		},
		{
			name:        "query with pagination returns data properly",
			req:         &types.QueryRewardsPlansRequest{Pagination: &query.PageRequest{Offset: 1, Limit: 1}},
			expectedErr: "",
			check: func(resp *types.QueryRewardsPlansResponse) {
				s.Require().Equal([]types.RewardsPlan{plan2}, resp.RewardsPlans)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.RewardsPlans(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_RewardsPlan() {
	service, _ := s.setupSampleServiceAndOperator()

	plan := s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	testCases := []struct {
		name        string
		req         *types.QueryRewardsPlanRequest
		expectedErr string
		check       func(resp *types.QueryRewardsPlanResponse)
	}{
		{
			name:        "success",
			req:         &types.QueryRewardsPlanRequest{PlanId: plan.ID},
			expectedErr: "",
			check: func(resp *types.QueryRewardsPlanResponse) {
				s.Require().Equal(plan, resp.RewardsPlan)
			},
		},
		{
			name:        "invalid plan ID returns error",
			req:         &types.QueryRewardsPlanRequest{PlanId: 0},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid plan id",
		},
		{
			name:        "plan not found",
			req:         &types.QueryRewardsPlanRequest{PlanId: 2},
			expectedErr: "rpc error: code = NotFound desc = plan not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.RewardsPlan(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_PoolOutstandingRewards() {
	service, _ := s.setupSampleServiceAndOperator()

	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	delAddr := utils.TestAddress(1)
	s.DelegatePool(utils.MustParseCoin("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		req         *types.QueryPoolOutstandingRewardsRequest
		expectedErr string
		check       func(resp *types.QueryPoolOutstandingRewardsResponse)
	}{
		{
			name:        "success",
			req:         &types.QueryPoolOutstandingRewardsRequest{PoolId: 1},
			expectedErr: "",
			check: func(resp *types.QueryPoolOutstandingRewardsResponse) {
				s.Require().Equal(types.OutstandingRewards{
					Rewards: types.DecPools{
						{
							Denom:    "umilk",
							DecCoins: utils.MustParseDecCoins("11574service"),
						},
					},
				}, resp.Rewards)
			},
		},
		{
			name:        "invalid pool ID returns error",
			req:         &types.QueryPoolOutstandingRewardsRequest{PoolId: 0},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid pool id",
		},
		{
			name:        "pool not found",
			req:         &types.QueryPoolOutstandingRewardsRequest{PoolId: 2},
			expectedErr: "pool not found: not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.PoolOutstandingRewards(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_OperatorOutstandingRewards() {
	service, operator := s.setupSampleServiceAndOperator()

	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	delAddr := utils.TestAddress(1)
	s.DelegateOperator(operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		req         *types.QueryOperatorOutstandingRewardsRequest
		expectedErr string
		check       func(resp *types.QueryOperatorOutstandingRewardsResponse)
	}{
		{
			name:        "success",
			req:         &types.QueryOperatorOutstandingRewardsRequest{OperatorId: operator.ID},
			expectedErr: "",
			check: func(resp *types.QueryOperatorOutstandingRewardsResponse) {
				s.Require().Equal(types.OutstandingRewards{
					Rewards: types.DecPools{
						{
							Denom:    "umilk",
							DecCoins: utils.MustParseDecCoins("11574service"),
						},
					},
				}, resp.Rewards)
			},
		},
		{
			name:        "invalid operator ID returns error",
			req:         &types.QueryOperatorOutstandingRewardsRequest{OperatorId: 0},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid operator id",
		},
		{
			name:        "operator not found",
			req:         &types.QueryOperatorOutstandingRewardsRequest{OperatorId: 2},
			expectedErr: "operator not found: not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.OperatorOutstandingRewards(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_ServiceOutstandingRewards() {
	service, _ := s.setupSampleServiceAndOperator()

	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	delAddr := utils.TestAddress(1)
	s.DelegateService(service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		req         *types.QueryServiceOutstandingRewardsRequest
		expectedErr string
		check       func(resp *types.QueryServiceOutstandingRewardsResponse)
	}{
		{
			name:        "success",
			req:         &types.QueryServiceOutstandingRewardsRequest{ServiceId: service.ID},
			expectedErr: "",
			check: func(resp *types.QueryServiceOutstandingRewardsResponse) {
				s.Require().Equal(types.OutstandingRewards{
					Rewards: types.DecPools{
						{
							Denom:    "umilk",
							DecCoins: utils.MustParseDecCoins("11574service"),
						},
					},
				}, resp.Rewards)
			},
		},
		{
			name:        "invalid service ID returns error",
			req:         &types.QueryServiceOutstandingRewardsRequest{ServiceId: 0},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid service id",
		},
		{
			name:        "service not found",
			req:         &types.QueryServiceOutstandingRewardsRequest{ServiceId: 2},
			expectedErr: "service not found: not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.ServiceOutstandingRewards(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_OperatorCommission() {
	service, operator := s.setupSampleServiceAndOperator()

	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	delAddr := utils.TestAddress(1)
	s.DelegateOperator(operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		req         *types.QueryOperatorCommissionRequest
		expectedErr string
		check       func(resp *types.QueryOperatorCommissionResponse)
	}{
		{
			name:        "success",
			req:         &types.QueryOperatorCommissionRequest{OperatorId: operator.ID},
			expectedErr: "",
			check: func(resp *types.QueryOperatorCommissionResponse) {
				s.Require().Equal(types.AccumulatedCommission{
					Commissions: types.DecPools{
						{
							Denom:    "umilk",
							DecCoins: utils.MustParseDecCoins("1157.4service"),
						},
					},
				}, resp.Commission)
			},
		},
		{
			name:        "invalid operator ID returns error",
			req:         &types.QueryOperatorCommissionRequest{OperatorId: 0},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid operator id",
		},
		{
			name:        "operator not found",
			req:         &types.QueryOperatorCommissionRequest{OperatorId: 2},
			expectedErr: "",
			check: func(resp *types.QueryOperatorCommissionResponse) {
				s.Require().Equal(types.AccumulatedCommission{}, resp.Commission)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.OperatorCommission(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_PoolDelegationRewards() {
	service, _ := s.setupSampleServiceAndOperator()

	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	delAddr := utils.TestAddress(1)
	s.DelegatePool(utils.MustParseCoin("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		req         *types.QueryPoolDelegationRewardsRequest
		expectedErr string
		check       func(resp *types.QueryPoolDelegationRewardsResponse)
	}{
		{
			name: "success",
			req: &types.QueryPoolDelegationRewardsRequest{
				DelegatorAddress: delAddr.String(),
				PoolId:           1,
			},
			expectedErr: "",
			check: func(resp *types.QueryPoolDelegationRewardsResponse) {
				s.Require().Equal(types.DecPools{
					{
						Denom:    "umilk",
						DecCoins: utils.MustParseDecCoins("11574service"),
					},
				}, resp.Rewards)
			},
		},
		{
			name: "invalid delegator address returns error",
			req:  &types.QueryPoolDelegationRewardsRequest{DelegatorAddress: "invalid", PoolId: 1},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid delegator address: decoding bech32 failed:" +
				" invalid bech32 string length 7",
		},
		{
			name: "invalid pool ID returns error",
			req: &types.QueryPoolDelegationRewardsRequest{
				DelegatorAddress: delAddr.String(),
				PoolId:           0,
			},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid pool id",
		},
		{
			name: "pool not found",
			req: &types.QueryPoolDelegationRewardsRequest{
				DelegatorAddress: delAddr.String(),
				PoolId:           2,
			},
			expectedErr: "pool not found: not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.PoolDelegationRewards(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_OperatorDelegationRewards() {
	service, operator := s.setupSampleServiceAndOperator()

	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	delAddr := utils.TestAddress(1)
	s.DelegateOperator(operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		req         *types.QueryOperatorDelegationRewardsRequest
		expectedErr string
		check       func(resp *types.QueryOperatorDelegationRewardsResponse)
	}{
		{
			name: "success",
			req: &types.QueryOperatorDelegationRewardsRequest{
				DelegatorAddress: delAddr.String(),
				OperatorId:       operator.ID,
			},
			expectedErr: "",
			check: func(resp *types.QueryOperatorDelegationRewardsResponse) {
				s.Require().Equal(types.DecPools{
					{
						Denom:    "umilk",
						DecCoins: utils.MustParseDecCoins("10416.6service"), // After deducting commission
					},
				}, resp.Rewards)
			},
		},
		{
			name: "invalid delegator address returns error",
			req: &types.QueryOperatorDelegationRewardsRequest{
				DelegatorAddress: "invalid",
				OperatorId:       operator.ID,
			},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid delegator address: decoding bech32 failed:" +
				" invalid bech32 string length 7",
		},
		{
			name: "invalid operator ID returns error",
			req: &types.QueryOperatorDelegationRewardsRequest{
				DelegatorAddress: delAddr.String(),
				OperatorId:       0,
			},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid operator id",
		},
		{
			name: "operator not found",
			req: &types.QueryOperatorDelegationRewardsRequest{
				DelegatorAddress: delAddr.String(),
				OperatorId:       2,
			},
			expectedErr: "operator not found: not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.OperatorDelegationRewards(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_ServiceDelegationRewards() {
	service, _ := s.setupSampleServiceAndOperator()

	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	delAddr := utils.TestAddress(1)
	s.DelegateService(service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		req         *types.QueryServiceDelegationRewardsRequest
		expectedErr string
		check       func(resp *types.QueryServiceDelegationRewardsResponse)
	}{
		{
			name: "success",
			req: &types.QueryServiceDelegationRewardsRequest{
				DelegatorAddress: delAddr.String(),
				ServiceId:        service.ID,
			},
			expectedErr: "",
			check: func(resp *types.QueryServiceDelegationRewardsResponse) {
				s.Require().Equal(types.DecPools{
					{
						Denom:    "umilk",
						DecCoins: utils.MustParseDecCoins("11574service"),
					},
				}, resp.Rewards)
			},
		},
		{
			name: "invalid delegator address returns error",
			req: &types.QueryServiceDelegationRewardsRequest{
				DelegatorAddress: "invalid",
				ServiceId:        service.ID,
			},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid delegator address: decoding bech32 failed:" +
				" invalid bech32 string length 7",
		},
		{
			name: "invalid service ID returns error",
			req: &types.QueryServiceDelegationRewardsRequest{
				DelegatorAddress: delAddr.String(),
				ServiceId:        0,
			},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid service id",
		},
		{
			name: "service not found",
			req: &types.QueryServiceDelegationRewardsRequest{
				DelegatorAddress: delAddr.String(),
				ServiceId:        2,
			},
			expectedErr: "service not found: not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.ServiceDelegationRewards(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_DelegationTotalRewards() {
	service, operator := s.setupSampleServiceAndOperator()

	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	delAddr := utils.TestAddress(1)
	s.DelegatePool(utils.MustParseCoin("100_000000umilk"), delAddr.String(), true)
	s.DelegateOperator(operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)
	s.DelegateService(service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		req         *types.QueryDelegationTotalRewardsRequest
		expectedErr string
		check       func(resp *types.QueryDelegationTotalRewardsResponse)
	}{
		{
			name: "success",
			req: &types.QueryDelegationTotalRewardsRequest{
				DelegatorAddress: delAddr.String(),
			},
			expectedErr: "",
			check: func(resp *types.QueryDelegationTotalRewardsResponse) {
				s.Require().Equal([]types.DelegationDelegatorReward{
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
						restakingtypes.DELEGATION_TYPE_OPERATOR, operator.ID,
						types.DecPools{
							{
								Denom:    "umilk",
								DecCoins: utils.MustParseDecCoins("3472.2service"),
							},
						},
					),
					types.NewDelegationDelegatorReward(
						restakingtypes.DELEGATION_TYPE_SERVICE, service.ID,
						types.DecPools{
							{
								Denom:    "umilk",
								DecCoins: utils.MustParseDecCoins("3858service"),
							},
						},
					),
				}, resp.Rewards)
				s.Require().Equal(types.DecPools{
					{
						Denom:    "umilk",
						DecCoins: utils.MustParseDecCoins("11188.2service"),
					},
				}, resp.Total)
			},
		},
		{
			name: "invalid delegator address returns error",
			req: &types.QueryDelegationTotalRewardsRequest{
				DelegatorAddress: "invalid",
			},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid delegator address: decoding bech32 failed:" +
				" invalid bech32 string length 7",
		},
		{
			name: "no delegations found",
			req: &types.QueryDelegationTotalRewardsRequest{
				DelegatorAddress: utils.TestAddress(2).String(),
			},
			expectedErr: "",
			check: func(resp *types.QueryDelegationTotalRewardsResponse) {
				s.Require().Empty(resp.Rewards)
				s.Require().Equal("", resp.Total.String())
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			resp, err := s.queryServer.DelegationTotalRewards(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}

func (s *KeeperTestSuite) TestQuerier_DelegatorWithdrawAddress() {
	service, operator := s.setupSampleServiceAndOperator()

	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("10000_000000service"),
	)

	delAddr := utils.TestAddress(1)
	s.DelegatePool(utils.MustParseCoin("100_000000umilk"), delAddr.String(), true)
	s.DelegateOperator(operator.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)
	s.DelegateService(service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(10 * time.Second)

	testCases := []struct {
		name        string
		store       func(ctx context.Context)
		req         *types.QueryDelegatorWithdrawAddressRequest
		expectedErr string
		check       func(resp *types.QueryDelegatorWithdrawAddressResponse)
	}{
		{
			name: "success",
			req: &types.QueryDelegatorWithdrawAddressRequest{
				DelegatorAddress: delAddr.String(),
			},
			expectedErr: "",
			check: func(resp *types.QueryDelegatorWithdrawAddressResponse) {
				s.Require().Equal(delAddr.String(), resp.WithdrawAddress)
			},
		},
		{
			name: "invalid delegator address returns error",
			req: &types.QueryDelegatorWithdrawAddressRequest{
				DelegatorAddress: "invalid",
			},
			expectedErr: "rpc error: code = InvalidArgument desc = invalid delegator address: decoding bech32 failed:" +
				" invalid bech32 string length 7",
		},
		{
			name: "different withdraw address set",
			store: func(ctx context.Context) {
				_, err := s.msgServer.SetWithdrawAddress(ctx, types.NewMsgSetWithdrawAddress(
					delAddr.String(), utils.TestAddress(2).String()))
				s.Require().NoError(err)
			},
			req: &types.QueryDelegatorWithdrawAddressRequest{
				DelegatorAddress: delAddr.String(),
			},
			expectedErr: "",
			check: func(resp *types.QueryDelegatorWithdrawAddressResponse) {
				s.Require().Equal(utils.TestAddress(2).String(), resp.WithdrawAddress)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx, _ := s.Ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}
			resp, err := s.queryServer.DelegatorWithdrawAddress(ctx, tc.req)
			if tc.expectedErr == "" {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check(resp)
				}
			} else {
				s.Require().EqualError(err, tc.expectedErr)
			}
		})
	}
}
