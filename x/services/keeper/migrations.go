package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "github.com/milkyway-labs/milkyway/x/services/legacy/v2"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

type Migrator struct {
	k  *Keeper
	pk types.PoolsKeeper
}

func NewMigrator(keeper *Keeper, pk types.PoolsKeeper) Migrator {
	return Migrator{
		k:  keeper,
		pk: pk,
	}
}

func (m Migrator) Migrate1To2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.k, m.pk)
}
