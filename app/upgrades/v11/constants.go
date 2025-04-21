package v11

import (
	_ "embed"

	storetypes "cosmossdk.io/store/types"
	hyperlanetypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	warptypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"

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
			hyperlanetypes.ModuleName,
			warptypes.ModuleName,
		},
		Deleted: []string{},
	},
}
