package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

func (suite *KeeperTestSuite) TestExportGenesis() {
	testCases := []struct {
		name       string
		store      func(ctx sdk.Context)
		shouldErr  bool
		expGenesis *types.GenesisState
	}{
		{
			name: "params are exported correctly",
			store: func(ctx sdk.Context) {
				err := suite.k.Params.Set(ctx, types.DefaultParams())
				suite.Require().NoError(err)
			},
			shouldErr: false,
			expGenesis: &types.GenesisState{
				Params: types.DefaultParams(),
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			genState, err := suite.k.ExportGenesis(ctx)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expGenesis, genState)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestInitGenesis() {
	testCases := []struct {
		name      string
		genesis   *types.GenesisState
		shouldErr bool
		check     func(ctx sdk.Context)
	}{
		{
			name: "genesis is initialized properly",
			genesis: types.NewGenesisState(
				types.DefaultParams(),
			),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				params, _ := suite.k.Params.Get(ctx)
				suite.Assert().Equal(types.DefaultParams(), params)
			},
		},
		{
			name:      "should block negative insurance fund percentage",
			genesis:   types.NewGenesisState(types.NewParams(math.LegacyNewDec(-1), nil, nil)),
			shouldErr: true,
		},
		{
			name:      "should block 0 insurance fund percentage",
			genesis:   types.NewGenesisState(types.NewParams(math.LegacyNewDec(0), nil, nil)),
			shouldErr: true,
		},
		{
			name:      "should allow 100 insurance fund percentage",
			genesis:   types.NewGenesisState(types.NewParams(math.LegacyNewDec(100), nil, nil)),
			shouldErr: false,
			check: func(ctx sdk.Context) {
				params, _ := suite.k.Params.Get(ctx)
				suite.Assert().Equal(math.LegacyNewDec(100), params.InsurancePercentage)
			},
		},
		{
			name:      "should block > 100 insurance fund percentage",
			genesis:   types.NewGenesisState(types.NewParams(math.LegacyNewDec(101), nil, nil)),
			shouldErr: true,
		},
		{
			name:      "should block invalid minter address",
			genesis:   types.NewGenesisState(types.NewParams(math.LegacyNewDec(2), nil, []string{"cosmos1fdsfd"})),
			shouldErr: true,
		},
		{
			name:      "should block invalid burners address",
			genesis:   types.NewGenesisState(types.NewParams(math.LegacyNewDec(2), []string{"cosmos1fdsfd"}, nil)),
			shouldErr: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()

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
