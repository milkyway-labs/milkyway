package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/pools/keeper"
	"github.com/milkyway-labs/milkyway/v12/x/pools/types"
)

func (suite *KeeperTestSuite) TestValidPoolsInvariant() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		expBroken bool
	}{
		{
			name: "not found next service id breaks invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "service with id equals to next service id breaks invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "service with id higher than next service id breaks invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SavePool(ctx, types.NewPool(2, "umilk"))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "invalid service breaks invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 2)
				suite.Require().NoError(err)

				err = suite.k.SavePool(ctx, types.NewPool(1, "invalid!"))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "valid data does not break invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 3)
				suite.Require().NoError(err)

				err = suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
				err = suite.k.SavePool(ctx, types.NewPool(2, "unit"))
				suite.Require().NoError(err)
			},
			expBroken: false,
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

			_, broken := keeper.ValidPoolsInvariant(suite.k)(ctx)
			suite.Require().Equal(tc.expBroken, broken)
		})
	}
}

func (suite *KeeperTestSuite) TestUniquePoolsInvariant() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		expBroken bool
	}{
		{
			name: "duplicated pools break invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
				err = suite.k.SavePool(ctx, types.NewPool(2, "umilk"))
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "valid data does not break invariant",
			store: func(ctx sdk.Context) {
				err := suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
				err = suite.k.SavePool(ctx, types.NewPool(2, "unit"))
				suite.Require().NoError(err)
			},
			expBroken: false,
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

			_, broken := keeper.UniquePoolsInvariant(suite.k)(ctx)
			suite.Require().Equal(tc.expBroken, broken)
		})
	}
}
