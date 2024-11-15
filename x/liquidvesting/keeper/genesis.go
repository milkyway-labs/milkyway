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
	var usersInsuranceFundState []types.UserInsuranceFundState
	err = k.insuranceFunds.Walk(ctx, nil,
		func(accAddr sdk.AccAddress, insuranceFund types.UserInsuranceFund) (stop bool, err error) {
			strAddr, err := k.accountKeeper.AddressCodec().BytesToString(accAddr)
			if err != nil {
				return true, err
			}
			usersInsuranceFundState = append(usersInsuranceFundState, types.NewUserInsuranceFundState(strAddr, insuranceFund))
			return false, nil
		})
	if err != nil {
		return nil, err
	}

	return types.NewGenesisState(params, k.GetAllBurnCoins(ctx), usersInsuranceFundState), nil
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
	for _, insuranceFund := range state.UserInsuranceFunds {
		accAddr := sdk.MustAccAddressFromBech32(insuranceFund.UserAddress)
		err = k.insuranceFunds.Set(ctx, accAddr, insuranceFund.InsuranceFund)
		if err != nil {
			return err
		}
		// Update the total coins in the insurance fund
		totalCoins = totalCoins.Add(insuranceFund.InsuranceFund.Balance...)

		// We compute the total amount of coins that should be used
		// in the insurance fund based on the assets that the user has restaked
		expectedUsed := sdk.NewCoins()
		userRestakedCoins, err := k.restakingKeeper.GetAllUserRestakedCoins(ctx, insuranceFund.UserAddress)
		if err != nil {
			return err
		}
		// Update the user restaked amount with also the tokens that are
		// in the unbonding state and will be received at completion
		for _, undelegation := range k.restakingKeeper.GetAllUserUnbondingDelegations(ctx, insuranceFund.UserAddress) {
			for _, entry := range undelegation.Entries {
				userRestakedCoins = userRestakedCoins.Add(sdk.NewDecCoinsFromCoins(entry.Balance...)...)
			}
		}

		for _, restakedCoin := range userRestakedCoins {
			if types.IsVestedRepresentationDenom(restakedCoin.Denom) {
				requiredAmount, err := k.getRequiredAmountInInsuranceFund(ctx,
					sdk.NewCoin(restakedCoin.Denom, restakedCoin.Amount.RoundInt()))
				if err != nil {
					return err
				}
				expectedUsed = expectedUsed.Add(requiredAmount)
			}
		}

		// Ensure that the amount that we have computed is the same as the
		// used amount in the user's insurance fund
		if !expectedUsed.Equal(insuranceFund.InsuranceFund.Used) {
			return fmt.Errorf("user: %s insurance fund used amount is incorrect, expected %s, got %s",
				insuranceFund.UserAddress, expectedUsed.String(), insuranceFund.InsuranceFund.Used.String())
		}
	}

	// Ensure that the balance of the liquid vesting module that coins the
	// insurance fund is equal to the sum of amount of the users' insurance fund
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
