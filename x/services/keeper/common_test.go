package keeper_test

import (
	"testing"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"go.uber.org/mock/gomock"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	bankkeeper "github.com/milkyway-labs/milkyway/x/bank/keeper"
	"github.com/milkyway-labs/milkyway/x/services/keeper"
	"github.com/milkyway-labs/milkyway/x/services/testutils"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	cdc            codec.Codec
	legacyAminoCdc *codec.LegacyAmino
	ctx            sdk.Context

	storeKey storetypes.StoreKey

	ak    authkeeper.AccountKeeper
	bk    bankkeeper.Keeper
	k     *keeper.Keeper
	hooks *testutils.MockHooks

	ctrl       *gomock.Controller
	poolKeeper *testutils.MockCommunityPoolKeeper
}

func (suite *KeeperTestSuite) SetupTest() {
	data := testutils.SetupKeeperTest(suite.T())

	suite.storeKey = data.StoreKey

	suite.ctx = data.Context
	suite.cdc, suite.legacyAminoCdc = data.Cdc, data.LegacyAmino

	suite.ctrl = data.MockCtrl

	// Build keepers
	suite.ak = data.AccountKeeper
	suite.bk = data.BankKeeper
	suite.k = data.Keeper
	suite.poolKeeper = data.PoolKeeper

	// Set hooks
	suite.hooks = data.Hooks
}

func (suite *KeeperTestSuite) TearDownTest() {
	suite.ctrl.Finish()
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
