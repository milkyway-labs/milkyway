# `x/restaking`

## Abstract

The following document specifies the restaking module.

This module contains the logic required for users to re-stake their assets.
These assets can be restaked within a pool, towards an operator, or toward an AVS.
In addition to managing delegation, this module allows operators to start and stop securing an AVS, 
as well as enabling AVS admins to determine which operator is authorized to secure their service and 
from which pool they may borrow security.

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

## State

### Unbonding ID

UnbondingID stores the ID of the latest unbonding operation. 
It enables creating unique IDs for unbonding operations, i.e., UnbondingID 
is incremented every time a new unbonding operation (pool undelegation, service undelegation and operator undelegation) is initiated.

* Unbonding ID: `0x01 -> uint64`

### Unbonding index

UnbondingIndex stores indexes that are used to obtain the key of a `UnbondingDelegation` given 
an `UnbondingID`.

* UnbondingIndex: `0x02 | UnbondingID -> UnbondingDelegationKey`  

### Unbonding type

UnbondingType stores the type of a `UnbondingDelegation`.
This allow to know to which target the `UnbondingDelegation` refers to.

* UnbondingType: `0x03 | UnbondingID -> TargetID`

 ### Operator joined services

OperatorJoinedServices stores which services an operator has joined.

* OperatorJoinedServices: `collections.NewIndexedMap(0x13)`

### Service operators allow list

ServiceOperatorsAllowList stores which operators are allowed to join a service.
It enables the creation of a whitelist of operators that can join a service.

* ServiceOperatorsAllowList: `collections.KeySet(0x14)`

### Service securing pools

ServiceSecuringPools stores from which pools a service is borrowing security.

* ServiceSecuringPools: `collections.KeySet(0x15)`

### Service joined by operator

