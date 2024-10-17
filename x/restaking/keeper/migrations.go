package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "github.com/milkyway-labs/milkyway/x/restaking/migrations/v2"
)

type Migrator struct {
	k *Keeper
}

func NewMigrator(k *Keeper) Migrator {
	return Migrator{
		k: k,
	}
}

func (m *Migrator) Migrate1To2(ctx sdk.Context) error {
	return v2.Migrate1To2(ctx, m.k.storeKey, m.k.cdc, m.k.operatorsKeeper)
}
