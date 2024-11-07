package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

// SetUserPreferences sets the given preferences for the user having the given address
func (k *Keeper) SetUserPreferences(ctx context.Context, userAddress string, preferences types.UserPreferences) error {
	return k.usersPreferences.Set(ctx, userAddress, preferences)
}

// GetUserPreferences returns the preferences of the user having the given address.
// If no custom preferences have been previously set by a user, the default ones will be returned instead.
func (k *Keeper) GetUserPreferences(ctx context.Context, userAddress string) (types.UserPreferences, error) {
	preferences, err := k.usersPreferences.Get(ctx, userAddress)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.DefaultUserPreferences(), nil
		}
		return types.UserPreferences{}, err
	}
	return preferences, nil
}
