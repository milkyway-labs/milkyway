package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v4/x/operators/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetParams() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		params    types.Params
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing params are set properly",
			params: types.NewParams(
				sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
				3*time.Hour,
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetParams(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
					3*time.Hour,
				), stored)
			},
		},
		{
			name: "existing params are overridden properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
					3*time.Hour,
				))
				suite.Require().NoError(err)
			},
			params: types.NewParams(
				sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(2000))),
				24*time.Hour,
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				stored, err := suite.k.GetParams(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(2000))),
					24*time.Hour,
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

			err := suite.k.SetParams(ctx, tc.params)
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

func (suite *KeeperTestSuite) TestKeeper_GetParams() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		expParams types.Params
	}{
		{
			name:      "non existing params return an error",
			shouldErr: true,
		},
		{
			name: "existing params are returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetParams(ctx, types.NewParams(
					sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
					3*time.Hour,
				))
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expParams: types.NewParams(
				sdk.NewCoins(sdk.NewCoin("uatom", sdkmath.NewInt(1000))),
				3*time.Hour,
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

			stored, err := suite.k.GetParams(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expParams, stored)
			}
		})
	}
}
