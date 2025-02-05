package types

import (
	"context"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	AddressCodec() address.Codec
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

// StakingKeeper defines the expected staking keeper.
type StakingKeeper interface {
	distrtypes.StakingKeeper

	GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error)
}

// DistrKeeper defines the expected distribution keeper.
type DistrKeeper interface {
	IncrementValidatorPeriod(ctx context.Context, val stakingtypes.ValidatorI) (uint64, error)
}
