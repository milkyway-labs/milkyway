package keeper_test

import (
	"github.com/milkyway-labs/milkyway/v6/utils"
)

func (suite *KeeperTestSuite) TestGetAssetAndPrice() {
	// Cache the context to avoid issues
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))

	_, price, err := suite.k.GetAssetAndPrice(ctx, "umilk")
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("2"), price)

	_, price, err = suite.k.GetAssetAndPrice(ctx, "uinit")
	suite.Require().NoError(err)
	suite.Require().True(price.IsZero())
}

func (suite *KeeperTestSuite) TestGetCoinValue() {
	// Cache the context to avoid issues
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))
	suite.RegisterCurrency(ctx, "afoo", "FOO", 18, utils.MustParseDec("0.53"))

	value, err := suite.k.GetCoinValue(ctx, utils.MustParseCoin("10_000000umilk"))
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("20"), value)

	value, err = suite.k.GetCoinValue(ctx, utils.MustParseCoin("10_000000uinit"))
	suite.Require().NoError(err)
	suite.Require().True(value.IsZero())

	value, err = suite.k.GetCoinValue(ctx, utils.MustParseCoin("1200000000000000afoo")) // 0.0012 $FOO
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("0.000636"), value)
}

func (suite *KeeperTestSuite) TestGetCoinsValue() {
	// Cache the context to avoid issues
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))
	suite.RegisterCurrency(ctx, "uinit", "INIT", 6, utils.MustParseDec("3"))

	value, err := suite.k.GetCoinsValue(ctx, utils.MustParseCoins("10_000000umilk,5_500000uinit"))
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("36.5"), value)

	value, err = suite.k.GetCoinsValue(ctx, utils.MustParseCoins("10_000000umilk,10_000000uatom"))
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDec("20"), value)
}
