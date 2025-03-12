package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/liquidvesting/types"
	restakingtypes "github.com/milkyway-labs/milkyway/v9/x/restaking/types"
)

func (k *Keeper) GetAllUsersInsuranceFundsEntries(ctx sdk.Context) ([]types.UserInsuranceFundEntry, error) {
	var usersInsuranceFundState []types.UserInsuranceFundEntry
	err := k.insuranceFunds.Walk(ctx, nil, func(userAddress string, insuranceFund types.UserInsuranceFund) (stop bool, err error) {
		usersInsuranceFundState = append(usersInsuranceFundState, types.NewUserInsuranceFundEntry(
			userAddress,
			insuranceFund.Balance,
		))
		return false, nil
	})
	return usersInsuranceFundState, err
}

// GetAllUserRestakedLockedRepresentations returns all restaked coins that are locked
// representation tokens for the provided user.
func (k *Keeper) GetAllUserRestakedLockedRepresentations(ctx context.Context, userAddress string) (sdk.DecCoins, error) {
	restakedCoins, err := k.restakingKeeper.GetAllUserRestakedCoins(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	lockedRepresentations := sdk.NewDecCoins()
	for _, coin := range restakedCoins {
		if types.IsLockedRepresentationDenom(coin.Denom) {
			lockedRepresentations = lockedRepresentations.Add(coin)
		}
	}

	return lockedRepresentations, nil
}

// GetAllUserUnbondingLockedRepresentations returns all the locked representation
// tokens that are currently unbonding for the provided user.
func (k *Keeper) GetAllUserUnbondingLockedRepresentations(ctx context.Context, userAddress string) sdk.Coins {
	lockedRepresentations := sdk.NewCoins()

	userUndelegations := k.restakingKeeper.GetAllUserUnbondingDelegations(ctx, userAddress)
	for _, undelegation := range userUndelegations {
		for _, entry := range undelegation.Entries {
			for _, coin := range entry.Balance {
				if types.IsLockedRepresentationDenom(coin.Denom) {
					lockedRepresentations = lockedRepresentations.Add(coin)
				}
			}
		}
	}

	return lockedRepresentations
}

// GetAllUserActiveLockedRepresentations gets all the locked representation tokens
// that are restaked or are currently unbonding for the provided user.
func (k *Keeper) GetAllUserActiveLockedRepresentations(ctx context.Context, userAddress string) (sdk.DecCoins, error) {
	restakedCoins, err := k.GetAllUserRestakedLockedRepresentations(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	// Get the locked representation tokens that are currently unbonding
	userUnbondingLockedRepresentations := k.GetAllUserUnbondingLockedRepresentations(ctx, userAddress)

	return restakedCoins.Add(sdk.NewDecCoinsFromCoins(userUnbondingLockedRepresentations...)...), nil
}

// --------------------------------------------------------------------------------------------------------------------

// GetTargetCoveredLockedShares returns the covered locked shares for a delegation target.
func (k *Keeper) GetTargetCoveredLockedShares(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) (sdk.DecCoins, error) {
	shares, err := k.TargetsCoveredLockedShares.Get(ctx, collections.Join(int32(delType), targetID))
	if err != nil && !errors.IsOf(err, collections.ErrNotFound) {
		return nil, err
	}
	return shares.Shares, nil
}

// SetTargetCoveredLockedShares sets the total locked shares for a target.
func (k *Keeper) SetTargetCoveredLockedShares(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, shares sdk.DecCoins) error {
	return k.TargetsCoveredLockedShares.Set(
		ctx,
		collections.Join(int32(delType), targetID),
		types.TargetCoveredLockedShares{Shares: shares},
	)
}

// IncrementTargetCoveredLockedShares increments the total locked shares for a target.
func (k *Keeper) IncrementTargetCoveredLockedShares(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, shares sdk.DecCoins) error {
	prevShares, err := k.GetTargetCoveredLockedShares(ctx, delType, targetID)
	if err != nil {
		return err
	}
	newShares := prevShares.Add(shares...)
	return k.SetTargetCoveredLockedShares(ctx, delType, targetID, newShares)
}

// DecrementTargetCoveredLockedShares decrements the total locked shares for a target.
// If the total locked shares become zero, the record is deleted instead.
func (k *Keeper) DecrementTargetCoveredLockedShares(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, shares sdk.DecCoins) error {
	prevShares, err := k.GetTargetCoveredLockedShares(ctx, delType, targetID)
	if err != nil {
		return err
	}
	newShares := prevShares.Sub(shares)
	if newShares.IsZero() {
		// Delete the shares record if it becomes zero
		return k.RemoveTargetCoveredLockedShares(ctx, delType, targetID)
	}
	return k.SetTargetCoveredLockedShares(ctx, delType, targetID, newShares)
}

// IterateTargetsCoveredLockedShares iterates over all the targets covered locked
// shares and calls cb.
func (k *Keeper) IterateTargetsCoveredLockedShares(ctx context.Context, cb func(delType restakingtypes.DelegationType, targetID uint32, shares sdk.DecCoins) (stop bool, err error)) error {
	err := k.TargetsCoveredLockedShares.Walk(ctx, nil, func(key collections.Pair[int32, uint32], shares types.TargetCoveredLockedShares) (stop bool, err error) {
		delType := restakingtypes.DelegationType(key.K1())
		targetID := key.K2()
		return cb(delType, targetID, shares.Shares)
	})
	return err
}

// RemoveTargetCoveredLockedShares removes a target covered locked shares record
// for the given delegation target.
func (k *Keeper) RemoveTargetCoveredLockedShares(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32) error {
	return k.TargetsCoveredLockedShares.Remove(ctx, collections.Join(int32(delType), targetID))
}

// --------------------------------------------------------------------------------------------------------------------

// SetLockedRepresentationDelegator marks the user as a locked representation
// delegator.
func (k *Keeper) SetLockedRepresentationDelegator(ctx context.Context, userAddress string) error {
	return k.LockedRepresentationDelegators.Set(ctx, userAddress)
}

// RemoveLockedRepresentationDelegator removes the user from the locked
// representation delegators list.
func (k *Keeper) RemoveLockedRepresentationDelegator(ctx context.Context, userAddress string) error {
	return k.LockedRepresentationDelegators.Remove(ctx, userAddress)
}

// GetAllLockedRepresentationDelegators returns all the locked representation
// delegators.
func (k *Keeper) GetAllLockedRepresentationDelegators(ctx context.Context) ([]string, error) {
	iter, err := k.LockedRepresentationDelegators.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	return iter.Keys()
}

// --------------------------------------------------------------------------------------------------------------------

// SetPreviousDelegationTokens caches the previous delegation's tokens for the
// user and the delegation target.
func (k *Keeper) SetPreviousDelegationTokens(
	ctx context.Context,
	user string,
	delType restakingtypes.DelegationType,
	targetID uint32,
	tokens sdk.DecCoins,
) error {
	return k.PreviousDelegationsTokens.Set(
		ctx,
		collections.Join3(user, int32(delType), targetID),
		types.PreviousDelegationTokens{Tokens: tokens},
	)
}

// GetPreviousDelegationTokens returns the cached previous delegation for the
// user and the delegation target.
func (k *Keeper) GetPreviousDelegationTokens(
	ctx context.Context,
	user string,
	delType restakingtypes.DelegationType,
	targetID uint32,
) (sdk.DecCoins, error) {
	tokens, err := k.PreviousDelegationsTokens.Get(ctx, collections.Join3(user, int32(delType), targetID))
	if err != nil {
		if errors.IsOf(err, collections.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return tokens.Tokens, nil
}

// RemovePreviousDelegationTokens removes the cache of the previous delegation.
func (k *Keeper) RemovePreviousDelegationTokens(ctx context.Context, user string, delType restakingtypes.DelegationType, targetID uint32) error {
	return k.PreviousDelegationsTokens.Remove(ctx, collections.Join3(user, int32(delType), targetID))
}
