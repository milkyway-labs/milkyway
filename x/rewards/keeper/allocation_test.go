package keeper_test

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/v12/app/testutil"
	"github.com/milkyway-labs/milkyway/v12/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/v12/x/operators/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v12/x/restaking/types"
	"github.com/milkyway-labs/milkyway/v12/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/v12/x/rewards/types"
)

func (suite *KeeperTestSuite) TestAllocateRewards_InactivePlan() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	// Inactive plans(current block time is out of their date range) don't allocate rewards.

	// Plan's start time is 2024-01-01 so set block time before that.
	ctx = ctx.WithBlockTime(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC))
	service, _ := suite.setupSampleServiceAndOperator(ctx)

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("100_000000service"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service"),
	)

	delAddr := testutil.TestAddress(1)
	suite.DelegateService(ctx, service.ID, utils.MustParseCoins("100_000000umilk"), delAddr.String(), true)

	// Allocate the rewards
	ctx = suite.allocateRewards(ctx, 3*time.Second)

	rewards, err := suite.keeper.GetServiceDelegationRewards(ctx, delAddr, service.ID)
	suite.Require().NoError(err)
	suite.Require().Empty(rewards)
}

func (suite *KeeperTestSuite) TestAllocateRewards_BasicScenario() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

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

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))
	suite.RegisterCurrency(ctx, "uinit", "INIT", 6, utils.MustParseDec("3"))
	suite.RegisterCurrency(ctx, "uusd", "MUSD", 6, utils.MustParseDec("1"))

	// Create services.
	serviceAdmin1 := testutil.TestAddress(10000)
	service1 := suite.CreateService(ctx, "Service1", serviceAdmin1.String())
	err := suite.servicesKeeper.SetServiceAccredited(ctx, service1.ID, true)
	suite.Require().NoError(err)

	serviceAdmin2 := testutil.TestAddress(10001)
	service2 := suite.CreateService(ctx, "Service2", serviceAdmin2.String())
	err = suite.servicesKeeper.SetServiceAccredited(ctx, service2.ID, true)
	suite.Require().NoError(err)

	serviceAdmin3 := testutil.TestAddress(10003)
	service3 := suite.CreateService(ctx, "Service3", serviceAdmin3.String())
	err = suite.servicesKeeper.SetServiceAccredited(ctx, service3.ID, false)
	suite.Require().NoError(err)

	// Create operators.
	operatorAdmin1 := testutil.TestAddress(10004)
	operator1 := suite.CreateOperator(ctx, "Operator1", operatorAdmin1.String())
	operatorAdmin2 := testutil.TestAddress(10005)
	operator2 := suite.CreateOperator(ctx, "Operator2", operatorAdmin2.String())
	operatorAdmin3 := testutil.TestAddress(10006)
	operator3 := suite.CreateOperator(ctx, "Operator3", operatorAdmin3.String())

	// Whitelist all pools.
	suite.AddPoolsToServiceSecuringPools(ctx, service1.ID, []uint32{1, 2, 3})
	// Whitelist only $MILK and $MUSD pools.
	suite.AddPoolsToServiceSecuringPools(ctx, service2.ID, []uint32{1, 3})
	// Whitelist only Operator2 and Operator3.
	suite.AddOperatorsToServiceAllowList(ctx, service3.ID, []uint32{operator2.ID, operator3.ID})

	suite.UpdateOperatorParams(ctx, operator1.ID, utils.MustParseDec("0.1"), []uint32{service1.ID, service2.ID, service3.ID})
	suite.UpdateOperatorParams(ctx, operator2.ID, utils.MustParseDec("0.05"), []uint32{service1.ID, service3.ID})
	suite.UpdateOperatorParams(ctx, operator3.ID, utils.MustParseDec("0.02"), []uint32{service2.ID, service3.ID})

	// Create active rewards plans.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateBasicRewardsPlan(
		ctx,
		service1.ID,
		utils.MustParseCoin("1000_000000service1"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service1"),
	)
	suite.CreateBasicRewardsPlan(
		ctx,
		service2.ID,
		utils.MustParseCoin("5000_000000service2"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service2"),
	)
	suite.CreateBasicRewardsPlan(
		ctx,
		service3.ID,
		utils.MustParseCoin("10000_000000service3"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service3"),
	)

	// Call AllocateRewards to set last rewards allocation time.
	err = suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	aliceAddr := testutil.TestAddress(1)
	suite.SetUserPreferences(ctx, aliceAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service1.ID, nil),
		restakingtypes.NewTrustedServiceEntry(service2.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("100_000000umilk"), aliceAddr.String(), true) // $200
	suite.DelegatePool(ctx, utils.MustParseCoin("100_000000uinit"), aliceAddr.String(), true) // $300
	suite.DelegatePool(ctx, utils.MustParseCoin("500_000000uusd"), aliceAddr.String(), true)  // $500

	bobAddr := testutil.TestAddress(2)
	suite.SetUserPreferences(ctx, bobAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service1.ID, nil),
		restakingtypes.NewTrustedServiceEntry(service2.ID, nil),
	})
	suite.DelegateService(ctx, service1.ID, utils.MustParseCoins("100_000000uinit"), bobAddr.String(), true) // $300
	suite.DelegateService(ctx, service2.ID, utils.MustParseCoins("200_000000uinit"), bobAddr.String(), true) // $600
	suite.DelegateService(ctx, service3.ID, utils.MustParseCoins("300_000000uinit"), bobAddr.String(), true) // $900

	charlieAddr := testutil.TestAddress(3)
	suite.SetUserPreferences(ctx, charlieAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service1.ID, nil),
		restakingtypes.NewTrustedServiceEntry(service2.ID, nil),
	})
	suite.DelegateOperator(ctx, operator1.ID, utils.MustParseCoins("1000_000000uusd"), charlieAddr.String(), true) // $1000
	suite.DelegateOperator(ctx, operator2.ID, utils.MustParseCoins("1000_000000uusd"), charlieAddr.String(), true) // $1000
	suite.DelegateOperator(ctx, operator3.ID, utils.MustParseCoins("500_000000uusd"), charlieAddr.String(), true)  // $500

	davidAddr := testutil.TestAddress(4)
	suite.SetUserPreferences(ctx, davidAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service1.ID, nil),
		restakingtypes.NewTrustedServiceEntry(service2.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000umilk"), davidAddr.String(), true)                    // $400
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000uinit"), davidAddr.String(), true)                    // $600
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000uusd"), davidAddr.String(), true)                     // $200
	suite.DelegateService(ctx, service1.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true)   // $400
	suite.DelegateService(ctx, service2.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true)   // $400
	suite.DelegateService(ctx, service3.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true)   // $400
	suite.DelegateOperator(ctx, operator1.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true) // $400
	suite.DelegateOperator(ctx, operator2.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true) // $400
	suite.DelegateOperator(ctx, operator3.ID, utils.MustParseCoins("200_000000umilk"), davidAddr.String(), true) // $400

	// Rewards plan 1 allocates 1000 * 10 / 86400 ~= 0.115741 $SERVICE1
	// Rewards plan 2 allocates 5000 * 10 / 86400 ~= 0.578704 $SERVICE1
	// Rewards plan 3 allocates 10000 * 10 / 86400 ~= 1.157407 $SERVICE1
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Alice receives:
	// - $200 / $5700 * 0.115741 ~= 0.004061 $SERVICE1 (from Pool1)
	// - $200 / $4600 * 0.578704 ~= 0.025161 $SERVICE2 (from Pool1)
	// - $300 / $5700 * 0.115741 ~= 0.006092 $SERVICE1 (from Pool2)
	// - $500 / $5700 * 0.115741 ~= 0.010153 $SERVICE1 (from Pool3)
	// - $500 / $4600 * 0.578704 ~= 0.062903 $SERVICE2 (from Pool3)
	rewards, err := suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("4061.052631578900000000service1,25161.000000000000000000service2", rewards.Sum().String())

	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 2)
	suite.Require().NoError(err)
	suite.Require().Equal("6091.578947368400000000service1", rewards.Sum().String())

	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 3)
	suite.Require().NoError(err)
	suite.Require().Equal("10152.631578947000000000service1,62902.500000000000000000service2", rewards.Sum().String())

	// Bob receives:
	// - $300 / $5700 * 0.115741 ~= 0.006092 $SERVICE1 (from Service1)
	// - $600 / $4600 * 0.578704 ~= 0.075483 $SERVICE2 (from Service2)
	// - $900 / $3600 * 1.157407 ~= 0.289352 $SERVICE3 (from Service3)
	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, bobAddr, service1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("6091.578947368400000000service1", rewards.Sum().String())
	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, bobAddr, service2.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("75483.000000000000000000service2", rewards.Sum().String())
	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, bobAddr, service3.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("289351.749999999900000000service3", rewards.Sum().String())

	// Charlie receives:
	// - $1000 / $5700 * 0.115741 * 0.9 ~= 0.018275 $SERVICE1 (from Operator1)
	// - $1000 / $4600 * 0.578704 * 0.9 ~= 0.113225 $SERVICE2 (from Operator1)
	// - $1000 / $5700 * 0.115741 * 0.95 ~= 0.019290 $SERVICE1 (from Operator2)
	// - $1000 / $3600 * 1.157407 * 0.95 ~= 0.305427 $SERVICE3 (from Operator2)
	// - $500 / $4600 * 0.578704 * 0.98 ~= 0.061645 $SERVICE2 (from Operator3)
	// - $500 / $3600 * 1.157407 * 0.98 ~= 0.157536 $SERVICE3 (from Operator3)
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, charlieAddr, operator1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("18274.736842105000000000service1,113224.500000000000000000service2", rewards.Sum().String())
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, charlieAddr, operator2.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("19290.000000000000000000service1,305426.847222222000000000service3", rewards.Sum().String())
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, charlieAddr, operator3.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("61644.450000000000000000service2,157535.952777777500000000service3", rewards.Sum().String())

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
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, davidAddr, 1)
	suite.Require().NoError(err)
	suite.Assert().Equal("8122.105263157800000000service1,50322.000000000000000000service2", rewards.Sum().String())

	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, davidAddr, 2)
	suite.Require().NoError(err)
	suite.Assert().Equal("12183.157894736800000000service1", rewards.Sum().String())

	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, davidAddr, 3)
	suite.Require().NoError(err)
	suite.Assert().Equal("4061.052631578800000000service1,25161.000000000000000000service2", rewards.Sum().String())

	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, davidAddr, service1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("8122.105263157800000000service1", rewards.Sum().String())

	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, davidAddr, service2.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("50322.000000000000000000service2", rewards.Sum().String())

	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, davidAddr, service3.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("128600.777777777600000000service3", rewards.Sum().String())

	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, davidAddr, operator1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("7309.894736842000000000service1,45289.800000000000000000service2", rewards.Sum().String())

	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, davidAddr, operator2.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("7716.000000000000000000service1,122170.738888888800000000service3", rewards.Sum().String())

	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, davidAddr, operator3.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("49315.560000000000000000service2,126028.762222222200000000service3", rewards.Sum().String())
}

