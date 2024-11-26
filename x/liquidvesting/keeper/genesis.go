package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	// Get the params
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	// Get the users' insurance fund
	insuranceFundsEntries, err := k.GetAllUsersInsuranceFundsEntries(ctx)
	if err != nil {
		return nil, err
	}

	return types.NewGenesisState(params, k.GetAllBurnCoins(ctx), insuranceFundsEntries), nil
}

// --------------------------------------------------------------------------------------------------------------------

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	err := state.Validate()
	if err != nil {
		return err
	}

	// Store the params
	err = k.SetParams(ctx, state.Params)
	if err != nil {
		return err
	}

	totalCoins := sdk.NewCoins()
	for _, entry := range state.UserInsuranceFunds {
		insuranceFund := types.NewInsuranceFund(entry.Balance)

		// Store the insurance fund
		err = k.insuranceFunds.Set(ctx, entry.UserAddress, insuranceFund)
		if err != nil {
			return err
		}

		// Update the total coins in the insurance fund
		totalCoins = totalCoins.Add(entry.Balance...)

		// Get the total locked representation that should be covered by the
		// insurance fund
		totalLockedRepresentations, err := k.GetAllUserActiveLockedRepresentations(ctx, entry.UserAddress)
		if err != nil {
			return err
		}

		// Check if the insurance fund can cover the restaked coins
		canCover, required, err := insuranceFund.CanCoverDecCoins(state.Params.InsurancePercentage, totalLockedRepresentations)
		if err != nil {
			return err
		}

		if !canCover {
			return fmt.Errorf("user: %s insurance fund amount is too low, expected %s, got %s",
				entry.UserAddress, required.String(), entry.Balance.String())
		}
	}

	// Ensure that the balance of the liquid vesting module is equal to the
	// sum of the users' insurance fund
	coins, err := k.GetInsuranceFundBalance(ctx)
	if err != nil {
		return err
	}
	if !coins.Equal(totalCoins) {
		return fmt.Errorf("the liquid vesting module doesn't have the coins specified in the users' insurance fund")
	}

	unbondingDelegations, err := k.restakingKeeper.GetAllUnbondingDelegations(ctx)
	if err != nil {
		return err
	}

	undelegateAmounts := make(map[string]sdk.Coins)
	for _, ud := range unbondingDelegations {
		balance, found := undelegateAmounts[ud.DelegatorAddress]
		if !found {
			balance = sdk.NewCoins()
		}
		// Compute the new amount
		for _, entry := range ud.Entries {
			balance = balance.Add(entry.Balance...)
		}
		// Store the newly computed undelegate amount
		undelegateAmounts[ud.DelegatorAddress] = balance
	}

	for _, burnCoins := range state.BurnCoins {
		userUndelegateAmount, found := undelegateAmounts[burnCoins.DelegatorAddress]
		if !found {
			return fmt.Errorf("%s don't have tokens that are being undelegated", burnCoins.DelegatorAddress)
		}
		if !userUndelegateAmount.IsAllGTE(burnCoins.Amount) {
			return fmt.Errorf("%s don't have enough tokens that are being undelegated", burnCoins.DelegatorAddress)
		}
		// Update the undelegate amounts that can be considered for this user
		undelegateAmounts[burnCoins.DelegatorAddress] = userUndelegateAmount.Sub(burnCoins.Amount...)
		err = k.InsertBurnCoinsToUnbondingQueue(ctx, burnCoins)
		if err != nil {
			return err
		}
	}

	return nil
}
