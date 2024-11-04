package app

import (
	v110 "github.com/milkyway-labs/milkyway/app/upgrades/v110"
	v122 "github.com/milkyway-labs/milkyway/app/upgrades/v122"
	v130 "github.com/milkyway-labs/milkyway/app/upgrades/v130"
)

// RegisterUpgradeHandlers returns upgrade handlers
func (app *MilkyWayApp) RegisterUpgradeHandlers() {
	app.registerUpgrade(v110.NewUpgrade(app.ModuleManager, app.Configurator(), app.appCodec, app.keys, app.RewardsKeeper))
	app.registerUpgrade(v122.NewUpgrade(app.ModuleManager, app.Configurator(), app.StakeIBCKeeper))
	app.registerUpgrade(v130.NewUpgrade(app.ModuleManager, app.Configurator(), app.LiquidVestingKeeper))
}
