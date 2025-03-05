package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/investors/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	investorsRewardRatio, err := k.InvestorsRewardRatio.Get(ctx)
	if err != nil {
		panic(err)
	}

	vestingInvestorsAddrs, err := k.GetAllVestingInvestorsAddresses(ctx)
	if err != nil {
		panic(err)
	}

	var validatorsInvestorsShares []types.ValidatorInvestorsShares
	err = k.ValidatorsInvestorsShares.Walk(ctx, nil, func(valAddr sdk.ValAddress, shares sdkmath.LegacyDec) (stop bool, err error) {
		validator, err := k.stakingKeeper.ValidatorAddressCodec().BytesToString(valAddr)
		if err != nil {
			return true, err
		}
		validatorsInvestorsShares = append(validatorsInvestorsShares, types.ValidatorInvestorsShares{
			ValidatorAddress: validator,
			InvestorsShares:  shares,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	return types.NewGenesisState(
		investorsRewardRatio,
		vestingInvestorsAddrs,
		validatorsInvestorsShares,
	)
}

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	err := k.InvestorsRewardRatio.Set(ctx, state.InvestorsRewardRatio)
	if err != nil {
		return err
	}

	for _, investor := range state.VestingInvestorsAddresses {
		investorAddr, err := k.accountKeeper.AddressCodec().StringToBytes(investor)
		if err != nil {
			return err
		}
		err = k.SetVestingInvestor(ctx, investorAddr)
		if err != nil {
			return err
		}
	}

	for _, shares := range state.ValidatorsInvestorsShares {
		valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(shares.ValidatorAddress)
		if err != nil {
			return err
		}
		err = k.ValidatorsInvestorsShares.Set(ctx, valAddr, shares.InvestorsShares)
		if err != nil {
			return err
		}
	}
	return nil
}
