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

// IncrementTargetCoveredLockedShares increments the total locked shares for a target.
func (k *Keeper) IncrementTargetCoveredLockedShares(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, shares sdk.DecCoins) error {
	prevShares, err := k.GetTargetCoveredLockedShares(ctx, delType, targetID)
	if err != nil {
		return err
	}
	newShares := prevShares.Add(shares...)
	return k.TargetsCoveredLockedShares.Set(
		ctx,
		collections.Join(int32(delType), targetID),
		types.CoveredLockedShares{Shares: newShares},
	)
}

// DecrementTargetCoveredLockedShares decrements the total locked shares for a target.
// If the total locked shares become zero, the record is deleted instead.
func (k *Keeper) DecrementTargetCoveredLockedShares(ctx context.Context, delType restakingtypes.DelegationType, targetID uint32, shares sdk.DecCoins) error {
	prevShares, err := k.GetTargetCoveredLockedShares(ctx, delType, targetID)
	if err != nil {
		return err
	}
	newShares := prevShares.Sub(shares)
	key := collections.Join(int32(delType), targetID)
	// Delete the shares record if it becomes zero
	if newShares.IsZero() {
		return k.TargetsCoveredLockedShares.Remove(ctx, key)
	}
	return k.TargetsCoveredLockedShares.Set(ctx, key, types.CoveredLockedShares{Shares: newShares})
}
