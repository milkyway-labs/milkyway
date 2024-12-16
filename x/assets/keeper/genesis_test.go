package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v7/x/assets/types"
)

func (suite *KeeperTestSuite) TestExportGenesis() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		shouldErr  bool
		expGenesis *types.GenesisState
	}{
		{
			name: "assets are exported correctly",
			store: func(ctx sdk.Context) {
				err := suite.keeper.SetAsset(ctx, types.NewAsset("umilk", "MILK", 6))
				suite.Require().NoError(err)

				err = suite.keeper.SetAsset(ctx, types.NewAsset("umilk2", "MILK", 6))
				suite.Require().NoError(err)

				err = suite.keeper.SetAsset(ctx, types.NewAsset("uatom", "ATOM", 6))
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expGenesis: &types.GenesisState{
				Assets: []types.Asset{
					types.NewAsset("uatom", "ATOM", 6),
					types.NewAsset("umilk", "MILK", 6),
					types.NewAsset("umilk2", "MILK", 6),
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			genState, err := suite.keeper.ExportGenesis(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expGenesis, genState)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestInitGenesis() {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "genesis is initialized properly",
			genesis: types.NewGenesisState(
				[]types.Asset{
					types.NewAsset("umilk", "MILK", 6),
					types.NewAsset("umilk2", "MILK", 6),
					types.NewAsset("uatom", "ATOM", 6),
				},
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				stored, err := suite.keeper.GetAsset(ctx, "umilk")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewAsset("umilk", "MILK", 6), stored)

				stored, err = suite.keeper.GetAsset(ctx, "umilk2")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewAsset("umilk2", "MILK", 6), stored)

				stored, err = suite.keeper.GetAsset(ctx, "uatom")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewAsset("uatom", "ATOM", 6), stored)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			err := suite.keeper.InitGenesis(ctx, tc.genesis)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}

}