ServiceJoinedByOperator stores the indexs to perform a reverse lookup of 
[Operator joined services](#operator-joined-services).
It allows to obtain the list of operators that have joined a specific service.

* ServiceJoinedByOperator: `indexes.NewReversePair(0x16)`

### Pool delegation

PoolDelegation stores the delegation made toward a pool by a user.
It allows to obtain all the user's delegation toward a pool and to quick check 
if a user has delegated toward a specific pool.

* PoolDelegation: `0xa1 | UserAddr | PoolID -> ProtocolBuffer(Delegation)`

### Pool delegations by pool ID

PoolDelegation stores the delegations made toward a specific pool.
It allows to perform a reverse lookup of [Pool delegation](#pool-delegation).

* PoolDelegationsByPoolID: `0xa2 | PoolID | UserAddr -> []byte{}`

### Pool unbonding delegation

PoolUnbondingDelegation stores the pools unbonding delegation made by a user.
It allows to obtain an user's pool unbonding delegation.

* PoolUnbondingDelegation: `0xa3 | UserAddr | PoolID -> ProtocolBuffer(UnbondingDelegation)`

### Operator delegation

OperatorDelegation stores the delegation made toward an operator by a user.
It allows to obtain all the user's delegation toward an operator and to quick check 
if a user has delegated toward a specific operator.

* PoolDelegation: `0xb1 | UserAddr | OperatorID -> ProtocolBuffer(Delegation)`

### Operator delegations by operator ID

OperatorDelegationsByOperatorID stores the delegations made toward a specific operator.
It allows to perform a reverse lookup of [Operator delegation](#operator-delegation).

* OperatorDelegationsByOperatorID: `0xb2 | OperatorID | UserAddr -> []byte{}`

### Operator unbonding delegation

OperatorUnbondingDelegation stores the users' unbonding delegation that are related to 
an operator.
It allows to obtain an user's operator unbonding delegation.

* OperatorUnbondingDelegation: `0xb3 | UserAddr | OperatorID -> ProtocolBuffer(UnbondingDelegation)`

### Service delegation

ServiceDelegation stores the delegation made toward a service by a user.
It allows to obtain all the user's delegation toward a service and to quick check 
if a user has delegated toward a specific service.

* ServiceDelegation: `0xc1 | UserAddr | ServiceID -> ProtocolBuffer(Delegation)`

### Service delegations by service ID

ServiceDelegationsByServiceID stores the delegations made toward a specific service.
It allows to perform a reverse lookup of [Service delegation](#delegations-delegation).

* ServiceDelegationsByServiceID: `0xc2 | ServiceID | UserAddr -> []byte{}`

### Service unbonding delegation

ServiceUnbondingDelegation stores the users' unbonding delegation that are related to 
a service.
It allows to obtain an user's service unbonding delegation.

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

* determine the delegators shares based on tokens delegated and the target's exchange rate
* creates a new `Delegation` or update an existing one with the computed shares
* move the tokens from the delegator account to the target's account
* store the `Delegation` object in the target's delegation store (`PoolDelegation`, `OperatorDelegation` or `ServiceDelegation`)
* store the index to perform the inverse look of all the delegations given the target id (`PoolDelegationByUser`, `OperatorDelegationByUser` or `ServiceDelegationByUser`)

### Begin unbonding

When a user perform an `UndelegatePool`, `UndelegateOperator` or `UndelegateService` 
the following operations occur:

* determines the shares based the amount of tokens that the user want to undelegate
* subtract the computed shares from the `Delegation` object. In case the final shares are 
0 the `Delegation` object will be removed
* subtract the undelegated tokens from the `Target` total delegated tokens
* computes the time when the undelegated tokens will be returned to the users
* increments the `UnbondingID`
* creates a new `UnbondingDelegation` or obtain an existing `UnbondingDelegation` in case 
the user was already undelegating some tokens from that target
* adds a new `UnbondingDelegationEntry` to the `UnbondingDelegation`
* stores the `UnbondingDelegation`, this is stored in `UserPoolUnbondingDelegation` for 
pools, `UserOperatorUnbondingDelegation` for operators and `UserServiceUnbondingDelegation` for services
* stores an index to retrive the `UnbondingDelegation` given an `UnbondingID`
* stores in the `UnbondingQueue` when the `UnbondingDelegation` completes

### Complete unbonding

When an `UnbondingDelegationEntry` matures the following operations occur:

* removes the index that associates the `UnbondigID` with the `UnbondingDelegation` object
* transfer the unbonded tokens from the target to the user
* removes the `UnbondingDelegationEntry` from the `UnbondingDelegation` object
* store the updated `UnbondingDelegation` or deletes in case the removed `UnbondingDelegationEntry` was the last one

## Messages

In this section we describe the processing of the restaking messages.

### MsgJoinService

It allows the operator's admin to start securing a AVS.

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

It allows the service admin to add an operator to the list of allowed operator to secure the service.
If the service didn't have an allow list after adding the operator to the allow list all the operators
not in the allow list will be stopped from securing the service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L115-L127
```

This message is expected to fail if:

* the `sender` is not the `Service` admin
* the `Operator` is already in the allow list

### MsgRemoveOperatorFromAllowlist

It allows the service admin to remove an operator to the list of allowed operator to secure the service.
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

It allows the service admin to remove a pool from the list of pools from which the service has chosen to borrow security.

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

* the user wants to undelegates an amount greater then amount delegated

### MsgDelegateOperator

It allows a user to delegate their assets to an operator.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L210-L229
```

This message is expected to fail if:

* don't exist an `Operator` with the given id
* the `Operator` is not active
* the denom of the coin that the user wants to delegate is not allowed
* the `sender` don't have the amount of coin that wants to delegate

### MsgUndelegateOperator

It allows a user to undelegate their assets from a restaking operator.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L298-L320
```

This message is expected to fail if:

* don't exist an `Operator` with the given id
* the user wants to undelegates an amount greater then amount delegated toward that operator

### MsgDelegateService

It allows a user to delegate their assets to a service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L234-L252
```

This message is expected to fail if:

* don't exist a `Service` with the given id
* the denom of the coin that the user wants to delegate is not allowed
* the `sender` don't have the amount of coin that wants to delegate

### MsgUndelegateService

It allows a user to undelegate their assets from a restaking service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L322-L344
```

This message is expected to fail if:

* don't exist a `Service` with the given id
* the user wants to undelegates an amount greater then amount delegated toward that operator

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

Allows to update the module parameters.  The params are updated through a 
governance proposal where the signer is the gov module account address.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/messages.proto#L257-L274
```

Params:

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/restaking/v1/params.proto#L9-L33
```

