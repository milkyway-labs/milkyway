package v710

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/milkyway-labs/milkyway/app/upgrades"
	assetstypes "github.com/milkyway-labs/milkyway/x/assets/types"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	rewardstypes "github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

var (
	_ upgrades.Upgrade = &Upgrade{}
)

// Upgrade represents the v1.1.0 upgrade
type Upgrade struct {
	mm           *module.Manager
	configurator module.Configurator
}

// NewUpgrade returns a new Upgrade instance
func NewUpgrade(mm *module.Manager, configurator module.Configurator) *Upgrade {
	return &Upgrade{
		mm:           mm,
		configurator: configurator,
	}
}

// Name implements upgrades.Upgrade
func (u *Upgrade) Name() string {
	return "v1.1.0"
}

// Handler implements upgrades.Upgrade
func (u *Upgrade) Handler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Set the consensus version to 1 in order to run migrations from 1 to 2 for the following modules
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
