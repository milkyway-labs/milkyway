package v9

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v10/app/upgrades"
	tokenfactorytypes "github.com/milkyway-labs/milkyway/v10/x/tokenfactory/types"
)

const UpgradeName = "v9"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added:   []string{tokenfactorytypes.StoreKey},
		Deleted: []string{},
	},
}
