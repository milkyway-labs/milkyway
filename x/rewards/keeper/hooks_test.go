package keeper_test

import (
	"github.com/milkyway-labs/milkyway/v7/app/testutil"
	"github.com/milkyway-labs/milkyway/v7/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/v7/x/restaking/types"
)

func (suite *KeeperTestSuite) TestPoolServiceTotalDelegatorShares() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "umilk", 6, utils.MustParseDec("1"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Alice by default trusts all services through all pools, so when Alice delegates
	// it increments pool-service total delegator shares of the pool.
	aliceAddr := testutil.TestAddress(1)
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), aliceAddr.String(), true)

	shares, err := suite.keeper.GetPoolServiceTotalDelegatorShares(ctx, 1, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDecCoins("10_000000pool/1/umilk"), shares)

	// Bob doesn't trust the service through the pool, so when Bob delegates it
	// doesn't increment pool-service total delegator shares of the pool.
	bobAddr := testutil.TestAddress(2)
	suite.SetUserPreferences(ctx, bobAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(1, []uint32{2}),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), bobAddr.String(), true)

	shares, err = suite.keeper.GetPoolServiceTotalDelegatorShares(ctx, 1, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDecCoins("10_000000pool/1/umilk"), shares)

	// But when Bob decides to trust the service through the pool, it increments
	// the total shares.

	suite.SetUserPreferences(ctx, bobAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(1, []uint32{service.ID}),
	})
	shares, err = suite.keeper.GetPoolServiceTotalDelegatorShares(ctx, 1, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDecCoins("20_000000pool/1/umilk"), shares)
}
