package testutil

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"

	milkywayapp "github.com/milkyway-labs/milkyway/app"
	assetstypes "github.com/milkyway-labs/milkyway/x/assets/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	rewardskeeper "github.com/milkyway-labs/milkyway/x/rewards/keeper"
	rewardstypes "github.com/milkyway-labs/milkyway/x/rewards/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// KeeperTestSuite is a base test suite for the keeper tests.
type KeeperTestSuite struct {
	suite.Suite

	App *milkywayapp.MilkyWayApp
	Ctx sdk.Context
}

// SetupTest creates a new MilkyWayApp and context for the test.
func (suite *KeeperTestSuite) SetupTest() {
	suite.App = milkywayapp.Setup(suite.T(), false)
	suite.Ctx = suite.App.NewContextLegacy(false, cmtproto.Header{
		Height: 1,
		Time:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	})
}

// FundAccount adds the given amount of coins to the account with the given address
func (suite *KeeperTestSuite) FundAccount(ctx sdk.Context, addr string, amt sdk.Coins) {
	// Mint the coins
	moduleAcc := suite.App.AccountKeeper.GetModuleAccount(ctx, authtypes.Minter)

	err := suite.App.BankKeeper.MintCoins(ctx, moduleAcc.GetName(), amt)
	suite.Require().NoError(err)

	// Send the amount to the user
	accAddr, err := sdk.AccAddressFromBech32(addr)
	suite.Require().NoError(err)

	err = suite.App.BankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAcc.GetName(), accAddr, amt)
	suite.Require().NoError(err)
}

// RegisterCurrency registers a currency with the given denomination, ticker
// and price. RegisterCurrency creates a market for the currency if not exists.
func (suite *KeeperTestSuite) RegisterCurrency(ctx sdk.Context, denom string, ticker string, exponent uint32, price math.LegacyDec) {
	// Create the market only if it doesn't exist.
	mmTicker := marketmaptypes.NewTicker(ticker, rewardstypes.USDTicker, math.LegacyPrecision, 0, true)
	hasMarket, err := suite.App.MarketMapKeeper.HasMarket(ctx, mmTicker.String())
	suite.Require().NoError(err)

	if !hasMarket {
		err = suite.App.MarketMapKeeper.CreateMarket(ctx, marketmaptypes.Market{Ticker: mmTicker})
		suite.Require().NoError(err)
	}

	// Set the price for the currency pair.
	err = suite.App.OracleKeeper.SetPriceForCurrencyPair(
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
	err = suite.App.AssetsKeeper.SetAsset(ctx, assetstypes.NewAsset(denom, ticker, exponent))
	suite.Require().NoError(err)
}

// CreateService creates an example service with the given service name and
// admin address. The service description and URLs related to the service are
// randomly chosen. The service is also activated.
func (suite *KeeperTestSuite) CreateService(ctx sdk.Context, name string, admin string) servicestypes.Service {
	servicesMsgServer := serviceskeeper.NewMsgServer(suite.App.ServicesKeeper)
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

	service, found := suite.App.ServicesKeeper.GetService(ctx, resp.NewServiceID)
	suite.Require().True(found, "service must be found")
	return service
}

// CreateOperator creates an example operator with the given operator name and
// admin address. The operator description and URLs related to the operator are
// randomly chosen.
func (suite *KeeperTestSuite) CreateOperator(ctx sdk.Context, name string, admin string) operatorstypes.Operator {
	// Register the operator
	operatorsMsgServer := operatorskeeper.NewMsgServer(suite.App.OperatorsKeeper)
	resp, err := operatorsMsgServer.RegisterOperator(ctx, operatorstypes.NewMsgRegisterOperator(
		name,
		"https://example.com",
		"https://example.com/picture.png",
		admin,
	))
	suite.Require().NoError(err)

	// Make sure the operator is found
	operator, found := suite.App.OperatorsKeeper.GetOperator(ctx, resp.NewOperatorID)
	suite.Require().True(found, "operator must be found")
	return operator
}

// UpdateOperatorParams updates the operator's params.
func (suite *KeeperTestSuite) UpdateOperatorParams(
	ctx sdk.Context,
	operatorID uint32,
	commissionRate math.LegacyDec,
	joinedServicesIDs []uint32,
) {
	// Make sure the operator is found
	_, found := suite.App.OperatorsKeeper.GetOperator(ctx, operatorID)
	suite.Require().True(found, "operator must be found")

	// Sets the operator commission rate
	err := suite.App.OperatorsKeeper.SaveOperatorParams(ctx, operatorID,
		operatorstypes.NewOperatorParams(commissionRate))
	suite.Require().NoError(err)

	// Make the operator join the service.
	joinedServices, err := suite.App.RestakingKeeper.GetOperatorJoinedServices(ctx, operatorID)
	suite.Require().NoError(err)
	for _, serviceID := range joinedServicesIDs {
		err = joinedServices.Add(serviceID)
		suite.Require().NoError(err)
	}
	suite.App.RestakingKeeper.SaveOperatorJoinedServices(ctx, operatorID, joinedServices)
}

// UpdateServiceParams updates the service's params.
func (suite *KeeperTestSuite) UpdateServiceParams(
	ctx sdk.Context,
	serviceID uint32,
	slashFraction math.LegacyDec,
	whitelistedPoolsIDs []uint32,
	whitelistedOperatorsIDs []uint32,
) {
	// Make sure the service is found
	service, found := suite.App.ServicesKeeper.GetService(ctx, serviceID)
	suite.Require().True(found, "service must be found")

	servicesMsgServer := serviceskeeper.NewMsgServer(suite.App.ServicesKeeper)
	serviceParams := servicestypes.NewServiceParams(slashFraction)
	_, err := servicesMsgServer.SetServiceParams(ctx, servicestypes.NewMsgSetServiceParams(
		serviceID,
		serviceParams,
		service.Admin,
	))
	suite.Require().NoError(err)

	// Make the operator join the service and set its commission rate to 10%.
	restakingMsgServer := restakingkeeper.NewMsgServer(suite.App.RestakingKeeper)
	_, err = restakingMsgServer.UpdateServiceParams(ctx, restakingtypes.NewMsgUpdateServiceParams(
		service.ID,
		restakingtypes.NewServiceParams(whitelistedPoolsIDs, whitelistedOperatorsIDs),
		service.Admin,
	))
	suite.Require().NoError(err)
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
	service, found := suite.App.ServicesKeeper.GetService(ctx, serviceID)
	suite.Require().True(found, "service must be found")

	rewardsMsgServer := rewardskeeper.NewMsgServer(suite.App.RewardsKeeper)
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
	plan, err := suite.App.RewardsKeeper.GetRewardsPlan(ctx, resp.NewRewardsPlanID)
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
	restakingMsgServer := restakingkeeper.NewMsgServer(suite.App.RestakingKeeper)
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
	restakingMsgServer := restakingkeeper.NewMsgServer(suite.App.RestakingKeeper)
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
	restakingMsgServer := restakingkeeper.NewMsgServer(suite.App.RestakingKeeper)
	_, err := restakingMsgServer.DelegatePool(ctx, restakingtypes.NewMsgDelegatePool(amt, delegator))
	suite.Require().NoError(err)
}
