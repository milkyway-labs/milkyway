package keeper_test

import (
	"time"

	"cosmossdk.io/math"

	"github.com/milkyway-labs/milkyway/utils"
)

func (s *KeeperTestSuite) TestAllocateRewards_InactivePlan() {
	// Inactive plans(current block time is out of their date range) don't
	// allocate rewards.

	// Plan's start time is 2024-01-01 so set block time before that.
	s.Ctx = s.Ctx.WithBlockTime(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC))
	service, _, _ := s.setupSimpleScenario()

	delAddr := utils.TestAddress(1)
	s.DelegateService(service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	s.advanceBlock(3 * time.Second)
	err := s.keeper.AllocateRewards(s.Ctx)
	s.Require().NoError(err)

	rewards, err := s.keeper.ServiceDelegationRewards(s.Ctx, delAddr, service.ID)
	s.Require().NoError(err)
	s.Require().Empty(rewards)
}

func (s *KeeperTestSuite) TestAllocateRewards_BasicScenario() {
	// - x/pools whitelists Service1, Service2
	// - Service3 whitelists Operator2, Operator3
	// - Service2 whitelists $MILK, $MUSD pool
	// - Operator1 joins Service1, Service2, Service3
	//   - but Service3 doesn't have Operator1 in its whitelist
	// - Operator2 joins Service1, Service3
	// - Operator3 joins Service2, Service3
	// - Operator1 sets 10% commission rate
	// - Operator2 sets 5% commission rate
	// - Operator3 sets 2% commission rate

	s.RegisterCurrency("umilk", "MILK", utils.MustParseDec("2"))
	s.RegisterCurrency("uinit", "INIT", utils.MustParseDec("3"))
	s.RegisterCurrency("uusd", "MUSD", utils.MustParseDec("1"))

	// Create services.
	serviceAdmin1 := utils.TestAddress(10000)
	service1 := s.CreateService("Service1", serviceAdmin1.String())
	serviceAdmin2 := utils.TestAddress(10001)
	service2 := s.CreateService("Service2", serviceAdmin2.String())
	serviceAdmin3 := utils.TestAddress(10003)
	service3 := s.CreateService("Service3", serviceAdmin3.String())

	// Add only Service1 and Service2 to the pools module's allowed list.
	poolsParams := s.App.PoolsKeeper.GetParams(s.Ctx)
	poolsParams.AllowedServiceIDs = []uint32{service1.ID, service2.ID}
	s.App.PoolsKeeper.SetParams(s.Ctx, poolsParams)

	// Create operators.
	operatorAdmin1 := utils.TestAddress(10004)
	operator1 := s.CreateOperator("Operator1", operatorAdmin1.String())
	operatorAdmin2 := utils.TestAddress(10005)
	operator2 := s.CreateOperator("Operator2", operatorAdmin2.String())
	operatorAdmin3 := utils.TestAddress(10006)
	operator3 := s.CreateOperator("Operator3", operatorAdmin3.String())

	s.UpdateOperatorParams(operator1.ID, utils.MustParseDec("0.1"), []uint32{service1.ID, service2.ID, service3.ID})
	s.UpdateOperatorParams(operator2.ID, utils.MustParseDec("0.05"), []uint32{service1.ID, service3.ID})
	s.UpdateOperatorParams(operator3.ID, utils.MustParseDec("0.02"), []uint32{service2.ID, service3.ID})

	// Create active rewards plans.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	plan1 := s.CreateBasicRewardsPlan(
		service1.ID, utils.MustParseCoins("1000_000000service1"),
		planStartTime, planEndTime, serviceAdmin1.String())
	plan2 := s.CreateBasicRewardsPlan(
		service2.ID, utils.MustParseCoins("5000_000000service2"),
		planStartTime, planEndTime, serviceAdmin1.String())
	plan3 := s.CreateBasicRewardsPlan(
		service3.ID, utils.MustParseCoins("10000_000000service3"),
		planStartTime, planEndTime, serviceAdmin1.String())
	s.FundAccount(plan1.RewardsPool, utils.MustParseCoins("1000000_000000service1"))
	s.FundAccount(plan2.RewardsPool, utils.MustParseCoins("1000000_000000service2"))
	s.FundAccount(plan3.RewardsPool, utils.MustParseCoins("1000000_000000service3"))

	// Call AllocateRewards to set last rewards allocation time.
	err := s.keeper.AllocateRewards(s.Ctx)
	s.Require().NoError(err)

	aliceAddr := utils.TestAddress(1)
	s.DelegatePool(utils.MustParseCoin("100_000000umilk"), aliceAddr.String(), true) // $300
	s.DelegatePool(utils.MustParseCoin("100_000000uinit"), aliceAddr.String(), true) // $200
	s.DelegatePool(utils.MustParseCoin("500_000000uusd"), aliceAddr.String(), true)  // $500

	// Whitelist only $MILK and $MUSD pools.
	s.UpdateServiceParams(service2.ID, math.LegacyZeroDec(), []uint32{1, 3}, nil)
	// Whitelist only Operator2 and Operator3.
	s.UpdateServiceParams(service3.ID, math.LegacyZeroDec(), nil, []uint32{operator2.ID, operator3.ID})

	bobAddr := utils.TestAddress(2)
	s.DelegateService(service1.ID, utils.MustParseCoins("100_000000uinit"), bobAddr.String(), true) // $300
	s.DelegateService(
		service2.ID, utils.MustParseCoins("200_000000uinit"), bobAddr.String(), true) // $600
	s.DelegateService(service3.ID, utils.MustParseCoins("300_000000uinit"), bobAddr.String(), true) // $900

	charlieAddr := utils.TestAddress(3)
	s.DelegateOperator(operator1.ID, utils.MustParseCoins("1000_000000uusd"), charlieAddr.String(), true) // $1000
	s.DelegateOperator(operator2.ID, utils.MustParseCoins("1000_000000uusd"), charlieAddr.String(), true) // $1000
	s.DelegateOperator(operator3.ID, utils.MustParseCoins("500_000000uusd"), charlieAddr.String(), true)  // $500

	davidAddr := utils.TestAddress(4)
	s.DelegatePool(utils.MustParseCoin("200_000000umilk"), davidAddr.String(), true)                    // $400
	s.DelegatePool(utils.MustParseCoin("200_000000uinit"), davidAddr.String(), true)                    // $600
	s.DelegatePool(utils.MustParseCoin("200_000000uusd"), davidAddr.String(), true)                     // $200
	s.DelegateService(service1.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true)   // $400
	s.DelegateService(service2.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true)   // $400
	s.DelegateService(service3.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true)   // $400
	s.DelegateOperator(operator1.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true) // $400
	s.DelegateOperator(operator2.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true) // $400
	s.DelegateOperator(operator3.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true) // $400

	s.advanceBlock(10 * time.Second)
	// Rewards plan 1 allocates 1000 * 10 / 86400 ~= 0.115741 $SERVICE1
	// Rewards plan 2 allocates 5000 * 10 / 86400 ~= 0.578704 $SERVICE1
	// Rewards plan 3 allocates 10000 * 10 / 86400 ~= 1.157407 $SERVICE1
	err = s.keeper.AllocateRewards(s.Ctx)
	s.Require().NoError(err)

	// Alice receives:
	// - $200 / $5700 * 0.115741 ~= 0.004061 $SERVICE1 (from Pool1)
	// - $200 / $4600 * 0.578704 ~= 0.025161 $SERVICE2 (from Pool1)
	// - $300 / $5700 * 0.115741 ~= 0.006092 $SERVICE1 (from Pool2)
	// - $500 / $5700 * 0.115741 ~= 0.010153 $SERVICE1 (from Pool3)
	// - $500 / $4600 * 0.578704 ~= 0.062903 $SERVICE2 (from Pool3)
	rewards, err := s.keeper.PoolDelegationRewards(s.Ctx, aliceAddr, 1)
	s.Require().NoError(err)
	s.Require().Equal("4061.078622482100000000service1,25161.030595813200000000service2", rewards.String())
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, aliceAddr, 2)
	s.Require().NoError(err)
	s.Require().Equal("6091.617933723100000000service1", rewards.String())
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, aliceAddr, 3)
	s.Require().NoError(err)
	s.Require().Equal("10152.696556205000000000service1,62902.576489533000000000service2", rewards.String())

	// Bob receives:
	// - $300 / $5700 * 0.115741 ~= 0.006092 $SERVICE1 (from Service1)
	// - $600 / $4600 * 0.578704 ~= 0.075483 $SERVICE2 (from Service2)
	// - $900 / $3600 * 1.157407 ~= 0.289352 $SERVICE3 (from Service3)
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, bobAddr, service1.ID)
	s.Require().NoError(err)
	s.Require().Equal("6091.617933723100000000service1", rewards.String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, bobAddr, service2.ID)
	s.Require().NoError(err)
	s.Require().Equal("75483.091787439600000000service2", rewards.String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, bobAddr, service3.ID)
	s.Require().NoError(err)
	s.Require().Equal("289351.851851851800000000service3", rewards.String())

	// Charlie receives:
	// - $1000 / $5700 * 0.115741 * 0.9 ~= 0.018275 $SERVICE1 (from Operator1)
	// - $1000 / $4600 * 0.578704 * 0.9 ~= 0.113225 $SERVICE2 (from Operator1)
	// - $1000 / $5700 * 0.115741 * 0.95 ~= 0.019290 $SERVICE1 (from Operator2)
	// - $1000 / $3600 * 1.157407 * 0.95 ~= 0.305427 $SERVICE3 (from Operator2)
	// - $500 / $4600 * 0.578704 * 0.98 ~= 0.061645 $SERVICE2 (from Operator3)
	// - $500 / $3600 * 1.157407 * 0.98 ~= 0.157536 $SERVICE3 (from Operator3)
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, charlieAddr, operator1.ID)
	s.Require().NoError(err)
	s.Require().Equal("18274.853801169000000000service1,113224.637681159000000000service2", rewards.String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, charlieAddr, operator2.ID)
	s.Require().NoError(err)
	s.Require().Equal("19290.123456790000000000service1,305426.954732510000000000service3", rewards.String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, charlieAddr, operator3.ID)
	s.Require().NoError(err)
	s.Require().Equal("61644.524959742000000000service2,157536.008230452500000000service3", rewards.String())

	// David receives:
	// - $400 / $5700 * 0.115741 ~= 0.008122 $SERVICE1 (from Pool1)
	// - $400 / $4600 * 0.578704 ~= 0.050322 $SERVICE2 (from Pool1)
	// - $600 / $5700 * 0.115741 ~= 0.012183 $SERVICE1 (from Pool2)
	// - $200 / $5700 * 0.115741 ~= 0.004061 $SERVICE1 (from Pool3)
	// - $200 / $4600 * 0.578704 ~= 0.025161 $SERVICE2 (from Pool3)
	// - $400 / $5700 * 0.115741 ~= 0.008122 $SERVICE1 (from Service1)
	// - $400 / $4600 * 0.578704 ~= 0.050322 $SERVICE2 (from Service2)
	// - $400 / $3600 * 1.157407 ~= 0.128601 $SERVICE3 (from Service3)
	// - $400 / $5700 * 0.115741 * 0.9 ~= 0.007310 $SERVICE1 (from Operator1)
	// - $400 / $4600 * 0.578704 * 0.9 ~= 0.045290 $SERVICE2 (from Operator1)
	// - $400 / $5700 * 0.115741 * 0.95 ~= 0.007716 $SERVICE1 (from Operator2)
	// - $400 / $3600 * 1.157407 * 0.95 ~= 0.122171 $SERVICE3 (from Operator2)
	// - $400 / $4600 * 0.578704 * 0.98 ~= 0.049316 $SERVICE2 (from Operator3)
	// - $400 / $3600 * 1.157407 * 0.98 ~= 0.126029 $SERVICE3 (from Operator3)
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, davidAddr, 1)
	s.Require().NoError(err)
	s.Require().Equal("8122.157244964200000000service1,50322.061191626400000000service2", rewards.String())
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, davidAddr, 2)
	s.Require().NoError(err)
	s.Require().Equal("12183.235867446200000000service1", rewards.String())
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, davidAddr, 3)
	s.Require().NoError(err)
	s.Require().Equal("4061.078622482000000000service1,25161.030595813200000000service2", rewards.String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, davidAddr, service1.ID)
	s.Require().NoError(err)
	s.Require().Equal("8122.157244964200000000service1", rewards.String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, davidAddr, service2.ID)
	s.Require().NoError(err)
	s.Require().Equal("50322.061191626400000000service2", rewards.String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, davidAddr, service3.ID)
	s.Require().NoError(err)
	s.Require().Equal("128600.823045267400000000service3", rewards.String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, davidAddr, operator1.ID)
	s.Require().NoError(err)
	s.Require().Equal("7309.941520467800000000service1,45289.855072463600000000service2", rewards.String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, davidAddr, operator2.ID)
	s.Require().NoError(err)
	s.Require().Equal("7716.049382716000000000service1,122170.781893004000000000service3", rewards.String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, davidAddr, operator3.ID)
	s.Require().NoError(err)
	s.Require().Equal("49315.619967793800000000service2,126028.806584362000000000service3", rewards.String())
}

func (s *KeeperTestSuite) TestAllocateRewards_MovingPrice() {
	panic("not implemented")
}
