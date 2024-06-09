package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/pools/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetNextPoolID() {
	testCases := []struct {
		name  string
		store func(ctx sdk.Context)
		id    uint32
		check func(ctx sdk.Context)
	}{
		{
			name: "next pool id is saved correctly",
			id:   1,
			check: func(ctx sdk.Context) {
				store := ctx.KVStore(suite.storeKey)
				serviceIDBz := store.Get(types.NextPoolIDKey)
				suite.Require().Equal(uint32(1), types.GetPoolIDFromBytes(serviceIDBz))
			},
		},
		{
			name: "next pool id is overridden properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextPoolID(ctx, 1)
			},
			id: 2,
			check: func(ctx sdk.Context) {
				store := ctx.KVStore(suite.storeKey)
				serviceIDBz := store.Get(types.NextPoolIDKey)
				suite.Require().Equal(uint32(2), types.GetPoolIDFromBytes(serviceIDBz))
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

			suite.k.SetNextPoolID(ctx, tc.id)
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
			name:      "non existing next pool id returns error",
			shouldErr: true,
		},
		{
			name: "exiting next pool id is returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SetNextPoolID(ctx, 1)
			},
			expNext: 1,
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
				suite.k.SetNextPoolID(ctx, 1)
			},
			shouldErr: true,
			pool:      types.NewPool(1, "uatom"),
			check: func(ctx sdk.Context) {
				pool, found := suite.k.GetPool(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPool(1, "uatom"), pool)
			},
		},
		{
			name: "existing pool is overridden properly",
			setup: func() {
				suite.k.SetNextPoolID(suite.ctx, 1)
				suite.k.SavePool(suite.ctx, types.NewPool(1, "uatom"))
			},
			pool: types.NewPool(1, "usdt"),
			check: func(ctx sdk.Context) {
				pool, found := suite.k.GetPool(ctx, 1)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPool(1, "usdt"), pool)
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

			suite.k.SavePool(ctx, tc.pool)

			if tc.check != nil {
				tc.check(ctx)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetPool() {
	testCases := []struct {
		name     string
		setup    func()
		store    func(ctx sdk.Context)
		poolID   uint32
		expFound bool
		expPool  types.Pool
		check    func(ctx sdk.Context)
	}{
		{
			name:     "not found pool returns error",
			poolID:   1,
			expFound: false,
		},
		{
			name: "found pool is returned properly",
			store: func(ctx sdk.Context) {
				suite.k.SavePool(ctx, types.NewPool(1, "uatom"))
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

			pool, found := suite.k.GetPool(ctx, tc.poolID)
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
