package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
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