func (suite *KeeperTestSuite) TestAllocateRewards_MovingPrice() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	// $MILK is $2 and $INIT is $3.
	service, _ := suite.setupSampleServiceAndOperator(ctx)

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("100_000000service"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service"),
	)

	delAddr1 := testutil.TestAddress(1)
	suite.DelegateService(ctx, service.ID, utils.MustParseCoins("100_000000umilk"), delAddr1.String(), true)
	delAddr2 := testutil.TestAddress(2)
	suite.DelegateService(ctx, service.ID, utils.MustParseCoins("100_000000uinit"), delAddr2.String(), true)

	// Allocate rewards.
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// They receive rewards by 1:2 ratio.
	rewards, err := suite.keeper.GetServiceDelegationRewards(ctx, delAddr1, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("4629.600000000000000000service", rewards.Sum().String())
	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, delAddr2, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("6944.400000000000000000service", rewards.Sum().String())

	// Now price changes.
	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("3"))
	suite.RegisterCurrency(ctx, "uinit", "INIT", 6, utils.MustParseDec("1"))

	// Allocate rewards again.
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Now they receive rewards by 3:1 ratio.
	// Note that already accumulated rewards are not affected.
	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, delAddr1, service.ID)
	suite.Require().NoError(err)

	// Delta: +8680.5555555555umilk
	suite.Require().Equal("13310.100000000000000000service", rewards.Sum().String())
	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, delAddr2, service.ID)
	suite.Require().NoError(err)

	// Delta: +2893.5185185185umilk
	suite.Require().Equal("9837.900000000000000000service", rewards.Sum().String())
}

