package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	collcodec "cosmossdk.io/collections/codec"
	corestoretypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/cosmos/gogoproto/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

type Keeper struct {
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService

	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	communityPoolKeeper types.CommunityPoolKeeper
	oracleKeeper        types.OracleKeeper
	poolsKeeper         types.PoolsKeeper
	operatorsKeeper     types.OperatorsKeeper
	servicesKeeper      types.ServicesKeeper
	restakingKeeper     types.RestakingKeeper
	tickersKeeper       types.TickersKeeper

	Schema                         collections.Schema
	Params                         collections.Item[types.Params]
	NextRewardsPlanID              collections.Item[uint64]
	RewardsPlans                   collections.Map[uint64, types.RewardsPlan]
	LastRewardsAllocationTime      collections.Item[gogotypes.Timestamp]
	DelegatorWithdrawAddrs         collections.Map[sdk.AccAddress, sdk.AccAddress]
	PoolDelegatorStartingInfos     collections.Map[collections.Pair[uint32, sdk.AccAddress], types.DelegatorStartingInfo]
	PoolHistoricalRewards          collections.Map[collections.Pair[uint32, uint64], types.HistoricalRewards]
	PoolCurrentRewards             collections.Map[uint32, types.CurrentRewards]
	PoolOutstandingRewards         collections.Map[uint32, types.OutstandingRewards]
	OperatorAccumulatedCommissions collections.Map[uint32, types.AccumulatedCommission]
	OperatorDelegatorStartingInfos collections.Map[collections.Pair[uint32, sdk.AccAddress], types.DelegatorStartingInfo]
	OperatorHistoricalRewards      collections.Map[collections.Pair[uint32, uint64], types.HistoricalRewards]
	OperatorCurrentRewards         collections.Map[uint32, types.CurrentRewards]
	OperatorOutstandingRewards     collections.Map[uint32, types.OutstandingRewards]
	ServiceDelegatorStartingInfos  collections.Map[collections.Pair[uint32, sdk.AccAddress], types.DelegatorStartingInfo]
	ServiceHistoricalRewards       collections.Map[collections.Pair[uint32, uint64], types.HistoricalRewards]
	ServiceCurrentRewards          collections.Map[uint32, types.CurrentRewards]
	ServiceOutstandingRewards      collections.Map[uint32, types.OutstandingRewards]

	authority string
}

