package app

import (
	v110 "github.com/milkyway-labs/milkyway/app/upgrades/v110"
)

// RegisterUpgradeHandlers returns upgrade handlers
func (app *MilkyWayApp) RegisterUpgradeHandlers() {
	app.registerUpgrade(v110.NewUpgrade(app.ModuleManager, app.Configurator()))
}
