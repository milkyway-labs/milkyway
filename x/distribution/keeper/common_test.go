package keeper_test

import (
	"context"
	"testing"

	corestoretypes "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	marketmapkeeper "github.com/skip-mev/connect/v2/x/marketmap/keeper"
	oraclekeeper "github.com/skip-mev/connect/v2/x/oracle/keeper"
	"github.com/stretchr/testify/suite"

	assetskeeper "github.com/milkyway-labs/milkyway/v7/x/assets/keeper"
	"github.com/milkyway-labs/milkyway/v7/x/distribution/keeper"
	"github.com/milkyway-labs/milkyway/v7/x/distribution/testutils"
	"github.com/milkyway-labs/milkyway/v7/x/distribution/types"
	operatorskeeper "github.com/milkyway-labs/milkyway/v7/x/operators/keeper"
	poolskeeper "github.com/milkyway-labs/milkyway/v7/x/pools/keeper"
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
	stakingKeeper   *stakingkeeper.Keeper
	pk              *poolskeeper.Keeper
	ok              *operatorskeeper.Keeper
	sk              *serviceskeeper.Keeper
	marketMapKeeper *marketmapkeeper.Keeper
	oracleKeeper    *oraclekeeper.Keeper
	assetsKeeper    *assetskeeper.Keeper
	k               keeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	data := testutils.NewKeeperTestData(suite.T())
	suite.ctx = data.Context
	suite.cdc = data.Cdc
	suite.storeService = data.StoreService

	// Build keepers
	suite.ak = data.AccountKeeper
	suite.bk = data.BankKeeper
	suite.stakingKeeper = data.StakingKeeper
	suite.pk = data.PoolsKeeper
	suite.ok = data.OperatorsKeeper
	suite.sk = data.ServicesKeeper
	suite.marketMapKeeper = data.MarketMapKeeper
	suite.oracleKeeper = data.OracleKeeper
	suite.assetsKeeper = data.AssetsKeeper
	suite.k = data.Keeper

	err := suite.stakingKeeper.SetParams(suite.ctx, stakingtypes.DefaultParams())
	suite.Require().NoError(err)

	// Reset to the default(50%) for every test
	types.VestingAccountRewardsRatio = sdkmath.LegacyNewDecWithPrec(5, 1)
}

// fundAccount adds the given amount of coins to the account with the given address
func (suite *KeeperTestSuite) fundAccount(ctx context.Context, address string, amount sdk.Coins) {
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

func (suite *KeeperTestSuite) createValidator(
	ctx context.Context,
	owner sdk.AccAddress,
	commissionRates stakingtypes.CommissionRates,
	value sdk.Coin,
	fund bool,
) stakingtypes.ValidatorI {
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	valAddr := sdk.ValAddress(owner)
	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	suite.Require().NoError(err)

	if fund {
		suite.fundAccount(ctx, owner.String(), sdk.NewCoins(value))
	}

	msg := &stakingtypes.MsgCreateValidator{
		Description:       stakingtypes.Description{Moniker: "Validator"},
		Commission:        commissionRates,
		MinSelfDelegation: sdkmath.OneInt(),
		ValidatorAddress:  valAddr.String(),
		Pubkey:            pkAny,
		Value:             value,
	}
	_, err = stakingkeeper.NewMsgServerImpl(suite.stakingKeeper).CreateValidator(ctx, msg)
	suite.Require().NoError(err)

	validator, err := suite.stakingKeeper.Validator(ctx, valAddr)
	suite.Require().NoError(err)
	return validator
}

func (suite *KeeperTestSuite) delegate(ctx context.Context, delegator, validator string, amount sdk.Coin, fund bool) {
	if fund {
		suite.fundAccount(ctx, delegator, sdk.NewCoins(amount))
	}

	_, err := stakingkeeper.NewMsgServerImpl(suite.stakingKeeper).Delegate(ctx, &stakingtypes.MsgDelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: validator,
		Amount:           amount,
	})
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) createVestingAccount(ctx context.Context, from, to string, amount sdk.Coins, endTime int64, delayed, fund bool) {
	if fund {
		suite.fundAccount(ctx, from, amount)
	}
	msg := &vestingtypes.MsgCreateVestingAccount{
		FromAddress: from,
		ToAddress:   to,
		Amount:      amount,
		EndTime:     endTime,
		Delayed:     delayed,
	}
	_, err := vesting.NewMsgServerImpl(suite.ak, suite.bk).CreateVestingAccount(ctx, msg)
	suite.Require().NoError(err)
}
