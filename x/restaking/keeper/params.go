package keeper

import (
	"errors"
	"time"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/utils"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// UnbondingTime returns the unbonding time.
func (k *Keeper) UnbondingTime(ctx sdk.Context) (time.Duration, error) {
	unbondingTimeNanos, err := k.unbondingTimeNanos.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.DefaultUnbondingTime, nil
		}
		return time.Duration(0), err
	}
	return time.Duration(unbondingTimeNanos) * time.Nanosecond, nil
}

// SetUnbondingTime sets the unbonding time.
func (k *Keeper) SetUnbondingTime(ctx sdk.Context, unbondingTime time.Duration) error {
	return k.unbondingTimeNanos.Set(ctx, unbondingTime.Nanoseconds())
}

// GetAllRestakbleAssets returns all the restakable assets.
func (k *Keeper) GetAllRestakbleAssets(ctx sdk.Context) ([]string, error) {
	var restakableAssets []string
	err := k.allowedRestakableDenoms.Walk(ctx, nil, func(denom string) (stop bool, err error) {
		restakableAssets = append(restakableAssets, denom)
		return false, nil
	})
	if err != nil {
		return nil, err
	}
	return restakableAssets, nil
}

// SetRestakableDenoms sets the denoms that are allowed to be restaked.
func (k *Keeper) SetRestakableDenoms(ctx sdk.Context, denoms []string) error {
	// Remove all the allowed restakable denoms from the store
	err := k.allowedRestakableDenoms.Clear(ctx, nil)
	if err != nil {
		return err
	}
	// Add the
	for _, denom := range denoms {
		err := k.allowedRestakableDenoms.Set(ctx, denom)
		if err != nil {
			return err
		}
	}

	return nil
}

// IsAssetRestakable checks if the asset with the provided denom is allowed
// to be restaked.
func (k *Keeper) IsAssetRestakable(ctx sdk.Context, denom string) (bool, error) {
	isEmpty, err := utils.IsKeySetEmpty(ctx, k.allowedRestakableDenoms, nil)
	if err != nil {
		return false, err
	}
	if isEmpty {
		return true, nil
	}

	return k.allowedRestakableDenoms.Has(ctx, denom)
}

// SetParams sets module parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	err := k.SetUnbondingTime(ctx, params.UnbondingTime)
	if err != nil {
		return err
	}

	err = k.SetRestakableDenoms(ctx, params.AllowedDenoms)
	if err != nil {
		return err
	}

	return nil
}

// GetParams returns the module parameters
func (k *Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	unbondingTime, err := k.UnbondingTime(ctx)
	if err != nil {
		return types.Params{}, err
	}
	restakableAssets, err := k.GetAllRestakbleAssets(ctx)
	if err != nil {
		return types.Params{}, nil
	}

	return types.NewParams(unbondingTime, restakableAssets), nil
}
