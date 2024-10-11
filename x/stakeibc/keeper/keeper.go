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

	icacallbackskeeper "github.com/milkyway-labs/milkyway/x/icacallbacks/keeper"
	icqkeeper "github.com/milkyway-labs/milkyway/x/interchainquery/keeper"
	recordsmodulekeeper "github.com/milkyway-labs/milkyway/x/records/keeper"
	"github.com/milkyway-labs/milkyway/x/stakeibc/types"
)

type (
	Keeper struct {
		cdc                   codec.Codec
		storeKey              storetypes.StoreKey
		memKey                storetypes.StoreKey
		authority             string
		icaControllerKeeper   icacontrollerkeeper.Keeper
		ibcKeeper             ibckeeper.Keeper
		bankKeeper            bankkeeper.Keeper
		accountKeeper         types.AccountKeeper
		interchainQueryKeeper icqkeeper.Keeper
		recordsKeeper         recordsmodulekeeper.Keeper
		icaCallbacksKeeper    icacallbackskeeper.Keeper
		hooks                 types.StakeIBCHooks
		rateLimitKeeper       types.RatelimitKeeper
		opChildKeeper         types.OPChildKeeper
		params                collections.Item[types.Params]
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	storeService corestoretypes.KVStoreService,
	authority string,
	accountKeeper types.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	icaControllerKeeper icacontrollerkeeper.Keeper,
	ibcKeeper ibckeeper.Keeper,
	interchainQueryKeeper icqkeeper.Keeper,
	RecordsKeeper recordsmodulekeeper.Keeper,
	icaCallbacksKeeper icacallbackskeeper.Keeper,
	rateLimitKeeper types.RatelimitKeeper,
	opChildKeeper types.OPChildKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	return Keeper{
		cdc:                   cdc,
		storeKey:              storeKey,
		memKey:                memKey,
		authority:             authority,
		accountKeeper:         accountKeeper,
		bankKeeper:            bankKeeper,
		icaControllerKeeper:   icaControllerKeeper,
		ibcKeeper:             ibcKeeper,
		interchainQueryKeeper: interchainQueryKeeper,
		recordsKeeper:         RecordsKeeper,
		icaCallbacksKeeper:    icaCallbacksKeeper,
		rateLimitKeeper:       rateLimitKeeper,
		opChildKeeper:         opChildKeeper,
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

// ValidateAdminAddress makes sure that the given address is the one of an account that can perform admin operations.
func (k Keeper) ValidateAdminAddress(ctx sdk.Context, address string) error {
	// The authority can always perform admin operations.
	if k.GetAuthority() == address {
		return nil
	}

	// The OpChild admin can always perform admin operations
	opChildParams, err := k.opChildKeeper.GetParams(ctx)
	if err != nil {
		return err
	}

	if opChildParams.Admin == address {
		return nil
	}

	return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "address (%s) is not an admin", address)
}

func (k Keeper) GetICATimeoutNanos(ctx sdk.Context, epochType string) (uint64, error) {
	epochTracker, found := k.GetEpochTracker(ctx, epochType)
	if !found {
		k.Logger(ctx).Error(fmt.Sprintf("Failed to get epoch tracker for %s", epochType))
		return 0, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "Failed to get epoch tracker for %s", epochType)
	}
	// BUFFER by 5% of the epoch length
	bufferSizeParam := k.GetParams(ctx).BufferSize
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
	params := k.GetParams(ctx)
	minSafetyThresholdInt := params.DefaultMinRedemptionRateThreshold
	minSafetyThreshold := sdkmath.LegacyNewDec(int64(minSafetyThresholdInt)).Quo(sdkmath.LegacyNewDec(100))

	if !zone.MinRedemptionRate.IsNil() && zone.MinRedemptionRate.IsPositive() {
		minSafetyThreshold = zone.MinRedemptionRate
	}

	maxSafetyThresholdInt := params.DefaultMaxRedemptionRateThreshold
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
