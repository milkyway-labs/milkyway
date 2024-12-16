package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/milkyway-labs/milkyway/v7/x/assets/keeper"
	"github.com/milkyway-labs/milkyway/v7/x/assets/types"
)

func (suite *KeeperTestSuite) TestQuerier_Assets() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		req       *types.QueryAssetsRequest
		shouldErr bool
		expAssets []types.Asset
	}{
		{
			name: "query without pagination returns all assets",
			store: func(ctx sdk.Context) {
				err := suite.keeper.SetAsset(ctx, types.NewAsset("umilk", "MILK", 6))
				suite.Require().NoError(err)

				err = suite.keeper.SetAsset(ctx, types.NewAsset("umilk2", "MILK", 6))
				suite.Require().NoError(err)

				err = suite.keeper.SetAsset(ctx, types.NewAsset("uatom", "ATOM", 6))
				suite.Require().NoError(err)
			},
			req:       types.NewQueryAssetsRequest("", nil),
			shouldErr: false,
			expAssets: []types.Asset{
				types.NewAsset("uatom", "ATOM", 6),
				types.NewAsset("umilk", "MILK", 6),
				types.NewAsset("umilk2", "MILK", 6),
			},
		},
		{
			name: "query with ticker returns assets with the given ticker",
			store: func(ctx sdk.Context) {
				err := suite.keeper.SetAsset(ctx, types.NewAsset("umilk", "MILK", 6))
				suite.Require().NoError(err)

				err = suite.keeper.SetAsset(ctx, types.NewAsset("umilk2", "MILK", 6))
				suite.Require().NoError(err)

				err = suite.keeper.SetAsset(ctx, types.NewAsset("uatom", "ATOM", 6))
				suite.Require().NoError(err)
			},
			req:       types.NewQueryAssetsRequest("MILK", nil),
			shouldErr: false,
			expAssets: []types.Asset{
				types.NewAsset("umilk", "MILK", 6),
				types.NewAsset("umilk2", "MILK", 6),
			},
		},
		{
			name: "query with ticker and pagination returns assets with the given ticker",
			store: func(ctx sdk.Context) {
				err := suite.keeper.SetAsset(ctx, types.NewAsset("umilk", "MILK", 6))
				suite.Require().NoError(err)

				err = suite.keeper.SetAsset(ctx, types.NewAsset("umilk2", "MILK", 6))
				suite.Require().NoError(err)

				err = suite.keeper.SetAsset(ctx, types.NewAsset("uatom", "ATOM", 6))
				suite.Require().NoError(err)
			},
			req: types.NewQueryAssetsRequest("MILK", &query.PageRequest{
				Limit:  1,
				Offset: 1,
			}),
			shouldErr: false,
			expAssets: []types.Asset{
				types.NewAsset("umilk2", "MILK", 6),
			},
		},
		{
			name:      "invalid ticker returns error",
			req:       types.NewQueryAssetsRequest("!@#$", nil),
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.Assets(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expAssets, res.Assets)
			}

		})
	}
}

func (suite *KeeperTestSuite) TestQuerier_Asset() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		req       *types.QueryAssetRequest
		shouldErr bool
		expAsset  types.Asset
	}{
		{
			name: "existing asset is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.keeper.SetAsset(ctx, types.NewAsset("umilk", "MILK", 6))
				suite.Require().NoError(err)
			},
			req:       types.NewQueryAssetRequest("umilk"),
			shouldErr: false,
			expAsset:  types.NewAsset("umilk", "MILK", 6),
		},
		{
			name:      "invalid denom returns error",
			req:       types.NewQueryAssetRequest("!@#$"),
			shouldErr: true,
		},
		{
			name:      "non-existing asset returns error",
			req:       types.NewQueryAssetRequest("umilk"),
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			queryServer := keeper.NewQueryServer(suite.keeper)
			res, err := queryServer.Asset(ctx, tc.req)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expAsset, res.Asset)
			}
		})
	}
}
