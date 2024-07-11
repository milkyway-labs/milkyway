package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/app/testutil"
	"github.com/milkyway-labs/milkyway/x/rewards/keeper"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

type KeeperTestSuite struct {
	testutil.KeeperTestSuite

	authority   string
	msgServer   types.MsgServer
	queryServer types.QueryServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.KeeperTestSuite.SetupTest()
	s.authority = s.App.TickersKeeper.GetAuthority()
	s.msgServer = keeper.NewMsgServer(s.App.RewardsKeeper)
	s.queryServer = keeper.NewQueryServer(s.App.RewardsKeeper)
}
