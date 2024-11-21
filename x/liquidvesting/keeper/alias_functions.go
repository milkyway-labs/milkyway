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

// GetAllUserRestakedVestedRepresentations returns all restaked coins that are vested
// representation tokens for the provided user.
func (k *Keeper) GetAllUserRestakedVestedRepresentations(ctx context.Context, userAddress string) (sdk.DecCoins, error) {
	restakedCoins, err := k.restakingKeeper.GetAllUserRestakedCoins(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	vestedRepresentations := sdk.NewDecCoins()
	for _, coin := range restakedCoins {
		if types.IsVestedRepresentationDenom(coin.Denom) {
			vestedRepresentations = vestedRepresentations.Add(coin)
		}
	}

	return vestedRepresentations, nil
}

// GetAllUserUnbondingVestedRepresentations returns all the vested representation
// tokens that are currently unbonding for the provided user.
func (k *Keeper) GetAllUserUnbondingVestedRepresentations(ctx context.Context, userAddress string) sdk.Coins {
	vestedRepresentations := sdk.NewCoins()

	userUndelegations := k.restakingKeeper.GetAllUserUnbondingDelegations(ctx, userAddress)
	for _, undelegation := range userUndelegations {
		for _, entry := range undelegation.Entries {
			for _, coin := range entry.Balance {
				if types.IsVestedRepresentationDenom(coin.Denom) {
					vestedRepresentations = vestedRepresentations.Add(coin)
				}
			}
		}
	}

	return vestedRepresentations
}

// GetAllUserActiveVestedRepresentations gets all the vested representation tokens
// that are restaked or are currently unbonding for the provided user.
func (k *Keeper) GetAllUserActiveVestedRepresentations(ctx context.Context, userAddress string) (sdk.DecCoins, error) {
	restakedCoins, err := k.GetAllUserRestakedVestedRepresentations(ctx, userAddress)
	if err != nil {
		return nil, err
	}

	// Get the vested representation tokens that are currently unbonding
	userUnbondingVestedRepresentations := k.GetAllUserUnbondingVestedRepresentations(ctx, userAddress)

	return restakedCoins.Add(sdk.NewDecCoinsFromCoins(userUnbondingVestedRepresentations...)...), nil
}
