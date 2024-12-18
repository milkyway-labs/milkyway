package v5_test

import (
	"testing"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/header"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	milkywayapp "github.com/milkyway-labs/milkyway/v6/app"
	"github.com/milkyway-labs/milkyway/v6/app/helpers"
	v4 "github.com/milkyway-labs/milkyway/v6/app/upgrades/v5"
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

func (suite *UpgradeTestSuite) TestUpgradeV4() {
	// Make sure the markets are set
	markets, err := suite.App.MarketMapKeeper.GetAllMarkets(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Len(markets, 63)

	// Make sure the currency pairs are set
	currencyPairs := suite.App.OracleKeeper.GetAllCurrencyPairs(suite.Ctx)
	suite.Require().Len(currencyPairs, 63)

	// Perform the upgrade
	suite.performUpgrade()

	// Make sure only the TIA/USD market is left
	markets, err = suite.App.MarketMapKeeper.GetAllMarkets(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Len(markets, 1)
	suite.Require().Equal("TIA/USD", markets["TIA/USD"].Ticker.String())

	// Make sure the only currency pair left is TIA/USD
	currencyPairs = suite.App.OracleKeeper.GetAllCurrencyPairs(suite.Ctx)
	suite.Require().Len(currencyPairs, 1)
	suite.Require().Equal("TIA/USD", currencyPairs[0].String())
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