func (suite *KeeperTestSuite) TestAllocateRewards_ZeroDelegations() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	// Test if AllocateRewards handles pool/operator/service distribution
	// correctly when the distribution info has weight specified but there's
	// no delegation yet.

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())
	// Whitelist all pools.
	suite.AddPoolsToServiceSecuringPools(ctx, service.ID, []uint32{1})

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("100_000000service"),
		planStartTime,
		planEndTime,
		types.NewBasicPoolsDistribution(1),
		types.NewBasicOperatorsDistribution(2),
		types.NewBasicUsersDistribution(3),
		utils.MustParseCoins("100000_000000service"),
	)

	// Create an operator.
	operatorAdmin := testutil.TestAddress(10001)
	operator := suite.CreateOperator(ctx, "Operator", operatorAdmin.String())

	// Make the operator join the service and set its commission rate to 10%.
	suite.UpdateOperatorParams(ctx, operator.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Try allocating rewards.
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// The rewards pool must be empty since no rewards have been allocated.
	rewardsPool := suite.accountKeeper.GetModuleAddress(types.RewardsPoolName)
	balances := suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().Empty(balances)

	// There must be no outstanding rewards allocated.
	target, err := suite.keeper.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operator.ID)
	suite.Require().NoError(err)
	rewards, err := suite.keeper.GetOutstandingRewardsCoins(ctx, target)
	suite.Require().NoError(err)
	suite.Require().Empty(rewards)

	target, err = suite.keeper.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, service.ID)
	suite.Require().NoError(err)
	rewards, err = suite.keeper.GetOutstandingRewardsCoins(ctx, target)
	suite.Require().NoError(err)
	suite.Require().Empty(rewards)

	// Two users delegate the same amount of $MILK to a pool and the service.
	delAddr1 := testutil.TestAddress(1)
	suite.SetUserPreferences(ctx, delAddr1.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), delAddr1.String(), true)
	delAddr2 := testutil.TestAddress(2)
	suite.DelegateService(ctx, service.ID, utils.MustParseCoins("10_000000umilk"), delAddr2.String(), true)

	// Allocate rewards.
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Still the operator has no rewards.
	target, err = suite.keeper.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operator.ID)
	suite.Require().NoError(err)
	rewards, err = suite.keeper.GetOutstandingRewardsCoins(ctx, target)
	suite.Require().NoError(err)
	suite.Require().Empty(rewards)

	// The pool and the service receive rewards by 1:3 ratio.
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, delAddr1, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("2893.500000000000000000service", rewards.Sum().String())

	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, delAddr2, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("8680.500000000000000000service", rewards.Sum().String())
}

func (suite *KeeperTestSuite) TestAllocateRewards_TransferRewardsOnlyWhenAllocatingRewards_Pool() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Create an active rewards plan.
	suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("100000_000000service"),
	)

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	rewardsPool := suite.accountKeeper.GetModuleAddress(types.RewardsPoolName)

	// There's no delegations yet, so no need to allocate rewards.
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances := suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().Empty(balances)

	// Alice delegates $MILK to the pool.
	aliceAddr := testutil.TestAddress(1)
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), aliceAddr.String(), true)

	// The service is not being secured by the $MILK pool so no rewards allocation
	// happens yet. Note that Alice trusts all services through all pools by default.
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances = suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().Empty(balances)

	// Now the service whitelists the pool, but Alice doesn't trust the service.
	// So no rewards allocation happens yet.
	suite.AddPoolsToServiceSecuringPools(ctx, service.ID, []uint32{1})
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances = suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().Empty(balances)

	// Finally, Alice trusts the service through the pool and rewards are allocated.
	suite.SetUserPreferences(ctx, aliceAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances = suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().NotEmpty(balances)
}

func (suite *KeeperTestSuite) TestAllocateRewards_TransferRewardsOnlyWhenAllocatingRewards_Operator() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Create an inactive operator.
	operatorAdmin := testutil.TestAddress(10001)
	operator := suite.CreateOperator(ctx, "Operator", operatorAdmin.String())

	// Create an active rewards plan.
	suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("100000_000000service"),
	)

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	rewardsPool := suite.accountKeeper.GetModuleAddress(types.RewardsPoolName)

	// There's no delegations yet, so no need to allocate rewards.
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances := suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().Empty(balances)

	// Alice delegates $MILK to the operator.
	aliceAddr := testutil.TestAddress(1)
	suite.DelegateOperator(ctx, operator.ID, utils.MustParseCoins("10_000000umilk"), aliceAddr.String(), true)

	// The operator hasn't joined the service so no rewards allocation happens yet.
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances = suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().Empty(balances)

	// Now the operator joins the service but at the same time it becomes inactive.
	// So no rewards allocation happens yet.
	suite.UpdateOperatorParams(ctx, operator.ID, utils.MustParseDec("0.1"), []uint32{service.ID})
	operator, err = suite.operatorsKeeper.GetOperator(ctx, operator.ID)
	suite.Require().NoError(err)
	operator.Status = operatorstypes.OPERATOR_STATUS_INACTIVE
	err = suite.operatorsKeeper.SaveOperator(ctx, operator)
	suite.Require().NoError(err)
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances = suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().Empty(balances)

	// Now the operator becomes active, but the service doesn't allow the operator to
	// validate it, so no rewards allocation happens yet.
	operator.Status = operatorstypes.OPERATOR_STATUS_ACTIVE
	err = suite.operatorsKeeper.SaveOperator(ctx, operator)
	suite.Require().NoError(err)
	suite.AddOperatorsToServiceAllowList(ctx, service.ID, []uint32{2})
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances = suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().Empty(balances)

	// Finally, the service allows the operator to validate it and rewards are
	// allocated.
	suite.AddOperatorsToServiceAllowList(ctx, service.ID, []uint32{operator.ID})
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances = suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().NotEmpty(balances)
}

