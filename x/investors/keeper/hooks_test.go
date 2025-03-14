package keeper_test

import (
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/milkyway-labs/milkyway/v10/app/testutil"
	distrkeeper "github.com/milkyway-labs/milkyway/v10/x/distribution/keeper"
)

func (suite *KeeperTestSuite) TestDistrHooks_AfterSetWithdrawAddress() {
	ctx, _ := suite.ctx.CacheContext()

	delAddr := testutil.TestAddress(1)
	withdrawAddr := testutil.TestAddress(2)

	// When the withdraw address is not set, the delegator address is the withdraw
	// address
	delegator, err := suite.k.GetDelegatorAddressByWithdrawAddress(ctx, delAddr.String())
	suite.Require().NoError(err)
	suite.Assert().Equal(delAddr.String(), delegator)

	msgServer := distrkeeper.NewMsgServerImpl(suite.dk)
	_, err = msgServer.SetWithdrawAddress(ctx, distrtypes.NewMsgSetWithdrawAddress(delAddr, withdrawAddr))
	suite.Require().NoError(err)

	// After setting the withdraw address, it should be possible to lookup the
	// delegator address by the withdraw address
	delegator, err = suite.k.GetDelegatorAddressByWithdrawAddress(ctx, withdrawAddr.String())
	suite.Require().NoError(err)
	suite.Assert().Equal(delAddr.String(), delegator)
}
