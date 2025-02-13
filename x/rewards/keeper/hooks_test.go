package keeper_test

import (
	"github.com/milkyway-labs/milkyway/v9/app/testutil"
	"github.com/milkyway-labs/milkyway/v9/utils"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

func (suite *KeeperTestSuite) TestPoolServiceTotalDelegatorShares() {
	// Cache the context to avoid test conflicts
	ctx, _ := suite.ctx.CacheContext()

	suite.RegisterCurrency(ctx, "umilk", "umilk", 6, utils.MustParseDec("1"))
	suite.RegisterCurrency(ctx, "utia", "TIA", 6, utils.MustParseDec("3"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Alice by default trusts all services through all pools, so when Alice delegates
	// it increments pool-service total delegator shares of the pool.
	aliceAddr := testutil.TestAddress(1)
	suite.SetUserPreferences(ctx, aliceAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, nil),
	})
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), aliceAddr.String(), true)

	shares, err := suite.keeper.GetPoolServiceTotalDelegatorShares(ctx, 1, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDecCoins("10_000000pool/1/umilk"), shares)

	// Bob doesn't trust the service through the $MILK pool, so when Bob delegates it
	// doesn't increment pool-service total delegator shares of the $MILK pool.
	bobAddr := testutil.TestAddress(2)
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000umilk"), bobAddr.String(), true)
	suite.DelegatePool(ctx, utils.MustParseCoin("10_000000utia"), bobAddr.String(), true)
	suite.SetUserPreferences(ctx, bobAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, []uint32{2}),
	})

	shares, err = suite.keeper.GetPoolServiceTotalDelegatorShares(ctx, 1, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDecCoins("10_000000pool/1/umilk"), shares)

	// But when Bob decides to trust the service through the $MILK pool, it
	// increments the total shares.

	suite.SetUserPreferences(ctx, bobAddr.String(), []restakingtypes.TrustedServiceEntry{
		restakingtypes.NewTrustedServiceEntry(service.ID, []uint32{1}),
	})
	shares, err = suite.keeper.GetPoolServiceTotalDelegatorShares(ctx, 1, service.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(utils.MustParseDecCoins("20_000000pool/1/umilk"), shares)
}
