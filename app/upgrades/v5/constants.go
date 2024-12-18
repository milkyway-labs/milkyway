package v5

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v5/app/upgrades"
)

const UpgradeName = "v5"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
