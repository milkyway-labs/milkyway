package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var _ stakingtypes.StakingHooks = Hooks{}

type Hooks struct {
	distrkeeper.Hooks

	k Keeper
}

func (k Keeper) Hooks() Hooks {
	return Hooks{
		Hooks: k.Keeper.Hooks(),
		k:     k,
	}
}

func (h Hooks) AfterValidatorCreated(ctx context.Context, valAddr sdk.ValAddress) error {
	// Clone the validator
	validator, err := h.k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}
	err = h.k.Validators.Set(ctx, valAddr, validator)
	if err != nil {
		return err
	}
	return h.Hooks.AfterValidatorCreated(ctx, valAddr)
}

func (h Hooks) AfterValidatorRemoved(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	// Delete the validator
	err := h.k.Validators.Remove(ctx, valAddr)
	if err != nil {
		return err
	}
	return h.Hooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
}

func (h Hooks) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	delegation, err := h.k.stakingKeeper.Delegation(ctx, delAddr, valAddr)
	if err != nil {
		return err
	}

	shares := delegation.GetShares()
	// If the account is a vesting account, halve the shares
	acc := h.k.accountKeeper.GetAccount(ctx, delAddr)
	_, isVestingAcc := acc.(vestingexported.VestingAccount)
	if isVestingAcc {
		shares = shares.QuoInt64(2) // 50%
	}

	validator, err := h.k.Validators.Get(ctx, valAddr)
	if err != nil {
		return err
	}
	validator, _ = validator.RemoveDelShares(shares)
	err = h.k.Validators.Set(ctx, valAddr, validator)
	if err != nil {
		return err
	}

	return h.Hooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
}

func (h Hooks) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	delegation, err := h.k.stakingKeeper.Delegation(ctx, delAddr, valAddr)
	if err != nil {
		return err
	}

	shares := delegation.GetShares()
	// If the account is a vesting account, halve the shares
	acc := h.k.accountKeeper.GetAccount(ctx, delAddr)
	_, isVestingAcc := acc.(vestingexported.VestingAccount)
	if isVestingAcc {
		shares = shares.QuoInt64(2)
	}

	validator, err := h.k.Validators.Get(ctx, valAddr)
	if err != nil {
		return err
	}
	var tokensToAdd sdkmath.Int
	if validator.DelegatorShares.IsZero() {
		tokensToAdd = shares.TruncateInt()
	} else {
		tokensToAdd = validator.TokensFromSharesTruncated(shares).TruncateInt()
	}
	validator, _ = validator.AddTokensFromDel(tokensToAdd)
	err = h.k.Validators.Set(ctx, valAddr, validator)
	if err != nil {
		return err
	}

	return h.Hooks.AfterDelegationModified(ctx, delAddr, valAddr)
}

func (h Hooks) BeforeValidatorSlashed(ctx context.Context, valAddr sdk.ValAddress, fraction sdkmath.LegacyDec) error {
	validator, err := h.k.Validators.Get(ctx, valAddr)
	if err != nil {
		return err
	}
	validator.Tokens = (sdkmath.LegacyOneDec().Sub(fraction)).MulInt(validator.Tokens).TruncateInt()
	err = h.k.Validators.Set(ctx, valAddr, validator)
	if err != nil {
		return err
	}

	return h.Hooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
}
