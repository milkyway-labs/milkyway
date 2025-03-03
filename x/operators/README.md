# `x/operators`

## Abstract

The following document specifies the operators module.

This module allows the registration and management of metadata related to restaking operators.

## Contents

* [Concepts](#concepts)
   * [Operator](#operator)
   * [Operator Params](#operator-params)
   * [Inactivating queue](#inactivating-queue)
* [State](#state)
   * [Params](#params)
   * [Next Operator ID](#next-operator-id)
   * [Operators](#operators)
   * [Inactivating queue](#inactivating-queue)
   * [Operator addresses](#operator-addresses)
   * [Operator params](#operator-params)
* [Messages](#messages)
   * [MsgRegisterOperator](#msgregisteroperator)
   * [MsgUpdateOperator](#msgupdateoperator)
   * [MsgDeactivateOperator](#msgdeactivateoperator)
   * [MsgReactivateOperator](#msgreactivateoperator)
   * [MsgDeleteOperator](#msgdeleteoperator)
   * [MsgSetOperatorParams](#msgsetoperatorparams)
   * [MsgTransferOperatorOwnership](#msgtransferoperatorownership)
* [Events](#events)
   * [BeginBlocker](#beginblocker)
   * [Handlers](#handlers)
      * [MsgRegisterOperator](#msgregisteroperator)
      * [MsgUpdateOperator](#msgupdateoperator)
      * [MsgDeactivateOperator](#msgdeactivateoperator)
      * [MsgReactivateOperator](#msgreactivateoperator)
      * [MsgDeleteOperator](#msgdeleteoperator)
      * [MsgSetOperatorParams](#msgsetoperatorparams)
      * [MsgTransferOperatorOwnership](#msgtransferoperatorownership)
* [Parameters](#parameters)

## Concepts

### Operator

An operator is the on-chain representation of an individual or company that is responsible for running off-chain
programs for each of the services that they are partaking. Operators are responsible for the uptime and the
correct operation of the services that they are running, and can be slashed if they are found to be acting
maliciously or if they are not providing the services that they are supposed to.

When registering a new operator, it is automatically assigned a new ID, which is incremental and unique across all
operators.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/models.proto#L31-L66
```

An operator can have one of the following statuses:

* `ACTIVE`: The operator is currently running the services that they are responsible for, receiving rewards for their
  services and being eligible for slashing.
* `INACTIVATING`: The operator has declared their intention of becoming inactive. In this state, an operator is no
  longer eligible for rewards, but is still eligible for slashing.
* `INACTIVE`: The operator is no longer running the services that they are responsible for, and is no longer eligible
  for rewards nor slashing.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/models.proto#L11-L29
```

### Operator Params

Each operator can set a series of parameters that only apply to them. These parameters are collectively called
`OperatorParams` and can be edited without any previous notice from the operator's admin.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/models.proto#L68-L77
```

### Inactivating queue

The inactivating queue is a list of operators that have declared their intention of becoming inactive. This queue is
used to keep track of the operators that are in the process of becoming inactive, and to ensure that they are not
eligible for rewards while they are in this state.

## State

### Params

The module params are stored using the `0x01` key:

* Params: `0x01 -> ProtocolBuffer(params)`

### Next Operator ID

The ID that will be assigned to the next registered operator is stored using the `0xa1` key:

* Next operator ID: `0xa1 -> uint32`

### Operators

Each operator is stored in state with the prefix of `0xa2`:

* Operator: `0xa2 | OperatorID -> ProtocolBuffer(Operator)`

### Inactivating queue

Each part of the inactivating queue is stored using the `0xa3` prefix:

* Inactivating queue: `0xa3 | InactivatingEndTime | OperatorID -> OperatorID`

### Operator addresses

In order to know more easily and faster if a particular address represents an operator, we store the set of operator
addresses using the `0xa4` prefix:

* Operator address: `0xa4 | OperatorAddress -> []byte{}`

### Operator params

Each operator's params are stored using the `0xa5` prefix:

* Operator params: `0xa5 | OperatorID -> ProtocolBuffer(OperatorParams)`

## Messages

### MsgRegisterOperator

The `MsgRegisterOperator` can be sent by anyone to register a new operator.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/messages.proto#L51-L83
```

The message will fail under the following conditions:

* The operator data are not valid
* The user registering for the operator has not enough funds to pay for the registration fee set inside the module's
  params

This message returns a `MsgRegisterOperatorResponse` that contains the ID of the newly registered operator.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/messages.proto#L85-L90
```

### MsgUpdateOperator

The `MsgUpdateOperator` can be sent by the operator's admin to update the operator's data.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/messages.proto#L92-L119
```

The message will fail under the following conditions:

* The operator data are not valid
* The user updating the operator is not the operator's admin

### MsgDeactivateOperator

The `MsgDeactivateOperator` can be sent by the operator's admin to declare the intention of deactivating the operator,
and initiate the inactivating process.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/messages.proto#L124-L138
```

The message will fail under the following conditions:

* The operator is already inactive
* The user deactivating the operator is not the operator's admin

### MsgReactivateOperator

The `MsgReactivateOperator` can be sent by the operator's admin to reactivate the operator, after it has been
deactivated.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/messages.proto#L143-L155
```

The message will fail under the following conditions:

* The operator is not inactive
* The user reactivating the operator is not the operator's admin

### MsgDeleteOperator

The `MsgDeleteOperator` can be sent by the operator's admin to delete the operator from the state.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/messages.proto#L160-L172
```

The message will fail under the following conditions:

* The operator is not inactive
* The user deleting the operator is not the operator's admin

### MsgSetOperatorParams

The `MsgSetOperatorParams` can be sent by the operator's admin to set the operator's params.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/messages.proto#L177-L189
```

The message will fail under the following conditions:

* The params are not valid
* The user setting the operator's params is not the operator's admin

### MsgTransferOperatorOwnership

The `MsgTransferOperatorOwnership` can be sent by the operator's admin to transfer the operator's ownership to another
address.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/messages.proto#L195-L210
```

The message will fail under the following conditions:

* The user transferring the operator's ownership is not the operator's admin

## Events

### BeginBlocker

| Type                             | Attribute Key | Attribute Value   |
|----------------------------------|---------------|-------------------|
| `complete_operator_inactivation` | `operator_id` | `{operatorID}`    |
| `complete_operator_inactivation` | `sender`      | `{senderAddress}` |

### Handlers

#### MsgRegisterOperator

| Type                | Attribute Key | Attribute Value   |
|---------------------|---------------|-------------------|
| `register_operator` | `operator_id` | `{operatorID}`    |
| `register_operator` | `sender`      | `{senderAddress}` |

#### MsgUpdateOperator

| Type              | Attribute Key | Attribute Value   |
|-------------------|---------------|-------------------|
| `update_operator` | `operator_id` | `{operatorID}`    |
| `update_operator` | `sender`      | `{senderAddress}` |

#### MsgDeactivateOperator

| Type                          | Attribute Key | Attribute Value   |
|-------------------------------|---------------|-------------------|
| `start_operator_inactivation` | `operator_id` | `{operatorID}`    |
| `start_operator_inactivation` | `sender`      | `{senderAddress}` |

#### MsgReactivateOperator

| Type                  | Attribute Key | Attribute Value   |
|-----------------------|---------------|-------------------|
| `reactivate_operator` | `operator_id` | `{operatorID}`    |
| `reactivate_operator` | `sender`      | `{senderAddress}` |

#### MsgDeleteOperator

| Type               | Attribute Key | Attribute Value   |
|--------------------|---------------|-------------------|
| `delete_opearator` | `operator_id` | `{operatorID}`    |
| `delete_opearator` | `sender`      | `{senderAddress}` |

#### MsgSetOperatorParams

| Type                  | Attribute Key | Attribute Value   |
|-----------------------|---------------|-------------------|
| `set_operator_params` | `operator_id` | `{operatorID}`    |
| `set_operator_params` | `sender`      | `{senderAddress}` |

## Parameters

The operators module contains the following parameters:

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/operators/v1/params.proto#L9-L25
```