func (suite *KeeperTestSuite) TestAllocateRewards_TransferRewardsOnlyWhenAllocatingRewards_Service() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Create an active rewards plan.
	suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("100_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("100000_000000service"),
	)

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	rewardsPool := suite.accountKeeper.GetModuleAddress(types.RewardsPoolName)

	// There's no delegations yet, so no need to allocate rewards.
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances := suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().Empty(balances)

	// Alice delegates $MILK to the service.
	aliceAddr := testutil.TestAddress(1)
	suite.DelegateService(ctx, service.ID, utils.MustParseCoins("10_000000umilk"), aliceAddr.String(), true)

	// The service allocates rewards right after a user delegates to it.
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	balances = suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	suite.Require().NotEmpty(balances)
}

func (suite *KeeperTestSuite) TestAllocateRewards_NoRewardsForNonTrustedPools() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("1"))
	suite.RegisterCurrency(ctx, "utia", "MILK", 6, utils.MustParseDec("3"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Create an active rewards plan.
	suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("1000_000000service"),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		utils.MustParseCoins("100000_000000service"),
	)
	suite.AddPoolsToServiceSecuringPools(ctx, service.ID, []uint32{1})

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Alice, Bob and Charlie delegate $MILK to the pool.
	aliceAddr := testutil.TestAddress(1)
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), aliceAddr.String(), true)
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000utia"), aliceAddr.String(), true)
	bobAddr := testutil.TestAddress(2)
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), bobAddr.String(), true)
	charlieAddr := testutil.TestAddress(3)
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), charlieAddr.String(), true)

	// Alice trusts the service through pool 2
	suite.SetUserPreferences(ctx, aliceAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, []uint32{2}),
	})
	// Bob doesn't trust the service at all.
	// Charlie trusts the service through pool 1
	suite.SetUserPreferences(ctx, charlieAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, []uint32{1}),
	})

	// Only Charlie receives rewards.
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	rewards, err := suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Empty(rewards)
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, bobAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Empty(rewards)
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, charlieAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("115740.000000000000000000service", rewards.Sum().String())
}

