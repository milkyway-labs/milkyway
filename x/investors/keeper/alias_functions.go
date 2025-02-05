package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v7/x/investors/types"
)

func (k *Keeper) GetValidatorInvestorsShares(ctx context.Context, valAddr sdk.ValAddress) (sdkmath.LegacyDec, error) {
	shares, err := k.ValidatorsInvestorsShares.Get(ctx, valAddr)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return sdkmath.LegacyZeroDec(), nil
		}
		return sdkmath.LegacyDec{}, err
	}
	return shares, nil
}

func (k *Keeper) IncrementValidatorInvestorsShares(ctx context.Context, valAddr sdk.ValAddress, shares sdkmath.LegacyDec) error {
	prevShares, err := k.GetValidatorInvestorsShares(ctx, valAddr)
	if err != nil {
		return err
	}
	return k.ValidatorsInvestorsShares.Set(ctx, valAddr, prevShares.Add(shares))
}

func (k *Keeper) DecrementValidatorInvestorsShares(ctx context.Context, valAddr sdk.ValAddress, shares sdkmath.LegacyDec) error {
	prevShares, err := k.GetValidatorInvestorsShares(ctx, valAddr)
	if err != nil {
		return err
	}
	newShares := prevShares.Sub(shares)
	if newShares.IsNegative() {
		panic("cannot set negative shares")
	} else if newShares.IsZero() {
		return k.ValidatorsInvestorsShares.Remove(ctx, valAddr)
	}
	return k.ValidatorsInvestorsShares.Set(ctx, valAddr, newShares)
}

// UpdateInvestorsRewardRatio updates the investors reward ratio. It also
// increments the period of all validators, since each validator's total tokens
// need to be adjusted after updating the ratio.
func (k *Keeper) UpdateInvestorsRewardRatio(ctx context.Context, ratio sdkmath.LegacyDec) error {
	err := types.ValidateInvestorsRewardRatio(ratio)
	if err != nil {
		return err
	}

	var innerErr error
	err = k.stakingKeeper.IterateValidators(ctx, func(_ int64, validator stakingtypes.ValidatorI) (stop bool) {
		var valAddr sdk.ValAddress
		valAddr, innerErr = k.stakingKeeper.ValidatorAddressCodec().StringToBytes(validator.GetOperator())
		if innerErr != nil {
			return true
		}
		validator, innerErr = k.GetAdjustedValidator(ctx, valAddr)
		if innerErr != nil {
			return true
		}
		_, innerErr = k.distrKeeper.IncrementValidatorPeriod(ctx, validator)
		return innerErr != nil
	})
	if err != nil {
		return err
	}
	if innerErr != nil {
		return innerErr
	}

	return k.InvestorsRewardRatio.Set(ctx, ratio)
}

// GetDelegation returns a delegation and a boolean flag indicating if the
// account is an investor.
func (k *Keeper) GetDelegation(
	ctx context.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
) (delegation stakingtypes.Delegation, isInvestor bool, err error) {
	delegation, err = k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return
	}

	acc := k.accountKeeper.GetAccount(ctx, delAddr)
	_, isInvestor = acc.(vestingexported.VestingAccount)
	return delegation, isInvestor, nil
}

func (k *Keeper) GetAdjustedValidator(ctx context.Context, valAddr sdk.ValAddress) (stakingtypes.Validator, error) {
	validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}
	investorsShares, err := k.GetValidatorInvestorsShares(ctx, valAddr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}
	if investorsShares.IsPositive() {
		rewardRatio, err := k.InvestorsRewardRatio.Get(ctx)
		if err != nil {
			return stakingtypes.Validator{}, err
		}
		oneMinusRewardRatio := sdkmath.LegacyOneDec().Sub(rewardRatio)
		validator, _ = validator.RemoveDelShares(investorsShares.MulRoundUp(oneMinusRewardRatio))
	}
	return validator, nil
}
