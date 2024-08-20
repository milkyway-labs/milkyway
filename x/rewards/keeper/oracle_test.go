package keeper_test

import (
	"github.com/milkyway-labs/milkyway/utils"
)

func (s *KeeperTestSuite) TestGetPrice() {
	s.RegisterCurrency("umilk", "MILK", utils.MustParseDec("2"))

	price, err := s.keeper.GetPrice(s.Ctx, "umilk")
	s.Require().NoError(err)
	s.Require().Equal(utils.MustParseDec("2"), price)
}

func (s *KeeperTestSuite) TestGetCoinValue() {
	// TODO: uncomment this after adding exponent to tickers module
	//s.RegisterCurrency("umilk", "MILK", utils.MustParseDec("2"))
	//
	//value, err := s.keeper.GetCoinValue(s.Ctx, utils.MustParseCoin("10_000000umilk"))
	//s.Require().NoError(err)
	//s.Require().Equal(utils.MustParseDec("20"), value)
}

func (s *KeeperTestSuite) TestGetCoinsValue() {
	// TODO: uncomment this after adding exponent to tickers module
	//s.RegisterCurrency("umilk", "MILK", utils.MustParseDec("2"))
	//s.RegisterCurrency("uinit", "INIT", utils.MustParseDec("3"))
	//
	//value, err := s.keeper.GetCoinsValue(s.Ctx, utils.MustParseCoins("10_000000umilk,5_5000000uinit"))
	//s.Require().NoError(err)
	//s.Require().Equal(utils.MustParseDec("36.5"), value)
}
