package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oraclekeeper "github.com/skip-mev/connect/v2/x/oracle/keeper"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/app/testutil"
	"github.com/milkyway-labs/milkyway/utils"
	assetskeeper "github.com/milkyway-labs/milkyway/x/assets/keeper"
	assetstypes "github.com/milkyway-labs/milkyway/x/assets/types"
	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/x/rewards/testutils"
	rewardstypes "github.com/milkyway-labs/milkyway/x/rewards/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type KeeperTestSuite struct {
	suite.Suite

	authority string

	ctx sdk.Context
	cdc codec.Codec

	accountKeeper   authkeeper.AccountKeeper
	bankKeeper      bankkeeper.Keeper
	marketMapKeeper *marketmapkeeper.Keeper
	oracleKeeper    oraclekeeper.Keeper
	assetsKeeper    *assetskeeper.Keeper
	poolsKeeper     *poolskeeper.Keeper
	servicesKeeper  *serviceskeeper.Keeper
	operatorsKeeper *operatorskeeper.Keeper
	restakingKeeper *restakingkeeper.Keeper

	keeper *keeper.Keeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())

	suite.authority = data.AuthorityAddress

	suite.cdc = data.Cdc
	suite.ctx = data.Context

	suite.accountKeeper = data.AccountKeeper
	suite.bankKeeper = data.BankKeeper
	suite.marketMapKeeper = data.MarketMapKeeper
	suite.oracleKeeper = data.OracleKeeper
	suite.assetsKeeper = data.AssetsKeeper
	suite.poolsKeeper = data.PoolsKeeper
	suite.servicesKeeper = data.ServicesKeeper
	suite.operatorsKeeper = data.OperatorsKeeper
	suite.restakingKeeper = data.RestakingKeeper
	suite.keeper = data.Keeper
}

// FundAccount adds the given amount of coins to the account with the given address
func (suite *KeeperTestSuite) FundAccount(ctx sdk.Context, addr string, amt sdk.Coins) {
	moduleAcc := suite.accountKeeper.GetModuleAccount(ctx, minttypes.ModuleName)

	// Mint the coins
	err := suite.bankKeeper.MintCoins(ctx, moduleAcc.GetName(), amt)
	suite.Require().NoError(err)

	// Send the amount to the user
	accAddr, err := sdk.AccAddressFromBech32(addr)
	suite.Require().NoError(err)

	err = suite.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAcc.GetName(), accAddr, amt)
	suite.Require().NoError(err)
}

// RegisterCurrency registers a currency with the given denomination, ticker
// and price. RegisterCurrency creates a market for the currency if not exists.
func (suite *KeeperTestSuite) RegisterCurrency(ctx sdk.Context, denom string, ticker string, exponent uint32, price math.LegacyDec) {
	// Create the market only if it doesn't exist.
	mmTicker := marketmaptypes.NewTicker(ticker, rewardstypes.USDTicker, math.LegacyPrecision, 0, true)
	hasMarket, err := suite.marketMapKeeper.HasMarket(ctx, mmTicker.String())
	suite.Require().NoError(err)

	if !hasMarket {
		err = suite.marketMapKeeper.CreateMarket(ctx, marketmaptypes.Market{Ticker: mmTicker})
		suite.Require().NoError(err)
	}

	// Set the price for the currency pair.
	err = suite.oracleKeeper.SetPriceForCurrencyPair(
		ctx,
		connecttypes.NewCurrencyPair(ticker, rewardstypes.USDTicker),
		oracletypes.QuotePrice{
			Price:          math.NewIntFromBigInt(price.BigInt()),
			BlockTimestamp: ctx.BlockTime(),
			BlockHeight:    uint64(ctx.BlockHeight()),
		},
	)
	suite.Require().NoError(err)

	// Register the currency.
	err = suite.assetsKeeper.SetAsset(ctx, assetstypes.NewAsset(denom, ticker, exponent))
	suite.Require().NoError(err)
}

// CreateService creates an example service with the given service name and
// admin address. The service description and URLs related to the service are
// randomly chosen. The service is also activated.
func (suite *KeeperTestSuite) CreateService(ctx sdk.Context, name string, admin string) servicestypes.Service {
	// Register the service
	servicesMsgServer := serviceskeeper.NewMsgServer(suite.servicesKeeper)
	resp, err := servicesMsgServer.CreateService(ctx, servicestypes.NewMsgCreateService(
		name,
		fmt.Sprintf("%s AVS", name),
		"https://example.com",
		"https://example.com/picture.png",
		admin,
	))
	suite.Require().NoError(err)

	// Also activate the service.
	_, err = servicesMsgServer.ActivateService(ctx, servicestypes.NewMsgActivateService(resp.NewServiceID, admin))
	suite.Require().NoError(err)

	service, found, err := suite.servicesKeeper.GetService(ctx, resp.NewServiceID)
	suite.Require().NoError(err)
	suite.Require().True(found, "service must be found")
	return service
}

