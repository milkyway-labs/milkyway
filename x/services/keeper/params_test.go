package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/services/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetParams() {
	testCases := []struct {
		name   string
		store  func(ctx sdk.Context)
		params types.Params
		check  func(ctx sdk.Context)
	}{
		{
			name: "non existing params are set properly",
			params: types.NewParams(
				sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
			),
			check: func(ctx sdk.Context) {
				stored := suite.k.GetParams(ctx)
				suite.Require().Equal(types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
				), stored)
			},
		},
		{
			name: "existing params are overridden properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
				))
			},
			params: types.NewParams(
				sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(2000))),
			),
			check: func(ctx sdk.Context) {
				stored := suite.k.GetParams(ctx)
				suite.Require().Equal(types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(2000))),
				), stored)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			suite.k.SetParams(ctx, tc.params)
			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetParams() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		expParams types.Params
	}{
		{
			name:      "non existing params are returned properly",
			expParams: types.Params{},
		},
		{
			name: "existing params are returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
				))
			},
			expParams: types.NewParams(
				sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			stored := suite.k.GetParams(ctx)
			suite.Require().Equal(tc.expParams, stored)
		})
	}
}