func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	communityPoolKeeper types.CommunityPoolKeeper,
	oracleKeeper types.OracleKeeper,
	poolsKeeper types.PoolsKeeper,
	operatorsKeeper types.OperatorsKeeper,
	servicesKeeper types.ServicesKeeper,
	restakingKeeper types.RestakingKeeper,
	tickersKeeper types.TickersKeeper,
	authority string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := &Keeper{
		cdc:                 cdc,
		storeService:        storeService,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		communityPoolKeeper: communityPoolKeeper,
		oracleKeeper:        oracleKeeper,
		poolsKeeper:         poolsKeeper,
		operatorsKeeper:     operatorsKeeper,
		servicesKeeper:      servicesKeeper,
		restakingKeeper:     restakingKeeper,
		tickersKeeper:       tickersKeeper,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		NextRewardsPlanID: collections.NewItem(
			sb, types.NextRewardsPlanIDKey, "next_rewards_plan_id", collections.Uint64Value),
		RewardsPlans: collections.NewMap(
			sb, types.RewardsPlanKeyPrefix, "rewards_plans",
			collections.Uint64Key, codec.CollValue[types.RewardsPlan](cdc)),
		LastRewardsAllocationTime: collections.NewItem(sb, types.LastRewardsAllocationTimeKey, "last_rewards_allocation_time",
			codec.CollValue[gogotypes.Timestamp](cdc)),
		DelegatorWithdrawAddrs: collections.NewMap(
			sb, types.DelegatorWithdrawAddrKeyPrefix, "delegator_withdraw_addrs",
			sdk.AccAddressKey,
			collcodec.KeyToValueCodec(sdk.AccAddressKey)),
		PoolDelegatorStartingInfos: collections.NewMap(
			sb, types.PoolDelegatorStartingInfoKeyPrefix, "pool_delegator_starting_infos",
			collections.PairKeyCodec(collections.Uint32Key, sdk.AccAddressKey),
			codec.CollValue[types.DelegatorStartingInfo](cdc)),
		PoolHistoricalRewards: collections.NewMap(
			sb, types.PoolHistoricalRewardsKeyPrefix, "pool_historical_rewards",
			collections.PairKeyCodec(collections.Uint32Key, collections.Uint64Key),
			codec.CollValue[types.HistoricalRewards](cdc)),
		PoolCurrentRewards: collections.NewMap(
			sb, types.PoolCurrentRewardsKeyPrefix, "pool_current_rewards",
			collections.Uint32Key, codec.CollValue[types.CurrentRewards](cdc)),
		PoolOutstandingRewards: collections.NewMap(
			sb, types.PoolOutstandingRewardsKeyPrefix, "pool_outstanding_rewards",
			collections.Uint32Key, codec.CollValue[types.OutstandingRewards](cdc)),
		OperatorAccumulatedCommissions: collections.NewMap(
			sb, types.OperatorAccumulatedCommissionKeyPrefix, "operator_accumulated_commissions",
			collections.Uint32Key, codec.CollValue[types.AccumulatedCommission](cdc)),
		OperatorDelegatorStartingInfos: collections.NewMap(
			sb, types.OperatorDelegatorStartingInfoKeyPrefix, "operator_delegator_starting_infos",
			collections.PairKeyCodec(collections.Uint32Key, sdk.AccAddressKey),
			codec.CollValue[types.DelegatorStartingInfo](cdc)),
		OperatorHistoricalRewards: collections.NewMap(
			sb, types.OperatorHistoricalRewardsKeyPrefix, "operator_historical_rewards",
			collections.PairKeyCodec(collections.Uint32Key, collections.Uint64Key),
			codec.CollValue[types.HistoricalRewards](cdc)),
		OperatorCurrentRewards: collections.NewMap(
			sb, types.OperatorCurrentRewardsKeyPrefix, "operator_current_rewards",
			collections.Uint32Key, codec.CollValue[types.CurrentRewards](cdc)),
		OperatorOutstandingRewards: collections.NewMap(
			sb, types.OperatorOutstandingRewardsKeyPrefix, "operator_outstanding_rewards",
			collections.Uint32Key, codec.CollValue[types.OutstandingRewards](cdc)),
		ServiceDelegatorStartingInfos: collections.NewMap(
			sb, types.ServiceDelegatorStartingInfoKeyPrefix, "service_delegator_starting_infos",
			collections.PairKeyCodec(collections.Uint32Key, sdk.AccAddressKey),
			codec.CollValue[types.DelegatorStartingInfo](cdc)),
		ServiceHistoricalRewards: collections.NewMap(
			sb, types.ServiceHistoricalRewardsKeyPrefix, "service_historical_rewards",
			collections.PairKeyCodec(collections.Uint32Key, collections.Uint64Key),
			codec.CollValue[types.HistoricalRewards](cdc)),
		ServiceCurrentRewards: collections.NewMap(
			sb, types.ServiceCurrentRewardsKeyPrefix, "service_current_rewards",
			collections.Uint32Key, codec.CollValue[types.CurrentRewards](cdc)),
		ServiceOutstandingRewards: collections.NewMap(
			sb, types.ServiceOutstandingRewardsKeyPrefix, "service_outstanding_rewards",
			collections.Uint32Key, codec.CollValue[types.OutstandingRewards](cdc)),

		authority: authority,
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

// GetAuthority returns the module's authority.
func (k *Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// SetWithdrawAddr sets a new address that will receive the rewards upon withdrawal
func (k *Keeper) SetWithdrawAddr(ctx context.Context, delegatorAddr, withdrawAddr sdk.AccAddress) error {
	if k.bankKeeper.BlockedAddr(withdrawAddr) {
		return errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive external funds", withdrawAddr)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetWithdrawAddress,
			sdk.NewAttribute(types.AttributeKeyWithdrawAddress, withdrawAddr.String()),
		),
	)

	err := k.DelegatorWithdrawAddrs.Set(ctx, delegatorAddr, withdrawAddr)
	if err != nil {
		return err
	}
	return nil
}

func (k *Keeper) WithdrawDelegationRewards(
	ctx context.Context, delAddr sdk.AccAddress, target *types.DelegationTarget) (types.Pools, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	del, found := k.GetDelegation(sdkCtx, target, delAddr)
	if !found {
		return nil, sdkerrors.ErrNotFound.Wrapf("delegation not found: %d, %s", target.GetID(), delAddr.String())
	}

	// withdraw rewards
	rewards, err := k.withdrawDelegationRewards(ctx, target, del)
	if err != nil {
		return nil, err
	}

	// reinitialize the delegation
	err = k.initializeDelegation(ctx, target.Type(), target.GetID(), delAddr)
	if err != nil {
		return nil, err
	}
	return rewards, nil
}

func (k *Keeper) WithdrawOperatorCommission(ctx context.Context, operatorID uint32) (types.Pools, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	operator, found := k.operatorsKeeper.GetOperator(sdkCtx, operatorID)
	if !found {
		return nil, operatorstypes.ErrOperatorNotFound
	}

	// fetch operator accumulated commission
	accumCommission, err := k.OperatorAccumulatedCommissions.Get(ctx, operatorID)
	if err != nil {
		return nil, err
	}
	if accumCommission.Commissions.IsEmpty() {
		return nil, types.ErrNoOperatorCommission
	}

	commissions, remainder := accumCommission.Commissions.TruncateDecimal()
	// leave remainder to withdraw later
	err = k.OperatorAccumulatedCommissions.Set(ctx, operatorID, types.AccumulatedCommission{
		Commissions: remainder,
	})
	if err != nil {
		return nil, err
	}

	// update outstanding
	outstanding, err := k.OperatorOutstandingRewards.Get(ctx, operatorID)
	if err != nil {
		return nil, err
	}

	err = k.OperatorOutstandingRewards.Set(ctx, operatorID, types.OutstandingRewards{
		Rewards: outstanding.Rewards.Sub(types.NewDecPoolsFromPools(commissions)),
	})
	if err != nil {
		return nil, err
	}

	commissionCoins := commissions.Sum()
	if !commissionCoins.IsZero() {
		adminAddr, err := k.accountKeeper.AddressCodec().StringToBytes(operator.Admin)
		if err != nil {
			return nil, err
		}
		withdrawAddr, err := k.GetDelegatorWithdrawAddr(ctx, adminAddr)
		if err != nil {
			return nil, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, commissionCoins)
		if err != nil {
			return nil, err
		}
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawCommission,
			sdk.NewAttribute(types.AttributeKeyDelegationType, restakingtypes.DELEGATION_TYPE_OPERATOR.String()),
			sdk.NewAttribute(types.AttributeKeyDelegationTargetID, fmt.Sprint(operatorID)),
			sdk.NewAttribute(sdk.AttributeKeyAmount, commissions.String()),
			sdk.NewAttribute(types.AttributeKeyAmountPerPool, commissions.String()),
		),
	)

	return commissions, nil
}
