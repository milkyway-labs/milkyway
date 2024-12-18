package v6_test

import (
	"testing"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/header"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/suite"

	milkywayapp "github.com/milkyway-labs/milkyway/v6/app"
	"github.com/milkyway-labs/milkyway/v6/app/helpers"
	v4 "github.com/milkyway-labs/milkyway/v6/app/upgrades/v6"
	"github.com/milkyway-labs/milkyway/v6/utils"
	liquidvestingtypes "github.com/milkyway-labs/milkyway/v6/x/liquidvesting/types"
)

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

type UpgradeTestSuite struct {
	suite.Suite

	App           *milkywayapp.MilkyWayApp
	Ctx           sdk.Context
	UpgradeModule appmodule.HasPreBlocker
}

func (suite *UpgradeTestSuite) SetupTest() {
	suite.App = helpers.Setup(suite.T())
	suite.Ctx = suite.App.NewUncachedContext(true, tmproto.Header{})
	suite.Ctx = suite.Ctx.WithHeaderInfo(header.Info{Height: 1})
	suite.UpgradeModule = upgrade.NewAppModule(suite.App.UpgradeKeeper, suite.App.AccountKeeper.AddressCodec())
}

func (suite *UpgradeTestSuite) TestUpgradeV6() {
	// Make sure the markets are set
	markets, err := suite.App.MarketMapKeeper.GetAllMarkets(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Len(markets, 63)

	// Make sure the currency pairs are set
	currencyPairs := suite.App.OracleKeeper.GetAllCurrencyPairs(suite.Ctx)
	suite.Require().Len(currencyPairs, 63)

	// In testnet v5, we have liquid vesting module account set as a BaseAccount.
	liquidVestingModuleAddr := suite.App.AccountKeeper.GetModuleAddress(liquidvestingtypes.ModuleName)
	acc := suite.App.AccountKeeper.GetAccount(suite.Ctx, liquidVestingModuleAddr)
	suite.App.AccountKeeper.RemoveAccount(suite.Ctx, acc)
	err = suite.App.BankKeeper.MintCoins(suite.Ctx, minttypes.ModuleName, utils.MustParseCoins("1utia"))
	suite.Require().NoError(err)
	err = suite.App.BankKeeper.SendCoinsFromModuleToAccount(
		suite.Ctx,
		minttypes.ModuleName,
		liquidVestingModuleAddr,
		utils.MustParseCoins("1utia"),
	)
	suite.Require().NoError(err)

	// Perform the upgrade
	suite.performUpgrade()

	// Make sure all the currency pairs are deleted
	currencyPairs = suite.App.OracleKeeper.GetAllCurrencyPairs(suite.Ctx)
	suite.Require().Empty(currencyPairs)

	acc = suite.App.AccountKeeper.GetAccount(suite.Ctx, liquidVestingModuleAddr)
	suite.Require().NotNil(acc)
	_, ok := acc.(sdk.ModuleAccountI)
	suite.Require().True(ok)
}

func (suite *UpgradeTestSuite) performUpgrade() {
	upgradeHeight := suite.Ctx.HeaderInfo().Height + 1

	// Schedule the upgrade
	err := suite.App.UpgradeKeeper.ScheduleUpgrade(suite.Ctx, upgradetypes.Plan{Name: v4.UpgradeName, Height: upgradeHeight})
	suite.Require().NoError(err)

	// Make sure the upgrade plan is set
	_, err = suite.App.UpgradeKeeper.GetUpgradePlan(suite.Ctx)
	suite.Require().NoError(err)

	// Fast-forward to the upgrade height
	suite.Ctx = suite.Ctx.WithHeaderInfo(header.Info{Height: upgradeHeight})

	// Run PreBlocker to trigger the upgrade
	_, err = suite.UpgradeModule.PreBlock(suite.Ctx)
	suite.Require().NoError(err)
}
