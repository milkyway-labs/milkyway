package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/app/testutil"
	"github.com/milkyway-labs/milkyway/x/assets/keeper"
	"github.com/milkyway-labs/milkyway/x/assets/types"
)

type KeeperTestSuite struct {
	testutil.KeeperTestSuite

	authority string

	keeper      *keeper.Keeper
	msgServer   types.MsgServer
	queryServer types.QueryServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.KeeperTestSuite.SetupTest()
	suite.authority = suite.App.AssetsKeeper.GetAuthority()

	suite.keeper = suite.App.AssetsKeeper
	suite.msgServer = keeper.NewMsgServer(suite.App.AssetsKeeper)
	suite.queryServer = keeper.NewQueryServer(suite.App.AssetsKeeper)
}
