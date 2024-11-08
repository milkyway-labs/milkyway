package app

import (
	v110 "github.com/milkyway-labs/milkyway/app/upgrades/v110"
	v122 "github.com/milkyway-labs/milkyway/app/upgrades/v122"
	v130 "github.com/milkyway-labs/milkyway/app/upgrades/v130"
	v140 "github.com/milkyway-labs/milkyway/app/upgrades/v140"
	v150 "github.com/milkyway-labs/milkyway/app/upgrades/v150"
)

// RegisterUpgradeHandlers returns upgrade handlers
func (app *MilkyWayApp) RegisterUpgradeHandlers() {
	app.registerUpgrade(v110.NewUpgrade(app.ModuleManager, app.Configurator(), app.appCodec, app.keys, app.RewardsKeeper))
	app.registerUpgrade(v122.NewUpgrade(app.ModuleManager, app.Configurator(), app.StakeIBCKeeper))
	app.registerUpgrade(v130.NewUpgrade(app.ModuleManager, app.Configurator(), app.LiquidVestingKeeper))
	app.registerUpgrade(v140.NewUpgrade(app.ModuleManager, app.Configurator(), app.LiquidVestingKeeper))
	app.registerUpgrade(v150.NewUpgrade(app.ModuleManager, app.Configurator(), app.LiquidVestingKeeper))
}
