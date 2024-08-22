package keeper_test

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/x/assets/types"
)

func (s *KeeperTestSuite) TestRegisterAsset() {
	_, err := s.msgServer.RegisterAsset(s.Ctx, &types.MsgRegisterAsset{
		Authority: s.authority,
		Asset:     types.NewAsset("umilk", "MILK", 6),
	})
	s.Require().NoError(err)

	resp, err := s.queryServer.Asset(s.Ctx, &types.QueryAssetRequest{Denom: "umilk"})
	s.Require().NoError(err)
	s.Require().Equal(types.NewAsset("umilk", "MILK", 6), resp.Asset)
}

func (s *KeeperTestSuite) TestDeregisterAsset() {
	_, err := s.msgServer.RegisterAsset(s.Ctx, &types.MsgRegisterAsset{
		Authority: s.authority,
		Asset:     types.NewAsset("umilk", "MILK", 6),
	})
	s.Require().NoError(err)

	_, err = s.msgServer.DeregisterAsset(s.Ctx, &types.MsgDeregisterAsset{
		Authority: s.authority,
		Denom:     "umilk",
	})
	s.Require().NoError(err)

	_, err = s.queryServer.Asset(s.Ctx, &types.QueryAssetRequest{Denom: "umilk"})
	s.Require().Equal(codes.NotFound, status.Code(err))
}
