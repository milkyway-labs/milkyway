# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [milkyway/rewards/v1/models.proto](#milkyway_rewards_v1_models-proto)
    - [AccumulatedCommission](#milkyway-rewards-v1-AccumulatedCommission)
    - [CurrentRewards](#milkyway-rewards-v1-CurrentRewards)
    - [DecPool](#milkyway-rewards-v1-DecPool)
    - [DelegationDelegatorReward](#milkyway-rewards-v1-DelegationDelegatorReward)
    - [DelegatorStartingInfo](#milkyway-rewards-v1-DelegatorStartingInfo)
    - [Distribution](#milkyway-rewards-v1-Distribution)
    - [DistributionTypeBasic](#milkyway-rewards-v1-DistributionTypeBasic)
    - [DistributionTypeEgalitarian](#milkyway-rewards-v1-DistributionTypeEgalitarian)
    - [DistributionTypeWeighted](#milkyway-rewards-v1-DistributionTypeWeighted)
    - [DistributionWeight](#milkyway-rewards-v1-DistributionWeight)
    - [HistoricalRewards](#milkyway-rewards-v1-HistoricalRewards)
    - [OutstandingRewards](#milkyway-rewards-v1-OutstandingRewards)
    - [Pool](#milkyway-rewards-v1-Pool)
    - [PoolServiceTotalDelegatorShares](#milkyway-rewards-v1-PoolServiceTotalDelegatorShares)
    - [RewardsPlan](#milkyway-rewards-v1-RewardsPlan)
    - [ServicePool](#milkyway-rewards-v1-ServicePool)
    - [UsersDistribution](#milkyway-rewards-v1-UsersDistribution)
    - [UsersDistributionTypeBasic](#milkyway-rewards-v1-UsersDistributionTypeBasic)
  
- [milkyway/rewards/v1/params.proto](#milkyway_rewards_v1_params-proto)
    - [Params](#milkyway-rewards-v1-Params)
  
- [milkyway/rewards/v1/genesis.proto](#milkyway_rewards_v1_genesis-proto)
    - [CurrentRewardsRecord](#milkyway-rewards-v1-CurrentRewardsRecord)
    - [DelegationTypeRecords](#milkyway-rewards-v1-DelegationTypeRecords)
    - [DelegatorStartingInfoRecord](#milkyway-rewards-v1-DelegatorStartingInfoRecord)
    - [DelegatorWithdrawInfo](#milkyway-rewards-v1-DelegatorWithdrawInfo)
    - [GenesisState](#milkyway-rewards-v1-GenesisState)
    - [HistoricalRewardsRecord](#milkyway-rewards-v1-HistoricalRewardsRecord)
    - [OperatorAccumulatedCommissionRecord](#milkyway-rewards-v1-OperatorAccumulatedCommissionRecord)
    - [OutstandingRewardsRecord](#milkyway-rewards-v1-OutstandingRewardsRecord)
  
- [milkyway/rewards/v1/messages.proto](#milkyway_rewards_v1_messages-proto)
    - [MsgCreateRewardsPlan](#milkyway-rewards-v1-MsgCreateRewardsPlan)
    - [MsgCreateRewardsPlanResponse](#milkyway-rewards-v1-MsgCreateRewardsPlanResponse)
    - [MsgEditRewardsPlan](#milkyway-rewards-v1-MsgEditRewardsPlan)
    - [MsgEditRewardsPlanResponse](#milkyway-rewards-v1-MsgEditRewardsPlanResponse)
    - [MsgSetWithdrawAddress](#milkyway-rewards-v1-MsgSetWithdrawAddress)
    - [MsgSetWithdrawAddressResponse](#milkyway-rewards-v1-MsgSetWithdrawAddressResponse)
    - [MsgUpdateParams](#milkyway-rewards-v1-MsgUpdateParams)
    - [MsgUpdateParamsResponse](#milkyway-rewards-v1-MsgUpdateParamsResponse)
    - [MsgWithdrawDelegatorReward](#milkyway-rewards-v1-MsgWithdrawDelegatorReward)
    - [MsgWithdrawDelegatorRewardResponse](#milkyway-rewards-v1-MsgWithdrawDelegatorRewardResponse)
    - [MsgWithdrawOperatorCommission](#milkyway-rewards-v1-MsgWithdrawOperatorCommission)
    - [MsgWithdrawOperatorCommissionResponse](#milkyway-rewards-v1-MsgWithdrawOperatorCommissionResponse)
  
    - [Msg](#milkyway-rewards-v1-Msg)
  
- [milkyway/rewards/v1/query.proto](#milkyway_rewards_v1_query-proto)
    - [QueryDelegatorTotalRewardsRequest](#milkyway-rewards-v1-QueryDelegatorTotalRewardsRequest)
    - [QueryDelegatorTotalRewardsResponse](#milkyway-rewards-v1-QueryDelegatorTotalRewardsResponse)
    - [QueryDelegatorWithdrawAddressRequest](#milkyway-rewards-v1-QueryDelegatorWithdrawAddressRequest)
    - [QueryDelegatorWithdrawAddressResponse](#milkyway-rewards-v1-QueryDelegatorWithdrawAddressResponse)
    - [QueryOperatorCommissionRequest](#milkyway-rewards-v1-QueryOperatorCommissionRequest)
    - [QueryOperatorCommissionResponse](#milkyway-rewards-v1-QueryOperatorCommissionResponse)
    - [QueryOperatorDelegationRewardsRequest](#milkyway-rewards-v1-QueryOperatorDelegationRewardsRequest)
    - [QueryOperatorDelegationRewardsResponse](#milkyway-rewards-v1-QueryOperatorDelegationRewardsResponse)
    - [QueryOperatorOutstandingRewardsRequest](#milkyway-rewards-v1-QueryOperatorOutstandingRewardsRequest)
    - [QueryOperatorOutstandingRewardsResponse](#milkyway-rewards-v1-QueryOperatorOutstandingRewardsResponse)
    - [QueryParamsRequest](#milkyway-rewards-v1-QueryParamsRequest)
    - [QueryParamsResponse](#milkyway-rewards-v1-QueryParamsResponse)
    - [QueryPoolDelegationRewardsRequest](#milkyway-rewards-v1-QueryPoolDelegationRewardsRequest)
    - [QueryPoolDelegationRewardsResponse](#milkyway-rewards-v1-QueryPoolDelegationRewardsResponse)
    - [QueryPoolOutstandingRewardsRequest](#milkyway-rewards-v1-QueryPoolOutstandingRewardsRequest)
    - [QueryPoolOutstandingRewardsResponse](#milkyway-rewards-v1-QueryPoolOutstandingRewardsResponse)
    - [QueryRewardsPlanRequest](#milkyway-rewards-v1-QueryRewardsPlanRequest)
    - [QueryRewardsPlanResponse](#milkyway-rewards-v1-QueryRewardsPlanResponse)
    - [QueryRewardsPlansRequest](#milkyway-rewards-v1-QueryRewardsPlansRequest)
    - [QueryRewardsPlansResponse](#milkyway-rewards-v1-QueryRewardsPlansResponse)
    - [QueryServiceDelegationRewardsRequest](#milkyway-rewards-v1-QueryServiceDelegationRewardsRequest)
    - [QueryServiceDelegationRewardsResponse](#milkyway-rewards-v1-QueryServiceDelegationRewardsResponse)
    - [QueryServiceOutstandingRewardsRequest](#milkyway-rewards-v1-QueryServiceOutstandingRewardsRequest)
    - [QueryServiceOutstandingRewardsResponse](#milkyway-rewards-v1-QueryServiceOutstandingRewardsResponse)
  
    - [Query](#milkyway-rewards-v1-Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="milkyway_rewards_v1_models-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/rewards/v1/models.proto



<a name="milkyway-rewards-v1-AccumulatedCommission"></a>

### AccumulatedCommission
AccumulatedCommission represents accumulated commission
for a delegation target kept as a running counter, can be withdrawn at any
time.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commissions | [DecPool](#milkyway-rewards-v1-DecPool) | repeated |  |






<a name="milkyway-rewards-v1-CurrentRewards"></a>

### CurrentRewards
CurrentRewards represents current rewards and current
period for a delegation target kept as a running counter and incremented
each block as long as the delegation target&#39;s tokens remain constant.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards | [ServicePool](#milkyway-rewards-v1-ServicePool) | repeated |  |
| period | [uint64](#uint64) |  |  |






<a name="milkyway-rewards-v1-DecPool"></a>

### DecPool
DecPool is a DecCoins wrapper with denom which represents the rewards pool
for the given denom. It is used to represent the rewards associated with the
denom.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| denom | [string](#string) |  |  |
| dec_coins | [cosmos.base.v1beta1.DecCoin](#cosmos-base-v1beta1-DecCoin) | repeated |  |






<a name="milkyway-rewards-v1-DelegationDelegatorReward"></a>

### DelegationDelegatorReward
DelegationDelegatorReward represents the properties of a delegator&#39;s
delegation reward. The delegator address implicit in the within the
query request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation_type | [milkyway.restaking.v1.DelegationType](#milkyway-restaking-v1-DelegationType) |  |  |
| delegation_target_id | [uint32](#uint32) |  |  |
| reward | [DecPool](#milkyway-rewards-v1-DecPool) | repeated |  |






<a name="milkyway-rewards-v1-DelegatorStartingInfo"></a>

### DelegatorStartingInfo
DelegatorStartingInfo represents the starting info for a delegator reward
period. It tracks the previous delegation target period, the delegation&#39;s
amount of staking token, and the creation height (to check later on if any
slashes have occurred). NOTE: Even though validators are slashed to whole
staking tokens, the delegators within the validator may be left with less
than a full token, thus sdk.Dec is used.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| previous_period | [uint64](#uint64) |  |  |
| stakes | [cosmos.base.v1beta1.DecCoin](#cosmos-base-v1beta1-DecCoin) | repeated |  |
| height | [uint64](#uint64) |  |  |






<a name="milkyway-rewards-v1-Distribution"></a>

### Distribution
Distribution represents distribution parameters for restaking
pools/operators.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation_type | [milkyway.restaking.v1.DelegationType](#milkyway-restaking-v1-DelegationType) |  | DelegationType is the type of delegation target which this distribution parameters are for. It can be one of DELEGATION_TYPE_POOL and DELEGATION_TYPE_OPERATOR. |
| weight | [uint32](#uint32) |  | Weight is the rewards distribution weight among other types of delegation targets. |
| type | [google.protobuf.Any](#google-protobuf-Any) |  | Type is one of basic/weighted/egalitarian distributions. |






<a name="milkyway-rewards-v1-DistributionTypeBasic"></a>

### DistributionTypeBasic
DistributionTypeBasic represents the simplest form of distribution.
Rewards are allocated to entities based on their delegation values.
For example, if there are three operators with delegation values of
$1000, $1500, and $2000, their rewards will be distributed in a
2:3:4 ratio.






<a name="milkyway-rewards-v1-DistributionTypeEgalitarian"></a>

### DistributionTypeEgalitarian
DistributionTypeEgalitarian is a distribution method where all entities
receive an equal share of rewards(a.k.a. egalitarian method).






<a name="milkyway-rewards-v1-DistributionTypeWeighted"></a>

### DistributionTypeWeighted
DistributionTypeWeighted is a type of distribution where the reward
weights for each entity are explicitly defined. Only the specified
delegation targets will receive rewards.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| weights | [DistributionWeight](#milkyway-rewards-v1-DistributionWeight) | repeated |  |






<a name="milkyway-rewards-v1-DistributionWeight"></a>

### DistributionWeight
DistributionWeight defines a delegation target and its assigned weight.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation_target_id | [uint32](#uint32) |  |  |
| weight | [uint32](#uint32) |  |  |






<a name="milkyway-rewards-v1-HistoricalRewards"></a>

### HistoricalRewards
HistoricalRewards represents historical rewards for a delegation target.
Height is implicit within the store key.
Cumulative reward ratio is the sum from the zeroeth period
until this period of rewards / tokens, per the spec.
The reference count indicates the number of objects
which might need to reference this historical entry at any point.
ReferenceCount =
   number of outstanding delegations which ended the associated period (and
   might need to read that record)
 &#43; number of slashes which ended the associated period (and might need to
 read that record)
 &#43; one per validator for the zeroeth period, set on initialization


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cumulative_reward_ratios | [ServicePool](#milkyway-rewards-v1-ServicePool) | repeated |  |
| reference_count | [uint32](#uint32) |  |  |






<a name="milkyway-rewards-v1-OutstandingRewards"></a>

### OutstandingRewards
OutstandingRewards represents outstanding (un-withdrawn) rewards
for a delegation target inexpensive to track, allows simple sanity checks.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards | [DecPool](#milkyway-rewards-v1-DecPool) | repeated |  |






<a name="milkyway-rewards-v1-Pool"></a>

### Pool
Pool is a Coins wrapper with denom which represents the rewards pool for the
given denom. It is used to represent the rewards associated with the denom.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| denom | [string](#string) |  |  |
| coins | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated |  |






<a name="milkyway-rewards-v1-PoolServiceTotalDelegatorShares"></a>

### PoolServiceTotalDelegatorShares
PoolServiceTotalDelegatorShares represents the total delegator shares for a
pool-service pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_id | [uint32](#uint32) |  |  |
| service_id | [uint32](#uint32) |  |  |
| shares | [cosmos.base.v1beta1.DecCoin](#cosmos-base-v1beta1-DecCoin) | repeated |  |






<a name="milkyway-rewards-v1-RewardsPlan"></a>

### RewardsPlan
RewardsPlan represents a rewards allocation plan.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint64](#uint64) |  | ID is the unique identifier of the plan. |
| description | [string](#string) |  | Description is the description of the plan. |
| service_id | [uint32](#uint32) |  | ServiceID is the service ID which the plan is related to. |
| amount_per_day | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | AmountPerDay is the amount of rewards to be distributed, per day. The rewards amount for every block will be calculated based on this. |
| start_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | StartTime is the starting time of the plan. |
| end_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | EndTime is the ending time of the plan. |
| rewards_pool | [string](#string) |  | RewardsPool is the address where rewards to be distributed are stored. If the rewards pool doesn&#39;t have enough funds to be distributed, then the rewards allocation for this plan will be skipped. |
| pools_distribution | [Distribution](#milkyway-rewards-v1-Distribution) |  | PoolsDistribution is the rewards distribution parameters for pools. |
| operators_distribution | [Distribution](#milkyway-rewards-v1-Distribution) |  | OperatorsDistribution is the rewards distribution parameters for operators. |
| users_distribution | [UsersDistribution](#milkyway-rewards-v1-UsersDistribution) |  | UsersDistribution is the rewards distribution parameters for users. |






<a name="milkyway-rewards-v1-ServicePool"></a>

### ServicePool
ServicePool represents the rewards pool for a service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  |  |
| dec_pools | [DecPool](#milkyway-rewards-v1-DecPool) | repeated |  |






<a name="milkyway-rewards-v1-UsersDistribution"></a>

### UsersDistribution
Distribution represents distribution parameters for delegators who directly
staked their tokens to the service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| weight | [uint32](#uint32) |  | Weight is the rewards distribution weight among other types of delegation targets. |
| type | [google.protobuf.Any](#google-protobuf-Any) |  | Type defines the rewards distribution method. Currently only the basic distribution is allowed. |






<a name="milkyway-rewards-v1-UsersDistributionTypeBasic"></a>

### UsersDistributionTypeBasic
UsersDistributionTypeBasic represents the simplest form of distribution.
Rewards are allocated to entities based on their delegation values.
For example, if there are three users with delegation values of
$1000, $1500, and $2000, their rewards will be distributed in a
2:3:4 ratio.





 

 

 

 



<a name="milkyway_rewards_v1_params-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/rewards/v1/params.proto



<a name="milkyway-rewards-v1-Params"></a>

### Params
Params defines the parameters for the module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards_plan_creation_fee | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | RewardsPlanCreationFee represents the fee that an account must pay in order to create a rewards plan. The fee is drawn from the MsgCreateRewardsPlan sender&#39;s account and transferred to the community pool. |





 

 

 

 



<a name="milkyway_rewards_v1_genesis-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/rewards/v1/genesis.proto



<a name="milkyway-rewards-v1-CurrentRewardsRecord"></a>

### CurrentRewardsRecord
CurrentRewardsRecord is used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation_target_id | [uint32](#uint32) |  | delegation_target_id is the ID of the delegation target. |
| rewards | [CurrentRewards](#milkyway-rewards-v1-CurrentRewards) |  | rewards defines the current rewards of the delegation target. |






<a name="milkyway-rewards-v1-DelegationTypeRecords"></a>

### DelegationTypeRecords
DelegationTypeRecords groups various genesis records under the same type
of delegation target.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| outstanding_rewards | [OutstandingRewardsRecord](#milkyway-rewards-v1-OutstandingRewardsRecord) | repeated | outstanding_rewards defines the outstanding rewards of all delegation targets with the same delegation type at genesis. |
| historical_rewards | [HistoricalRewardsRecord](#milkyway-rewards-v1-HistoricalRewardsRecord) | repeated | historical_rewards defines the historical rewards of all delegation targets with the same delegation type at genesis. |
| current_rewards | [CurrentRewardsRecord](#milkyway-rewards-v1-CurrentRewardsRecord) | repeated | current_rewards defines the current rewards of all delegation targets with the same delegation type at genesis. |
| delegator_starting_infos | [DelegatorStartingInfoRecord](#milkyway-rewards-v1-DelegatorStartingInfoRecord) | repeated | delegator_starting_infos defines the delegator starting infos of all delegation targets with the same delegation type at genesis. |






<a name="milkyway-rewards-v1-DelegatorStartingInfoRecord"></a>

### DelegatorStartingInfoRecord
DelegatorStartingInfoRecord used for import / export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | delegator_address is the address of the delegator. |
| delegation_target_id | [uint32](#uint32) |  | delegation_target_id is the ID of the delegation target. |
| starting_info | [DelegatorStartingInfo](#milkyway-rewards-v1-DelegatorStartingInfo) |  | starting_info defines the starting info of a delegator. |






<a name="milkyway-rewards-v1-DelegatorWithdrawInfo"></a>

### DelegatorWithdrawInfo
DelegatorWithdrawInfo is the address for where delegation rewards are
withdrawn to by default this struct is only used at genesis to feed in
default withdraw addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | delegator_address is the address of the delegator. |
| withdraw_address | [string](#string) |  | withdraw_address is the address to withdraw the delegation rewards to. |






<a name="milkyway-rewards-v1-GenesisState"></a>

### GenesisState
GenesisState defines the module&#39;s genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-rewards-v1-Params) |  | Params defines the parameters of the module. |
| next_rewards_plan_id | [uint64](#uint64) |  | NextRewardsPlanID represents the id to be used when creating the next rewards plan. |
| rewards_plans | [RewardsPlan](#milkyway-rewards-v1-RewardsPlan) | repeated | RewardsPlans defines the list of rewards plans. |
| last_rewards_allocation_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | last_rewards_allocation_time is the last time rewards were allocated. |
| delegator_withdraw_infos | [DelegatorWithdrawInfo](#milkyway-rewards-v1-DelegatorWithdrawInfo) | repeated | delegator_withdraw_infos defines the delegator withdraw infos at genesis. |
| pools_records | [DelegationTypeRecords](#milkyway-rewards-v1-DelegationTypeRecords) |  | pools_records defines a group of genesis records of all pools at genesis. |
| operators_records | [DelegationTypeRecords](#milkyway-rewards-v1-DelegationTypeRecords) |  | operators_records defines a group of genesis records of all operators at genesis. |
| services_records | [DelegationTypeRecords](#milkyway-rewards-v1-DelegationTypeRecords) |  | services_records defines a group of genesis records of all services at genesis. |
| operator_accumulated_commissions | [OperatorAccumulatedCommissionRecord](#milkyway-rewards-v1-OperatorAccumulatedCommissionRecord) | repeated | operator_accumulated_commissions defines the accumulated commissions of all operators at genesis. |
| pool_service_total_delegator_shares | [PoolServiceTotalDelegatorShares](#milkyway-rewards-v1-PoolServiceTotalDelegatorShares) | repeated | pool_service_total_delegator_shares defines the total delegator shares at genesis. |






<a name="milkyway-rewards-v1-HistoricalRewardsRecord"></a>

### HistoricalRewardsRecord
HistoricalRewardsRecord is used for import / export via genesis
json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation_target_id | [uint32](#uint32) |  | delegation_target_id is the ID of the delegation target. |
| period | [uint64](#uint64) |  | period defines the period the historical rewards apply to. |
| rewards | [HistoricalRewards](#milkyway-rewards-v1-HistoricalRewards) |  | rewards defines the historical rewards of the delegation target. |






<a name="milkyway-rewards-v1-OperatorAccumulatedCommissionRecord"></a>

### OperatorAccumulatedCommissionRecord
OperatorAccumulatedCommissionRecord contains the data about the accumulated commission of an operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  |  |
| accumulated | [AccumulatedCommission](#milkyway-rewards-v1-AccumulatedCommission) |  | accumulated is the accumulated commission of an operator. |






<a name="milkyway-rewards-v1-OutstandingRewardsRecord"></a>

### OutstandingRewardsRecord
OutstandingRewardsRecord is used for import/export via genesis json.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation_target_id | [uint32](#uint32) |  | delegation_target_id is the ID of the delegation target. |
| outstanding_rewards | [DecPool](#milkyway-rewards-v1-DecPool) | repeated | outstanding_rewards represents the outstanding rewards of the delegation target. |





 

 

 

 



<a name="milkyway_rewards_v1_messages-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/rewards/v1/messages.proto



<a name="milkyway-rewards-v1-MsgCreateRewardsPlan"></a>

### MsgCreateRewardsPlan
MsgCreateRewardsPlan defines the message structure for the
CreateRewardsPlan gRPC service method. It allows an account to create a
new rewards plan. It requires a sender address as well as the details of
the plan to be created.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user creating the rewards plan |
| description | [string](#string) |  |  |
| service_id | [uint32](#uint32) |  |  |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Amount is the amount of rewards to be distributed. |
| start_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | StartTime is the starting time of the plan. |
| end_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | EndTime is the ending time of the plan. |
| pools_distribution | [Distribution](#milkyway-rewards-v1-Distribution) |  | PoolsDistribution is the rewards distribution parameters for pools. |
| operators_distribution | [Distribution](#milkyway-rewards-v1-Distribution) |  | OperatorsDistribution is the rewards distribution parameters for operators. |
| users_distribution | [UsersDistribution](#milkyway-rewards-v1-UsersDistribution) |  | UsersDistribution is the rewards distribution parameters for users who delegated directly to the service. |
| fee_amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | FeeAmount represents the fees that are going to be paid to create the rewards plan. These should always be greater or equals of any of the coins specified inside the RewardsPlanCreationFee field of the modules params. If no fees are specified inside the module parameters, this field can be omitted. |






<a name="milkyway-rewards-v1-MsgCreateRewardsPlanResponse"></a>

### MsgCreateRewardsPlanResponse
MsgCreateRewardsPlanResponse is the return value of
MsgCreateRewardsPlan. It returns the newly created plan ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| new_rewards_plan_id | [uint64](#uint64) |  | NewRewardsPlanID is the ID of the newly created rewards plan |






<a name="milkyway-rewards-v1-MsgEditRewardsPlan"></a>

### MsgEditRewardsPlan
MsgEditRewardsPlan defines the message structure for the
EditRewardsPlan gRPC service method. It allows an account to edit a
previously created rewards plan.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user editing the rewards plan. |
| id | [uint64](#uint64) |  | ID is the ID of the rewards plan to be edited. |
| description | [string](#string) |  |  |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Amount is the amount of rewards to be distributed. |
| start_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | StartTime is the starting time of the plan. |
| end_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | EndTime is the ending time of the plan. |
| pools_distribution | [Distribution](#milkyway-rewards-v1-Distribution) |  | PoolsDistribution is the rewards distribution parameters for pools. |
| operators_distribution | [Distribution](#milkyway-rewards-v1-Distribution) |  | OperatorsDistribution is the rewards distribution parameters for operators. |
| users_distribution | [UsersDistribution](#milkyway-rewards-v1-UsersDistribution) |  | UsersDistribution is the rewards distribution parameters for users who delegated directly to the service. |






<a name="milkyway-rewards-v1-MsgEditRewardsPlanResponse"></a>

### MsgEditRewardsPlanResponse
MsgEditRewardsPlanResponse is the return value of
MsgEditRewardsPlan.






<a name="milkyway-rewards-v1-MsgSetWithdrawAddress"></a>

### MsgSetWithdrawAddress
MsgSetWithdrawAddress sets the withdraw address for a delegator(or an
operator when withdrawing commission).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  |  |
| withdraw_address | [string](#string) |  |  |






<a name="milkyway-rewards-v1-MsgSetWithdrawAddressResponse"></a>

### MsgSetWithdrawAddressResponse
MsgSetWithdrawAddressResponse defines the Msg/SetWithdrawAddress response
type.






<a name="milkyway-rewards-v1-MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams defines the message structure for the UpdateParams gRPC
service method. It allows the authority to update the module parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authority | [string](#string) |  | Authority is the address that controls the module (defaults to x/gov unless overwritten). |
| params | [Params](#milkyway-rewards-v1-Params) |  | Params define the parameters to update.

NOTE: All parameters must be supplied. |






<a name="milkyway-rewards-v1-MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse is the return value of MsgUpdateParams.






<a name="milkyway-rewards-v1-MsgWithdrawDelegatorReward"></a>

### MsgWithdrawDelegatorReward
MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator
from a single delegation target.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  |  |
| delegation_type | [milkyway.restaking.v1.DelegationType](#milkyway-restaking-v1-DelegationType) |  |  |
| delegation_target_id | [uint32](#uint32) |  |  |






<a name="milkyway-rewards-v1-MsgWithdrawDelegatorRewardResponse"></a>

### MsgWithdrawDelegatorRewardResponse
MsgWithdrawDelegatorRewardResponse defines the Msg/WithdrawDelegatorReward
response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated |  |






<a name="milkyway-rewards-v1-MsgWithdrawOperatorCommission"></a>

### MsgWithdrawOperatorCommission
MsgWithdrawOperatorCommission withdraws the full commission to the operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  |  |
| operator_id | [uint32](#uint32) |  |  |






<a name="milkyway-rewards-v1-MsgWithdrawOperatorCommissionResponse"></a>

### MsgWithdrawOperatorCommissionResponse
MsgWithdrawOperatorCommissionResponse defines the
Msg/WithdrawOperatorCommission response type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Since: cosmos-sdk 0.46 |





 

 

 


<a name="milkyway-rewards-v1-Msg"></a>

### Msg
Msg defines the services module&#39;s gRPC message service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateRewardsPlan | [MsgCreateRewardsPlan](#milkyway-rewards-v1-MsgCreateRewardsPlan) | [MsgCreateRewardsPlanResponse](#milkyway-rewards-v1-MsgCreateRewardsPlanResponse) | CreateRewardsPlan defines the operation for creating a new rewards plan. |
| EditRewardsPlan | [MsgEditRewardsPlan](#milkyway-rewards-v1-MsgEditRewardsPlan) | [MsgEditRewardsPlanResponse](#milkyway-rewards-v1-MsgEditRewardsPlanResponse) | EditRewardsPlan defines the operation to edit an existing rewards plan. |
| SetWithdrawAddress | [MsgSetWithdrawAddress](#milkyway-rewards-v1-MsgSetWithdrawAddress) | [MsgSetWithdrawAddressResponse](#milkyway-rewards-v1-MsgSetWithdrawAddressResponse) | SetWithdrawAddress defines a method to change the withdraw address for a delegator(or an operator, when withdrawing commission). |
| WithdrawDelegatorReward | [MsgWithdrawDelegatorReward](#milkyway-rewards-v1-MsgWithdrawDelegatorReward) | [MsgWithdrawDelegatorRewardResponse](#milkyway-rewards-v1-MsgWithdrawDelegatorRewardResponse) | WithdrawDelegatorReward defines a method to withdraw rewards of delegator from a single delegation target. |
| WithdrawOperatorCommission | [MsgWithdrawOperatorCommission](#milkyway-rewards-v1-MsgWithdrawOperatorCommission) | [MsgWithdrawOperatorCommissionResponse](#milkyway-rewards-v1-MsgWithdrawOperatorCommissionResponse) | WithdrawOperatorCommission defines a method to withdraw the full commission to the operator. |
| UpdateParams | [MsgUpdateParams](#milkyway-rewards-v1-MsgUpdateParams) | [MsgUpdateParamsResponse](#milkyway-rewards-v1-MsgUpdateParamsResponse) | UpdateParams defines a (governance) operation for updating the module parameters. The authority defaults to the x/gov module account. |

 



<a name="milkyway_rewards_v1_query-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/rewards/v1/query.proto



<a name="milkyway-rewards-v1-QueryDelegatorTotalRewardsRequest"></a>

### QueryDelegatorTotalRewardsRequest
QueryDelegatorTotalRewardsRequest is the request type for the
Query/DelegatorTotalRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="milkyway-rewards-v1-QueryDelegatorTotalRewardsResponse"></a>

### QueryDelegatorTotalRewardsResponse
QueryDelegatorTotalRewardsResponse is the response type for the
Query/DelegatorTotalRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards | [DelegationDelegatorReward](#milkyway-rewards-v1-DelegationDelegatorReward) | repeated | rewards defines all the rewards accrued by a delegator. |
| total | [DecPool](#milkyway-rewards-v1-DecPool) | repeated | total defines the sum of all the rewards. |






<a name="milkyway-rewards-v1-QueryDelegatorWithdrawAddressRequest"></a>

### QueryDelegatorWithdrawAddressRequest
QueryDelegatorWithdrawAddressRequest is the request type for the
Query/DelegatorWithdrawAddress RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | delegator_address defines the delegator address to query for. |






<a name="milkyway-rewards-v1-QueryDelegatorWithdrawAddressResponse"></a>

### QueryDelegatorWithdrawAddressResponse
QueryDelegatorWithdrawAddressResponse is the response type for the
Query/DelegatorWithdrawAddress RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| withdraw_address | [string](#string) |  | withdraw_address defines the delegator address to query for. |






<a name="milkyway-rewards-v1-QueryOperatorCommissionRequest"></a>

### QueryOperatorCommissionRequest
QueryOperatorCommissionRequest is the request type for the
Query/OperatorCommission RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | operator_id defines the validator address to query for. |






<a name="milkyway-rewards-v1-QueryOperatorCommissionResponse"></a>

### QueryOperatorCommissionResponse
QueryOperatorCommissionResponse is the response type for the
Query/OperatorCommission RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commission | [AccumulatedCommission](#milkyway-rewards-v1-AccumulatedCommission) |  | commission defines the commission the operator received. |






<a name="milkyway-rewards-v1-QueryOperatorDelegationRewardsRequest"></a>

### QueryOperatorDelegationRewardsRequest
QueryOperatorDelegationRewardsRequest is the request type for the
Query/OperatorDelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | delegator_address defines the delegator address to query for. |
| operator_id | [uint32](#uint32) |  | operator_id defines the operator ID to query for. |






<a name="milkyway-rewards-v1-QueryOperatorDelegationRewardsResponse"></a>

### QueryOperatorDelegationRewardsResponse
QueryOperatorDelegationRewardsResponse is the response type for the
Query/OperatorDelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards | [DecPool](#milkyway-rewards-v1-DecPool) | repeated | rewards defines the rewards accrued by a delegation. |






<a name="milkyway-rewards-v1-QueryOperatorOutstandingRewardsRequest"></a>

### QueryOperatorOutstandingRewardsRequest
QueryOperatorOutstandingRewardsRequest is the request type for the
Query/OperatorOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | operator_id defines the operator ID to query for. |






<a name="milkyway-rewards-v1-QueryOperatorOutstandingRewardsResponse"></a>

### QueryOperatorOutstandingRewardsResponse
QueryOperatorOutstandingRewardsResponse is the response type for the
Query/OperatorOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards | [OutstandingRewards](#milkyway-rewards-v1-OutstandingRewards) |  |  |






<a name="milkyway-rewards-v1-QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="milkyway-rewards-v1-QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-rewards-v1-Params) |  |  |






<a name="milkyway-rewards-v1-QueryPoolDelegationRewardsRequest"></a>

### QueryPoolDelegationRewardsRequest
QueryPoolDelegationRewardsRequest is the request type for the
Query/PoolDelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | delegator_address defines the delegator address to query for. |
| pool_id | [uint32](#uint32) |  | pool_id defines the pool ID to query for. |






<a name="milkyway-rewards-v1-QueryPoolDelegationRewardsResponse"></a>

### QueryPoolDelegationRewardsResponse
QueryPoolDelegationRewardsResponse is the response type for the
Query/PoolDelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards | [DecPool](#milkyway-rewards-v1-DecPool) | repeated | rewards defines the rewards accrued by a delegation. |






<a name="milkyway-rewards-v1-QueryPoolOutstandingRewardsRequest"></a>

### QueryPoolOutstandingRewardsRequest
QueryPoolOutstandingRewardsRequest is the request type for the
Query/PoolOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_id | [uint32](#uint32) |  | pool_id defines the pool ID to query for. |






<a name="milkyway-rewards-v1-QueryPoolOutstandingRewardsResponse"></a>

### QueryPoolOutstandingRewardsResponse
QueryPoolOutstandingRewardsResponse is the response type for the
Query/PoolOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards | [OutstandingRewards](#milkyway-rewards-v1-OutstandingRewards) |  |  |






<a name="milkyway-rewards-v1-QueryRewardsPlanRequest"></a>

### QueryRewardsPlanRequest
QueryRewardsPlanRequest is the request type for the Query/RewardsPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| plan_id | [uint64](#uint64) |  |  |






<a name="milkyway-rewards-v1-QueryRewardsPlanResponse"></a>

### QueryRewardsPlanResponse
QueryRewardsPlanResponse is the response type for the Query/RewardsPlan RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards_plan | [RewardsPlan](#milkyway-rewards-v1-RewardsPlan) |  |  |






<a name="milkyway-rewards-v1-QueryRewardsPlansRequest"></a>

### QueryRewardsPlansRequest
QueryRewardsPlansRequest is the request type for the Query/RewardsPlans RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="milkyway-rewards-v1-QueryRewardsPlansResponse"></a>

### QueryRewardsPlansResponse
QueryRewardsPlansResponse is the response type for the Query/RewardsPlans
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards_plans | [RewardsPlan](#milkyway-rewards-v1-RewardsPlan) | repeated |  |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | pagination defines the pagination in the response. |






<a name="milkyway-rewards-v1-QueryServiceDelegationRewardsRequest"></a>

### QueryServiceDelegationRewardsRequest
QueryServiceDelegationRewardsRequest is the request type for the
Query/ServiceDelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | delegator_address defines the delegator address to query for. |
| service_id | [uint32](#uint32) |  | service_id defines the service ID to query for. |






<a name="milkyway-rewards-v1-QueryServiceDelegationRewardsResponse"></a>

### QueryServiceDelegationRewardsResponse
QueryServiceDelegationRewardsResponse is the response type for the
Query/ServiceDelegationRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards | [DecPool](#milkyway-rewards-v1-DecPool) | repeated | rewards defines the rewards accrued by a delegation. |






<a name="milkyway-rewards-v1-QueryServiceOutstandingRewardsRequest"></a>

### QueryServiceOutstandingRewardsRequest
QueryServiceOutstandingRewardsRequest is the request type for the
Query/ServiceOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | service_id defines the service ID to query for. |






<a name="milkyway-rewards-v1-QueryServiceOutstandingRewardsResponse"></a>

### QueryServiceOutstandingRewardsResponse
QueryServiceOutstandingRewardsResponse is the response type for the
Query/ServiceOutstandingRewards RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rewards | [OutstandingRewards](#milkyway-rewards-v1-OutstandingRewards) |  |  |





 

 

 


<a name="milkyway-rewards-v1-Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Params | [QueryParamsRequest](#milkyway-rewards-v1-QueryParamsRequest) | [QueryParamsResponse](#milkyway-rewards-v1-QueryParamsResponse) | Params defines a gRPC query method that returns the parameters of the module. |
| RewardsPlans | [QueryRewardsPlansRequest](#milkyway-rewards-v1-QueryRewardsPlansRequest) | [QueryRewardsPlansResponse](#milkyway-rewards-v1-QueryRewardsPlansResponse) | RewardsPlans queries all rewards plans. |
| RewardsPlan | [QueryRewardsPlanRequest](#milkyway-rewards-v1-QueryRewardsPlanRequest) | [QueryRewardsPlanResponse](#milkyway-rewards-v1-QueryRewardsPlanResponse) | RewardsPlan queries a specific rewards plan by its ID. |
| PoolOutstandingRewards | [QueryPoolOutstandingRewardsRequest](#milkyway-rewards-v1-QueryPoolOutstandingRewardsRequest) | [QueryPoolOutstandingRewardsResponse](#milkyway-rewards-v1-QueryPoolOutstandingRewardsResponse) | PoolOutstandingRewards queries rewards of a pool. |
| OperatorOutstandingRewards | [QueryOperatorOutstandingRewardsRequest](#milkyway-rewards-v1-QueryOperatorOutstandingRewardsRequest) | [QueryOperatorOutstandingRewardsResponse](#milkyway-rewards-v1-QueryOperatorOutstandingRewardsResponse) | OperatorOutstandingRewards queries rewards of an operator. |
| ServiceOutstandingRewards | [QueryServiceOutstandingRewardsRequest](#milkyway-rewards-v1-QueryServiceOutstandingRewardsRequest) | [QueryServiceOutstandingRewardsResponse](#milkyway-rewards-v1-QueryServiceOutstandingRewardsResponse) | ServiceOutstandingRewards queries rewards of a service. |
| OperatorCommission | [QueryOperatorCommissionRequest](#milkyway-rewards-v1-QueryOperatorCommissionRequest) | [QueryOperatorCommissionResponse](#milkyway-rewards-v1-QueryOperatorCommissionResponse) | OperatorCommission queries accumulated commission for an operator. |
| PoolDelegationRewards | [QueryPoolDelegationRewardsRequest](#milkyway-rewards-v1-QueryPoolDelegationRewardsRequest) | [QueryPoolDelegationRewardsResponse](#milkyway-rewards-v1-QueryPoolDelegationRewardsResponse) | PoolDelegationRewards queries the total rewards accrued by a pool delegation. |
| OperatorDelegationRewards | [QueryOperatorDelegationRewardsRequest](#milkyway-rewards-v1-QueryOperatorDelegationRewardsRequest) | [QueryOperatorDelegationRewardsResponse](#milkyway-rewards-v1-QueryOperatorDelegationRewardsResponse) | OperatorDelegationRewards queries the total rewards accrued by a operator delegation. |
| ServiceDelegationRewards | [QueryServiceDelegationRewardsRequest](#milkyway-rewards-v1-QueryServiceDelegationRewardsRequest) | [QueryServiceDelegationRewardsResponse](#milkyway-rewards-v1-QueryServiceDelegationRewardsResponse) | ServiceDelegationRewards queries the total rewards accrued by a service delegation. |
| DelegatorTotalRewards | [QueryDelegatorTotalRewardsRequest](#milkyway-rewards-v1-QueryDelegatorTotalRewardsRequest) | [QueryDelegatorTotalRewardsResponse](#milkyway-rewards-v1-QueryDelegatorTotalRewardsResponse) | DelegatorTotalRewards queries the total rewards accrued by a single delegator |
| DelegatorWithdrawAddress | [QueryDelegatorWithdrawAddressRequest](#milkyway-rewards-v1-QueryDelegatorWithdrawAddressRequest) | [QueryDelegatorWithdrawAddressResponse](#milkyway-rewards-v1-QueryDelegatorWithdrawAddressResponse) | DelegatorWithdrawAddress queries withdraw address of a delegator. |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

