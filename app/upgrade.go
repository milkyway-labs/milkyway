package app

import (
	v160 "github.com/milkyway-labs/milkyway/app/upgrades/v160"
)

// RegisterUpgradeHandlers returns upgrade handlers
func (app *MilkyWayApp) RegisterUpgradeHandlers() {
	app.registerUpgrade(v160.NewUpgrade(app.ModuleManager, app.Configurator(), app.StakeIBCKeeper))
}
