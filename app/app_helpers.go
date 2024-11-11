package milkyway

import (
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctestingtypes "github.com/cosmos/ibc-go/v8/testing/types"
	icstest "github.com/cosmos/interchain-security/v6/testutil/integration"
	ibcproviderkeeper "github.com/cosmos/interchain-security/v6/x/ccv/provider/keeper"
)

// ProviderApp interface implementations for icstest tests

// GetProviderKeeper implements the ProviderApp interface.
func (app *MilkyWayApp) GetProviderKeeper() ibcproviderkeeper.Keeper { //nolint:nolintlint
	return app.ProviderKeeper
}

// GetStakingKeeper implements the TestingApp interface. Needed for ICS.
func (app *MilkyWayApp) GetStakingKeeper() ibctestingtypes.StakingKeeper { //nolint:nolintlint
	return app.StakingKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *MilkyWayApp) GetIBCKeeper() *ibckeeper.Keeper { //nolint:nolintlint
	return app.IBCKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *MilkyWayApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper { //nolint:nolintlint
	return app.ScopedIBCKeeper
}

// GetTestStakingKeeper implements the ProviderApp interface.
func (app *MilkyWayApp) GetTestStakingKeeper() icstest.TestStakingKeeper { //nolint:nolintlint
	return app.StakingKeeper
}

// GetTestBankKeeper implements the ProviderApp interface.
func (app *MilkyWayApp) GetTestBankKeeper() icstest.TestBankKeeper { //nolint:nolintlint
	return app.BankKeeper
}

// GetTestSlashingKeeper implements the ProviderApp interface.
func (app *MilkyWayApp) GetTestSlashingKeeper() icstest.TestSlashingKeeper { //nolint:nolintlint
	return app.SlashingKeeper
}

// GetTestDistributionKeeper implements the ProviderApp interface.
func (app *MilkyWayApp) GetTestDistributionKeeper() icstest.TestDistributionKeeper { //nolint:nolintlint
	return app.DistrKeeper
}

func (app *MilkyWayApp) GetTestAccountKeeper() icstest.TestAccountKeeper { //nolint:nolintlint
	return app.AccountKeeper
}
