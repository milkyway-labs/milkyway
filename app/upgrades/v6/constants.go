package v6

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v9/app/upgrades"
)

const UpgradeName = "v6"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
