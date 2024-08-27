package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/app/testutil"
	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type KeeperTestSuite struct {
	testutil.KeeperTestSuite

	authority   string
	keeper      *keeper.Keeper
	msgServer   types.MsgServer
	queryServer types.QueryServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.KeeperTestSuite.SetupTest()
	suite.authority = suite.App.RewardsKeeper.GetAuthority()
	suite.keeper = suite.App.RewardsKeeper
	suite.msgServer = keeper.NewMsgServer(suite.keeper)
	suite.queryServer = keeper.NewQueryServer(suite.keeper)
}

func (suite *KeeperTestSuite) allocateRewards(ctx sdk.Context, duration time.Duration) sdk.Context {
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(duration)).WithBlockHeight(ctx.BlockHeight() + 1)
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)
	return ctx
}

func (suite *KeeperTestSuite) setupSampleServiceAndOperator(ctx sdk.Context) (servicestypes.Service, operatorstypes.Operator) {
	// This helper method:
	// - registers $MILK, $INIT
	// - creates a service named "MilkyWay"
	// - creates a rewards plan with basic distribution types
	//   - it distributes 100 $MILK every day
	// - creates an operator named "MilkyWay Operator"
	//   - it has 10% commission rate
	//   - it joins the newly created service

	// Register $MILK and $INIT.
	suite.RegisterCurrency(ctx, "umilk", "MILK", 6, utils.MustParseDec("2"))
	suite.RegisterCurrency(ctx, "uinit", "INIT", 6, utils.MustParseDec("3"))

	// Create a service.
	serviceAdmin := testutil.TestAddress(10000)
	service := suite.CreateService(ctx, "Service", serviceAdmin.String())

	// Add the created service ID to the pools module's allowed list.
	poolsParams := suite.App.PoolsKeeper.GetParams(ctx)
	poolsParams.AllowedServicesIDs = []uint32{service.ID}
	suite.App.PoolsKeeper.SetParams(ctx, poolsParams)

	// Create an operator.
	operatorAdmin := testutil.TestAddress(10001)
	operator := suite.CreateOperator(ctx, "Operator", operatorAdmin.String())

	// Make the operator join the service and set its commission rate to 10%.
	suite.UpdateOperatorParams(ctx, operator.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Call AllocateRewards to set last rewards allocation time.
	err := suite.keeper.AllocateRewards(ctx)
	suite.Require().NoError(err)

	return service, operator
}
