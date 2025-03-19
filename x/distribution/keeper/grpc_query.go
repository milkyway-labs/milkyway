package keeper

import (
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
)

func NewQuerier(keeper Keeper) distrkeeper.Querier {
	return distrkeeper.Querier{Keeper: keeper.Keeper}
}
