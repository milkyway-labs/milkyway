package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/milkyway-labs/milkyway/v10/x/assets/types"
)

// SetAsset stores an asset in the store
func (k *Keeper) SetAsset(ctx context.Context, asset types.Asset) error {
	if err := k.Assets.Set(ctx, asset.Denom, asset); err != nil {
		return err
	}
	if err := k.TickerIndexes.Set(ctx, collections.Join(asset.Ticker, asset.Denom)); err != nil {
		return err
	}
	return nil
}

// GetAsset retrieves an asset from the store.
// If the asset does not exist, a collections.ErrNotFound error is returned.
func (k *Keeper) GetAsset(ctx context.Context, denom string) (types.Asset, error) {
	return k.Assets.Get(ctx, denom)
}

// RemoveAsset removes the asset having the given denom from the store.
// If the asset does not exist, a collections.ErrNotFound error is returned.
func (k *Keeper) RemoveAsset(ctx context.Context, denom string) error {
	asset, err := k.GetAsset(ctx, denom)
	if err != nil {
		return err
	}
	if err := k.Assets.Remove(ctx, denom); err != nil {
		return err
	}
	return k.TickerIndexes.Remove(ctx, collections.Join(asset.Ticker, denom))
}
