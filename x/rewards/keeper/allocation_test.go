package keeper_test

import (
	"time"

	"cosmossdk.io/math"

	"github.com/milkyway-labs/milkyway/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func (s *KeeperTestSuite) TestAllocateRewards_InactivePlan() {
	// Inactive plans(current block time is out of their date range) don't
	// allocate rewards.

	// Plan's start time is 2024-01-01 so set block time before that.
	s.Ctx = s.Ctx.WithBlockTime(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC))
	service, _ := s.setupSampleServiceAndOperator()

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"), planStartTime, planEndTime,
		utils.MustParseCoins("100000_000000service"))

	delAddr := utils.TestAddress(1)
	s.DelegateService(service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	s.allocateRewards(3 * time.Second)

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
	poolsParams.AllowedServicesIDs = []uint32{service1.ID, service2.ID}
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
	s.CreateBasicRewardsPlan(
		service1.ID, utils.MustParseCoins("1000_000000service1"),
		planStartTime, planEndTime,
		utils.MustParseCoins("100000_000000service1"))
	s.CreateBasicRewardsPlan(
		service2.ID, utils.MustParseCoins("5000_000000service2"),
		planStartTime, planEndTime,
		utils.MustParseCoins("100000_000000service2"))
	s.CreateBasicRewardsPlan(
		service3.ID, utils.MustParseCoins("10000_000000service3"),
		planStartTime, planEndTime,
		utils.MustParseCoins("100000_000000service3"))

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

	// Rewards plan 1 allocates 1000 * 10 / 86400 ~= 0.115741 $SERVICE1
	// Rewards plan 2 allocates 5000 * 10 / 86400 ~= 0.578704 $SERVICE1
	// Rewards plan 3 allocates 10000 * 10 / 86400 ~= 1.157407 $SERVICE1
	s.allocateRewards(10 * time.Second)

	// Alice receives:
	// - $200 / $5700 * 0.115741 ~= 0.004061 $SERVICE1 (from Pool1)
	// - $200 / $4600 * 0.578704 ~= 0.025161 $SERVICE2 (from Pool1)
	// - $300 / $5700 * 0.115741 ~= 0.006092 $SERVICE1 (from Pool2)
	// - $500 / $5700 * 0.115741 ~= 0.010153 $SERVICE1 (from Pool3)
	// - $500 / $4600 * 0.578704 ~= 0.062903 $SERVICE2 (from Pool3)
	rewards, err := s.keeper.PoolDelegationRewards(s.Ctx, aliceAddr, 1)
	s.Require().NoError(err)
	s.Require().Equal("4061.052631578900000000service1,25161.000000000000000000service2", rewards.Sum().String())
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, aliceAddr, 2)
	s.Require().NoError(err)
	s.Require().Equal("6091.578947368400000000service1", rewards.Sum().String())
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, aliceAddr, 3)
	s.Require().NoError(err)
	s.Require().Equal("10152.631578947000000000service1,62902.500000000000000000service2", rewards.Sum().String())

	// Bob receives:
	// - $300 / $5700 * 0.115741 ~= 0.006092 $SERVICE1 (from Service1)
	// - $600 / $4600 * 0.578704 ~= 0.075483 $SERVICE2 (from Service2)
	// - $900 / $3600 * 1.157407 ~= 0.289352 $SERVICE3 (from Service3)
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, bobAddr, service1.ID)
	s.Require().NoError(err)
	s.Assert().Equal("6091.578947368400000000service1", rewards.Sum().String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, bobAddr, service2.ID)
	s.Require().NoError(err)
	s.Assert().Equal("75483.000000000000000000service2", rewards.Sum().String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, bobAddr, service3.ID)
	s.Require().NoError(err)
	s.Assert().Equal("289351.749999999900000000service3", rewards.Sum().String())

	// Charlie receives:
	// - $1000 / $5700 * 0.115741 * 0.9 ~= 0.018275 $SERVICE1 (from Operator1)
	// - $1000 / $4600 * 0.578704 * 0.9 ~= 0.113225 $SERVICE2 (from Operator1)
	// - $1000 / $5700 * 0.115741 * 0.95 ~= 0.019290 $SERVICE1 (from Operator2)
	// - $1000 / $3600 * 1.157407 * 0.95 ~= 0.305427 $SERVICE3 (from Operator2)
	// - $500 / $4600 * 0.578704 * 0.98 ~= 0.061645 $SERVICE2 (from Operator3)
	// - $500 / $3600 * 1.157407 * 0.98 ~= 0.157536 $SERVICE3 (from Operator3)
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, charlieAddr, operator1.ID)
	s.Require().NoError(err)
	s.Assert().Equal("18274.736842105000000000service1,113224.500000000000000000service2", rewards.Sum().String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, charlieAddr, operator2.ID)
	s.Require().NoError(err)
	s.Assert().Equal("19290.000000000000000000service1,305426.847222222000000000service3", rewards.Sum().String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, charlieAddr, operator3.ID)
	s.Require().NoError(err)
	s.Assert().Equal("61644.450000000000000000service2,157535.952777777500000000service3", rewards.Sum().String())

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
	s.Assert().Equal("8122.105263157800000000service1,50322.000000000000000000service2", rewards.Sum().String())
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, davidAddr, 2)
	s.Require().NoError(err)
	s.Assert().Equal("12183.157894736800000000service1", rewards.Sum().String())
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, davidAddr, 3)
	s.Require().NoError(err)
	s.Assert().Equal("4061.052631578800000000service1,25161.000000000000000000service2", rewards.Sum().String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, davidAddr, service1.ID)
	s.Require().NoError(err)
	s.Assert().Equal("8122.105263157800000000service1", rewards.Sum().String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, davidAddr, service2.ID)
	s.Require().NoError(err)
	s.Assert().Equal("50322.000000000000000000service2", rewards.Sum().String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, davidAddr, service3.ID)
	s.Require().NoError(err)
	s.Assert().Equal("128600.777777777600000000service3", rewards.Sum().String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, davidAddr, operator1.ID)
	s.Require().NoError(err)
	s.Assert().Equal("7309.894736842000000000service1,45289.800000000000000000service2", rewards.Sum().String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, davidAddr, operator2.ID)
	s.Require().NoError(err)
	s.Assert().Equal("7716.000000000000000000service1,122170.738888888800000000service3", rewards.Sum().String())
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, davidAddr, operator3.ID)
	s.Require().NoError(err)
	s.Assert().Equal("49315.560000000000000000service2,126028.762222222200000000service3", rewards.Sum().String())
}

