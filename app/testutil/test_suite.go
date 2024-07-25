package testutil

import (
	"time"

	"cosmossdk.io/core/header"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	milkywayapp "github.com/milkyway-labs/milkyway/app"
)

type KeeperTestSuite struct {
	suite.Suite

	App *milkywayapp.MilkyWayApp
	Ctx sdk.Context
}

func (s *KeeperTestSuite) SetupTest() {
	s.App = milkywayapp.SetupWithGenesisAccounts(nil, nil)
	s.Ctx = s.App.NewContext(true).WithHeaderInfo(header.Info{
		Height: s.App.LastBlockHeight(),
		Time:   time.Time{},
	})
}
