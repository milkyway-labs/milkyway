package v9

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v7/app/upgrades"
)

const UpgradeName = "v9"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