func (suite *KeeperTestSuite) TestAllocateRewards_WeightedDistributions() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	// Register $MILK and $INIT. For simple calculation, set both currencies'
	// price $1.
	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("1"))
	suite.RegisterCurrency(ctx, "uinit", "INIT", 6, utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())
	// Whitelist all pools.
	suite.AddPoolsToServiceSecuringPools(ctx, service.ID, []uint32{1, 2})

	// Create operators.
	operatorAdmin1 := testutil.TestAddress(10001)
	operator1 := suite.CreateOperator(ctx, "Operator1", operatorAdmin1.String())
	suite.UpdateOperatorParams(ctx, operator1.ID, utils.MustParseDec("0.1"), []uint32{service.ID})
	operatorAdmin2 := testutil.TestAddress(10002)
	operator2 := suite.CreateOperator(ctx, "Operator2", operatorAdmin2.String())
	suite.UpdateOperatorParams(ctx, operator2.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Delegate to $MILK pool.
	delAddr1 := testutil.TestAddress(1)
	suite.SetUserPreferences(ctx, delAddr1.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("300_000000umilk"), delAddr1.String(), true)
	delAddr2 := testutil.TestAddress(2)
	suite.SetUserPreferences(ctx, delAddr2.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000umilk"), delAddr2.String(), true)
	// Delegate to $INIT pool.
	delAddr3 := testutil.TestAddress(3)
	suite.SetUserPreferences(ctx, delAddr3.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000uinit"), delAddr3.String(), true)
	delAddr4 := testutil.TestAddress(4)
	suite.SetUserPreferences(ctx, delAddr4.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("300_000000uinit"), delAddr4.String(), true)
	// Delegate to Operator1.
	delAddr5 := testutil.TestAddress(5)
	suite.DelegateOperator(ctx, operator1.ID, utils.MustParseCoins("100_000000umilk"), delAddr5.String(), true)
	delAddr6 := testutil.TestAddress(6)
	suite.DelegateOperator(ctx, operator1.ID, utils.MustParseCoins("200_000000uinit"), delAddr6.String(), true)
	// Delegate to Operator2.
	delAddr7 := testutil.TestAddress(7)
	suite.DelegateOperator(ctx, operator2.ID, utils.MustParseCoins("200_000000umilk"), delAddr7.String(), true)
	delAddr8 := testutil.TestAddress(8)
	suite.DelegateOperator(ctx, operator2.ID, utils.MustParseCoins("100_000000uinit"), delAddr8.String(), true)

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("100_000000service"),
		planStartTime,
		planEndTime,
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

	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// delAddr1 receives 3/4 * 2/3 * $300 / $500 * 100 * (10 / 86400) ~= 0.003472 $SERVICE
	rewards, err := suite.keeper.GetPoolDelegationRewards(ctx, delAddr1, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("3472.200000000000000000service", rewards.Sum().String())

	// delAddr2 receives 3/4 * 2/3 * $200 / $500 * 100 * (10 / 86400) ~= 0.002315 $SERVICE
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, delAddr2, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("2314.800000000000000000service", rewards.Sum().String())

	// delAddr3 receives 3/4 * 1/3 * $200 / $500 * 100 * (10 / 86400) ~= 0.001157 $SERVICE
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, delAddr3, 2)
	suite.Require().NoError(err)
	suite.Require().Equal("1157.400000000000000000service", rewards.Sum().String())

	// delAddr4 receives 3/4 * 1/3 * $300 / $500 * 100 * (10 / 86400) ~= 0.001736 $SERVICE
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, delAddr4, 2)
	suite.Require().NoError(err)
	suite.Require().Equal("1736.100000000000000000service", rewards.Sum().String())

	// Note that operators take commission from rewards.

	// delAddr5 receives 1/4 * 2/5 * $100 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000347 $SERVICE
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, delAddr5, operator1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("347.220000000000000000service", rewards.Sum().String())

	// delAddr6 receives 1/4 * 2/5 * $200 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000694 $SERVICE
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, delAddr6, operator1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("694.440000000000000000service", rewards.Sum().String())

	// delAddr7 receives 1/4 * 3/5 * $200 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.001042 $SERVICE
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, delAddr7, operator2.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("1041.660000000000000000service", rewards.Sum().String())

	// delAddr8 receives 1/4 * 3/5 * $100 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000521 $SERVICE
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, delAddr8, operator2.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("520.830000000000000000service", rewards.Sum().String())
}

func (suite *KeeperTestSuite) TestAllocateRewards_EgalitarianDistributions() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	// Register $MILK and $INIT. For simple calculation, set both currencies'
	// price $1.
	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("1"))
	suite.RegisterCurrency(ctx, "uinit", "INIT", 6, utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())
	// Whitelist all pools.
	suite.AddPoolsToServiceSecuringPools(ctx, service.ID, []uint32{1, 2})

	// Create operators.
	operatorAdmin1 := testutil.TestAddress(10001)
	operator1 := suite.CreateOperator(ctx, "Operator1", operatorAdmin1.String())
	suite.UpdateOperatorParams(ctx, operator1.ID, utils.MustParseDec("0.1"), []uint32{service.ID})
	operatorAdmin2 := testutil.TestAddress(10002)
	operator2 := suite.CreateOperator(ctx, "Operator2", operatorAdmin2.String())
	suite.UpdateOperatorParams(ctx, operator2.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("100_000000service"),
		planStartTime,
		planEndTime,
		types.NewEgalitarianPoolsDistribution(3),
		types.NewEgalitarianOperatorsDistribution(1),
		types.NewBasicUsersDistribution(0), // No user rewards
		utils.MustParseCoins("100000_000000service"),
	)

	// Delegate to $MILK pool.
	delAddr1 := testutil.TestAddress(1)
	suite.SetUserPreferences(ctx, delAddr1.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("300_000000umilk"), delAddr1.String(), true)
	delAddr2 := testutil.TestAddress(2)
	suite.SetUserPreferences(ctx, delAddr2.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000umilk"), delAddr2.String(), true)

	// Delegate to $INIT pool.
	delAddr3 := testutil.TestAddress(3)
	suite.SetUserPreferences(ctx, delAddr3.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000uinit"), delAddr3.String(), true)
	delAddr4 := testutil.TestAddress(4)
	suite.SetUserPreferences(ctx, delAddr4.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("300_000000uinit"), delAddr4.String(), true)

	// Delegate to Operator1.
	delAddr5 := testutil.TestAddress(5)
	suite.DelegateOperator(ctx, operator1.ID, utils.MustParseCoins("100_000000umilk"), delAddr5.String(), true)
	delAddr6 := testutil.TestAddress(6)
	suite.DelegateOperator(ctx, operator1.ID, utils.MustParseCoins("200_000000uinit"), delAddr6.String(), true)

	// Delegate to Operator2.
	delAddr7 := testutil.TestAddress(7)
	suite.DelegateOperator(ctx, operator2.ID, utils.MustParseCoins("200_000000umilk"), delAddr7.String(), true)
	delAddr8 := testutil.TestAddress(8)
	suite.DelegateOperator(ctx, operator2.ID, utils.MustParseCoins("100_000000uinit"), delAddr8.String(), true)

	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// delAddr1 receives 3/4 * 1/2 * $300 / $500 * 100 * (10 / 86400) ~= 0.002604 $SERVICE
	rewards, err := suite.keeper.GetPoolDelegationRewards(ctx, delAddr1, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("2604.150000000000000000service", rewards.Sum().String())

	// delAddr2 receives 3/4 * 1/2 * $200 / $500 * 100 * (10 / 86400) ~= 0.001736 $SERVICE
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, delAddr2, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("1736.100000000000000000service", rewards.Sum().String())

	// delAddr3 receives 3/4 * 1/2 * $200 / $500 * 100 * (10 / 86400) ~= 0.001736 $SERVICE
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, delAddr3, 2)
	suite.Require().NoError(err)
	suite.Require().Equal("1736.100000000000000000service", rewards.Sum().String())

	// delAddr4 receives 3/4 * 1/2 * $300 / $500 * 100 * (10 / 86400) ~= 0.002604 $SERVICE
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, delAddr4, 2)
	suite.Require().NoError(err)
	suite.Require().Equal("2604.150000000000000000service", rewards.Sum().String())

	// Note that operators take commission from rewards.

	// delAddr5 receives 1/4 * 1/2 * $100 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000434 $SERVICE
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, delAddr5, operator1.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("434.025000000000000000service", rewards.Sum().String())

	// delAddr6 receives 1/4 * 1/2 * $200 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000868 $SERVICE
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, delAddr6, operator1.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("868.050000000000000000service", rewards.Sum().String())

	// delAddr7 receives 1/4 * 1/2 * $200 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000868 $SERVICE
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, delAddr7, operator2.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("868.050000000000000000service", rewards.Sum().String())

	// delAddr8 receives 1/4 * 1/2 * $100 / $300 * 100 * (10 / 86400) * 0.9 ~= 0.000434 $SERVICE
	rewards, err = suite.keeper.GetOperatorDelegationRewards(ctx, delAddr8, operator2.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("434.025000000000000000service", rewards.Sum().String())
}

func (suite *KeeperTestSuite) TestAllocateRewards_TrustedServices() {
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))

	// Create services.
	serviceAdmin1 := testutil.TestAddress(10000)
	service1 := suite.CreateService(ctx, "Service1", serviceAdmin1.String())
	// Whitelist all pools.
	suite.AddPoolsToServiceSecuringPools(ctx, service1.ID, []uint32{1})
	serviceAdmin2 := testutil.TestAddress(10001)
	service2 := suite.CreateService(ctx, "Service2", serviceAdmin2.String())
	// Whitelist all pools.
	suite.AddPoolsToServiceSecuringPools(ctx, service2.ID, []uint32{1})

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Create active rewards plans.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateBasicRewardsPlan(
		ctx,
		service1.ID,
		utils.MustParseCoin("1000_000000service1"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service1"),
	)
	suite.CreateBasicRewardsPlan(
		ctx,
		service2.ID,
		utils.MustParseCoin("5000_000000service2"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service2"),
	)

	// Delegate to $MILK pool.
	aliceAddr := testutil.TestAddress(1)
	suite.SetUserPreferences(ctx, aliceAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service1.ID, nil),
		restakingtypes.NewTrustedServiceEntry(service2.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("300_000000umilk"), aliceAddr.String(), true)

	bobAddr := testutil.TestAddress(2)
	suite.SetUserPreferences(ctx, bobAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service2.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000umilk"), bobAddr.String(), true)

	// Rewards plan 1 allocates 1000 * 10 / 86400 ~= 0.115741 $SERVICE1
	// Rewards plan 2 allocates 5000 * 10 / 86400 ~= 0.578704 $SERVICE1
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Alice receives:
	// - $600 / $600 * 0.115741 ~= 0.115741 $SERVICE1
	// - $600 / $1000 * 0.578704 ~= 0.347222 $SERVICE2
	rewards, err := suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("115740.000000000000000000service1,347221.800000000000000000service2", rewards.Sum().String())

	// Bob receives:
	// - $400 / $1000 * 0.578704 ~= 0.231482 $SERVICE2
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, bobAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("231481.200000000000000000service2", rewards.Sum().String())

	// Bob trusts the service as well.
	suite.SetUserPreferences(ctx, bobAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service1.ID, nil),
		restakingtypes.NewTrustedServiceEntry(service2.ID, nil),
	})

	// Withdraw all rewards to make calculation easier.
	_, err = keeper.NewMsgServer(suite.keeper).WithdrawDelegatorReward(
		ctx,
		types.NewMsgWithdrawDelegatorReward(restakingtypes.DELEGATION_TYPE_POOL, 1, aliceAddr.String()),
	)
	suite.Require().NoError(err)
	_, err = keeper.NewMsgServer(suite.keeper).WithdrawDelegatorReward(
		ctx,
		types.NewMsgWithdrawDelegatorReward(restakingtypes.DELEGATION_TYPE_POOL, 1, bobAddr.String()),
	)
	suite.Require().NoError(err)

	// Rewards plan 1 allocates 1000 * 10 / 86400 ~= 0.115741 $SERVICE1
	// Rewards plan 2 allocates 5000 * 10 / 86400 ~= 0.578704 $SERVICE1
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Alice receives:
	// - $600 / $1000 * 0.115741 ~= 0.069445 $SERVICE1
	// - $600 / $1000 * 0.578704 ~= 0.347222 $SERVICE2
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("69444.000000000000000000service1,347221.800000000000000000service2", rewards.Sum().String())

	// Bob receives:
	// - $400 / $1000 * 0.115741 ~= 0.046296 $SERVICE1
	// - $400 / $1000 * 0.578704 ~= 0.231482 $SERVICE2
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, bobAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("46296.000000000000000000service1,231481.200000000000000000service2", rewards.Sum().String())

	// Now Alice decides to not trust service 2.
	// This will make Alice's rewards for the pool to be withdrawn automatically.
	suite.SetUserPreferences(ctx, aliceAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service1.ID, nil),
	})

	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("", rewards.Sum().String())
}

