package keeper_test

import (
	"testing"
	"time"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"

	milkywayapp "github.com/milkyway-labs/milkyway/v7/app"
	"github.com/milkyway-labs/milkyway/v7/x/tokenfactory/keeper"
	"github.com/milkyway-labs/milkyway/v7/x/tokenfactory/testutils"
	"github.com/milkyway-labs/milkyway/v7/x/tokenfactory/types"
)

type KeeperTestSuite struct {
	suite.Suite

	App      *milkywayapp.MilkyWayApp
	Ctx      sdk.Context
	TestAccs []sdk.AccAddress

	queryClient    types.QueryClient
	msgServer      types.MsgServer
	contractKeeper wasmtypes.ContractOpsKeeper
	bankMsgServer  banktypes.MsgServer

	// defaultDenom is on the suite, as it depends on the creator test address.
	defaultDenom string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type SudoAuthorizationPolicy struct{}

func (p SudoAuthorizationPolicy) CanCreateCode(chainAccesscoConfig wasmtypes.ChainAccessConfigs, actor sdk.AccAddress, config wasmtypes.AccessConfig) bool {
	return true
}

func (p SudoAuthorizationPolicy) CanInstantiateContract(config wasmtypes.AccessConfig, actor sdk.AccAddress) bool {
	return true
}

func (p SudoAuthorizationPolicy) CanModifyContract(admin, actor sdk.AccAddress) bool {
	return true
}

func (p SudoAuthorizationPolicy) CanModifyCodeAccessConfig(creator, actor sdk.AccAddress, isSubset bool) bool {
	return true
}

func (s *KeeperTestSuite) SetupTest() {
	s.App = testutils.Setup(false, s.T().TempDir())
	s.Ctx = s.App.BaseApp.NewContextLegacy(false, tmproto.Header{
		Height: 1,
		Time:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	s.TestAccs = sims.CreateRandomAccounts(3)

	// Fund every TestAcc with two denoms, one of which is the denom creation fee
	s.contractKeeper = wasmkeeper.NewGovPermissionKeeper(s.App.WasmKeeper)
	s.queryClient = types.NewQueryClient(&baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	})
	s.msgServer = keeper.NewMsgServerImpl(*s.App.TokenFactoryKeeper)
	s.bankMsgServer = bankkeeper.NewMsgServerImpl(s.App.BankKeeper)
}

func (s *KeeperTestSuite) SetupTestForInitGenesis() {
	// Setting to True, leads to init genesis not running
	s.App = testutils.Setup(true, s.T().TempDir())
	s.Ctx = s.App.BaseApp.NewContextLegacy(true, tmproto.Header{})
}

func (s *KeeperTestSuite) CreateDefaultDenom() {
	res, _ := s.msgServer.CreateDenom(s.Ctx, types.NewMsgCreateDenom(s.TestAccs[0].String(), "bitcoin"))
	s.defaultDenom = res.GetNewTokenDenom()
}

// FundAcc funds target address with specified amount.
func (s *KeeperTestSuite) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := banktestutil.FundAccount(s.Ctx, s.App.BankKeeper, acc, amounts)
	s.Require().NoError(err)
}

// FundModuleAcc funds target modules with specified amount.
func (s *KeeperTestSuite) FundModuleAcc(moduleName string, amounts sdk.Coins) {
	err := banktestutil.FundModuleAccount(s.Ctx, s.App.BankKeeper, moduleName, amounts)
	s.Require().NoError(err)
}

// AssertEventEmitted asserts that ctx's event manager has emitted the given number of events
// of the given type.
func (s *KeeperTestSuite) AssertEventEmitted(ctx sdk.Context, eventTypeExpected string, numEventsExpected int) {
	allEvents := ctx.EventManager().Events()
	// filter out other events
	actualEvents := make([]sdk.Event, 0)
	for _, event := range allEvents {
		if event.Type == eventTypeExpected {
			actualEvents = append(actualEvents, event)
		}
	}
	s.Require().Equal(numEventsExpected, len(actualEvents))
}

func (s *KeeperTestSuite) TestCreateModuleAccount() {
	app := s.App

	// setup new next account number
	nextAccountNumber := app.AccountKeeper.NextAccountNumber(s.Ctx)

	// remove module account
	tokenfactoryModuleAccount := app.AccountKeeper.GetAccount(s.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	app.AccountKeeper.RemoveAccount(s.Ctx, tokenfactoryModuleAccount)

	// ensure module account was removed
	s.Ctx = app.BaseApp.NewContextLegacy(false, tmproto.Header{})
	tokenfactoryModuleAccount = app.AccountKeeper.GetAccount(s.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().Nil(tokenfactoryModuleAccount)

	// create module account
	app.TokenFactoryKeeper.CreateModuleAccount(s.Ctx)

	// check that the module account is now initialized
	tokenfactoryModuleAccount = app.AccountKeeper.GetAccount(s.Ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	s.Require().NotNil(tokenfactoryModuleAccount)

	// check that the account number of the module account is now initialized correctly
	s.Require().Equal(nextAccountNumber+1, tokenfactoryModuleAccount.GetAccountNumber())
}
