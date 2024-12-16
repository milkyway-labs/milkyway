package keeper_test

import (
	"testing"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	oraclekeeper "github.com/skip-mev/connect/v2/x/oracle/keeper"
	oracletypes "github.com/skip-mev/connect/v2/x/oracle/types"
	"github.com/stretchr/testify/suite"

	assetskeeper "github.com/milkyway-labs/milkyway/v7/x/assets/keeper"
	assetstypes "github.com/milkyway-labs/milkyway/v7/x/assets/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/v7/x/operators/keeper"
	poolskeeper "github.com/milkyway-labs/milkyway/v7/x/pools/keeper"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/keeper"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/testutils"
	rewardstypes "github.com/milkyway-labs/milkyway/v7/x/rewards/types"
	serviceskeeper "github.com/milkyway-labs/milkyway/v7/x/services/keeper"
)

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

type KeeperTestSuite struct {
	suite.Suite

	cdc codec.Codec
	ctx sdk.Context

	storeService corestoretypes.KVStoreService

	ak              authkeeper.AccountKeeper
	bk              bankkeeper.BaseKeeper
	pk              *poolskeeper.Keeper
	ok              *operatorskeeper.Keeper
	sk              *serviceskeeper.Keeper
	marketMapKeeper *marketmapkeeper.Keeper
	oracleKeeper    *oraclekeeper.Keeper
	assetsKeeper    *assetskeeper.Keeper
	k               *keeper.Keeper
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
	suite.marketMapKeeper = data.MarketMapKeeper
	suite.oracleKeeper = data.OracleKeeper
	suite.assetsKeeper = data.AssetsKeeper
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

// This code snippet is copied from x/rewards/keeper/common_test.go
// TODO: remove redundant code
// RegisterCurrency registers a currency with the given denomination, ticker
// and price. RegisterCurrency creates a market for the currency if not exists.
func (suite *KeeperTestSuite) RegisterCurrency(ctx sdk.Context, denom string, ticker string, exponent uint32, price math.LegacyDec) {
	// Create the market only if it doesn't exist.
	mmTicker := marketmaptypes.NewTicker(ticker, rewardstypes.USDTicker, math.LegacyPrecision, 0, true)
	hasMarket, err := suite.marketMapKeeper.HasMarket(ctx, mmTicker.String())
	suite.Require().NoError(err)

	if !hasMarket {
		err = suite.marketMapKeeper.CreateMarket(ctx, marketmaptypes.Market{Ticker: mmTicker})
		suite.Require().NoError(err)
	}

	// Set the price for the currency pair.
	err = suite.oracleKeeper.SetPriceForCurrencyPair(
		ctx,
		connecttypes.NewCurrencyPair(ticker, rewardstypes.USDTicker),
		oracletypes.QuotePrice{
			Price:          math.NewIntFromBigInt(price.BigInt()),
			BlockTimestamp: ctx.BlockTime(),
			BlockHeight:    uint64(ctx.BlockHeight()),
		},
	)
	suite.Require().NoError(err)

	// Register the currency.
	err = suite.assetsKeeper.SetAsset(ctx, assetstypes.NewAsset(denom, ticker, exponent))
	suite.Require().NoError(err)
}
