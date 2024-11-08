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
	operatorskeeper "github.com/milkyway-labs/milkyway/x/operators/keeper"
	poolskeeper "github.com/milkyway-labs/milkyway/x/pools/keeper"
	"github.com/milkyway-labs/milkyway/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/x/restaking/testutils"
	serviceskeeper "github.com/milkyway-labs/milkyway/x/services/keeper"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	cdc codec.Codec
	ctx sdk.Context

	storeKey storetypes.StoreKey

	ak authkeeper.AccountKeeper
	bk bankkeeper.Keeper
	pk *poolskeeper.Keeper
	ok *operatorskeeper.Keeper
	sk *serviceskeeper.Keeper
	k  *keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	// Build the base data
	data := testutils.NewKeeperTestData(suite.T())

	// Define store keys
	suite.storeKey = data.StoreKey

	// Define the codec and context
	suite.ctx = data.Context
	suite.cdc = data.Cdc

	// Build keepers
	suite.ak = data.AccountKeeper
	suite.bk = data.BankKeeper
	suite.pk = data.PoolsKeeper
	suite.ok = data.OperatorsKeeper
	suite.sk = data.ServicesKeeper
	suite.k = data.Keeper
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

// --------------------------------------------------------------------------------------------------------------------
