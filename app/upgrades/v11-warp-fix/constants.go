package v11warpfix

import (
	_ "embed"

	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v11/app/upgrades"
)

const UpgradeName = "v11-warp-fix"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
