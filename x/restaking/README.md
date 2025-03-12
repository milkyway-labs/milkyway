# `x/restaking`

## Abstract

The following document specifies the restaking module.

This module contains the logic that allows users to restake their assets.

Assets can be restaked within a pool, towards an operator, or towards a specific Actively Validated Service (AVS).

In addition to managing the delegation, this module allows operators to start and stop securing an AVS, as well as AVS
admins to determine which operator is authorized to secure their service and from which pool they may borrow security.

## Contents

* [State](#state)
    * [Unbonding ID](#unbonding-id)
    * [Unbonding index](#unbonding-index)
    * [Unbonding type](#unbonding-type)
    * [Operator joined services](#operator-joined-services)
    * [Service operators allow list](#service-operators-allow-list)
    * [Service securing pools](#service-securing-pools)
    * [Service joined by operator](#service-joined-by-operator)
    * [Pool delegation](#pool-delegation)
    * [Pool delegation by pool ID](#pool-delegation-by-pool-id)
    * [Pool unbonding delegation](#pool-unbonding-delegation)
    * [Operator delegation](#operator-delegation)
    * [Operator delegations by operator ID](#operator-delegations-by-operator-id)
    * [Operator unbonding delegation](#operator-unbonding-delegation)
    * [Service delegation](#service-delegation)
    * [Service delegations by service ID](#service-delegations-by-service-id)
    * [Service unbonding delegation](#service-unbonding-delegation)
    * [Unbonding queue](#unbonding-queue)
    * [User preferences](#user-preferences)
* [State transitions](#state-transitions)
    * [Delegations](#delegations)
    * [Begin unbonding](#begin-unbonding)
    * [Complete unbonding](#complete-unbonding)
* [Messages](#messages)
    * [MsgJoinService](#msgjoinservice)
    * [MsgLeaveService](#msgleaveservice)
    * [MsgAddOperatorToAllowList](#msgaddoperatortoallowlist)
    * [MsgRemoveOperatorFromAllowlist](#msgremoveoperatorfromallowlist)
    * [MsgBorrowPoolSecurity](#msgborrowpoolsecurity)
    * [MsgCeasePoolSecurityBorrow](#msgceasepoolsecurityborrow)
    * [MsgDelegatePool](#msgdelegatepool)
    * [MsgUndelegatePool](#msgundelegatepool)
    * [MsgDelegateOperator](#msgdelegateoperator)
    * [MsgUndelegateOperator](#msgundelegateoperator)
    * [MsgDelegateService](#msgdelegateservice)
    * [MsgUndelegateService](#msgundelegateservice)
    * [MsgSetUserPreferences](#msgsetuserpreferences)
    * [MsgUpdateParams](#msgupdateparams)
* [End-Block](#end-block)
    * [Unbonding Delegations](#unbonding-delegations)
* [Hooks](#hooks)
* [Events](#events)
* [Parameters](#parameters)

## State

### Unbonding ID

UnbondingID stores the ID of the latest unbonding operation.
It enables creating unique IDs for unbonding operations, i.e., UnbondingID
is incremented every time a new unbonding operation (pool undelegation, service undelegation and operator undelegation)
is initiated.

* Unbonding ID: `0x01 -> uint64`

### Unbonding index

`UnbondingIndex` maintains indexes used to acquire the key of a `UnbondingDelegation` based on its unique identifier,
`UnbondingID`.

* UnbondingIndex: `0x02 | UnbondingID -> UnbondingDelegationKey`

### Unbonding type

UnbondingType represents the type of `UnbondingDelegation`,
indicating to which entity (pool, operator, or service) the unbonding delegation is related.  
This information can be used to determine the appropriate actions to take when a `UnbondingDelegation`
has completed and the funds should be returned to the user's account.

* UnbondingType: `0x03 | UnbondingID -> TargetID`

### Operator joined services

OperatorJoinedServices maintains a record of the services that an operator has joined.

* OperatorJoinedServices: `collections.NewIndexedMap(0x13)`

### Service operators allow list

ServiceOperatorsAllowList stores which operators are allowed to join a service.
It allows obtaining the list of operators that have joined a specific service.

* ServiceOperatorsAllowList: `collections.KeySet(0x14)`

### Service securing pools

ServiceSecuringPools stores from which pools a service is borrowing security.

* ServiceSecuringPools: `collections.KeySet(0x15)`

### Service joined by operator

ServiceJoinedByOperator stores the indexes to perform a reverse lookup of
[Operator joined services](#operator-joined-services).
It allows obtaining the list of operators that have joined a specific service.

* ServiceJoinedByOperator: `indexes.NewReversePair(0x16)`

### Pool delegation

PoolDelegation stores the delegation made by a user toward a pool.
It allows obtaining all of a user’s delegations for a pool and quickly checking if a delegation exists.

* PoolDelegation: `0xa1 | UserAddr | PoolID -> ProtocolBuffer(Delegation)`

### Pool delegations by pool ID

PoolDelegation stores the delegations made toward a specific pool.
It supports a reverse lookup for [Pool delegation](#pool-delegation).

* PoolDelegationsByPoolID: `0xa2 | PoolID | UserAddr -> []byte{}`

### Pool unbonding delegation

PoolUnbondingDelegation stores the pools unbonding delegation made by a user.
It allows obtaining a user's pool unbonding delegation.

* PoolUnbondingDelegation: `0xa3 | UserAddr | PoolID -> ProtocolBuffer(UnbondingDelegation)`

### Operator delegation

OperatorDelegation stores the delegation made by a user toward an operator.
It allows obtaining all of a user’s delegations to an operator.

* PoolDelegation: `0xb1 | UserAddr | OperatorID -> ProtocolBuffer(Delegation)`

### Operator delegations by operator ID

OperatorDelegationsByOperatorID stores the delegations made toward a specific operator.
It supports a reverse lookup for [Operator delegation](#operator-delegation).

* OperatorDelegationsByOperatorID: `0xb2 | OperatorID | UserAddr -> []byte{}`

### Operator unbonding delegation

OperatorUnbondingDelegation stores the users' unbonding delegation that are related to
an operator.
It allows to obtain a user's operator unbonding delegation.

* OperatorUnbondingDelegation: `0xb3 | UserAddr | OperatorID -> ProtocolBuffer(UnbondingDelegation)`

### Service delegation

ServiceDelegation stores the delegation made by a user toward a service.
It allows obtaining all of a user’s delegations to a service.

* ServiceDelegation: `0xc1 | UserAddr | ServiceID -> ProtocolBuffer(Delegation)`

### Service delegations by service ID

ServiceDelegationsByServiceID stores the delegations made toward a specific service.
It supports a reverse lookup for [Service delegation](#delegations-delegation).

* ServiceDelegationsByServiceID: `0xc2 | ServiceID | UserAddr -> []byte{}`

### Service unbonding delegation

ServiceUnbondingDelegation stores the users' unbonding delegation that are related to
a service.
It allows obtaining a user's service unbonding delegation.

* ServiceUnbondingDelegation: `0xc3 | UserAddr | ServiceID -> ProtocolBuffer(UnbondingDelegation)`

### Unbonding queue

UnbondingQueue stores the timestamp at which a list of `UnbodingDelegation` completes
and the delegated funds should return to the user's account.

* UnbondingQueue: `0xd1 | Timestamp -> ProtocolBuffer(DTDataList)`

### User preferences

UserPreferences stores the users' preferences.

* UserPreferences: `collections.NewMap[string, UserPreferences](0xe1)`

## State transitions

### Delegations

When a delegation occurs the target object is affected together with the
internal state of the module.
The target object can be a `Pool`, `Operator` or `Service`.

* determine the delegators shares based on the tokens delegated and the total amount of tokens delegated toward the
  target
* creates a new `Delegation` or update an existing one with the computed shares
* move the tokens from the delegator account to the target's account
* store the `Delegation` object in the target's delegation store (`PoolDelegation`, `OperatorDelegation` or
  `ServiceDelegation`)
* store the index to perform the inverse look of all the delegations given the target id (`PoolDelegationByUser`,
  `OperatorDelegationByUser` or `ServiceDelegationByUser`)

### Begin unbonding

When a user perform a `UndelegatePool`, `UndelegateOperator` or `UndelegateService` 
the following operations occur:

* determine the shares based on the amount of tokens that the user want to undelegate
* subtract the computed shares from the `Delegation` object. In case the final shares are 
0 the `Delegation` object will be removed
* subtract the undelegated tokens from the `Target` total delegated tokens
* computes the time when the undelegated tokens will be returned to the user
* increments the `UnbondingID`
* creates a new `UnbondingDelegation` or obtain an existing `UnbondingDelegation` in case
  the user was already undelegating some tokens from the target
* adds a new `UnbondingDelegationEntry` to the `UnbondingDelegation`
* stores the `UnbondingDelegation`, this is stored in `UserPoolUnbondingDelegation` for
  pools, `UserOperatorUnbondingDelegation` for operators and `UserServiceUnbondingDelegation` for services
* stores an index to retrive the `UnbondingDelegation` given an `UnbondingID`
* stores in the `UnbondingQueue` when the `UnbondingDelegation` completes

### Complete unbonding

When an `UnbondingDelegationEntry` matures the following operations occur:

* remove the index that associates the `UnbondingID` with the `UnbondingDelegation` object
* transfer the unbonded tokens from the target to the user
* removes the `UnbondingDelegationEntry` from the `UnbondingDelegation` object
* store the updated `UnbondingDelegation` or deletes in case the removed `UnbondingDelegationEntry` was the last one

## Messages

In this section we describe the processing of the restaking messages.

### MsgJoinService

It allows the operator's admin to start securing an AVS.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L81-L93
```

This message is expected to fail if:

* the `sender` is not the `Operator` admin
* the service that the admin wants to join is no logger active
* the operator is not in the service's operators allow list

### MsgLeaveService

It allows the operator's admin to stop securing an AVS.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L98-L110
```

This message is expected to fail if:

* the `sender` is not the `Operator` admin
* the `Operator` didn't join the service

### MsgAddOperatorToAllowList

It allows the service admin to add an operator to the allowed list for securing the service.
If the service didn't have an allow list after adding the operator to the allow list all the operators
not in the allow list will be stopped from securing the service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L115-L127
```

This message is expected to fail if:

* the `sender` is not the `Service` admin
* the `Operator` is already in the allow list

### MsgRemoveOperatorFromAllowlist

It allows the service admin to remove an operator from the allowed list.
If the operator was securing the service will be stopped from securing it.
When the last operator is removed from the allow list the service will allow all operators to join.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L133-L146
```

This message is expected to fail if:

* the `sender` is not the `Service` admin
* the `Operator` is not in the `Service`'s allow list

### MsgBorrowPoolSecurity

It allows the service admin to add a pool to the list of pools from which the service has chosen to borrow security.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L152-L165
```

This message is expected to fail if:

* the `sender` is not the `Service` admin
* the `Service` is inactive
* the `Service` is already secured by the provided `Pool`

### MsgCeasePoolSecurityBorrow

It allows the service admin to remove a pool from the list of pools from which the service has chosen to borrow
security.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L170-L183
```

This message is expected to fail if:

* the `sender` is not the `Service` admin
* the `Service` is not secured by the provided `Pool`

### MsgDelegatePool

It allows a user to put their assets into a restaking pool that will later be
used to provide cryptoeconomic security to services that choose it.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L189-L205
```

This message is expected to fail if:

* the denom of the coin that the user wants to delegate is not allowed
* the `sender` don't have the amount of coin that wants to delegate

### MsgUndelegatePool

It allows a user to undelegate their assets from a restaking pool.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L279-L296
```

This message is expected to fail if:

* the user wants to undelegate an amount greater than the delegated amount

### MsgDelegateOperator

It allows a user to delegate their assets to an operator.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L210-L229
```

This message is expected to fail if:

* no `Operator` with the given id exists
* the `Operator` is not active
* the denom of the coin that the user wants to delegate is not allowed
* the `sender` don't have the amount of coin that wants to delegate

### MsgUndelegateOperator

It allows a user to undelegate their assets from a restaking operator.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L298-L320
```

This message is expected to fail if:

* no `Operator` with the given id exists
* the user wants to undelegate an amount greater than the delegated amount

### MsgDelegateService

It allows a user to delegate their assets to a service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L234-L252
```

This message is expected to fail if:

* no `Service` with the given id exists
* the denom of the coin that the user wants to delegate is not allowed
* the `sender` don't have the amount of coin that wants to delegate

### MsgUndelegateService

It allows a user to undelegate their assets from a restaking service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L322-L344
```

This message is expected to fail if:

* no `Service` with the given id exists
* the user wants to undelegate an amount greater than the delegated amount

### MsgSetUserPreferences

It allows a user to set their preferences for the restaking module.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L357-L369
```

User preferences:

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/models.proto#L150-L170
```

This message is expected to fail if:

* one or more service specified in the preferences don't exists
* one or more pool specified in the preferences don't exists

### MsgUpdateParams

Allows to update the module parameters. Parameters are updated via a 
governance proposal, with the gov module account as the signer.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L257-L274
```

Params:

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/params.proto#L9-L33
```

## End-Block

In this section we describe the operation that this modules executes during each abci end block call.

### Unbonding delegations

Complete the unbonding of all mature `UnbondingDelegations.Entries` within the `UnbondingQueue` with the following
procedure:

* removes the index that associates the `UnbondigID` with the `UnbondingDelegation` object
* transfer the unbonded tokens from the target to the user
* removes the `UnbondingDelegationEntry` from the `UnbondingDelegation` object
* store the updated `UnbondingDelegation` or deletes in case the removed `UnbondingDelegationEntry` was the last one

## Hooks

Other modules may register operations to execute when a certain event has occurred within the restaking module.
The following hooks can registered with restaking:

* `BeforePoolDelegationCreated(ctx context.Context, poolID uint32, delegator string) error`
   * called before a new `Delegation` object associated to a `Pool` is created
* `BeforePoolDelegationSharesModified(ctx context.Context, poolID uint32, delegator string) error`
   * called before the shares of `Delegation` object associated a `Pool` are modified
* `AfterPoolDelegationModified(ctx context.Context, poolID uint32, delegator string) error`
   * called after a `Delegation` object associated to a `Pool` is modified
* `BeforePoolDelegationRemoved(ctx context.Context, poolID uint32, delegator string) error`
   * called after a `Delegation` object associated to a `Pool` is removed
* `BeforeOperatorDelegationCreated(ctx context.Context, operatorID uint32, delegator string) error`
   * called before a new `Delegation` object associated to a `Operator` is created
* `BeforeOperatorDelegationSharesModified(ctx context.Context, operatorID uint32, delegator string) error`
   * called before the shares of `Delegation` object associated a `Operator` are modified
* `AfterOperatorDelegationModified(ctx context.Context, operatorID uint32, delegator string) error`
   * called after a `Delegation` object associated to a `Operator` is modified
* `BeforeOperatorDelegationRemoved(ctx context.Context, operatorID uint32, delegator string) error`
   * called after a `Delegation` object associated to a `Operator` is removed
* `BeforeServiceDelegationCreated(ctx context.Context, serviceID uint32, delegator string) error`
   * called before a new `Delegation` object associated to a `Service` is created
* `BeforeServiceDelegationSharesModified(ctx context.Context, serviceID uint32, delegator string) error`
   * called before the shares of `Delegation` object associated a `Service` are modified
* `AfterServiceDelegationModified(ctx context.Context, serviceID uint32, delegator string) error`
   * called after a `Delegation` object associated to a `Service` is modified
* `BeforeServiceDelegationRemoved(ctx context.Context, serviceID uint32, delegator string) error`
   * called after a `Delegation` object associated to a `Service` is removed
* `AfterUnbondingInitiated(ctx context.Context, unbondingDelegationID uint64) error`
   * called after a new `UnbondingDelegation` is created
*
`AfterUserPreferencesModified(ctx context.Context, userAddress string, oldPreferences, newPreferences UserPreferences) error`
   * called after an user's preferences are modified

## Events

The restaking module emits the following events:

### EndBlocker

| Type                  | Attribute Key        | Attribute Value                 |
|-----------------------|----------------------|---------------------------------|
| complete_unbonding    | amount               | {totalUnbondingAmount}          |
| complete_unbonding    | unbonding_delegation | {unbondingDelegationTargetType} |
| complete_redelegation | target_id            | {delegationTargetId}            |
| complete_unbonding    | delegator            | {delegatorAddress}              |

### MsgJoinService

| Type         | Attribute Key | Attribute Value |
|--------------|---------------|-----------------|
| join_service | operator_id   | {operatorId}    |
| join_service | service_id    | {serviceId}     |

### MsgLeaveService

| Type          | Attribute Key | Attribute Value |
|---------------|---------------|-----------------|
| leave_service | operator_id   | {operatorId}    |
| leave_service | service_id    | {serviceId}     |

### MsgAddOperatorToAllowList

| Type           | Attribute Key | Attribute Value |
|----------------|---------------|-----------------|
| allow_operator | operator_id   | {operatorId}    |
| allow_operator | service_id    | {serviceId}     |

### MsgRemoveOperatorFromAllowlist

| Type                    | Attribute Key | Attribute Value |
|-------------------------|---------------|-----------------|
| remove_allowed_operator | operator_id   | {operatorId}    |
| remove_allowed_operator | service_id    | {serviceId}     |

### MsgBorrowPoolSecurity

| Type                 | Attribute Key | Attribute Value |
|----------------------|---------------|-----------------|
| borrow_pool_security | service_id    | {serviceId}     |
| borrow_pool_security | pool_id       | {poolId}        |

### MsgCeasePoolSecurityBorrow

| Type                       | Attribute Key | Attribute Value |
|----------------------------|---------------|-----------------|
| cease_pool_security_borrow | service_id    | {serviceId}     |
| cease_pool_security_borrow | pool_id       | {poolId}        |

### MsgDelegatePool

| Type          | Attribute Key | Attribute Value    |
|---------------|---------------|--------------------|
| delegate_pool | delegator     | {delegatorAddress} |
| delegate_pool | amount        | {delegationAmount} |
| delegate_pool | new_shares    | {newShares}        |

### MsgUndelegatePool

| Type        | Attribute Key   | Attribute Value    |
|-------------|-----------------|--------------------|
| unbond_pool | amount          | {unbondAmount}     |
| unbond_pool | delegator       | {delegatorAddress} |
| unbond_pool | completion_time | {completionTime}   |

### MsgDelegateOperator

| Type              | Attribute Key | Attribute Value    |
|-------------------|---------------|--------------------|
| delegate_operator | delegator     | {delegatorAddress} |
| delegate_operator | operator_id   | {operatorID}       |
| delegate_operator | amount        | {delegationAmount} |
| delegate_operator | new_shares    | {newShares}        |

### MsgUndelegateOperator

| Type            | Attribute Key   | Attribute Value    |
|-----------------|-----------------|--------------------|
| unbond_operator | amount          | {unbondAmount}     |
| unbond_operator | delegator       | {delegatorAddress} |
| unbond_operator | operator_id     | {operatorID}       |
| unbond_operator | completion_time | {completionTime}   |

### MsgDelegateService

| Type             | Attribute Key | Attribute Value    |
|------------------|---------------|--------------------|
| delegate_service | delegator     | {delegatorAddress} |
| delegate_service | service_id    | {serviceID}        |
| delegate_service | amount        | {delegationAmount} |
| delegate_service | new_shares    | {newShares}        |

### MsgUndelegateService

| Type           | Attribute Key   | Attribute Value    |
|----------------|-----------------|--------------------|
| unbond_service | amount          | {unbondAmount}     |
| unbond_service | delegator       | {delegatorAddress} |
| unbond_service | service_id      | {serviceID}        |
| unbond_service | completion_time | {completionTime}   |

### MsgSetUserPreferences

| Type                 | Attribute Key | Attribute Value |
|----------------------|---------------|-----------------|
| set_user_preferences | user          | {userAddress}   |

### Parameters

The restaking module contains the following parameters:

| Key           | Type              | Example           |
|---------------|-------------------|-------------------|
| UnbondingTime | string (time ns)  | "259200000000000" |
| AllowedDenoms | []string          | {"utia", "uusdc"} |
| RestakingCap  | sdkmath.LegacyDec | "259200000000000" |
| MaxEntries    | uint32            | 7                 |

