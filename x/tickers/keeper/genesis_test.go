package keeper_test

import (
	"github.com/milkyway-labs/milkyway/x/tickers/types"
)

func (s *KeeperTestSuite) TestExportGenesis() {
	_, err := s.msgServer.RegisterTicker(s.Ctx, types.NewMsgRegisterTicker(s.authority, "umilk", "MILK"))
	s.Require().NoError(err)

	_, err = s.msgServer.RegisterTicker(s.Ctx, types.NewMsgRegisterTicker(s.authority, "umilk2", "MILK"))
	s.Require().NoError(err)

	_, err = s.msgServer.RegisterTicker(s.Ctx, types.NewMsgRegisterTicker(s.authority, "uatom", "ATOM"))
	s.Require().NoError(err)

	genState := s.App.TickersKeeper.ExportGenesis(s.Ctx)

	s.Require().Equal(types.Tickers{
		"umilk":  "MILK",
		"umilk2": "MILK",
		"uatom":  "ATOM",
	}, genState.Tickers)
}

func (s *KeeperTestSuite) TestInitGenesis() {
	genState := types.NewGenesisState(types.DefaultParams(), types.Tickers{
		"umilk":  "MILK",
		"umilk2": "MILK",
		"uatom":  "ATOM",
	})

	s.App.TickersKeeper.InitGenesis(s.Ctx, genState)

	resp, err := s.queryServer.Ticker(s.Ctx, &types.QueryTickerRequest{Denom: "umilk"})
	s.Require().NoError(err)
	s.Require().Equal("MILK", resp.Ticker)
}
