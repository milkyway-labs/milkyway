package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/milkyway-labs/milkyway/v10/x/investors/types"
)

// SendRestrictionFn function that can be used in the x/bank module to intercept
// the transfer from the distribution module account to the delegator address and
// redirect the rewards to the investors module account. These redirected rewards
// are then redistributed to the investors in the end blocker, after deducting
// the amount based on the InvestorsRewardRatio parameter.
func (k *Keeper) SendRestrictionFn(ctx context.Context, from sdk.AccAddress, to sdk.AccAddress, amount sdk.Coins) (sdk.AccAddress, error) {
	// If the sender is not the distribution module account, skip.
	distrModuleAddr := k.accountKeeper.GetModuleAddress(distrtypes.ModuleName)
	if !from.Equals(distrModuleAddr) {
		return to, nil
	}

	// Get the delegator address from the withdraw address
	toAddrStr, err := k.accountKeeper.AddressCodec().BytesToString(to)
	if err != nil {
		return nil, err
	}
	delegator, err := k.GetDelegatorAddressByWithdrawAddress(ctx, toAddrStr)
	if err != nil {
		return nil, err
	}

	// If the delegator is not a vesting investor, skip.
	isVestingInvestor, err := k.IsVestingInvestor(ctx, delegator)
	if err != nil {
		return nil, err
	}
	if !isVestingInvestor {
		return to, nil
	}

	investorsRewardRatio, err := k.GetInvestorsRewardRatio(ctx)
	if err != nil {
		return nil, err
	}
	shared, _ := sdk.NewDecCoinsFromCoins(amount...).MulDecTruncate(investorsRewardRatio).TruncateDecimal()
	if !shared.IsZero() {
		err = k.IncrementVestingInvestorRewards(ctx, delegator, shared)
		if err != nil {
			return nil, err
		}
	}

	// Redirect the rewards to the module account
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	return moduleAddr, nil
}