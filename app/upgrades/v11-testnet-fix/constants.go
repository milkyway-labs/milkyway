package v11_testnet_fix

import (
	_ "embed"

	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v11/app/upgrades"
)

const UpgradeName = "v11-testnet-fix"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
