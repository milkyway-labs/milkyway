package v122

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/app/upgrades"
	stakeibckeeper "github.com/milkyway-labs/milkyway/x/stakeibc/keeper"
	stakeibctypes "github.com/milkyway-labs/milkyway/x/stakeibc/types"
)

var (
	_ upgrades.Upgrade = &Upgrade{}
)

// Upgrade represents the v1.2.2 upgrade
type Upgrade struct {
	mm             *module.Manager
	configurator   module.Configurator
	stakeIBCKeeper stakeibckeeper.Keeper
}

// NewUpgrade returns a new Upgrade instance
func NewUpgrade(
	mm *module.Manager,
	configurator module.Configurator,
	rewardsKeeper stakeibckeeper.Keeper,
) *Upgrade {
	return &Upgrade{
		mm:             mm,
		configurator:   configurator,
		stakeIBCKeeper: rewardsKeeper,
	}
}

// Name implements upgrades.Upgrade
func (u *Upgrade) Name() string {
	return "v1.2.2"
}

// Handler implements upgrades.Upgrade
func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)

		// Update the unbonding period of the host zone to 21 days.
		hostZone, found := u.stakeIBCKeeper.GetHostZone(sdkCtx, "initiation-2")
		if !found {
			return nil, stakeibctypes.ErrHostZoneNotFound
		}
		hostZone.UnbondingPeriod = 21
		u.stakeIBCKeeper.SetHostZone(sdkCtx, hostZone)

		// This upgrade does not require any migration, so we can simply return the current version map
		return u.mm.RunMigrations(ctx, u.configurator, fromVM)
	}
}

// StoreUpgrades implements upgrades.Upgrade
func (u *Upgrade) StoreUpgrades() *storetypes.StoreUpgrades {
	return &storetypes.StoreUpgrades{
		Added:   nil,
		Renamed: nil,
		Deleted: nil,
	}
}
