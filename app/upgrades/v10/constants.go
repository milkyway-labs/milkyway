package v10

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v9/app/upgrades"
	investorstypes "github.com/milkyway-labs/milkyway/v9/x/investors/types"
)

const UpgradeName = "v10"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{
			investorstypes.StoreKey,
		},
		Deleted: []string{},
	},
}
