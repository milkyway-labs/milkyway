package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/milkyway-labs/milkyway/v7/x/distribution/types"
)

type Keeper struct {
	distrkeeper.Keeper

	accountKeeper distrtypes.AccountKeeper
	stakingKeeper types.StakingKeeper

	Validators collections.Map[sdk.ValAddress, stakingtypes.Validator]
}

func NewKeeper(
	cdc codec.BinaryCodec, storeService store.KVStoreService,
	ak distrtypes.AccountKeeper, bk distrtypes.BankKeeper, sk types.StakingKeeper,
	feeCollectorName, authority string,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	stakingKeeper := &StakingKeeper{
		StakingKeeper: sk,
	}
	k := Keeper{
		Keeper:        distrkeeper.NewKeeper(cdc, storeService, ak, bk, stakingKeeper, feeCollectorName, authority),
		accountKeeper: ak,
		stakingKeeper: sk,
		Validators:    collections.NewMap(sb, types.ValidatorsKeyPrefix, "validators", sdk.ValAddressKey, codec.CollValue[stakingtypes.Validator](cdc)),
	}
	stakingKeeper.k = &k

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// StakingKeeper is a wrapper around the staking keeper that returns modified values of
// validator's total tokens, total shares and delegation's shares.
type StakingKeeper struct {
	types.StakingKeeper

	k *Keeper
}

func (sk StakingKeeper) Validator(ctx context.Context, address sdk.ValAddress) (stakingtypes.ValidatorI, error) {
	return sk.k.Validators.Get(ctx, address)
}

func (sk StakingKeeper) Delegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.DelegationI, error) {
	delegation, err := sk.k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return nil, err
	}
	// If the account is a vesting account, halve the shares
	acc := sk.k.accountKeeper.GetAccount(ctx, delAddr)
	_, isVestingAcc := acc.(vestingexported.VestingAccount)
	if isVestingAcc {
		delegation.Shares = delegation.Shares.MulTruncate(types.VestingAccountRewardsRatio)
	}
	return delegation, nil
}
