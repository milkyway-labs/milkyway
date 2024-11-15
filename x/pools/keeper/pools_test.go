package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetNextPoolID() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		id        uint32
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "next pool id is saved correctly",
			id:        1,
			shouldErr: false,
			check: func(ctx sdk.Context) {
				nextPoolID, err := suite.k.GetNextPoolID(ctx)
				suite.Require().NoError(err)
				suite.Require().EqualValues(1, nextPoolID)
			},
		},
		{
			name: "next pool id is overridden properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)
			},
			id:        2,
			shouldErr: false,
			check: func(ctx sdk.Context) {
				nextPoolID, err := suite.k.GetNextPoolID(ctx)
				suite.Require().NoError(err)
				suite.Require().EqualValues(2, nextPoolID)
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

			err := suite.k.SetNextPoolID(ctx, tc.id)
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

func (suite *KeeperTestSuite) TestKeeper_GetNextPoolID() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		shouldErr bool
		expNext   uint32
	}{
		{
			name:      "non existing next pool id does not error",
			shouldErr: false,
			expNext:   1,
		},
		{
			name: "exiting next pool id is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 2)
				suite.Require().NoError(err)
			},
			expNext: 2,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			next, err := suite.k.GetNextPoolID(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expNext, next)
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *KeeperTestSuite) TestKeeper_SavePool() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		pool      types.Pool
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "non existing pool is saved properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)
			},
			shouldErr: false,
			pool:      types.NewPool(1, "uatom"),
			check: func(ctx sdk.Context) {
				// Make sure the pool is saved properly
				pool, found, err := suite.k.GetPool(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPool(1, "uatom"), pool)

				// Make sure the pool account is created
				hasAccount := suite.ak.HasAccount(ctx, types.GetPoolAddress(1))
				suite.Require().True(hasAccount)
			},
		},
		{
			name: "existing pool is overridden properly",
			setup: func() {
				err := suite.k.SetNextPoolID(suite.ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SavePool(suite.ctx, types.NewPool(1, "uatom"))
				suite.Require().NoError(err)
			},
			pool:      types.NewPool(1, "usdt"),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				// Make sure the pool is saved properly
				pool, found, err := suite.k.GetPool(ctx, 1)
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPool(1, "usdt"), pool)

				// Make sure the pool account is created
				hasAccount := suite.ak.HasAccount(ctx, types.GetPoolAddress(1))
				suite.Require().True(hasAccount)
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

			err := suite.k.SavePool(ctx, tc.pool)
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

func (suite *KeeperTestSuite) TestKeeper_GetPool() {
	testCases := []struct {
		name      string
		setup     func()
		store     func(ctx sdk.Context)
		poolID    uint32
		shouldErr bool
		expFound  bool
		expPool   types.Pool
		check     func(ctx sdk.Context)
	}{
		{
			name:     "not found pool returns error",
			poolID:   1,
			expFound: false,
		},
		{
			name: "found pool is returned properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SavePool(ctx, types.NewPool(1, "uatom"))
				suite.Require().NoError(err)
			},
			poolID:   1,
			expFound: true,
			expPool:  types.NewPool(1, "uatom"),
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

			pool, found, err := suite.k.GetPool(ctx, tc.poolID)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expFound, found)
				if tc.expFound {
					suite.Require().Equal(tc.expPool, pool)
				}
			}

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}
