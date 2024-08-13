package keeper_test

import (
	"github.com/milkyway-labs/milkyway/x/assets/types"
)

func (s *KeeperTestSuite) TestExportGenesis() {
	_, err := s.msgServer.RegisterAsset(s.Ctx, &types.MsgRegisterAsset{
		Authority: s.authority,
		Asset:     types.NewAsset("umilk", "MILK", 6),
	})
	s.Require().NoError(err)

	_, err = s.msgServer.RegisterAsset(s.Ctx, &types.MsgRegisterAsset{
		Authority: s.authority,
		Asset:     types.NewAsset("umilk2", "MILK", 6),
	})
	s.Require().NoError(err)

	_, err = s.msgServer.RegisterAsset(s.Ctx, &types.MsgRegisterAsset{
		Authority: s.authority,
		Asset:     types.NewAsset("uatom", "ATOM", 6),
	})
	s.Require().NoError(err)

	genState := s.App.AssetsKeeper.ExportGenesis(s.Ctx)

	s.Require().Equal([]types.Asset{
		types.NewAsset("uatom", "ATOM", 6),
		types.NewAsset("umilk", "MILK", 6),
		types.NewAsset("umilk2", "MILK", 6),
	}, genState.Assets)
}

func (s *KeeperTestSuite) TestInitGenesis() {
	genState := types.NewGenesisState(types.DefaultParams(), []types.Asset{
		types.NewAsset("umilk", "MILK", 6),
		types.NewAsset("umilk2", "MILK", 6),
		types.NewAsset("uatom", "ATOM", 6),
	})

	s.App.AssetsKeeper.InitGenesis(s.Ctx, genState)

	resp, err := s.queryServer.Asset(s.Ctx, &types.QueryAssetRequest{Denom: "umilk"})
	s.Require().NoError(err)
	s.Require().Equal(types.NewAsset("umilk", "MILK", 6), resp.Asset)

	resp, err = s.queryServer.Asset(s.Ctx, &types.QueryAssetRequest{Denom: "umilk2"})
	s.Require().NoError(err)
	s.Require().Equal(types.NewAsset("umilk2", "MILK", 6), resp.Asset)

	resp, err = s.queryServer.Asset(s.Ctx, &types.QueryAssetRequest{Denom: "uatom"})
	s.Require().NoError(err)
	s.Require().Equal(types.NewAsset("uatom", "ATOM", 6), resp.Asset)
}
