package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/milkyway/v2/x/services/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/services/testutils"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	storeKey storetypes.StoreKey

	ak    authkeeper.AccountKeeper
	bk    bankkeeper.BaseKeeper
	k     *keeper.Keeper
	hooks *testutils.MockHooks
}

func (suite *KeeperTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())
	suite.ctx = data.Context
	suite.storeKey = data.StoreKey

	// Build keepers
	suite.ak = data.AccountKeeper
	suite.bk = data.BankKeeper
	suite.k = data.Keeper

	// Set hooks
	suite.hooks = data.Hooks
}

// --------------------------------------------------------------------------------------------------------------------

// fundAccount adds the given amount of coins to the account with the given address
func (suite *KeeperTestSuite) fundAccount(ctx sdk.Context, address string, amount sdk.Coins) {
	moduleAcc := suite.ak.GetModuleAccount(ctx, minttypes.ModuleName)

	// Mint the coins
	err := suite.bk.MintCoins(ctx, moduleAcc.GetName(), amount)
	suite.Require().NoError(err)

	// Get the amount to the user
	userAddress, err := sdk.AccAddressFromBech32(address)
	suite.Require().NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(ctx, moduleAcc.GetName(), userAddress, amount)
	suite.Require().NoError(err)
}
