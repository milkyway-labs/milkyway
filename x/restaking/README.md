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

UnbondingQueueKey stores the timestamp at which a list of `UnbodingDelegation` completes 
and the delegated funds should return to the user's account.

* UnbondingQueueKey: `0xd1 | Timestamp -> ProtocolBuffer(DTDataList)`

### User preferences

UserPreferences stores the users' preferences.  

* UserPreferences: `collections.NewMap[string, UserPreferences](0xe1)`

## State transitions