func (s *KeeperTestSuite) TestAllocateRewards_MovingPrice() {
	// $MILK is $2 and $INIT is $3.
	service, _ := s.setupSampleServiceAndOperator()

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	s.CreateBasicRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"), planStartTime, planEndTime,
		utils.MustParseCoins("100000_000000service"))

	delAddr1 := utils.TestAddress(1)
	s.DelegateService(service.ID, utils.MustParseCoins("100_000000umilk"), delAddr1.String(), true)
	delAddr2 := utils.TestAddress(2)
	s.DelegateService(service.ID, utils.MustParseCoins("100_000000uinit"), delAddr2.String(), true)

	// Allocate rewards.
	s.allocateRewards(10 * time.Second)

	// They receive rewards by 1:2 ratio.
	rewards, err := s.keeper.ServiceDelegationRewards(s.Ctx, delAddr1, service.ID)
	s.Require().NoError(err)
	s.Require().Equal("4629.600000000000000000service", rewards.Sum().String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, delAddr2, service.ID)
	s.Require().NoError(err)
	s.Require().Equal("6944.400000000000000000service", rewards.Sum().String())

	// Now price changes.
	s.RegisterCurrency("umilk", "MILK", utils.MustParseDec("3"))
	s.RegisterCurrency("uinit", "INIT", utils.MustParseDec("1"))

	// Allocate rewards again.
	s.allocateRewards(10 * time.Second)

	// Now they receive rewards by 3:1 ratio.
	// Note that already accumulated rewards are not affected.
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, delAddr1, service.ID)
	s.Require().NoError(err)
	// Delta: +8680.5555555555umilk
	s.Require().Equal("13310.100000000000000000service", rewards.Sum().String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, delAddr2, service.ID)
	s.Require().NoError(err)
	// Delta: +2893.5185185185umilk
	s.Require().Equal("9837.900000000000000000service", rewards.Sum().String())
}

