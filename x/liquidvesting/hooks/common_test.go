package hooks_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	ibchooks "github.com/initia-labs/initia/x/ibc-hooks"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/testutils"
)

func TestIBCHooksTestSuite(t *testing.T) {
	suite.Run(t, new(IBCHooksTestSuite))
}

type IBCHooksTestSuite struct {
	suite.Suite

	ctx sdk.Context

	ak   authkeeper.AccountKeeper
	ibcm ibchooks.IBCMiddleware

	k *keeper.Keeper
}

func (suite *IBCHooksTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())

	// Context and codecs
	suite.ctx = data.Context

	// Keepers
	suite.ak = data.AccountKeeper
	suite.k = data.Keeper
	suite.ibcm = data.IBCMiddleware
}