// CreateOperator creates an example operator with the given operator name and
// admin address. The operator description and URLs related to the operator are
// randomly chosen.
func (suite *KeeperTestSuite) CreateOperator(ctx sdk.Context, name string, admin string) operatorstypes.Operator {
	// Register the operator
	operatorsMsgServer := operatorskeeper.NewMsgServer(suite.operatorsKeeper)
	resp, err := operatorsMsgServer.RegisterOperator(ctx, operatorstypes.NewMsgRegisterOperator(
		name,
		"https://example.com",
		"https://example.com/picture.png",
		nil,
		admin,
	))
	suite.Require().NoError(err)

	// Make sure the operator is found
	operator, found, err := suite.operatorsKeeper.GetOperator(ctx, resp.NewOperatorID)
	suite.Require().NoError(err)
	suite.Require().True(found, "operator must be found")
	return operator
}

// UpdateOperatorParams updates the operator's params.
// TODO: split functionalities
func (suite *KeeperTestSuite) UpdateOperatorParams(
	ctx sdk.Context,
	operatorID uint32,
	commissionRate math.LegacyDec,
	joinedServicesIDs []uint32,
) {
	// Make sure the operator is found
	_, found, err := suite.operatorsKeeper.GetOperator(ctx, operatorID)
	suite.Require().NoError(err)
	suite.Require().True(found, "operator must be found")

	// Sets the operator commission rate
	err = suite.operatorsKeeper.SaveOperatorParams(ctx, operatorID, operatorstypes.NewOperatorParams(commissionRate))
	suite.Require().NoError(err)

	// Make the operator join the service.
	for _, serviceID := range joinedServicesIDs {
		err = suite.restakingKeeper.AddServiceToOperatorJoinedServices(ctx, operatorID, serviceID)
		suite.Require().NoError(err)
	}
}

// SetUserPreferences sets the user's preferences.
func (suite *KeeperTestSuite) SetUserPreferences(
	ctx sdk.Context,
	userAddress string,
	trustNonAccreditedServices,
	trustAccreditedServices bool,
	trustedServicesIDs []uint32,
) {
	err := suite.restakingKeeper.SetUserPreferences(
		ctx,
		userAddress,
		restakingtypes.NewUserPreferences(trustNonAccreditedServices, trustAccreditedServices, trustedServicesIDs),
	)
	suite.Require().NoError(err)
}

// AddPoolsToServiceSecuringPools adds the provided pools the list of
// pools from which the service can borrow security.
func (suite *KeeperTestSuite) AddPoolsToServiceSecuringPools(
	ctx sdk.Context,
	serviceID uint32,
	whitelistedPoolsIDs []uint32,
) {
	// Make sure the service is found
	_, found, err := suite.servicesKeeper.GetService(ctx, serviceID)
	suite.Require().NoError(err)
	suite.Require().True(found, "service must be found")

	for _, poolID := range whitelistedPoolsIDs {
		err := suite.restakingKeeper.AddPoolToServiceSecuringPools(ctx, serviceID, poolID)
		suite.Require().NoError(err)
	}
}

// AddOperatorsToServiceAllowList adds the given operators to the list of
// operators allowed to secure the service.
func (suite *KeeperTestSuite) AddOperatorsToServiceAllowList(
	ctx sdk.Context,
	serviceID uint32,
	allowedOperatorsID []uint32,
) {
	// Make sure the service is found
	_, found, err := suite.servicesKeeper.GetService(ctx, serviceID)
	suite.Require().NoError(err)
	suite.Require().True(found, "service must be found")

	for _, operatorID := range allowedOperatorsID {
		err := suite.restakingKeeper.AddOperatorToServiceAllowList(ctx, serviceID, operatorID)
		suite.Require().NoError(err)
	}
}

// CreateRewardsPlan creates a rewards plan with the given parameters.
// The plan's name is chosen randomly. The plan is also funded with the given
// initial rewards.
func (suite *KeeperTestSuite) CreateRewardsPlan(
	ctx sdk.Context,
	serviceID uint32,
	amtPerDay sdk.Coins,
	startTime time.Time,
	endTime time.Time,
	poolsDistr rewardstypes.Distribution,
	operatorsDistr rewardstypes.Distribution,
	usersDistr rewardstypes.UsersDistribution,
	initialRewards sdk.Coins,
) rewardstypes.RewardsPlan {
	service, found, err := suite.servicesKeeper.GetService(ctx, serviceID)
	suite.Require().NoError(err)
	suite.Require().True(found, "service must be found")

	rewardsMsgServer := keeper.NewMsgServer(suite.keeper)
	resp, err := rewardsMsgServer.CreateRewardsPlan(ctx, rewardstypes.NewMsgCreateRewardsPlan(
		serviceID,
		"Rewards Plan",
		amtPerDay,
		startTime,
		endTime,
		poolsDistr,
		operatorsDistr,
		usersDistr,
		service.Admin,
	))
	suite.Require().NoError(err)

	// Return the newly created plan.
	plan, err := suite.keeper.GetRewardsPlan(ctx, resp.NewRewardsPlanID)
	suite.Require().NoError(err)

	if initialRewards.IsAllPositive() {
		suite.FundAccount(ctx, plan.RewardsPool, initialRewards)
	}

	return plan
}

