package keeper_test

import (
	"github.com/milkyway-labs/milkyway/utils"
)

func (s *KeeperTestSuite) TestGetAssetAndPrice() {
	s.RegisterCurrency("umilk", "MILK", 6, utils.MustParseDec("2"))

	_, price, err := s.keeper.GetAssetAndPrice(s.Ctx, "umilk")
	s.Require().NoError(err)
	s.Require().Equal(utils.MustParseDec("2"), price)

	_, price, err = s.keeper.GetAssetAndPrice(s.Ctx, "uinit")
	s.Require().NoError(err)
	s.Require().True(price.IsZero())
}

func (s *KeeperTestSuite) TestGetCoinValue() {
	s.RegisterCurrency("umilk", "MILK", 6, utils.MustParseDec("2"))
	s.RegisterCurrency("afoo", "FOO", 18, utils.MustParseDec("0.53"))

	value, err := s.keeper.GetCoinValue(s.Ctx, utils.MustParseCoin("10_000000umilk"))
	s.Require().NoError(err)
	s.Require().Equal(utils.MustParseDec("20"), value)

	value, err = s.keeper.GetCoinValue(s.Ctx, utils.MustParseCoin("10_000000uinit"))
	s.Require().NoError(err)
	s.Require().True(value.IsZero())

	value, err = s.keeper.GetCoinValue(s.Ctx, utils.MustParseCoin("1200000000000000afoo")) // 0.0012 $FOO
	s.Require().NoError(err)
	s.Require().Equal(utils.MustParseDec("0.000636"), value)
}

func (s *KeeperTestSuite) TestGetCoinsValue() {
	s.RegisterCurrency("umilk", "MILK", 6, utils.MustParseDec("2"))
	s.RegisterCurrency("uinit", "INIT", 6, utils.MustParseDec("3"))

	value, err := s.keeper.GetCoinsValue(s.Ctx, utils.MustParseCoins("10_000000umilk,5_500000uinit"))
	s.Require().NoError(err)
	s.Require().Equal(utils.MustParseDec("36.5"), value)

	value, err = s.keeper.GetCoinsValue(s.Ctx, utils.MustParseCoins("10_000000umilk,10_000000uatom"))
	s.Require().NoError(err)
	s.Require().Equal(utils.MustParseDec("20"), value)
}
