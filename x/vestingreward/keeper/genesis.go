package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v7/x/vestingreward/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	vestingAccountsRewardRatio, err := k.VestingAccountsRewardRatio.Get(ctx)
	if err != nil {
		panic(err)
	}

	var validatorsVestingAccountsShares []types.ValidatorVestingAccountsShares
	err = k.ValidatorsVestingAccountsShares.Walk(ctx, nil, func(valAddr sdk.ValAddress, shares sdkmath.LegacyDec) (stop bool, err error) {
		validator, err := k.stakingKeeper.ValidatorAddressCodec().BytesToString(valAddr)
		if err != nil {
			return true, err
		}
		validatorsVestingAccountsShares = append(validatorsVestingAccountsShares, types.ValidatorVestingAccountsShares{
			ValidatorAddress:      validator,
			VestingAccountsShares: shares,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	return types.NewGenesisState(
		vestingAccountsRewardRatio,
		validatorsVestingAccountsShares,
	)
}

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	err := k.VestingAccountsRewardRatio.Set(ctx, state.VestingAccountsRewardRatio)
	if err != nil {
		return err
	}

	for _, shares := range state.ValidatorsVestingAccountsShares {
		valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(shares.ValidatorAddress)
		if err != nil {
			return err
		}
		err = k.ValidatorsVestingAccountsShares.Set(ctx, valAddr, shares.VestingAccountsShares)
		if err != nil {
			return err
		}
	}
	return nil
}