func (s *KeeperTestSuite) TestAllocateRewards_ZeroDelegations() {
	// Test if AllocateRewards handles pool/operator/service distribution
	// correctly when the distribution info has weight specified but there's
	// no delegation yet.

	s.RegisterCurrency("umilk", "MILK", utils.MustParseDec("2"))

	// Create a service.
	serviceAdmin := utils.TestAddress(10000)
	service := s.CreateService("Service", serviceAdmin.String())

	// Add the created service ID to the pools module's allowed list.
	poolsParams := s.App.PoolsKeeper.GetParams(s.Ctx)
	poolsParams.AllowedServicesIDs = []uint32{service.ID}
	s.App.PoolsKeeper.SetParams(s.Ctx, poolsParams)

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	s.CreateRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"), planStartTime, planEndTime,
		types.NewBasicPoolsDistribution(1), types.NewBasicOperatorsDistribution(2), types.NewBasicUsersDistribution(3),
		utils.MustParseCoins("100000_000000service"))

	// Create an operator.
	operatorAdmin := utils.TestAddress(10001)
	operator := s.CreateOperator("Operator", operatorAdmin.String())
	// Make the operator join the service and set its commission rate to 10%.
	s.UpdateOperatorParams(operator.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Call AllocateRewards to set last rewards allocation time.
	err := s.keeper.AllocateRewards(s.Ctx)
	s.Require().NoError(err)

	// Try allocating rewards.
	s.allocateRewards(10 * time.Second)

	// There must be no outstanding rewards allocated.
	target, err := s.keeper.GetDelegationTarget(s.Ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operator.ID)
	s.Require().NoError(err)
	rewards, err := s.keeper.GetOutstandingRewardsCoins(s.Ctx, target)
	s.Require().NoError(err)
	s.Require().Empty(rewards)
	target, err = s.keeper.GetDelegationTarget(s.Ctx, restakingtypes.DELEGATION_TYPE_SERVICE, service.ID)
	s.Require().NoError(err)
	rewards, err = s.keeper.GetOutstandingRewardsCoins(s.Ctx, target)
	s.Require().NoError(err)
	s.Require().Empty(rewards)

	// Two users delegate the same amount of $MILK to a pool and the service.
	delAddr1 := utils.TestAddress(1)
	s.DelegatePool(utils.MustParseCoin("10_000000umilk"), delAddr1.String(), true)
	delAddr2 := utils.TestAddress(2)
	s.DelegateService(service.ID, utils.MustParseCoins("10_000000umilk"), delAddr2.String(), true)

	// Allocate rewards.
	s.allocateRewards(10 * time.Second)

	// Still the operator has no rewards.
	target, err = s.keeper.GetDelegationTarget(s.Ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operator.ID)
	s.Require().NoError(err)
	rewards, err = s.keeper.GetOutstandingRewardsCoins(s.Ctx, target)
	s.Require().NoError(err)
	s.Require().Empty(rewards)

	// The pool and the service receive rewards by 1:3 ratio.
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, delAddr1, 1)
	s.Require().NoError(err)
	s.Require().Equal("2893.500000000000000000service", rewards.Sum().String())
	rewards, err = s.keeper.ServiceDelegationRewards(s.Ctx, delAddr2, service.ID)
	s.Require().NoError(err)
	s.Require().Equal("8680.500000000000000000service", rewards.Sum().String())
}

