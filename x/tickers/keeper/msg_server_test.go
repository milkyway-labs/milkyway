package keeper_test

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

func (s *KeeperTestSuite) TestRegisterTicker() {
	_, err := s.msgServer.RegisterTicker(s.Ctx, types.NewMsgRegisterTicker(s.authority, "umilk", "MILK"))
	s.Require().NoError(err)

	resp, err := s.queryServer.Ticker(s.Ctx, &types.QueryTickerRequest{Denom: "umilk"})
	s.Require().NoError(err)
	s.Require().Equal("MILK", resp.Ticker)
}

func (s *KeeperTestSuite) TestDeregisterTicker() {
	_, err := s.msgServer.RegisterTicker(s.Ctx, types.NewMsgRegisterTicker(s.authority, "umilk", "MILK"))
	s.Require().NoError(err)

	_, err = s.msgServer.DeregisterTicker(s.Ctx, types.NewMsgDeregisterTicker(s.authority, "umilk"))
	s.Require().NoError(err)

	_, err = s.queryServer.Ticker(s.Ctx, &types.QueryTickerRequest{Denom: "umilk"})
	s.Require().Equal(codes.NotFound, status.Code(err))
}
