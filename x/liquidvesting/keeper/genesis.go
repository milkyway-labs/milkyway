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
		k.insuranceFunds.Set(ctx, accAddr, insuranceFund.InsuranceFund)
		totalCoins = totalCoins.Add(insuranceFund.InsuranceFund.Balance...)
	}
	coins, err := k.GetInsuranceFundBalance(ctx)
	if err != nil {
		return err
	}
	if !coins.Equal(totalCoins) {
		return fmt.Errorf("the liquid vesting module doesn't have same coins as genesis state")
	}

	undelegateAmounts := make(map[string]sdk.Coins)
	for _, ud := range k.restakingKeeper.GetAllUnbondingDelegations(ctx) {
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
		k.InsertBurnCoinsQueue(ctx, burnCoins, burnCoins.CompletionTime)
	}

	return nil
}
