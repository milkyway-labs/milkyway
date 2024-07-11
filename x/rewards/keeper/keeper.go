package keeper

import (
	"context"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"

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
	PoolDelegatorStartingInfos     collections.Map[collections.Pair[uint32, string], types.DelegatorStartingInfo]
	PoolHistoricalRewards          collections.Map[collections.Pair[uint32, uint64], types.HistoricalRewards]
	PoolCurrentRewards             collections.Map[uint32, types.CurrentRewards]
	PoolOutstandingRewards         collections.Map[uint32, types.OutstandingRewards]
	OperatorAccumulatedCommissions collections.Map[uint32, types.AccumulatedCommission]
	OperatorDelegatorStartingInfos collections.Map[collections.Pair[uint32, string], types.MultiDelegatorStartingInfo]
	OperatorHistoricalRewards      collections.Map[collections.Pair[uint32, uint64], types.MultiHistoricalRewards]
	OperatorCurrentRewards         collections.Map[uint32, types.MultiCurrentRewards]
	OperatorOutstandingRewards     collections.Map[uint32, types.MultiOutstandingRewards]
	ServiceDelegatorStartingInfos  collections.Map[collections.Pair[uint32, string], types.MultiDelegatorStartingInfo]
	ServiceHistoricalRewards       collections.Map[collections.Pair[uint32, uint64], types.MultiHistoricalRewards]
	ServiceCurrentRewards          collections.Map[uint32, types.MultiCurrentRewards]
	ServiceOutstandingRewards      collections.Map[uint32, types.MultiOutstandingRewards]

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
		PoolDelegatorStartingInfos: collections.NewMap(
			sb, types.PoolDelegatorStartingInfoKeyPrefix, "pool_delegator_starting_infos",
			collections.PairKeyCodec(collections.Uint32Key, collections.StringKey),
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
			collections.PairKeyCodec(collections.Uint32Key, collections.StringKey),
			codec.CollValue[types.MultiDelegatorStartingInfo](cdc)),
		OperatorHistoricalRewards: collections.NewMap(
			sb, types.OperatorHistoricalRewardsKeyPrefix, "operator_historical_rewards",
			collections.PairKeyCodec(collections.Uint32Key, collections.Uint64Key),
			codec.CollValue[types.MultiHistoricalRewards](cdc)),
		OperatorCurrentRewards: collections.NewMap(
			sb, types.OperatorCurrentRewardsKeyPrefix, "operator_current_rewards",
			collections.Uint32Key, codec.CollValue[types.MultiCurrentRewards](cdc)),
		OperatorOutstandingRewards: collections.NewMap(
			sb, types.OperatorOutstandingRewardsKeyPrefix, "operator_outstanding_rewards",
			collections.Uint32Key, codec.CollValue[types.MultiOutstandingRewards](cdc)),
		ServiceDelegatorStartingInfos: collections.NewMap(
			sb, types.ServiceDelegatorStartingInfoKeyPrefix, "service_delegator_starting_infos",
			collections.PairKeyCodec(collections.Uint32Key, collections.StringKey),
			codec.CollValue[types.MultiDelegatorStartingInfo](cdc)),
		ServiceHistoricalRewards: collections.NewMap(
			sb, types.ServiceHistoricalRewardsKeyPrefix, "service_historical_rewards",
			collections.PairKeyCodec(collections.Uint32Key, collections.Uint64Key),
			codec.CollValue[types.MultiHistoricalRewards](cdc)),
		ServiceCurrentRewards: collections.NewMap(
			sb, types.ServiceCurrentRewardsKeyPrefix, "service_current_rewards",
			collections.Uint32Key, codec.CollValue[types.MultiCurrentRewards](cdc)),
		ServiceOutstandingRewards: collections.NewMap(
			sb, types.ServiceOutstandingRewardsKeyPrefix, "service_outstanding_rewards",
			collections.Uint32Key, codec.CollValue[types.MultiOutstandingRewards](cdc)),

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
