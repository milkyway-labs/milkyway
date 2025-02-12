package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/v9/x/assets/keeper"
	"github.com/milkyway-labs/milkyway/v9/x/assets/testutils"
)

type KeeperTestSuite struct {
	suite.Suite

	authority string

	ctx    sdk.Context
	keeper *keeper.Keeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())

	suite.authority = data.AuthorityAddress
	suite.ctx = data.Context

	suite.keeper = data.Keeper
}
