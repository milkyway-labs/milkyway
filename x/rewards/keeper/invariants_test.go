package keeper_test

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v4/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v4/x/pools/types"
	"github.com/milkyway-labs/milkyway/v4/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/v4/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/v4/x/services/types"
)

func (suite *KeeperTestSuite) TestInvariants_ReferenceCountInvariant() {
	testCases := []struct {
		name      string
		store     func(ctx sdk.Context)
		expBroken bool
	}{
		{
			name: "default genesis does not return errors",
			store: func(ctx sdk.Context) {
				// Store the genesis data for all the modules involved
				err := suite.poolsKeeper.InitGenesis(ctx, poolstypes.DefaultGenesis())
				suite.Require().NoError(err)
				err = suite.servicesKeeper.InitGenesis(ctx, servicestypes.DefaultGenesis())
				suite.NoError(err)
				err = suite.operatorsKeeper.InitGenesis(ctx, operatorstypes.DefaultGenesis())
				suite.NoError(err)
				err = suite.keeper.InitGenesis(ctx, types.DefaultGenesis())
				suite.NoError(err)
			},
			expBroken: false,
		},
		{
			name: "initializing services and operators does not break the invariant",
			store: func(ctx sdk.Context) {
				suite.CreateOperator(ctx, "MilkyWay", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
				suite.CreateService(ctx, "MilkyWay AVS", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")
			},
			expBroken: false,
		},
		{
			name: "invalid pools historical rewards reference count breaks the invariant",
			store: func(ctx sdk.Context) {
				pool, err := suite.poolsKeeper.CreateOrGetPoolByDenom(ctx, "stake")
				suite.NoError(err)

				// Create an invalid number of historical rewards
				historicalRewards := types.NewHistoricalRewards(types.ServicePools{
					types.NewServicePool(1, types.DecPools{types.NewDecPool("umilk", nil)}),
				}, 1)
				err = suite.keeper.PoolHistoricalRewards.Set(ctx, collections.Join[uint32, uint64](pool.ID, 1), historicalRewards)
				err = suite.keeper.PoolHistoricalRewards.Set(ctx, collections.Join[uint32, uint64](pool.ID, 2), historicalRewards)
				err = suite.keeper.PoolHistoricalRewards.Set(ctx, collections.Join[uint32, uint64](pool.ID, 3), historicalRewards)
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "invalid service historical rewards reference count breaks the invariant",
			store: func(ctx sdk.Context) {
				service := suite.CreateService(ctx, "MilkyWay AVS", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")

				// Create an invalid number of historical rewards
				historicalRewards := types.NewHistoricalRewards(types.ServicePools{
					types.NewServicePool(1, types.DecPools{types.NewDecPool("umilk", nil)}),
				}, 1)
				err := suite.keeper.ServiceHistoricalRewards.Set(ctx, collections.Join[uint32, uint64](service.ID, 1), historicalRewards)
				err = suite.keeper.ServiceHistoricalRewards.Set(ctx, collections.Join[uint32, uint64](service.ID, 2), historicalRewards)
				err = suite.keeper.ServiceHistoricalRewards.Set(ctx, collections.Join[uint32, uint64](service.ID, 3), historicalRewards)
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
		{
			name: "invalid operator historical rewards reference count breaks the invariant",
			store: func(ctx sdk.Context) {
				service := suite.CreateOperator(ctx, "MilkyWay", "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4")

				// Create an invalid number of historical rewards
				historicalRewards := types.NewHistoricalRewards(types.ServicePools{
					types.NewServicePool(1, types.DecPools{types.NewDecPool("umilk", nil)}),
				}, 1)
				err := suite.keeper.OperatorHistoricalRewards.Set(ctx, collections.Join[uint32, uint64](service.ID, 1), historicalRewards)
				err = suite.keeper.OperatorHistoricalRewards.Set(ctx, collections.Join[uint32, uint64](service.ID, 2), historicalRewards)
				err = suite.keeper.OperatorHistoricalRewards.Set(ctx, collections.Join[uint32, uint64](service.ID, 3), historicalRewards)
				suite.Require().NoError(err)
			},
			expBroken: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()

			ctx, _ := suite.ctx.CacheContext()
			if tc.store != nil {
				tc.store(ctx)
			}

			res, broken := keeper.ReferenceCountInvariant(suite.keeper)(ctx)
			suite.Equal(tc.expBroken, broken, res)
		})
	}
}
