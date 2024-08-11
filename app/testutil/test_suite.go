package testutil

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"

	milkywayapp "github.com/milkyway-labs/milkyway/app"
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	restakingkeeper "github.com/milkyway-labs/milkyway/x/restaking/keeper"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	rewardskeeper "github.com/milkyway-labs/milkyway/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
	tickerstypes "github.com/milkyway-labs/milkyway/x/tickers/types"
)

// KeeperTestSuite is a base test suite for the keeper tests.
type KeeperTestSuite struct {
	suite.Suite

	App *milkywayapp.MilkyWayApp
	Ctx sdk.Context
}

// SetupTest creates a new MilkyWayApp and context for the test.
func (s *KeeperTestSuite) SetupTest() {
	s.App = milkywayapp.Setup(false)
	s.Ctx = s.App.NewContextLegacy(false, cmtproto.Header{
		Height: 1,
		Time:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	})
}

// FundAccount adds the given amount of coins to the account with the given address
func (s *KeeperTestSuite) FundAccount(addr string, amt sdk.Coins) {
	// Mint the coins
	moduleAcc := s.App.AccountKeeper.GetModuleAccount(s.Ctx, authtypes.Minter)

	err := s.App.BankKeeper.MintCoins(s.Ctx, moduleAcc.GetName(), amt)
	s.Require().NoError(err)

	// Get the amount to the user
	accAddr, err := sdk.AccAddressFromBech32(addr)
	s.Require().NoError(err)
	err = s.App.BankKeeper.SendCoinsFromModuleToAccount(s.Ctx, moduleAcc.GetName(), accAddr, amt)
	s.Require().NoError(err)
}

// RegisterCurrency registers a currency with the given denomination, ticker
// and price. RegisterCurrency creates a market for the currency if not exists.
func (s *KeeperTestSuite) RegisterCurrency(denom, ticker string, exponent uint32, price math.LegacyDec) {
	// Create market only if it doesn't exist.
	mmTicker := marketmaptypes.NewTicker(ticker, types.USDTicker, math.LegacyPrecision, 0, true)
	hasMarket, err := s.App.MarketMapKeeper.HasMarket(s.Ctx, mmTicker.String())
	s.Require().NoError(err)
	if !hasMarket {
		err = s.App.MarketMapKeeper.CreateMarket(s.Ctx, marketmaptypes.Market{Ticker: mmTicker})
		s.Require().NoError(err)
	}
	err = s.App.OracleKeeper.SetPriceForCurrencyPair(
		s.Ctx, slinkytypes.NewCurrencyPair(ticker, types.USDTicker),
		oracletypes.QuotePrice{
			Price:          math.NewIntFromBigInt(price.BigInt()),
			BlockTimestamp: s.Ctx.BlockTime(),
			BlockHeight:    uint64(s.Ctx.BlockHeight()),
		})
	s.Require().NoError(err)
	err = s.App.TickersKeeper.SetAsset(s.Ctx, tickerstypes.NewAsset(denom, ticker, exponent))
	s.Require().NoError(err)
}

// CreateService creates an example service with the given service name and
// admin address. The service description and URLs related to the service are
// randomly chosen. The service is also activated.
func (s *KeeperTestSuite) CreateService(name, admin string) servicestypes.Service {
	servicesMsgServer := serviceskeeper.NewMsgServer(s.App.ServicesKeeper)
	resp, err := servicesMsgServer.CreateService(s.Ctx, servicestypes.NewMsgCreateService(
		name,
		fmt.Sprintf("%s AVS", name),
		"https://example.com",
		"https://example.com/picture.png",
		admin,
	))
	s.Require().NoError(err)
	// Also activate the service.
	_, err = servicesMsgServer.ActivateService(s.Ctx, servicestypes.NewMsgActivateService(resp.NewServiceID, admin))
	s.Require().NoError(err)
	service, found := s.App.ServicesKeeper.GetService(s.Ctx, resp.NewServiceID)
	s.Require().True(found, "service must be found")
	return service
}

// CreateOperator creates an example operator with the given operator name and
// admin address. The operator description and URLs related to the operator are
// randomly chosen.
func (s *KeeperTestSuite) CreateOperator(name, admin string) operatorstypes.Operator {
	operatorsMsgServer := operatorskeeper.NewMsgServer(s.App.OperatorsKeeper)
	resp, err := operatorsMsgServer.RegisterOperator(s.Ctx, operatorstypes.NewMsgRegisterOperator(
		name,
		"https://example.com",
		"https://example.com/picture.png",
		admin,
	))
	s.Require().NoError(err)
	operator, found := s.App.OperatorsKeeper.GetOperator(s.Ctx, resp.NewOperatorID)
	s.Require().True(found, "operator must be found")
	return operator
}

