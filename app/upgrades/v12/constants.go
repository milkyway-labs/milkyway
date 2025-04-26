package v12

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v12/app/upgrades"
)

const UpgradeName = "v12"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
