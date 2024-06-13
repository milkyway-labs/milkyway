package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

func (suite *KeeperTestSuite) TestKeeper_GetPoolForDenom() {
	testCases := []struct {
		name     string
		setup    func()
		store    func(ctx sdk.Context)
		denom    string
		expFound bool
		expPool  types.Pool
		check    func(ctx sdk.Context)
	}{
		{
			name:     "non exiting pool returns error",
			denom:    "denom",
			expFound: false,
		},
		{
			name: "existing pool is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
			},
			denom:    "umilk",
			expFound: true,
			expPool:  types.NewPool(1, "umilk"),
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

			pool, found := suite.k.GetPoolByDenom(ctx, tc.denom)
			suite.Require().Equal(tc.expFound, found)
			if tc.expFound {
				suite.Require().Equal(tc.expPool, pool)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_CreateOrGetPoolByDenom() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		denom     string
		shouldErr bool
		expPool   types.Pool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "invalid next pool id returns error",
			denom:     "umilk",
			shouldErr: true,
		},
		{
			name: "invalid pool returns error",
			store: func(ctx sdk.Context) {
				suite.k.SetNextPoolID(ctx, 1)
			},
			denom:     "invalid!",
			shouldErr: true,
		},
		{
			name: "existing pool is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
			},
			denom:     "umilk",
			shouldErr: false,
			expPool:   types.NewPool(1, "umilk"),
		},
		{
			name: "non existing pool is created properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextPoolID(ctx, 1)
			},
			denom:     "umilk",
			shouldErr: false,
			expPool:   types.NewPool(1, "umilk"),
			check: func(ctx sdk.Context) {
				// Make sure the pool is stored properly
				pool, found := suite.k.GetPoolByDenom(ctx, "umilk")
				suite.Require().True(found)
				suite.Require().Equal(types.NewPool(1, "umilk"), pool)

				// Make sure the next pool id has been incremented
				nextPoolID, err := suite.k.GetNextPoolID(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(uint32(2), nextPoolID)
			},
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

			pool, err := suite.k.CreateOrGetPoolByDenom(ctx, tc.denom)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expPool, pool)
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
