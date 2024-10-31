package v122_test

import (
	"testing"

	"cosmossdk.io/core/header"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/app/testutil"
	v122 "github.com/milkyway-labs/milkyway/app/upgrades/v122"
	stakeibctypes "github.com/milkyway-labs/milkyway/x/stakeibc/types"
)

type KeeperTestSuite struct {
	testutil.KeeperTestSuite
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.KeeperTestSuite.SetupTest()
}

func (suite *KeeperTestSuite) TestUpgradeV122() {
	suite.Ctx = suite.Ctx.WithHeaderInfo(header.Info{Height: 1})

	// Set the host zone.
	suite.App.StakeIBCKeeper.SetHostZone(suite.Ctx, stakeibctypes.HostZone{
		ChainId:           "initiation-2",
		Bech32Prefix:      "init",
		ConnectionId:      "connection-0",
		TransferChannelId: "channel-0",
		IbcDenom:          "ibc/37A3FB4FED4CA04ED6D9E5DA36C6D27248645F0E22F585576A1488B8A89C5A50",
		HostDenom:         "uinit",
		UnbondingPeriod:   7,
	})

	upgradeHeight := suite.Ctx.HeaderInfo().Height + 1

	upgrade := v122.NewUpgrade(suite.App.ModuleManager, suite.App.Configurator(), suite.App.StakeIBCKeeper)
	plan := upgradetypes.Plan{Name: upgrade.Name(), Height: upgradeHeight}
	err := suite.App.UpgradeKeeper.ScheduleUpgrade(suite.Ctx, plan)
	suite.Require().NoError(err)

	_, err = suite.App.UpgradeKeeper.GetUpgradePlan(suite.Ctx)
	suite.Require().NoError(err)

	suite.Ctx = suite.Ctx.WithHeaderInfo(header.Info{Height: upgradeHeight})
	// PreBlocker triggers the upgrade
	_, err = suite.App.PreBlocker(suite.Ctx, nil)
	suite.Require().NoError(err)

	// Check if the unbonding period has been updated to 21 days.
	hostZone, found := suite.App.StakeIBCKeeper.GetHostZone(suite.Ctx, "initiation-2")
	suite.Require().True(found)
	suite.Require().EqualValues(21, hostZone.UnbondingPeriod)
}
