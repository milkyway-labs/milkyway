package keeper_test

import (
	"github.com/milkyway-labs/milkyway/utils"
)

func (suite *KeeperTestSuite) TestServicesHooks_AfterServiceDeleted() {
	ctx := suite.ctx

	service, _ := suite.setupSampleServiceAndOperator(ctx)

	err := suite.keeper.IncrementPoolServiceTotalDelegatorShares(
		ctx,
		1,
		service.ID,
		utils.MustParseDecCoins("100000000umilk"),
	)
	suite.Require().NoError(err)
	err = suite.keeper.IncrementPoolServiceTotalDelegatorShares(
		ctx,
		2,
		service.ID,
		utils.MustParseDecCoins("100000000umilk"),
	)
	suite.Require().NoError(err)

	hooks := suite.keeper.ServicesHooks()
	// Calling AfterServiceDeleted will delete all the pool service total delegator
	// shares records for the service
	err = hooks.AfterServiceDeleted(ctx, 1)
	suite.NoError(err)

	shares, err := suite.keeper.GetPoolServiceTotalDelegatorShares(ctx, 1, service.ID)
	suite.Require().NoError(err)
	suite.Require().Nil(shares)
	shares, err = suite.keeper.GetPoolServiceTotalDelegatorShares(ctx, 2, service.ID)
	suite.Require().NoError(err)
	suite.Require().Nil(shares)
}
