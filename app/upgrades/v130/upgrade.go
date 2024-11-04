package v122

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/app/upgrades"
	liquidvestingkeeper "github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
)

var (
	_ upgrades.Upgrade = &Upgrade{}
)

// Upgrade represents the v1.3.0 upgrade
type Upgrade struct {
	mm           *module.Manager
	configurator module.Configurator

	liquidVestingKeeper *liquidvestingkeeper.Keeper
}

// NewUpgrade returns a new Upgrade instance
func NewUpgrade(
	mm *module.Manager,
	configurator module.Configurator,
	liquidVestingKeeper *liquidvestingkeeper.Keeper,
) *Upgrade {
	return &Upgrade{
		mm:           mm,
		configurator: configurator,

		liquidVestingKeeper: liquidVestingKeeper,
	}
}

// Name implements upgrades.Upgrade
func (u *Upgrade) Name() string {
	return "v1.3.0"
}

// Handler implements upgrades.Upgrade
func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
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
