package keeper_test

import (
	"context"
	"testing"
	"time"

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
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"

	distrkeeper "github.com/milkyway-labs/milkyway/v9/x/distribution/keeper"
	"github.com/milkyway-labs/milkyway/v9/x/investors"
	"github.com/milkyway-labs/milkyway/v9/x/investors/keeper"
	"github.com/milkyway-labs/milkyway/v9/x/investors/testutils"
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
	sk *stakingkeeper.Keeper
	dk distrkeeper.Keeper
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
	suite.sk = data.StakingKeeper
	suite.dk = data.DistrKeeper
	suite.k = data.Keeper
}

// fundAccount adds the given amount of coins to the account with the given
// address.
func (suite *KeeperTestSuite) fundAccount(ctx context.Context, address string, amount sdk.Coins) {
	// Mint the coins
	err := suite.bk.MintCoins(ctx, minttypes.ModuleName, amount)
	suite.Require().NoError(err)

	// Send the amount to the user
	userAddress, err := suite.ak.AddressCodec().StringToBytes(address)
	suite.Require().NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, userAddress, amount)
	suite.Require().NoError(err)
}

// fundModuleAccount adds the given amount of coins to the module account with
// the given name.
func (suite *KeeperTestSuite) fundModuleAccount(ctx context.Context, moduleName string, amount sdk.Coins) {
	// Mint the coins
	err := suite.bk.MintCoins(ctx, minttypes.ModuleName, amount)
	suite.Require().NoError(err)

	// Send the amount to the module
	err = suite.bk.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, moduleName, amount)
	suite.Require().NoError(err)
}

// allocateTokensToValidator allocates the given amount of tokens to the given
// validator.
func (suite *KeeperTestSuite) allocateTokensToValidator(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	tokens sdk.DecCoins,
	fund bool,
) sdk.Context {
	if fund {
		// We need to ceil the dec coins to fund
		coins := sdk.NewCoins()
		for _, token := range tokens {
			coins = coins.Add(sdk.NewCoin(token.Denom, token.Amount.Ceil().TruncateInt()))
		}
		suite.fundModuleAccount(ctx, distrtypes.ModuleName, coins)
	}
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(5 * time.Second))
	// We get 'adjusted' validator since that is what the distribution module uses
	validator, err := suite.sk.Validator(ctx, valAddr)
	suite.Require().NoError(err)
	err = suite.dk.AllocateTokensToValidator(ctx, validator, tokens)
	suite.Require().NoError(err)
	return ctx
}

// createValidator creates a validator with the given owner, commission rates,
// initial delegation value.
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
	_, err = stakingkeeper.NewMsgServerImpl(suite.sk).CreateValidator(ctx, msg)
	suite.Require().NoError(err)

	validator, err := suite.sk.Validator(ctx, valAddr)
	suite.Require().NoError(err)
	return validator
}

// delegate creates a delegation from the given delegator to the given validator.
func (suite *KeeperTestSuite) delegate(ctx context.Context, delegator, validator string, amount sdk.Coin, fund bool) {
	if fund {
		suite.fundAccount(ctx, delegator, sdk.NewCoins(amount))
	}

	_, err := stakingkeeper.NewMsgServerImpl(suite.sk).Delegate(
		ctx,
		stakingtypes.NewMsgDelegate(delegator, validator, amount),
	)
	suite.Require().NoError(err)
}

// createVestingAccount creates a vesting account with the given parameters.
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

func (suite *KeeperTestSuite) withdrawRewardsAndIncrementBlockHeight(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, sdk.Context) {
	balancesBefore := suite.bk.GetAllBalances(ctx, delAddr)
	_, err := suite.dk.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	suite.Require().NoError(err)
	err = investors.EndBlocker(ctx, suite.k)
	suite.Require().NoError(err)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1).WithBlockTime(ctx.BlockTime().Add(5 * time.Second))
	balancesAfter := suite.bk.GetAllBalances(ctx, delAddr)
	rewards := balancesAfter.Sub(balancesBefore...)
	return rewards, ctx
}
