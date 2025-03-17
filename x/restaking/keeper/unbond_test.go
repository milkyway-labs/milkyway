package keeper_test

import (
	"time"

	"github.com/milkyway-labs/milkyway/v10/utils"
	"github.com/milkyway-labs/milkyway/v10/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v10/x/restaking/types"
)

func (suite *KeeperTestSuite) TestKeeper_MaxUnbondingEntries() {
	ctx, _ := suite.ctx.CacheContext()

	params, err := suite.k.GetParams(ctx)
	suite.Require().NoError(err)
	params.MaxEntries = 2
	err = suite.k.SetParams(ctx, params)
	suite.Require().NoError(err)

	delegator := "cosmos167x6ehhple8gwz5ezy9x0464jltvdpzl6qfdt4"
	suite.fundAccount(ctx, delegator, utils.MustParseCoins("10000_000000umilk"))

	// First delegate to pool
	msgServer := keeper.NewMsgServer(suite.k)
	_, err = msgServer.DelegatePool(ctx, types.NewMsgDelegatePool(utils.MustParseCoin("10000_000000umilk"), delegator))
	suite.Require().NoError(err)

	// Unbonding from pool for the first two times should be successful
	_, err = msgServer.UndelegatePool(ctx, types.NewMsgUndelegatePool(utils.MustParseCoin("1000_000000umilk"), delegator))
	suite.Require().NoError(err)
	// Increase the block height by 1 so that a separate unbonding entry is created
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(5 * time.Second))
	_, err = msgServer.UndelegatePool(ctx, types.NewMsgUndelegatePool(utils.MustParseCoin("1000_000000umilk"), delegator))
	suite.Require().NoError(err)

	// But it should fail for the third time, since it exceeds the max entries
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(5 * time.Second))
	_, err = msgServer.UndelegatePool(ctx, types.NewMsgUndelegatePool(utils.MustParseCoin("1000_000000umilk"), delegator))
	suite.Require().ErrorIs(err, types.ErrMaxUnbondingDelegationEntries)
}
