package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/pools/types"
)

func (suite *KeeperTestSuite) TestKeeper_ExportGenesis() {
	testCases := []struct {
		name       string
		setup      func()
		store      func(ctx sdk.Context)
		expGenesis *types.GenesisState
	}{
		{
			name: "next pool id is exported properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 10)
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				NextPoolID: 10,
			},
		},
		{
			name: "pools are exported properly",
			store: func(ctx sdk.Context) {
				err := suite.k.SetNextPoolID(ctx, 1)
				suite.Require().NoError(err)

				err = suite.k.SavePool(ctx, types.NewPool(1, "umilk"))
				suite.Require().NoError(err)
				err = suite.k.SavePool(ctx, types.NewPool(2, "uatom"))
				suite.Require().NoError(err)
			},
			expGenesis: &types.GenesisState{
				NextPoolID: 1,
				Pools: []types.Pool{
					types.NewPool(1, "umilk"),
					types.NewPool(2, "uatom"),
				},
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

			genesis := suite.k.ExportGenesis(ctx)
			suite.Require().Equal(tc.expGenesis, genesis)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_InitGenesis() {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name:      "default genesis is initialized properly",
			genesis:   types.DefaultGenesis(),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				nextPoolID, err := suite.k.GetNextPoolID(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(uint32(1), nextPoolID)
			},
		},
		{
			name: "genesis with pools is initialized properly",
			genesis: types.NewGenesis(
				10,
				[]types.Pool{
					types.NewPool(1, "umilk"),
					types.NewPool(2, "uatom"),
				},
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				nextPoolID, err := suite.k.GetNextPoolID(ctx)
				suite.Require().NoError(err)
				suite.Require().Equal(uint32(10), nextPoolID)

				pool, found, err := suite.k.GetPoolByDenom(ctx, "umilk")
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPool(1, "umilk"), pool)

				pool, found, err = suite.k.GetPoolByDenom(ctx, "uatom")
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().Equal(types.NewPool(2, "uatom"), pool)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			ctx, _ := suite.ctx.CacheContext()

			err := suite.k.InitGenesis(ctx, tc.genesis)
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
