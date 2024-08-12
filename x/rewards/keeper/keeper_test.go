package keeper_test

import (
	"testing"
	"time"

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

func (s *KeeperTestSuite) SetupTest() {
	s.KeeperTestSuite.SetupTest()
	s.authority = s.App.TickersKeeper.GetAuthority()
	s.keeper = s.App.RewardsKeeper
	s.msgServer = keeper.NewMsgServer(s.keeper)
	s.queryServer = keeper.NewQueryServer(s.keeper)
}

func (s *KeeperTestSuite) allocateRewards(duration time.Duration) {
	s.Ctx = s.Ctx.WithBlockTime(s.Ctx.BlockTime().Add(duration)).WithBlockHeight(s.Ctx.BlockHeight() + 1)
	err := s.keeper.AllocateRewards(s.Ctx)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) setupSampleServiceAndOperator() (servicestypes.Service, operatorstypes.Operator) {
	// This helper method:
	// - registers $MILK, $INIT
	// - creates a service named "MilkyWay"
	// - creates a rewards plan with basic distribution types
	//   - it distributes 100 $MILK every day
	// - creates an operator named "MilkyWay Operator"
	//   - it has 10% commission rate
	//   - it joins the newly created service

	// Register $MILK and $INIT.
	s.RegisterCurrency("umilk", "MILK", utils.MustParseDec("2"))
	s.RegisterCurrency("uinit", "INIT", utils.MustParseDec("3"))

	// Create a service.
	serviceAdmin := utils.TestAddress(10000)
	service := s.CreateService("Service", serviceAdmin.String())

	// Add the created service ID to the pools module's allowed list.
	poolsParams := s.App.PoolsKeeper.GetParams(s.Ctx)
	poolsParams.AllowedServicesIDs = []uint32{service.ID}
	s.App.PoolsKeeper.SetParams(s.Ctx, poolsParams)

	// Create an operator.
	operatorAdmin := utils.TestAddress(10001)
	operator := s.CreateOperator("Operator", operatorAdmin.String())
	// Make the operator join the service and set its commission rate to 10%.
	s.UpdateOperatorParams(operator.ID, utils.MustParseDec("0.1"), []uint32{service.ID})

	// Call AllocateRewards to set last rewards allocation time.
	err := s.keeper.AllocateRewards(s.Ctx)
	s.Require().NoError(err)

	return service, operator
}
