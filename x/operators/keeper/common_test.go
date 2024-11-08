package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"

	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	"github.com/milkyway-labs/milkyway/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/x/operators/testutils"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	cdc codec.Codec
	ctx sdk.Context

	storeKey storetypes.StoreKey

	ak    authkeeper.AccountKeeper
	bk    bankkeeper.Keeper
	k     *keeper.Keeper
	hooks *testutils.MockHooks
}

func (suite *KeeperTestSuite) SetupTest() {
	// Build the testing data
	data := testutils.NewKeeperTestData(suite.T())

	// Define the codec and context
	suite.cdc = data.Cdc
	suite.ctx = data.Context

	// Define store keys
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
	// Mint the coins
	moduleAcc := suite.ak.GetModuleAccount(ctx, authtypes.Minter)

	err := suite.bk.MintCoins(ctx, moduleAcc.GetName(), amount)
	suite.Require().NoError(err)

	// Get the amount to the user
	userAddress, err := sdk.AccAddressFromBech32(address)
	suite.Require().NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(ctx, moduleAcc.GetName(), userAddress, amount)
	suite.Require().NoError(err)
}