// CreateBasicRewardsPlan creates a rewards plan with basic distribution for
// all restaking entities. Weights among the entities are set to 0, which means
// that the rewards are distributed based on the total delegation values among
// all entities.
func (suite *KeeperTestSuite) CreateBasicRewardsPlan(
	ctx sdk.Context,
	serviceID uint32,
	amtPerDay sdk.Coins,
	startTime,
	endTime time.Time,
	initialRewards sdk.Coins,
) rewardstypes.RewardsPlan {
	return suite.CreateRewardsPlan(
		ctx,
		serviceID,
		amtPerDay,
		startTime,
		endTime,
		rewardstypes.NewBasicPoolsDistribution(0),
		rewardstypes.NewBasicOperatorsDistribution(0),
		rewardstypes.NewBasicUsersDistribution(0),
		initialRewards,
	)
}

// DelegateOperator delegates the given amount of coins to the operator. If fund
// is true, the delegator's account is funded with the given amount of coins.
func (suite *KeeperTestSuite) DelegateOperator(
	ctx sdk.Context,
	operatorID uint32,
	amt sdk.Coins,
	delegator string,
	fund bool,
) {
	// Fund the delegator's account if needed.
	if fund {
		suite.FundAccount(ctx, delegator, amt)
	}

	// Delegate the coins to the operator.
	restakingMsgServer := restakingkeeper.NewMsgServer(suite.restakingKeeper)
	_, err := restakingMsgServer.DelegateOperator(ctx, restakingtypes.NewMsgDelegateOperator(
		operatorID,
		amt,
		delegator,
	))
	suite.Require().NoError(err)
}

// DelegateService delegates the given amount of coins to the service. If fund
// is true, the delegator's account is funded with the given amount of coins.
func (suite *KeeperTestSuite) DelegateService(ctx sdk.Context, serviceID uint32, amt sdk.Coins, delegator string, fund bool) {
	// Fund the delegator's account if needed.
	if fund {
		suite.FundAccount(ctx, delegator, amt)
	}

	// Delegate the coins to the service.
	restakingMsgServer := restakingkeeper.NewMsgServer(suite.restakingKeeper)
	_, err := restakingMsgServer.DelegateService(ctx, restakingtypes.NewMsgDelegateService(
		serviceID,
		amt,
		delegator,
	))
	suite.Require().NoError(err)
}

// DelegatePool delegates the given amount of coins to the pool. If fund is
// true, the delegator's account is funded with the given amount of coins.
func (suite *KeeperTestSuite) DelegatePool(ctx sdk.Context, amt sdk.Coin, delegator string, fund bool) {
	// Fund the delegator's account if needed.
	if fund {
		suite.FundAccount(ctx, delegator, sdk.NewCoins(amt))
	}

	// Delegate the coins to the pool.
	restakingMsgServer := restakingkeeper.NewMsgServer(suite.restakingKeeper)
	_, err := restakingMsgServer.DelegatePool(ctx, restakingtypes.NewMsgDelegatePool(amt, delegator))
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) allocateRewards(ctx sdk.Context, duration time.Duration) sdk.Context {
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(duration)).WithBlockHeight(ctx.BlockHeight() + 1)
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)
	return ctx
}

func (suite *KeeperTestSuite) setupSampleServiceAndOperator(ctx sdk.Context) (servicestypes.Service, operatorstypes.Operator) {
	// This helper method:
	// - registers $MILK, $INIT
	// - creates a service named "MilkyWay"
	// - creates a rewards plan with basic distribution types
	//   - it distributes 100 $MILK every day
	// - creates an operator named "MilkyWay Operator"
	//   - it has 10% commission rate
	//   - it joins the newly created service

	// Register $MILK and $INIT.
	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))
	suite.RegisterCurrency(ctx, "uinit", "INIT", 6, utils.MustParseDec("3"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Create an operator.
	operatorAdmin := testutil.TestAddress(10001)
	operator := suite.CreateOperator(ctx, "Operator", operatorAdmin.String())

	// Make the operator join the service and set its commission rate to 10%.
	suite.UpdateOperatorParams(ctx, operator.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	return service, operator
}
