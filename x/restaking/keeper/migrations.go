package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "github.com/milkyway-labs/milkyway/x/restaking/migrations/v2"
	v3 "github.com/milkyway-labs/milkyway/x/restaking/migrations/v3"
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
	return v2.Migrate1To2(ctx,
		m.k.storeKey,
		m.k.cdc,
		m.k,
		m.k.operatorsKeeper,
		m.k.servicesKeeper,
	)
}

func (m *Migrator) Migrate2To3(ctx sdk.Context) error {
	return v3.Migrate2To3(ctx,
		m.k.storeKey,
		m.k.cdc,
		m.k,
	)
}
