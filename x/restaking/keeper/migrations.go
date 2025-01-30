package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "github.com/milkyway-labs/milkyway/v8/x/restaking/migrations/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper *Keeper
}

// NewMigrator creates a new instance of Migrator.
func NewMigrator(keeper *Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

// Migrate1To2 migrates from version 1 to 2.
func (m Migrator) Migrate1To2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.keeper, m.keeper.storeService, m.keeper.cdc)
}
