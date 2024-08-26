package keeper_test

import (
	"github.com/milkyway-labs/milkyway/utils"
)

func (suite *KeeperTestSuite) TestGetAssetAndPrice() {
	suite.RegisterCurrency("umilk", "MILK", 6, utils.MustParseDec("2"))

	_, price, err := suite.keeper.GetAssetAndPrice(suite.Ctx, "umilk")
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("2"), price)

	_, price, err = suite.keeper.GetAssetAndPrice(suite.Ctx, "uinit")
	suite.Require().NoError(err)
	suite.Require().True(price.IsZero())
}

func (suite *KeeperTestSuite) TestGetCoinValue() {
	suite.RegisterCurrency("umilk", "MILK", 6, utils.MustParseDec("2"))
	suite.RegisterCurrency("afoo", "FOO", 18, utils.MustParseDec("0.53"))

	value, err := suite.keeper.GetCoinValue(suite.Ctx, utils.MustParseCoin("10_000000umilk"))
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("20"), value)

	value, err = suite.keeper.GetCoinValue(suite.Ctx, utils.MustParseCoin("10_000000uinit"))
	suite.Require().NoError(err)
	suite.Require().True(value.IsZero())

	value, err = suite.keeper.GetCoinValue(suite.Ctx, utils.MustParseCoin("1200000000000000afoo")) // 0.0012 $FOO
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("0.000636"), value)
}

func (suite *KeeperTestSuite) TestGetCoinsValue() {
	suite.RegisterCurrency("umilk", "MILK", 6, utils.MustParseDec("2"))
	suite.RegisterCurrency("uinit", "INIT", 6, utils.MustParseDec("3"))

	value, err := suite.keeper.GetCoinsValue(suite.Ctx, utils.MustParseCoins("10_000000umilk,5_500000uinit"))
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("36.5"), value)

	value, err = suite.keeper.GetCoinsValue(suite.Ctx, utils.MustParseCoins("10_000000umilk,10_000000uatom"))
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("20"), value)
}
