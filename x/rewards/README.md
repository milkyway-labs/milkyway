# `x/rewards`

## Abstract

`x/rewards` is a Cosmos SDK module that manages restakig rewards.

## Contents

* [Concepts](#concepts)
    * [F1 Distribution](#f1-distribution)
    * [Rewards Plans](#rewards-plans)
        * [Basic Distribution](#basic-distribution)
        * [Weighted Distribution](#weighted-distribution)
        * [Egalitarian Distribution](#egalitarian-distribution)
* [State](#state)
    * [Params](#params)
    * [NextRewardsPlanID](#nextrewardsplanid)
    * [RewardsPlans](#rewardsplans)
    * [LastRewardsAllocationTime](#lastrewardsallocationtime)
    * [DelegatorWithdrawAddrs](#delegatorwithdrawaddrs)
    * [DelegatorStartingInfos](#delegatorstartinginfos)
    * [HistoricalRewards](#historicalrewards)
    * [CurrentRewards](#currentrewards)
    * [OutstandingRewards](#outstandingrewards)
    * [PoolServiceTotalDelegatorShares](#poolservicetotaldelegatorshares)
    * [OperatorAccumulatedCommissions](#operatoraccumulatedcommissions)
* [Messages](#messages)
    * [MsgCreateRewardsPlan](#msgcreaterewardsplan)
    * [MsgEditRewardsPlan](#msgeditrewardsplan)
    * [MsgSetWithdrawAddress](#msgsetwithdrawaddress)
    * [MsgWithdrawDelegatorReward](#msgwithdrawdelegatorreward)
    * [MsgWithdrawOperatorCommission](#msgwithdrawoperatorcommission)
* [Events](#events)
    * [BeginBlocker](#beginblocker)
    * [Handlers](#handlers)
        * [MsgCreateRewardsPlan](#msgcreaterewardsplan-1)
        * [MsgEditRewardsPlan](#msgeditrewardsplan-1)
        * [MsgSetWithdrawAddress](#msgsetwithdrawaddress-1)
        * [MsgWithdrawDelegatorReward](#msgwithdrawdelegatorreward-1)
        * [MsgWithdrawOperatorCommission](#msgwithdrawoperatorcommission-1)
* [Parameters](#parameters)

## Concepts

### F1 Distribution

x/rewards uses the F1 distribution algorithm which the Cosmos SDK's
x/distribution module also uses.
For more information on the F1 distribution algorithm, refer to the
[F1 Fee Distribution paper](https://github.com/cosmos/cosmos-sdk/blob/main/docs/spec/fee_distribution/f1_fee_distr.pdf)
and the [x/distribution module](https://github.com/cosmos/cosmos-sdk/blob/v0.50.11/x/distribution/README.md).

### Rewards Plans

Rewards plans are created by service admins to reward restakers who have
restaked their assets to the service.
A rewards plan consists of the following parameters:

* Description: A description of the rewards plan
* Amount per day: The amount of rewards to be distributed per day
    * Rewards are distributed per block, based on the previous block's duration
* Start time: The time when the rewards plan starts
* End time: The time when the rewards plan ends
* Pools distribution: The distribution method of rewards toward pools
* Operators distribution: The distribution method of rewards toward operators
* Users distribution: The distribution method of rewards toward delegators who
  have delegated to the service directly

Each distribution method object also has `weight` field which is used to
determine the distribution ratio among the pools, operators and users.

There are three types of distribution methods, while users distribution can only
be basic:

* Basic distribution
* Weighted distribution
* Egalitarian distribution

#### Basic Distribution

Rewards are distributed to all pools(or operators) proportionally based on the
USD value of restaked assets.

Example:

* Pool A value: \$30M
* Pool B value: \$50M
* Total value: \$80M
* Pool A's rewards: \$30M / \$80M * This block's rewards
* Pool B's rewards: \$50M / \$80M * This block's rewards

#### Weighted Distribution

Rewards are distributed to the predefined list of pools(or operators) with the
specified weights.

Example:

* Operator A weight: 4
* Operator B weight: 7
* Operator C is not in the list
* Total weight: 11
* Operator A's rewards: 4 / 11 * This block's rewards
* Operator B's rewards: 7 / 11 * This block's rewards
* Operator C's rewards: 0

#### Egalitarian Distribution

Rewards are distributed equally among all pools(or operators).

Example:

* Pool A's rewards: This block's rewards / Total number of pools
* Pool B's rewards: This block's rewards / Total number of pools
* ...

## State

### Params

The module parameters are stored under the `0x01` key:

* Params: `0x01 -> ProtocolBuffer(Params)`

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/params.proto#L10-L20
```

### NextRewardsPlanID

NextRewardsPlanID stores the ID of the next rewards plan to be created.

* NextRewardsPlanID: `0xa1 -> uint64`

### RewardsPlans

All the rewards plan are stored under the `0xa2` key prefix.

* RewardsPlans: `0xa2 | RewardsPlanID -> ProtocolBuffer(RewardsPlan)`

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/models.proto#L15-L59
```

### LastRewardsAllocationTime

LastRewardsAllocationTime stores the timestamp of the last rewards allocation.

* LastRewardsAllocationTime: `0xa3 -> Timestamp`

### DelegatorWithdrawAddrs

DelegatorWithdrawAddrs stores the withdraw address of each delegator.

* DelegatorWithdrawAddrs: `0xa4 | DelegatorAddr -> WithdrawAddr`

### DelegatorStartingInfos

DelegatorStartingInfos stores the delegation starting information of each
delegator.
There are three DelegatorStartingInfos for each delegation type:

* PoolDelegatorStartingInfos: `0xb1 | PoolID | DelegatorAddr -> ProtocolBuffer(DelegatorStartingInfo)`
* OperatorDelegatorStartingInfos: `0xc2 | OperatorID | DelegatorAddr -> ProtocolBuffer(DelegatorStartingInfo)`
* ServiceDelegatorStartingInfos: `0xd1 | ServiceID | DelegatorAddr -> ProtocolBuffer(DelegatorStartingInfo)`

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/models.proto#L187-L204
```

### HistoricalRewards

HistoricalRewards stores the historical rewards for each delegation target.
It keeps track of cumulative rewards ratios in order to calculate rewards lazily
using the F1 distribution logic.
There are three HistoricalRewards for each delegation type:

* PoolHistoricalRewards: `0xb2 | PoolID | Period -> ProtocolBuffer(HistoricalRewards)`
* OperatorHistoricalRewards: `0xc3 | OperatorID | Period -> ProtocolBuffer(HistoricalRewards)`
* ServiceHistoricalRewards: `0xd2 | ServiceID | Period -> ProtocolBuffer(HistoricalRewards)`

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/models.proto#L133-L152
```

### CurrentRewards

CurrentRewards stores the current rewards and the current period of each
delegation target.
There are three CurrentRewards for each delegation type:

* PoolCurrentRewards: `0xb3 | PoolID -> ProtocolBuffer(CurrentRewards)`
* OperatorCurrentRewards: `0xc4 | OperatorID -> ProtocolBuffer(CurrentRewards)`
* ServiceCurrentRewards: `0xd3 | ServiceID -> ProtocolBuffer(CurrentRewards)`

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/models.proto#L154-L164
```

### OutstandingRewards

OutstandingRewards stores the outstanding(un-withdrawn) rewards for each
delegation target.
There are three OutstandingRewards for each delegation type:

* PoolOutstandingRewards: `0xb4 | PoolID -> ProtocolBuffer(OutstandingRewards)`
* OperatorOutstandingRewards: `0xc5 | OperatorID -> ProtocolBuffer(OutstandingRewards)`
* ServiceOutstandingRewards: `0xd4 | ServiceID -> ProtocolBuffer(OutstandingRewards)`

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/models.proto#L166-L174
```

### PoolServiceTotalDelegatorShares

PoolServiceTotalDelegatorShares stores the total trusted delegator shares of
each pool-service pair.
The total trusted delegator shares means the sum of delegation shares of all
delegators who has delegated to the pool and trusts the service with it.

* PoolServiceTotalDelegatorShares: `0xb5 | PoolID | ServiceID -> ProtocolBuffer(PoolServiceTotalDelegatorShares)`

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/models.proto#L225-L234
```

### OperatorAccumulatedCommissions

The total accumulated commissions of each operator are stored under the `0xc1`
key prefix:

* OperatorAccumulatedCommissions: `0xc1 | OperatorID -> ProtocolBuffer(AccumulatedCommission)`

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/models.proto#L176-L185
```

## Messages

### MsgCreateRewardsPlan

The `MsgCreateRewardsPlan` can be sent by service admins to create a new rewards
plan.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/messages.proto#L46-L98
```

The message will fail under the following conditions:

* The service doesn't exist
* The sender is not the service admin
* The supplied rewards plan creation fee is insufficient
* The plan is invalid

### MsgEditRewardsPlan

The `MsgEditRewardsPlan` can be sent by service admins to edit an existing
rewards plan.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/messages.proto#L107-L149
```

The message will fail under the following conditions:

* The service doesn't exist
* The sender is not the service admin
* The plan is invalid

### MsgSetWithdrawAddress

The `MsgSetWithdrawAddress` can be sent by anyone to set the withdraw address.
By default, the withdraw address is the delegator address.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/messages.proto#L155-L166
```

The message will fail under the following conditions:

* The withdraw address is blocked from receiving funds by the bank module

### MsgWithdrawDelegatorReward

The `MsgWithdrawDelegatorReward` can be sent by anyone to withdraw the rewards.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/messages.proto#L172-L186
```

The message will fail under the following conditions:

* The delegation target doesn't exist
* The delegation doesn't exist

### MsgWithdrawOperatorCommission

The `MsgWithdrawOperatorCommission` can be sent by operator admins to withdraw
the accumulated commission.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/messages.proto#L199-L210
```

The message will fail under the following conditions:

* The operator doesn't exist
* The sender is not the operator admin
* There's no accumulated commission to withdraw

## Events

### BeginBlocker

| Type         | Attribute Key          | Attribute Value                       |
|--------------|------------------------|---------------------------------------|
| `commission` | `operator_id`          | `{operatorID}`                        |
| `commission` | `pool`                 | `{denomOfRestakedAssetBeingRewarded}` |
| `commission` | `amount`               | `{commissionAmount}`                  |
| `rewards`    | `delegation_type`      | `{delegationType}`                    |
| `rewards`    | `delegation_target_id` | `{delegationTargetID}`                |
| `rewards`    | `pool`                 | `{restakedAssetBeingRewarded}`        |
| `rewards`    | `amount`               | `{rewardAmount}`                      |

### Handlers

#### MsgCreateRewardsPlan

| Type                  | Attribute Key     | Attribute Value   |
|-----------------------|-------------------|-------------------|
| `create_rewards_plan` | `rewards_plan_id` | `{rewardsPlanID}` |
| `create_rewards_plan` | `service_id`      | `{serviceID}`     |
| `create_rewards_plan` | `sender`          | `{senderAddress}` |

#### MsgEditRewardsPlan

| Type                | Attribute Key     | Attribute Value   |
|---------------------|-------------------|-------------------|
| `edit_rewards_plan` | `rewards_plan_id` | `{rewardsPlanID}` |
| `edit_rewards_plan` | `service_id`      | `{serviceID}`     |
| `edit_rewards_plan` | `sender`          | `{senderAddress}` |

#### MsgSetWithdrawAddress

| Type                   | Attribute Key      | Attribute Value     |
|------------------------|--------------------|---------------------|
| `set_withdraw_address` | `sender`           | `{senderAddress}`   |
| `set_withdraw_address` | `withdraw_address` | `{withdrawAddress}` |

#### MsgWithdrawDelegatorReward

| Type               | Attribute Key          | Attribute Value                           |
|--------------------|------------------------|-------------------------------------------|
| `withdraw_rewards` | `delegation_type`      | `{delegationType}`                        |
| `withdraw_rewards` | `delegation_target_id` | `{delegationTargetID}`                    |
| `withdraw_rewards` | `delegator`            | `{delegatorAddress}`                      |
| `withdraw_rewards` | `amount`               | `{withdrawnRewardAmount}`                 |
| `withdraw_rewards` | `amount_per_pool`      | `{withdrawnRewardAmountPerRestakedAsset}` |

#### MsgWithdrawOperatorCommission

| Type                  | Attribute Key     | Attribute Value                               |
|-----------------------|-------------------|-----------------------------------------------|
| `withdraw_commission` | `operator_id`     | `{operatorID}`                                |
| `withdraw_commission` | `amount`          | `{withdrawnCommissionAmount}`                 |
| `withdraw_commission` | `amount_per_pool` | `{withdrawnCommissionAmountPerRestakedAsset}` |

## Parameters

The rewards module contains the following parameters:

| Key                       | Type          | Example                                 |
|---------------------------|---------------|-----------------------------------------|
| rewards_plan_creation_fee | array (coins) | [{"denom":"umilk","amount":"10000000"}] |

Note that the `rewards_plan_creation_fee` represents the OR(not AND) conditions
of the fee which means that a plan creator can pay the fee with any coin
denomination and amount specified inside the `rewards_plan_creation_fee`.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v9.0.0/proto/milkyway/rewards/v1/params.proto#L10-L20
```