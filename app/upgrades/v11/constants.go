package v11

import (
	_ "embed"

	storetypes "cosmossdk.io/store/types"

	"github.com/milkyway-labs/milkyway/v10/app/upgrades"
	investorstypes "github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

const UpgradeName = "v11"

//go:embed data.json
var dataBz []byte

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
