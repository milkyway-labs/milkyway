package v710

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/app/upgrades"
	assetstypes "github.com/milkyway-labs/milkyway/x/assets/types"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	rewardskeeper "github.com/milkyway-labs/milkyway/x/rewards/keeper"
	rewardstypes "github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

var _ upgrades.Upgrade = &Upgrade{}

// Upgrade represents the v1.1.0 upgrade
type Upgrade struct {
	mm           *module.Manager
	configurator module.Configurator

	cdc           codec.BinaryCodec
	keys          map[string]*storetypes.KVStoreKey
	rewardsKeeper *rewardskeeper.Keeper
}

// NewUpgrade returns a new Upgrade instance
func NewUpgrade(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.BinaryCodec,
	keys map[string]*storetypes.KVStoreKey,
	rewardsKeeper *rewardskeeper.Keeper,
) *Upgrade {
	return &Upgrade{
		mm:            mm,
		configurator:  configurator,
		cdc:           cdc,
		keys:          keys,
		rewardsKeeper: rewardsKeeper,
	}
}

// Name implements upgrades.Upgrade
func (u *Upgrade) Name() string {
	return "v1.1.0"
}

// Handler implements upgrades.Upgrade
func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// Set the operators module params
		store := sdkCtx.KVStore(u.keys[operatorstypes.ModuleName])
		operatorsParams := operatorstypes.DefaultParams()
		store.Set(operatorstypes.ParamsKey, u.cdc.MustMarshal(&operatorsParams))

		// Set the pools params
		store = sdkCtx.KVStore(u.keys[poolstypes.ModuleName])
		poolsParams := poolstypes.DefaultParams()
		store.Set(poolstypes.ParamsKey, u.cdc.MustMarshal(&poolsParams))

		// Set the restaking params
		store = sdkCtx.KVStore(u.keys[restakingtypes.ModuleName])
		restakingParams := restakingtypes.DefaultParams()
		store.Set(restakingtypes.LegacyParamsKey, u.cdc.MustMarshal(&restakingParams))

		// Set the rewards params
		if err := u.rewardsKeeper.Params.Set(ctx, rewardstypes.DefaultParams()); err != nil {
			return nil, err
		}

		// Set the services params
		store = sdkCtx.KVStore(u.keys[servicestypes.ModuleName])
		servicesParams := servicestypes.DefaultParams()
		store.Set(servicestypes.ParamsKey, u.cdc.MustMarshal(&servicesParams))

		// Set the module versions
		fromVM[assetstypes.ModuleName] = 1
		fromVM[operatorstypes.ModuleName] = 1
		fromVM[poolstypes.ModuleName] = 1
		fromVM[restakingtypes.ModuleName] = 1
		fromVM[rewardstypes.ModuleName] = 1
		fromVM[servicestypes.ModuleName] = 1

		// This upgrade does not require any migration, so we can simply return the current version map
		return u.mm.RunMigrations(ctx, u.configurator, fromVM)
	}
}

// StoreUpgrades implements upgrades.Upgrade
func (u *Upgrade) StoreUpgrades() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Added: []string{
			assetstypes.ModuleName,
			operatorstypes.ModuleName,
			poolstypes.ModuleName,
			restakingtypes.ModuleName,
			rewardstypes.ModuleName,
			servicestypes.ModuleName,
		},
		Renamed: nil,
		Deleted: nil,
	}
}
