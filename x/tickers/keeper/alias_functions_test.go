package keeper_test

import (
	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

func (s *KeeperTestSuite) TestRemoveTicker() {
	err := s.App.TickersKeeper.SetTicker(s.Ctx, "umilk", "MILK")
	s.Require().NoError(err)

	err = s.App.TickersKeeper.RemoveTicker(s.Ctx, "umilk")
	s.Require().NoError(err)

	resp, err := s.queryServer.Denoms(s.Ctx, &types.QueryDenomsRequest{Ticker: "MILK"})
	s.Require().NoError(err)
	s.Require().Empty(resp.Denoms)
}
