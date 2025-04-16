package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/milkyway-labs/milkyway/v10/app/testutil"
	"github.com/milkyway-labs/milkyway/v10/utils"
	"github.com/milkyway-labs/milkyway/v10/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v10/x/restaking/types"
)

func (suite *KeeperTestSuite) TestRestakingHooks_UpdateTargetCoveredLockedSharesAfterDelegation() {
	ctx, _ := suite.ctx.CacheContext()

	err := suite.k.SetParams(ctx, types.NewParams(sdkmath.LegacyNewDec(1), nil, nil, nil)) // 1%
	suite.Require().NoError(err)

	suite.createService(ctx, 1)
	suite.createService(ctx, 2)

	delAddr := testutil.TestAddress(1)
	suite.mintLockedRepresentation(ctx, delAddr.String(), utils.MustParseCoins("1000_000000stake"))
	suite.fundAccountInsuranceFund(ctx, delAddr.String(), utils.MustParseCoins("1_000000stake")) // 1%

	_, err = suite.restakingKeeper.DelegateToService(ctx, 1, utils.MustParseCoins("100_000000locked/stake"), delAddr.String())
	suite.Require().NoError(err)
	_, err = suite.restakingKeeper.DelegateToService(ctx, 2, utils.MustParseCoins("200_000000locked/stake"), delAddr.String())
	suite.Require().NoError(err)

	targetCoveredLockedShares, err := suite.k.GetTargetCoveredLockedShares(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, 1)
	suite.Require().NoError(err)
	suite.Assert().Equal("33333333.000000000000000000service/1/locked/stake", targetCoveredLockedShares.String())
	targetCoveredLockedShares, err = suite.k.GetTargetCoveredLockedShares(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, 2)
	suite.Require().NoError(err)
	suite.Assert().Equal("66666666.000000000000000000service/2/locked/stake", targetCoveredLockedShares.String())

	ark := suite.k.AdjustedRestakingKeeper(suite.restakingKeeper)
	target, err := ark.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, 1)
	suite.Require().NoError(err)
	suite.Assert().Equal("33333333locked/stake", target.GetTokens().String())
	suite.Assert().Equal("33333333.000000000000000000service/1/locked/stake", target.GetDelegatorShares().String())
	del, _, err := ark.GetDelegation(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, 1, delAddr.String())
	suite.Require().NoError(err)
	suite.Assert().Equal("33333333.000000000000000000service/1/locked/stake", del.Shares.String())

	target, err = ark.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, 2)
	suite.Require().NoError(err)
	suite.Assert().Equal("66666666locked/stake", target.GetTokens().String())
	suite.Assert().Equal("66666666.000000000000000000service/2/locked/stake", target.GetDelegatorShares().String())
	del, _, err = ark.GetDelegation(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, 2, delAddr.String())
	suite.Require().NoError(err)
	suite.Assert().Equal("66666666.000000000000000000service/2/locked/stake", del.Shares.String())
}
