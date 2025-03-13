package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

// ExportGenesis returns the GenesisState associated with the given context
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	investorsRewardRatio, err := k.GetInvestorsRewardRatio(ctx)
	if err != nil {
		panic(err)
	}

	vestingInvestorsAddrs, err := k.GetAllVestingInvestorsAddresses(ctx)
	if err != nil {
		panic(err)
	}

	return types.NewGenesisState(investorsRewardRatio, vestingInvestorsAddrs)
}

// InitGenesis initializes the state from a GenesisState
func (k *Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) error {
	err := k.SetInvestorsRewardRatio(ctx, state.InvestorsRewardRatio)
	if err != nil {
		return err
	}

	for _, investor := range state.VestingInvestorsAddresses {
		err = k.SetVestingInvestor(ctx, investor)
		if err != nil {
			return err
		}
	}

	// Initialize lookup table for delegator addresses
	k.distrKeeper.IterateDelegatorWithdrawAddrs(ctx, func(del, addr sdk.AccAddress) (stop bool) {
		var delegator, withdrawAddr string
		delegator, err = k.accountKeeper.AddressCodec().BytesToString(del)
		if err != nil {
			return true
		}
		withdrawAddr, err = k.accountKeeper.AddressCodec().BytesToString(addr)
		if err != nil {
			return true
		}
		err = k.Delegators.Set(ctx, withdrawAddr, delegator)
		return err != nil
	})
	if err != nil {
		return err
	}

	// Create the module account if it doesn't exist
	k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)

	return nil
}
