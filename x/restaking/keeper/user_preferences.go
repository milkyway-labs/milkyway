package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
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

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	k.servicesKeeper.IterateServices(sdkCtx, func(service servicestypes.Service) bool {
		trustedBefore := oldPreferences.IsServiceTrusted(service.ID, service.Accredited)
		trustedAfter := preferences.IsServiceTrusted(service.ID, service.Accredited)
		if trustedBefore != trustedAfter {
			err = k.AfterUserTrustedServiceUpdated(sdkCtx, userAddress, service.ID, trustedAfter)
			if err != nil {
				return true
			}
		}
		return false
	})
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
