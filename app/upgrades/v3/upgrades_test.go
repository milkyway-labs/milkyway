package v3_test

import (
	"testing"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/header"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	milkywayapp "github.com/milkyway-labs/milkyway/v3/app"
	"github.com/milkyway-labs/milkyway/v3/app/helpers"
	v3 "github.com/milkyway-labs/milkyway/v3/app/upgrades/v3"
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

func (suite *UpgradeTestSuite) TestUpgradeV3() {
	operatorsParams, err := suite.App.OperatorsKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	operatorsParams.OperatorRegistrationFee = sdk.NewCoins(
		sdk.NewInt64Coin("ibc/16065EE5282C5217685C8F084FC44864C25C706AC37356B0D62811D50B96920F", 1000000),
		sdk.NewInt64Coin("ibc/6C349F0EB135C5FA99301758F35B87DB88403D690E5E314AB080401FEE4066E5", 1000000),
		sdk.NewInt64Coin("ibc/84FBEC4BBB48BD7CC534ED7518F339CCF6C45529DC00C7BFB8605C9EE7D68AFC", 1000000),
		sdk.NewInt64Coin("ibc/8D4FC51F696E03711B9B37A5787FB89BD2DDBAF788813478B002D552A12F9157", 1000000),
		sdk.NewInt64Coin("ibc/9ACD338BC3B488E0F50A54DE9A844C8326AF0739D917922A9CE04D42AD66017E", 1000000),
	)
	err = suite.App.OperatorsKeeper.SetParams(suite.Ctx, operatorsParams)
	suite.Require().NoError(err)

	servicesParams, err := suite.App.ServicesKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	servicesParams.ServiceRegistrationFee = sdk.NewCoins(
		sdk.NewInt64Coin("ibc/16065EE5282C5217685C8F084FC44864C25C706AC37356B0D62811D50B96920F", 1000000),
		sdk.NewInt64Coin("ibc/6C349F0EB135C5FA99301758F35B87DB88403D690E5E314AB080401FEE4066E5", 1000000),
		sdk.NewInt64Coin("ibc/84FBEC4BBB48BD7CC534ED7518F339CCF6C45529DC00C7BFB8605C9EE7D68AFC", 1000000),
		sdk.NewInt64Coin("ibc/8D4FC51F696E03711B9B37A5787FB89BD2DDBAF788813478B002D552A12F9157", 1000000),
		sdk.NewInt64Coin("ibc/9ACD338BC3B488E0F50A54DE9A844C8326AF0739D917922A9CE04D42AD66017E", 1000000),
	)
	err = suite.App.ServicesKeeper.SetParams(suite.Ctx, servicesParams)
	suite.Require().NoError(err)

	restakingParams, err := suite.App.RestakingKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	restakingParams.AllowedDenoms = []string{
		"ibc/16065EE5282C5217685C8F084FC44864C25C706AC37356B0D62811D50B96920F",
		"ibc/6C349F0EB135C5FA99301758F35B87DB88403D690E5E314AB080401FEE4066E5",
		"ibc/84FBEC4BBB48BD7CC534ED7518F339CCF6C45529DC00C7BFB8605C9EE7D68AFC",
		"ibc/8D4FC51F696E03711B9B37A5787FB89BD2DDBAF788813478B002D552A12F9157",
		"ibc/9ACD338BC3B488E0F50A54DE9A844C8326AF0739D917922A9CE04D42AD66017E",
		"locked/ibc/F1183DB3D428313A6FD329DF18219F9D6B83257D07D292EA9EC1D877E89EC2B0",
	}
	err = suite.App.RestakingKeeper.SetParams(suite.Ctx, restakingParams)
	suite.Require().NoError(err)

	rewardsParams, err := suite.App.RewardsKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	rewardsParams.RewardsPlanCreationFee = sdk.NewCoins(
		sdk.NewInt64Coin("ibc/16065EE5282C5217685C8F084FC44864C25C706AC37356B0D62811D50B96920F", 1000000),
		sdk.NewInt64Coin("ibc/6C349F0EB135C5FA99301758F35B87DB88403D690E5E314AB080401FEE4066E5", 1000000),
		sdk.NewInt64Coin("ibc/84FBEC4BBB48BD7CC534ED7518F339CCF6C45529DC00C7BFB8605C9EE7D68AFC", 1000000),
		sdk.NewInt64Coin("ibc/8D4FC51F696E03711B9B37A5787FB89BD2DDBAF788813478B002D552A12F9157", 1000000),
		sdk.NewInt64Coin("ibc/9ACD338BC3B488E0F50A54DE9A844C8326AF0739D917922A9CE04D42AD66017E", 1000000),
	)
	err = suite.App.RewardsKeeper.SetParams(suite.Ctx, rewardsParams)
	suite.Require().NoError(err)

	suite.performUpgrade()

	operatorsParams, err = suite.App.OperatorsKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins(
		sdk.NewInt64Coin("ibc/16065EE5282C5217685C8F084FC44864C25C706AC37356B0D62811D50B96920F", 1000000),
		sdk.NewInt64Coin("ibc/6C349F0EB135C5FA99301758F35B87DB88403D690E5E314AB080401FEE4066E5", 1000000),
		sdk.NewInt64Coin("ibc/8D4FC51F696E03711B9B37A5787FB89BD2DDBAF788813478B002D552A12F9157", 1000000),
	), operatorsParams.OperatorRegistrationFee)

	servicesParams, err = suite.App.ServicesKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins(
		sdk.NewInt64Coin("ibc/16065EE5282C5217685C8F084FC44864C25C706AC37356B0D62811D50B96920F", 1000000),
		sdk.NewInt64Coin("ibc/6C349F0EB135C5FA99301758F35B87DB88403D690E5E314AB080401FEE4066E5", 1000000),
		sdk.NewInt64Coin("ibc/8D4FC51F696E03711B9B37A5787FB89BD2DDBAF788813478B002D552A12F9157", 1000000),
	), servicesParams.ServiceRegistrationFee)

	restakingParams, err = suite.App.RestakingKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal([]string{
		"ibc/16065EE5282C5217685C8F084FC44864C25C706AC37356B0D62811D50B96920F",
		"ibc/6C349F0EB135C5FA99301758F35B87DB88403D690E5E314AB080401FEE4066E5",
		"ibc/8D4FC51F696E03711B9B37A5787FB89BD2DDBAF788813478B002D552A12F9157",
		"locked/ibc/F1183DB3D428313A6FD329DF18219F9D6B83257D07D292EA9EC1D877E89EC2B0",
	}, restakingParams.AllowedDenoms)

	rewardsParams, err = suite.App.RewardsKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.NewCoins(
		sdk.NewInt64Coin("ibc/16065EE5282C5217685C8F084FC44864C25C706AC37356B0D62811D50B96920F", 1000000),
		sdk.NewInt64Coin("ibc/6C349F0EB135C5FA99301758F35B87DB88403D690E5E314AB080401FEE4066E5", 1000000),
		sdk.NewInt64Coin("ibc/8D4FC51F696E03711B9B37A5787FB89BD2DDBAF788813478B002D552A12F9157", 1000000),
	), rewardsParams.RewardsPlanCreationFee)

	liquidVestingParams, err := suite.App.LiquidVestingKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal([]string{
		"celestia1nyk8qsfkrplvzex5yc4l5kghdr60nj7mutnw6z",
		"celestia1vre9kvtw6w9lxkl3620zpzs5lpczhvjc0s6r2a",
	}, liquidVestingParams.TrustedDelegates)

	slashingParams, err := suite.App.SlashingKeeper.GetParams(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(time.Hour, slashingParams.DowntimeJailDuration)
}

func (suite *UpgradeTestSuite) performUpgrade() {
	upgradeHeight := suite.Ctx.HeaderInfo().Height + 1
	plan := upgradetypes.Plan{Name: v3.UpgradeName, Height: upgradeHeight}
	err := suite.App.UpgradeKeeper.ScheduleUpgrade(suite.Ctx, plan)
	suite.Require().NoError(err)
	_, err = suite.App.UpgradeKeeper.GetUpgradePlan(suite.Ctx)
	suite.Require().NoError(err)

	suite.Ctx = suite.Ctx.WithHeaderInfo(header.Info{Height: upgradeHeight})
	// PreBlocker triggers the upgrade
	_, err = suite.UpgradeModule.PreBlock(suite.Ctx)
	suite.Require().NoError(err)
}
