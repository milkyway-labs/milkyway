package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/spf13/cast"

	icacallbackskeeper "github.com/milkyway-labs/milk/x/icacallbacks/keeper"
	icqkeeper "github.com/milkyway-labs/milk/x/interchainquery/keeper"
	recordsmodulekeeper "github.com/milkyway-labs/milk/x/records/keeper"
	"github.com/milkyway-labs/milk/x/stakeibc/types"
)

type (
	Keeper struct {
		// *cosmosibckeeper.Keeper
		cdc                   codec.Codec
		storeService          corestoretypes.KVStoreService
		storeKey              storetypes.StoreKey
		memKey                storetypes.StoreKey
		authority             string
		ICAControllerKeeper   icacontrollerkeeper.Keeper
		IBCKeeper             ibckeeper.Keeper
		bankKeeper            bankkeeper.Keeper
		AccountKeeper         types.AccountKeeper
		InterchainQueryKeeper icqkeeper.Keeper
		RecordsKeeper         recordsmodulekeeper.Keeper
		ICACallbacksKeeper    icacallbackskeeper.Keeper
		hooks                 types.StakeIBCHooks
		RatelimitKeeper       types.RatelimitKeeper
		params                collections.Item[types.Params]
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	storeKey,
	memKey storetypes.StoreKey,
	authority string,
	accountKeeper types.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	icacontrollerkeeper icacontrollerkeeper.Keeper,
	ibcKeeper ibckeeper.Keeper,
	interchainQueryKeeper icqkeeper.Keeper,
	RecordsKeeper recordsmodulekeeper.Keeper,
	ICACallbacksKeeper icacallbackskeeper.Keeper,
	RatelimitKeeper types.RatelimitKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	return Keeper{
		cdc:                   cdc,
		storeService:          storeService,
		storeKey:              storeKey,
		memKey:                memKey,
		authority:             authority,
		AccountKeeper:         accountKeeper,
		bankKeeper:            bankKeeper,
		ICAControllerKeeper:   icacontrollerkeeper,
		IBCKeeper:             ibcKeeper,
		InterchainQueryKeeper: interchainQueryKeeper,
		RecordsKeeper:         RecordsKeeper,
		ICACallbacksKeeper:    ICACallbacksKeeper,
		RatelimitKeeper:       RatelimitKeeper,
		params:                collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks sets the hooks for ibc staking
func (k *Keeper) SetHooks(gh types.StakeIBCHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set ibc staking hooks twice")
	}

	k.hooks = gh

	return k
}

// GetAuthority returns the x/stakeibc module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) GetICATimeoutNanos(ctx sdk.Context, epochType string) (uint64, error) {
	epochTracker, found := k.GetEpochTracker(ctx, epochType)
	if !found {
		k.Logger(ctx).Error(fmt.Sprintf("Failed to get epoch tracker for %s", epochType))
		return 0, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Failed to get epoch tracker for %s", epochType)
	}
	// BUFFER by 5% of the epoch length
	bufferSizeParam := k.GetParam(ctx, types.KeyBufferSize)
	bufferSize := epochTracker.Duration / bufferSizeParam
	// buffer size should not be negative or longer than the epoch duration
	if bufferSize > epochTracker.Duration {
		k.Logger(ctx).Error(fmt.Sprintf("Invalid buffer size %d", bufferSize))
		return 0, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid buffer size %d", bufferSize)
	}
	timeoutNanos := epochTracker.NextEpochStartTime - bufferSize
	timeoutNanosUint64, err := cast.ToUint64E(timeoutNanos)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("Failed to convert timeoutNanos to uint64, error: %s", err.Error()))
		return 0, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Failed to convert timeoutNanos to uint64, error: %s", err.Error())
	}
	return timeoutNanosUint64, nil
}

func (k Keeper) GetOuterSafetyBounds(ctx sdk.Context, zone types.HostZone) (sdkmath.LegacyDec, sdkmath.LegacyDec) {
	// Fetch the wide bounds
	minSafetyThresholdInt := k.GetParam(ctx, types.KeyDefaultMinRedemptionRateThreshold)
	minSafetyThreshold := sdkmath.LegacyNewDec(int64(minSafetyThresholdInt)).Quo(sdkmath.LegacyNewDec(100))

	if !zone.MinRedemptionRate.IsNil() && zone.MinRedemptionRate.IsPositive() {
		minSafetyThreshold = zone.MinRedemptionRate
	}

	maxSafetyThresholdInt := k.GetParam(ctx, types.KeyDefaultMaxRedemptionRateThreshold)
	maxSafetyThreshold := sdkmath.LegacyNewDec(int64(maxSafetyThresholdInt)).Quo(sdkmath.LegacyNewDec(100))

	if !zone.MaxRedemptionRate.IsNil() && zone.MaxRedemptionRate.IsPositive() {
		maxSafetyThreshold = zone.MaxRedemptionRate
	}

	return minSafetyThreshold, maxSafetyThreshold
}

func (k Keeper) GetInnerSafetyBounds(ctx sdk.Context, zone types.HostZone) (sdkmath.LegacyDec, sdkmath.LegacyDec) {
	// Fetch the inner bounds
	minSafetyThreshold, maxSafetyThreshold := k.GetOuterSafetyBounds(ctx, zone)

	if !zone.MinInnerRedemptionRate.IsNil() && zone.MinInnerRedemptionRate.IsPositive() && zone.MinInnerRedemptionRate.GT(minSafetyThreshold) {
		minSafetyThreshold = zone.MinInnerRedemptionRate
	}
	if !zone.MaxInnerRedemptionRate.IsNil() && zone.MaxInnerRedemptionRate.IsPositive() && zone.MaxInnerRedemptionRate.LT(maxSafetyThreshold) {
		maxSafetyThreshold = zone.MaxInnerRedemptionRate
	}

	return minSafetyThreshold, maxSafetyThreshold
}
