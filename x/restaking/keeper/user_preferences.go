package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	cosmoserrors "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/v3/x/restaking/types"
)

// SetUserPreferences sets the given preferences for the user having the given address
func (k *Keeper) SetUserPreferences(ctx context.Context, userAddress string, preferences types.UserPreferences) error {
	err := preferences.Validate()
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid preferences: %s", err)
	}

	oldPreferences, err := k.GetUserPreferences(ctx, userAddress)
	if err != nil {
		return err
	}

	err = k.AfterUserPreferencesModified(ctx, userAddress, oldPreferences, preferences)
	if err != nil {
		return err
	}

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
