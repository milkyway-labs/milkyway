package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/pools/types"
)

func (suite *KeeperTestSuite) TestKeeper_GetPoolForDenom() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		denom     string
		shouldErr bool
		expFound  bool
		expPool   types.Pool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "non exiting pool returns error",
			denom:     "denom",
			expFound:  false,
			shouldErr: false,
		},
		{
			name: "existing pool is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
			},
			denom:     "umilk",
			shouldErr: false,
			expFound:  true,
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

			pool, found, err := suite.k.GetPoolByDenom(ctx, tc.denom)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expFound, found)
				if tc.expFound {
					suite.Require().Equal(tc.expPool, pool)
				}

				if tc.check != nil {
					tc.check(ctx)
				}
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
			name: "invalid pool returns error",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)
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
				err := suite.k.SetNextPoolID(ctx, 2)
				suite.Require().NoError(err)

				err = suite.k.SavePool(ctx, types.NewPool(1, "unit"))
				suite.Require().NoError(err)
			},
			denom:     "umilk",
			shouldErr: false,
			expPool:   types.NewPool(2, "umilk"),
			check: func(ctx sdk.Context) {
				// Make sure the pool is stored properly
				pool, found, err := suite.k.GetPoolByDenom(ctx, "umilk")
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPool(2, "umilk"), pool)

				// Make sure the next pool id has been incremented
				nextPoolID, err := suite.k.GetNextPoolID(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(uint32(3), nextPoolID)
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
