package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/assets/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetAsset() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		asset     types.Asset
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:  "non existing asset is set properly",
			asset: types.NewAsset("umilk", "MILK", 6),
			check: func(ctx sdk.Context) {
				asset, err := suite.keeper.GetAsset(ctx, "umilk")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewAsset("umilk", "MILK", 6), asset)
			},
		},
		{
			name: "existing asset is updated properly",
			store: func(ctx sdk.Context) {
				err := suite.keeper.SetAsset(ctx, types.NewAsset("umilk", "MILK", 6))
				suite.Require().NoError(err)
			},
			asset: types.NewAsset("umilk", "MILK", 8),
			check: func(ctx sdk.Context) {
				asset, err := suite.keeper.GetAsset(ctx, "umilk")
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewAsset("umilk", "MILK", 8), asset)
			},
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

			err := suite.keeper.SetAsset(ctx, tc.asset)
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

func (suite *KeeperTestSuite) TestKeeper_GetAsset() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		ticker    string
		shouldErr bool
		expAsset  types.Asset
	}{
		{
			name:      "not found asset returns error",
			ticker:    "MILK",
			shouldErr: true,
		},
		{
			name: "found asset is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.keeper.SetAsset(ctx, types.NewAsset("umilk", "MILK", 6))
				suite.Require().NoError(err)
			},
			ticker:   "umilk",
			expAsset: types.NewAsset("umilk", "MILK", 6),
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

			asset, err := suite.keeper.GetAsset(ctx, tc.ticker)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expAsset, asset)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestRemoveAsset() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		ticker    string
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "not found asset returns error",
			ticker:    "MILK",
			shouldErr: true,
		},
		{
			name: "found asset is removed properly",
			store: func(ctx sdk.Context) {
				err := suite.keeper.SetAsset(ctx, types.NewAsset("umilk", "MILK", 6))
				suite.Require().NoError(err)
			},
			ticker: "umilk",
			check: func(ctx sdk.Context) {
				_, err := suite.keeper.GetAsset(ctx, "umilk")
				suite.Require().Error(err)
			},
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

			err := suite.keeper.RemoveAsset(ctx, tc.ticker)
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
