package keeper_test

import (
	"testing"

	corestoretypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/suite"

	operatorskeeper "github.com/milkyway-labs/milkyway/v2/x/operators/keeper"
	poolskeeper "github.com/milkyway-labs/milkyway/v2/x/pools/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v2/x/restaking/testutils"
	serviceskeeper "github.com/milkyway-labs/milkyway/v2/x/services/keeper"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	cdc codec.Codec
	ctx sdk.Context

	storeService corestoretypes.KVStoreService

	ak authkeeper.AccountKeeper
	bk bankkeeper.BaseKeeper
	pk *poolskeeper.Keeper
	ok *operatorskeeper.Keeper
	sk *serviceskeeper.Keeper
	k  *keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())
	suite.ctx = data.Context
	suite.cdc = data.Cdc
	suite.storeService = data.StoreService

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