func (s *KeeperTestSuite) TestAllocateRewards_WeightedDistributions() {
	// Register $MILK and $INIT. For simple calculation, set both currencies'
	// price $1.
	s.RegisterCurrency("umilk", "MILK", utils.MustParseDec("1"))
	s.RegisterCurrency("uinit", "INIT", utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := utils.TestAddress(10000)
	service := s.CreateService("Service", serviceAdmin.String())

	// Add the created service ID to the pools module's allowed list.
	poolsParams := s.App.PoolsKeeper.GetParams(s.Ctx)
	poolsParams.AllowedServicesIDs = []uint32{service.ID}
	s.App.PoolsKeeper.SetParams(s.Ctx, poolsParams)

	// Create operators.
	operatorAdmin1 := utils.TestAddress(10001)
	operator1 := s.CreateOperator("Operator1", operatorAdmin1.String())
	s.UpdateOperatorParams(operator1.ID, utils.MustParseDec("0.1"), []uint32{service.ID})
	operatorAdmin2 := utils.TestAddress(10002)
	operator2 := s.CreateOperator("Operator2", operatorAdmin2.String())
	s.UpdateOperatorParams(operator2.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Call AllocateRewards to set last rewards allocation time.
	err := s.keeper.AllocateRewards(s.Ctx)
	s.Require().NoError(err)

	// Delegate to $MILK pool.
	delAddr1 := utils.TestAddress(1)
	s.DelegatePool(utils.MustParseCoin("300_000000umilk"), delAddr1.String(), true)
	delAddr2 := utils.TestAddress(2)
	s.DelegatePool(utils.MustParseCoin("200_000000umilk"), delAddr2.String(), true)
	// Delegate to $INIT pool.
	delAddr3 := utils.TestAddress(3)
	s.DelegatePool(utils.MustParseCoin("200_000000uinit"), delAddr3.String(), true)
	delAddr4 := utils.TestAddress(4)
	s.DelegatePool(utils.MustParseCoin("300_000000uinit"), delAddr4.String(), true)
	// Delegate to Operator1.
	delAddr5 := utils.TestAddress(5)
	s.DelegateOperator(operator1.ID, utils.MustParseCoins("100_000000umilk"), delAddr5.String(), true)
	delAddr6 := utils.TestAddress(6)
	s.DelegateOperator(operator1.ID, utils.MustParseCoins("200_000000uinit"), delAddr6.String(), true)
	// Delegate to Operator2.
	delAddr7 := utils.TestAddress(7)
	s.DelegateOperator(operator2.ID, utils.MustParseCoins("200_000000umilk"), delAddr7.String(), true)
	delAddr8 := utils.TestAddress(8)
	s.DelegateOperator(operator2.ID, utils.MustParseCoins("100_000000uinit"), delAddr8.String(), true)

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	s.CreateRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"), planStartTime, planEndTime,
		types.NewWeightedPoolsDistribution(3, []types.DistributionWeight{
			types.NewDistributionWeight(1, 2),
			types.NewDistributionWeight(2, 1),
		}),
		types.NewWeightedOperatorsDistribution(1, []types.DistributionWeight{
			types.NewDistributionWeight(operator1.ID, 2),
			types.NewDistributionWeight(operator2.ID, 3),
		}),
		types.NewBasicUsersDistribution(0), // No user rewards
		utils.MustParseCoins("100000_000000service"),
	)

	s.allocateRewards(10 * time.Second)

	// delAddr1 receives 3/4 * 2/3 * $300 / $500 * 100 * (10 / 86400) ~= 0.003472 $SERVICE
	rewards, err := s.keeper.PoolDelegationRewards(s.Ctx, delAddr1, 1)
	s.Require().NoError(err)
	s.Require().Equal("3472.200000000000000000service", rewards.Sum().String())

	// delAddr2 receives 3/4 * 2/3 * $200 / $500 * 100 * (10 / 86400) ~= 0.002315 $SERVICE
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, delAddr2, 1)
	s.Require().NoError(err)
	s.Require().Equal("2314.800000000000000000service", rewards.Sum().String())

	// delAddr3 receives 3/4 * 1/3 * $200 / $500 * 100 * (10 / 86400) ~= 0.001157 $SERVICE
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, delAddr3, 2)
	s.Require().NoError(err)
	s.Require().Equal("1157.400000000000000000service", rewards.Sum().String())

	// delAddr4 receives 3/4 * 1/3 * $300 / $500 * 100 * (10 / 86400) ~= 0.001736 $SERVICE
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, delAddr4, 2)
	s.Require().NoError(err)
	s.Require().Equal("1736.100000000000000000service", rewards.Sum().String())

	// Note that operators take commission from rewards.

	// delAddr5 receives 1/4 * 2/5 * $100 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000347 $SERVICE
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, delAddr5, operator1.ID)
	s.Require().NoError(err)
	s.Assert().Equal("347.220000000000000000service", rewards.Sum().String())

	// delAddr6 receives 1/4 * 2/5 * $200 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000694 $SERVICE
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, delAddr6, operator1.ID)
	s.Require().NoError(err)
	s.Assert().Equal("694.440000000000000000service", rewards.Sum().String())

	// delAddr7 receives 1/4 * 3/5 * $200 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.001042 $SERVICE
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, delAddr7, operator2.ID)
	s.Require().NoError(err)
	s.Assert().Equal("1041.660000000000000000service", rewards.Sum().String())

	// delAddr8 receives 1/4 * 3/5 * $100 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000521 $SERVICE
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, delAddr8, operator2.ID)
	s.Require().NoError(err)
	s.Assert().Equal("520.830000000000000000service", rewards.Sum().String())
}

