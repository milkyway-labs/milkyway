package keeper_test

import (
	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

func (s *KeeperTestSuite) TestRemoveAsset() {
	err := s.App.TickersKeeper.SetAsset(s.Ctx, types.NewAsset("umilk", "MILK", 6))
	s.Require().NoError(err)

	err = s.App.TickersKeeper.RemoveAsset(s.Ctx, "umilk")
	s.Require().NoError(err)

	resp, err := s.queryServer.Assets(s.Ctx, &types.QueryAssetsRequest{Ticker: "MILK"})
	s.Require().NoError(err)
	s.Require().Empty(resp.Assets)
}