// UpdateOperatorParams updates the operator's params.
func (s *KeeperTestSuite) UpdateOperatorParams(
	operatorID uint32, commissionRate math.LegacyDec, joinedServicesIDs []uint32) {
	operator, found := s.App.OperatorsKeeper.GetOperator(s.Ctx, operatorID)
	s.Require().True(found, "operator must be found")
	// Make the operator join the service and set its commission rate to 10%.
	restakingMsgServer := restakingkeeper.NewMsgServer(s.App.RestakingKeeper)
	_, err := restakingMsgServer.UpdateOperatorParams(
		s.Ctx, restakingtypes.NewMsgUpdateOperatorParams(
			operator.ID,
			restakingtypes.NewOperatorParams(commissionRate, joinedServicesIDs),
			operator.Admin))
	s.Require().NoError(err)
}

// UpdateServiceParams updates the service's params.
func (s *KeeperTestSuite) UpdateServiceParams(
	serviceID uint32, slashFraction math.LegacyDec, whitelistedPoolsIDs, whitelistedOperatorsIDs []uint32) {
	service, found := s.App.ServicesKeeper.GetService(s.Ctx, serviceID)
	s.Require().True(found, "service must be found")
	// Make the operator join the service and set its commission rate to 10%.
	restakingMsgServer := restakingkeeper.NewMsgServer(s.App.RestakingKeeper)
	_, err := restakingMsgServer.UpdateServiceParams(
		s.Ctx, restakingtypes.NewMsgUpdateServiceParams(
			service.ID,
			restakingtypes.NewServiceParams(slashFraction, whitelistedPoolsIDs, whitelistedOperatorsIDs),
			service.Admin))
	s.Require().NoError(err)
}

// CreateRewardsPlan creates a rewards plan with the given parameters.
// The plan's name is chosen randomly. The plan is also funded with the given
// initial rewards.
func (s *KeeperTestSuite) CreateRewardsPlan(
	serviceID uint32, amtPerDay sdk.Coins, startTime, endTime time.Time,
	poolsDistr types.Distribution, operatorsDistr types.Distribution,
	usersDistr types.UsersDistribution, initialRewards sdk.Coins,
) types.RewardsPlan {
	service, found := s.App.ServicesKeeper.GetService(s.Ctx, serviceID)
	s.Require().True(found, "service must be found")
	rewardsMsgServer := rewardskeeper.NewMsgServer(s.App.RewardsKeeper)
	resp, err := rewardsMsgServer.CreateRewardsPlan(s.Ctx, types.NewMsgCreateRewardsPlan(
		service.Admin, "Rewards Plan", serviceID, amtPerDay, startTime, endTime,
		poolsDistr, operatorsDistr, usersDistr))
	s.Require().NoError(err)
	// Return the newly created plan.
	plan, err := s.App.RewardsKeeper.GetRewardsPlan(s.Ctx, resp.NewRewardsPlanID)
	s.Require().NoError(err)
	if initialRewards.IsAllPositive() {
		s.FundAccount(plan.RewardsPool, initialRewards)
	}
	return plan
}

// CreateBasicRewardsPlan creates a rewards plan with basic distribution for
// all restaking entities. Weights among the entities are set to 0, which means
// that the rewards are distributed based on the total delegation values among
// all entities.
func (s *KeeperTestSuite) CreateBasicRewardsPlan(
	serviceID uint32, amtPerDay sdk.Coins, startTime, endTime time.Time, initialRewards sdk.Coins,
) types.RewardsPlan {
	return s.CreateRewardsPlan(
		serviceID, amtPerDay, startTime, endTime,
		types.NewBasicPoolsDistribution(0), types.NewBasicOperatorsDistribution(0), types.NewBasicUsersDistribution(0),
		initialRewards)
}

// DelegateOperator delegates the given amount of coins to the operator. If fund
// is true, the delegator's account is funded with the given amount of coins.
func (s *KeeperTestSuite) DelegateOperator(operatorID uint32, amt sdk.Coins, delegator string, fund bool) {
	if fund {
		s.FundAccount(delegator, amt)
	}
	restakingMsgServer := restakingkeeper.NewMsgServer(s.App.RestakingKeeper)
	_, err := restakingMsgServer.DelegateOperator(
		s.Ctx, restakingtypes.NewMsgDelegateOperator(operatorID, amt, delegator))
	s.Require().NoError(err)
}

// DelegateService delegates the given amount of coins to the service. If fund
// is true, the delegator's account is funded with the given amount of coins.
func (s *KeeperTestSuite) DelegateService(serviceID uint32, amt sdk.Coins, delegator string, fund bool) {
	if fund {
		s.FundAccount(delegator, amt)
	}
	restakingMsgServer := restakingkeeper.NewMsgServer(s.App.RestakingKeeper)
	_, err := restakingMsgServer.DelegateService(s.Ctx, restakingtypes.NewMsgDelegateService(serviceID, amt, delegator))
	s.Require().NoError(err)
}

// DelegatePool delegates the given amount of coins to the pool. If fund is
// true, the delegator's account is funded with the given amount of coins.
func (s *KeeperTestSuite) DelegatePool(amt sdk.Coin, delegator string, fund bool) {
	if fund {
		s.FundAccount(delegator, sdk.NewCoins(amt))
	}
	restakingMsgServer := restakingkeeper.NewMsgServer(s.App.RestakingKeeper)
	_, err := restakingMsgServer.DelegatePool(s.Ctx, restakingtypes.NewMsgDelegatePool(amt, delegator))
	s.Require().NoError(err)
}