func (suite *KeeperTestSuite) TestAllocateRewards_UserTrustedServiceUpdated() {
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))

	// Create services.
	serviceAdmin1 := testutil.TestAddress(10000)
	service1 := suite.CreateService(ctx, "Service1", serviceAdmin1.String())
	// Whitelist all pools.
	suite.AddPoolsToServiceSecuringPools(ctx, service1.ID, []uint32{1})
	serviceAdmin2 := testutil.TestAddress(10001)
	service2 := suite.CreateService(ctx, "Service2", serviceAdmin2.String())
	// Whitelist all pools.
	suite.AddPoolsToServiceSecuringPools(ctx, service2.ID, []uint32{1})

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Create active rewards plans.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateBasicRewardsPlan(
		ctx,
		service1.ID,
		utils.MustParseCoin("1000_000000service1"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service1"),
	)
	suite.CreateBasicRewardsPlan(
		ctx,
		service2.ID,
		utils.MustParseCoin("5000_000000service2"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service2"),
	)

	// Delegate to $MILK pool.
	aliceAddr := testutil.TestAddress(1)
	suite.SetUserPreferences(ctx, aliceAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service1.ID, nil),
		restakingtypes.NewTrustedServiceEntry(service2.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("300_000000umilk"), aliceAddr.String(), true)

	bobAddr := testutil.TestAddress(2)
	suite.SetUserPreferences(ctx, bobAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service2.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000umilk"), bobAddr.String(), true)

	// Rewards plan 1 allocates 1000 * 10 / 86400 ~= 0.115741 $SERVICE1
	// Rewards plan 2 allocates 5000 * 10 / 86400 ~= 0.578704 $SERVICE1
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Alice receives:
	// - $600 / $600 * 0.115741 ~= 0.115741 $SERVICE1
	// - $600 / $1000 * 0.578704 ~= 0.347222 $SERVICE2
	rewards, err := suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("115740.000000000000000000service1,347221.800000000000000000service2", rewards.Sum().String())

	// Bob receives:
	// - $400 / $1000 * 0.578704 ~= 0.231482 $SERVICE2
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, bobAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("231481.200000000000000000service2", rewards.Sum().String())

	// Charlie doesn't trust any services.
	charlieAddr := testutil.TestAddress(3)
	suite.DelegatePool(ctx, utils.MustParseCoin("200_000000umilk"), charlieAddr.String(), true)

	// Withdraw all rewards to make calculation easier.
	_, err = keeper.NewMsgServer(suite.keeper).WithdrawDelegatorReward(
		ctx,
		types.NewMsgWithdrawDelegatorReward(restakingtypes.DELEGATION_TYPE_POOL, 1, aliceAddr.String()),
	)
	suite.Require().NoError(err)
	_, err = keeper.NewMsgServer(suite.keeper).WithdrawDelegatorReward(
		ctx,
		types.NewMsgWithdrawDelegatorReward(restakingtypes.DELEGATION_TYPE_POOL, 1, bobAddr.String()),
	)
	suite.Require().NoError(err)

	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Rewards amount must not be changed since Charlie doesn't trust any services,
	// thus receives no rewards.
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("115740.000000000000000000service1,347221.800000000000000000service2", rewards.Sum().String())
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, bobAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("231481.200000000000000000service2", rewards.Sum().String())
}

func (suite *KeeperTestSuite) TestAllocateRewards_InactiveService() {
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))

	// Create services. They are active by default because CreateService helper
	// activated them automatically.
	serviceAdmin1 := testutil.TestAddress(10001)
	service1 := suite.CreateService(ctx, "Service1", serviceAdmin1.String())
	serviceAdmin2 := testutil.TestAddress(10002)
	service2 := suite.CreateService(ctx, "Service2", serviceAdmin2.String())

	// Create rewards plans.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateBasicRewardsPlan(
		ctx,
		service1.ID,
		utils.MustParseCoin("1000_000000service1"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service1"),
	)
	suite.CreateBasicRewardsPlan(
		ctx,
		service2.ID,
		utils.MustParseCoin("2000_000000service2"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service2"),
	)

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Alice delegates to the services.
	aliceAddr := testutil.TestAddress(1)
	suite.DelegateService(ctx, service1.ID, utils.MustParseCoins("300_000000umilk"), aliceAddr.String(), true)
	suite.DelegateService(ctx, service2.ID, utils.MustParseCoins("100_000000umilk"), aliceAddr.String(), true)

	// Try allocating rewards.
	// Service 1 allocates 1000 * 10 / 86400 ~= 0.115741 $SERVICE1
	// Service 2 allocates 2000 * 10 / 86400 ~= 0.231481 $SERVICE2
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Both services allocated rewards.
	rewards, err := suite.keeper.GetDelegationRewards(ctx, aliceAddr, restakingtypes.DELEGATION_TYPE_SERVICE, service1.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("115740.000000000000000000service1", rewards.Sum().String())
	rewards, err = suite.keeper.GetDelegationRewards(ctx, aliceAddr, restakingtypes.DELEGATION_TYPE_SERVICE, service2.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("231481.000000000000000000service2", rewards.Sum().String())

	// Withdraw rewards from services to make calculation easier.
	_, err = keeper.NewMsgServer(suite.keeper).WithdrawDelegatorReward(
		ctx,
		types.NewMsgWithdrawDelegatorReward(restakingtypes.DELEGATION_TYPE_SERVICE, service1.ID, aliceAddr.String()),
	)
	suite.Require().NoError(err)
	_, err = keeper.NewMsgServer(suite.keeper).WithdrawDelegatorReward(
		ctx,
		types.NewMsgWithdrawDelegatorReward(restakingtypes.DELEGATION_TYPE_SERVICE, service2.ID, aliceAddr.String()),
	)
	suite.Require().NoError(err)

	// Service 1 becomes inactive.
	err = suite.servicesKeeper.DeactivateService(ctx, service1.ID)
	suite.Require().NoError(err)

	// Try allocating rewards again.
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// There's no rewards allocated by service 1 because it was inactive,
	rewards, err = suite.keeper.GetDelegationRewards(ctx, aliceAddr, restakingtypes.DELEGATION_TYPE_SERVICE, service1.ID)
	suite.Require().NoError(err)
	suite.Require().True(rewards.IsEmpty())
	// but service 2 allocated rewards.
	rewards, err = suite.keeper.GetDelegationRewards(ctx, aliceAddr, restakingtypes.DELEGATION_TYPE_SERVICE, service2.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("231481.000000000000000000service2", rewards.Sum().String())
}

func (suite *KeeperTestSuite) TestAllocateRewards_InactiveOperator() {
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10001)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Create operators. They are active by default when creating.
	operatorAdmin1 := testutil.TestAddress(10002)
	operator1 := suite.CreateOperator(ctx, "Operator1", operatorAdmin1.String())
	operatorAdmin2 := testutil.TestAddress(10003)
	operator2 := suite.CreateOperator(ctx, "Operator2", operatorAdmin2.String())

	// Operators set their commission rate to 10% and join the service.
	suite.UpdateOperatorParams(ctx, operator1.ID, utils.MustParseDec("0.1"), []uint32{service.ID})
	suite.UpdateOperatorParams(ctx, operator2.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Create a rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("1000_000000service"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service"),
	)

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Alice delegates to both operators.
	aliceAddr := testutil.TestAddress(1)
	suite.DelegateOperator(ctx, operator1.ID, utils.MustParseCoins("300_000000umilk"), aliceAddr.String(), true)
	suite.DelegateOperator(ctx, operator2.ID, utils.MustParseCoins("100_000000umilk"), aliceAddr.String(), true)

	// Try allocating rewards.
	// The service allocates 1000 * 10 / 86400 ~= 0.115741 $SERVICE
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Both operators received rewards.
	// Alice receives $600 / $800 * 0.115741 * 0.9 ~= 0.078125 $SERVICE from operator 1.
	rewards, err := suite.keeper.GetDelegationRewards(ctx, aliceAddr, restakingtypes.DELEGATION_TYPE_OPERATOR, operator1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("78124.500000000000000000service", rewards.Sum().String())
	// Alice receives $200 / $800 * 0.115741 * 0.9 ~= 0.026042 $SERVICE from operator 1.
	rewards, err = suite.keeper.GetDelegationRewards(ctx, aliceAddr, restakingtypes.DELEGATION_TYPE_OPERATOR, operator2.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("26041.500000000000000000service", rewards.Sum().String())

	// Withdraw rewards from operators to make calculation easier.
	_, err = keeper.NewMsgServer(suite.keeper).WithdrawDelegatorReward(
		ctx,
		types.NewMsgWithdrawDelegatorReward(restakingtypes.DELEGATION_TYPE_OPERATOR, operator1.ID, aliceAddr.String()),
	)
	suite.Require().NoError(err)
	_, err = keeper.NewMsgServer(suite.keeper).WithdrawDelegatorReward(
		ctx,
		types.NewMsgWithdrawDelegatorReward(restakingtypes.DELEGATION_TYPE_OPERATOR, operator2.ID, aliceAddr.String()),
	)
	suite.Require().NoError(err)

	// Refresh the updated state of operator 2.
	operator2, err = suite.operatorsKeeper.GetOperator(ctx, operator2.ID)
	suite.Require().NoError(err)
	// Operator 2 becomes inactive.
	err = suite.operatorsKeeper.StartOperatorInactivation(ctx, operator2)
	suite.Require().NoError(err)

	// Try allocating rewards again.
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// This time Alice receives $600 / $600 * 0.115741 * 0.9 ~= 0.104167 $SERVICE
	// from operator 1.
	rewards, err = suite.keeper.GetDelegationRewards(ctx, aliceAddr, restakingtypes.DELEGATION_TYPE_OPERATOR, operator1.ID)
	suite.Require().NoError(err)
	suite.Assert().Equal("104166.000000000000000000service", rewards.Sum().String())
	// There's no rewards allocated to operator 2 because it was inactive.
	rewards, err = suite.keeper.GetDelegationRewards(ctx, aliceAddr, restakingtypes.DELEGATION_TYPE_OPERATOR, operator2.ID)
	suite.Require().NoError(err)
	suite.Assert().True(rewards.IsEmpty())
}

func (suite *KeeperTestSuite) TestAllocateRewards_ExcludeNonTrustedTokens() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	// Register $MILK
	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("1000_000000service"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service"),
	)

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Alice trusts the service and delegates $MILK to the pool.
	aliceAddr := testutil.TestAddress(1)
	suite.SetUserPreferences(ctx, aliceAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), aliceAddr.String(), true)

	// Add the pool to the service's securing pools list so that rewards can be
	// allocated to the pool.
	suite.AddPoolsToServiceSecuringPools(ctx, service.ID, []uint32{1})

	// Bob doesn't trust the service but still delegates $MILK to the pool.
	bobAddr := testutil.TestAddress(2)
	suite.DelegatePool(ctx, utils.MustParseCoin("1000_000000umilk"), bobAddr.String(), true)

	// Also, Alice delegates $MILK to the service directly.
	suite.DelegateService(ctx, service.ID, utils.MustParseCoins("10_000000umilk"), aliceAddr.String(), true)

	// Try allocating rewards.
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Alice receives the same amount of rewards from both the pool and the service
	// delegations.
	rewards, err := suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	// Rewards plan allocates total of 1000 * 10 / 86400 ~= 0.115741 $SERVICE and
	// Alice receives half of it as a pool restaker, and half of it as a service restaker
	suite.Require().Equal("57870.000000000000000000service", rewards.Sum().String())
	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, aliceAddr, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("57870.000000000000000000service", rewards.Sum().String())

	// Bob didn't receive rewards from the service since he doesn't trust it.
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, bobAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Empty(rewards)
}

func (suite *KeeperTestSuite) TestAllocateRewards_NoRewardsAfterUnbonding() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	// Register $MILK
	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Create an active rewards plan.
	planStartTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	planEndTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	plan := suite.CreateBasicRewardsPlan(
		ctx,
		service.ID,
		utils.MustParseCoin("1000_000000service"),
		planStartTime,
		planEndTime,
		utils.MustParseCoins("100000_000000service"),
	)

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	// Alice trusts the service and delegates $MILK to the pool.
	aliceAddr := testutil.TestAddress(1)
	suite.SetUserPreferences(ctx, aliceAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), aliceAddr.String(), true)

	// Add the pool to the service's securing pools list so that rewards can be
	// allocated to the pool.
	suite.AddPoolsToServiceSecuringPools(ctx, service.ID, []uint32{1})

	// Bob doesn't trust the service but still delegates $MILK to the pool.
	bobAddr := testutil.TestAddress(2)
	suite.DelegatePool(ctx, utils.MustParseCoin("1000_000000umilk"), bobAddr.String(), true)

	// Also, Alice delegates $MILK to the service directly.
	suite.DelegateService(ctx, service.ID, utils.MustParseCoins("10_000000umilk"), aliceAddr.String(), true)

	// Try allocating rewards.
	ctx = suite.allocateRewards(ctx, 10*time.Second)

	// Alice receives the same amount of rewards from both the pool and the service
	// delegations.
	rewards, err := suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	// Rewards plan allocates total of 1000 * 10 / 86400 ~= 0.115741 $SERVICE and
	// Alice receives half of it as a pool restaker, and half of it as a service restaker
	suite.Require().Equal("57870.000000000000000000service", rewards.Sum().String())
	rewards, err = suite.keeper.GetServiceDelegationRewards(ctx, aliceAddr, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal("57870.000000000000000000service", rewards.Sum().String())

	// Bob didn't receive rewards from the service since he doesn't trust it.
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, bobAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Empty(rewards)

	// Now Alice unbonds all tokens from the service.
	_, err = suite.restakingKeeper.UndelegateFromService(
		ctx,
		service.ID,
		utils.MustParseCoins("10_000000umilk"),
		aliceAddr.String(),
	)
	suite.Require().NoError(err)
	// Also make sure Alice withdraws all her rewards frm the pool to make the
	// calculation easier.
	_, err = keeper.NewMsgServer(suite.keeper).WithdrawDelegatorReward(
		ctx,
		types.NewMsgWithdrawDelegatorReward(restakingtypes.DELEGATION_TYPE_POOL, 1, aliceAddr.String()),
	)
	suite.Require().NoError(err)

	// Alice doesn't receive rewards from the service immediately after unbonding.
	// Instead, all rewards now go to the pool.
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	_, err = suite.keeper.GetServiceDelegationRewards(ctx, aliceAddr, service.ID)
	suite.Require().ErrorIs(err, sdkerrors.ErrNotFound)
	rewards, err = suite.keeper.GetPoolDelegationRewards(ctx, aliceAddr, 1)
	suite.Require().NoError(err)
	suite.Require().Equal("115740.000000000000000000service", rewards.Sum().String())

	// Alice also unbonds from the pool, making the rewards allocation to stop
	// entirely.
	_, err = suite.restakingKeeper.UndelegateFromPool(ctx, utils.MustParseCoin("10_000000umilk"), aliceAddr.String())
	suite.Require().NoError(err)

	// Make sure no rewards were distributed.
	rewardsPool := suite.accountKeeper.GetModuleAddress(types.RewardsPoolName)
	moduleBalancesBefore := suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	poolBalancesBefore := suite.bankKeeper.GetAllBalances(
		ctx,
		plan.MustGetRewardsPoolAddress(suite.accountKeeper.AddressCodec()),
	)
	ctx = suite.allocateRewards(ctx, 10*time.Second)
	moduleBalancesAfter := suite.bankKeeper.GetAllBalances(ctx, rewardsPool)
	poolBalancesAfter := suite.bankKeeper.GetAllBalances(
		ctx,
		plan.MustGetRewardsPoolAddress(suite.accountKeeper.AddressCodec()),
	)
	suite.Require().Equal(moduleBalancesBefore, moduleBalancesAfter)
	suite.Require().Equal(poolBalancesBefore, poolBalancesAfter)
}
