package keeper_test

import (
	"github.com/milkyway-labs/milkyway/x/assets/types"
)

func (s *KeeperTestSuite) TestRemoveAsset() {
	err := s.App.AssetsKeeper.SetAsset(s.Ctx, types.NewAsset("umilk", "MILK", 6))
	s.Require().NoError(err)

	err = s.App.AssetsKeeper.RemoveAsset(s.Ctx, "umilk")
	s.Require().NoError(err)

	resp, err := s.queryServer.Assets(s.Ctx, &types.QueryAssetsRequest{Ticker: "MILK"})
	s.Require().NoError(err)
	s.Require().Empty(resp.Assets)
}
