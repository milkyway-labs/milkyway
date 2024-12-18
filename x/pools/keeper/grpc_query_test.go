package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v6/x/pools/types"
)

func (suite *KeeperTestSuite) TestQueryServer_PoolByID() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		request   *types.QueryPoolByIdRequest
		shouldErr bool
		expPool   types.Pool
	}{
		{
			name: "not found pool returns error",
			request: &types.QueryPoolByIdRequest{
				PoolId: 1,
			},
			shouldErr: true,
		},
		{
			name: "found pool is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
			},
			request: &types.QueryPoolByIdRequest{
				PoolId: 1,
			},
			shouldErr: false,
			expPool:   types.NewPool(1, "umilk"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			res, err := suite.k.PoolByID(ctx, tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expPool, res.Pool)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryServer_PoolByDenom() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		request   *types.QueryPoolByDenomRequest
		shouldErr bool
		expPool   types.Pool
	}{
		{
			name: "not found pool returns error",
			request: &types.QueryPoolByDenomRequest{
				Denom: "umilk",
			},
			shouldErr: true,
		},
		{
			name: "found pool is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
			},
			request: &types.QueryPoolByDenomRequest{
				Denom: "umilk",
			},
			shouldErr: false,
			expPool:   types.NewPool(1, "umilk"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.setup != nil {
				tc.setup()
			}
			if tc.store != nil {
				tc.store(ctx)
			}

			res, err := suite.k.PoolByDenom(ctx, tc.request)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expPool, res.Pool)
			}
		})
	}
}
