package keeper

import (
	"context"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/investors/types"
)

type Keeper struct {
	cdc codec.Codec

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	distrKeeper   types.DistrKeeper

	Schema               collections.Schema
	InvestorsRewardRatio collections.Item[sdkmath.LegacyDec]
	// (vesting end time(in unix seconds), investor address)
	InvestorsVestingQueue collections.KeySet[collections.Pair[int64, string]]
	// Set of investors that are still in their vesting period
	VestingInvestors collections.KeySet[string]

	// authority represents the address capable of executing a governance message.
	// Typically, this should be the x/gov module account.
	authority string
}

// NewKeeper creates a new keeper
func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distrKeeper types.DistrKeeper,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
		distrKeeper:   distrKeeper,
		authority:     authority,

		InvestorsRewardRatio: collections.NewItem(sb, types.InvestorsRewardRatioKey, "investors_reward_ratio", sdk.LegacyDecValue),
		InvestorsVestingQueue: collections.NewKeySet(
			sb,
			types.InvestorsVestingQueueKeyPrefix,
			"investors_vesting_queue",
			collections.PairKeyCodec(collections.Int64Key, collections.StringKey),
		),
		VestingInvestors: collections.NewKeySet(
			sb,
			types.VestingInvestorsKeyPrefix,
			"vesting_investors",
			collections.StringKey,
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}
