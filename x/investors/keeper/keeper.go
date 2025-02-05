package keeper

import (
	"context"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v7/x/investors/types"
)

type Keeper struct {
	cdc codec.BinaryCodec

	accountKeeper types.AccountKeeper
	stakingKeeper types.StakingKeeper
	distrKeeper   types.DistrKeeper

	Schema                    collections.Schema
	InvestorsRewardRatio      collections.Item[sdkmath.LegacyDec]
	ValidatorsInvestorsShares collections.Map[sdk.ValAddress, sdkmath.LegacyDec]

	// authority represents the address capable of executing a governance message.
	// Typically, this should be the x/gov module account.
	authority string
}

// NewKeeper creates a new keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	stakingKeeper types.StakingKeeper,
	distrKeeper types.DistrKeeper,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		cdc:           cdc,
		accountKeeper: accountKeeper,
		stakingKeeper: stakingKeeper,
		distrKeeper:   distrKeeper,
		authority:     authority,

		InvestorsRewardRatio: collections.NewItem(
			sb,
			types.InvestorsRewardRatioKey,
			"investors_reward_ratio",
			sdk.LegacyDecValue,
		),
		ValidatorsInvestorsShares: collections.NewMap(
			sb,
			types.ValidatorsInvestorSharesKeyPrefix,
			"validators_investors_shares",
			sdk.ValAddressKey,
			sdk.LegacyDecValue,
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
