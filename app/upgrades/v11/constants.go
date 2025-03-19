package v11

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v10/app/upgrades"
	investorstypes "github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

const UpgradeName = "v11"

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