func (s *KeeperTestSuite) TestAllocateRewards_EgalitarianDistributions() {
	// Register $MILK and $INIT. For simple calculation, set both currencies'
	// price $1.
	s.RegisterCurrency("umilk", "MILK", utils.MustParseDec("1"))
	s.RegisterCurrency("uinit", "INIT", utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := utils.TestAddress(10000)
	service := s.CreateService("Service", serviceAdmin.String())

	// Add the created service ID to the pools module's allowed list.
	poolsParams := s.App.PoolsKeeper.GetParams(s.Ctx)
	poolsParams.AllowedServicesIDs = []uint32{service.ID}
	s.App.PoolsKeeper.SetParams(s.Ctx, poolsParams)

	// Create operators.
	operatorAdmin1 := utils.TestAddress(10001)
	operator1 := s.CreateOperator("Operator1", operatorAdmin1.String())
	s.UpdateOperatorParams(operator1.ID, utils.MustParseDec("0.1"), []uint32{service.ID})
	operatorAdmin2 := utils.TestAddress(10002)
	operator2 := s.CreateOperator("Operator2", operatorAdmin2.String())
	s.UpdateOperatorParams(operator2.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Call AllocateRewards to set last rewards allocation time.
	err := s.keeper.AllocateRewards(s.Ctx)
	s.Require().NoError(err)

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	s.CreateRewardsPlan(
		service.ID, utils.MustParseCoins("100_000000service"), planStartTime, planEndTime,
		types.NewEgalitarianPoolsDistribution(3),
		types.NewEgalitarianOperatorsDistribution(1),
		types.NewBasicUsersDistribution(0), // No user rewards
		utils.MustParseCoins("100000_000000service"),
	)

	// Delegate to $MILK pool.
	delAddr1 := utils.TestAddress(1)
	s.DelegatePool(utils.MustParseCoin("300_000000umilk"), delAddr1.String(), true)
	delAddr2 := utils.TestAddress(2)
	s.DelegatePool(utils.MustParseCoin("200_000000umilk"), delAddr2.String(), true)
	// Delegate to $INIT pool.
	delAddr3 := utils.TestAddress(3)
	s.DelegatePool(utils.MustParseCoin("200_000000uinit"), delAddr3.String(), true)
	delAddr4 := utils.TestAddress(4)
	s.DelegatePool(utils.MustParseCoin("300_000000uinit"), delAddr4.String(), true)
	// Delegate to Operator1.
	delAddr5 := utils.TestAddress(5)
	s.DelegateOperator(operator1.ID, utils.MustParseCoins("100_000000umilk"), delAddr5.String(), true)
	delAddr6 := utils.TestAddress(6)
	s.DelegateOperator(operator1.ID, utils.MustParseCoins("200_000000uinit"), delAddr6.String(), true)
	// Delegate to Operator2.
	delAddr7 := utils.TestAddress(7)
	s.DelegateOperator(operator2.ID, utils.MustParseCoins("200_000000umilk"), delAddr7.String(), true)
	delAddr8 := utils.TestAddress(8)
	s.DelegateOperator(operator2.ID, utils.MustParseCoins("100_000000uinit"), delAddr8.String(), true)

	s.allocateRewards(10 * time.Second)

	// delAddr1 receives 3/4 * 1/2 * $300 / $500 * 100 * (10 / 86400) ~= 0.002604 $SERVICE
	rewards, err := s.keeper.PoolDelegationRewards(s.Ctx, delAddr1, 1)
	s.Require().NoError(err)
	s.Require().Equal("2604.150000000000000000service", rewards.Sum().String())

	// delAddr2 receives 3/4 * 1/2 * $200 / $500 * 100 * (10 / 86400) ~= 0.001736 $SERVICE
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, delAddr2, 1)
	s.Require().NoError(err)
	s.Require().Equal("1736.100000000000000000service", rewards.Sum().String())

	// delAddr3 receives 3/4 * 1/2 * $200 / $500 * 100 * (10 / 86400) ~= 0.001736 $SERVICE
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, delAddr3, 2)
	s.Require().NoError(err)
	s.Require().Equal("1736.100000000000000000service", rewards.Sum().String())

	// delAddr4 receives 3/4 * 1/2 * $300 / $500 * 100 * (10 / 86400) ~= 0.002604 $SERVICE
	rewards, err = s.keeper.PoolDelegationRewards(s.Ctx, delAddr4, 2)
	s.Require().NoError(err)
	s.Require().Equal("2604.150000000000000000service", rewards.Sum().String())

	// Note that operators take commission from rewards.

	// delAddr5 receives 1/4 * 1/2 * $100 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000434 $SERVICE
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, delAddr5, operator1.ID)
	s.Require().NoError(err)
	s.Require().Equal("434.025000000000000000service", rewards.Sum().String())

	// delAddr6 receives 1/4 * 1/2 * $200 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000868 $SERVICE
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, delAddr6, operator1.ID)
	s.Require().NoError(err)
	s.Require().Equal("868.050000000000000000service", rewards.Sum().String())

	// delAddr7 receives 1/4 * 1/2 * $200 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000868 $SERVICE
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, delAddr7, operator2.ID)
	s.Require().NoError(err)
	s.Require().Equal("868.050000000000000000service", rewards.Sum().String())

	// delAddr8 receives 1/4 * 1/2 * $100 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000434 $SERVICE
	rewards, err = s.keeper.OperatorDelegationRewards(s.Ctx, delAddr8, operator2.ID)
	s.Require().NoError(err)
	s.Require().Equal("434.025000000000000000service", rewards.Sum().String())
}
