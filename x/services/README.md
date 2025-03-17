# `x/services`

## Abstract

The following document specifies the services module.

This module allows the registration and management of metadata related to Actively Validated Services (AVS) within the
MilkyWay restaking protocol.

## Contents

## Concepts

### Service

A service is the on-chain representation of an application that requires operators to run its software off-chain to
provide it with security, in exchange for rewards. Services are responsible for coding the off-chain program that need
to be run by operators, and to design their rewards distributions in order to make sure those operators are incentivized
to run their software. Also, service administrators are responsible for monitoring the off-chain execution of their
software, and to slash any operator that misbehaves.

When registering a new service, it is automatically assigned a new ID, which is incremental and unique across all
services.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/models.proto#L29-L78
```

A service can have one of the following statuses:

* `CREATED`: the service has been created but is not yet active.
* `ACTIVE`: the service is currently running and accepting for operators to validate it. This means that it's
  distributing rewards to operators and is monitoring their behavior.
* `INACTIVE`: the service is not currently running and is not accepting for operators to validate it. this means that
  it's not distributing rewards to operators and is not monitoring their behavior.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/models.proto#L11-L27
```

### Service params

Each service can set a series of parameters that only apply to them. These are collectively known as the service params
and can be edited without any notice from the service admin.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/models.proto#L80-L85
```

## State

### Params

The module params are stored using the `0x01` key:

* Params: `0x01 -> ProtocolBuffer(params)`

### Next Service ID

The ID that will be assigned to the next registered service is stored using the `0xa1` key:

* Next service ID: `0xa1 -> uint32`

### Services

Each service is stored in state with the prefix of `0xa2`:

* Service: `0xa2 | ServiceID -> ProtocolBuffer(Service)`

### Service addresses

In order to know more easily and faster if a particular address represents a service, we store the set of service
addresses using the `0xa3` prefix:

* Service address: `0xa3 | ServiceAddress -> []byte{}`

### Service params

Each service's params are stored using the `0xa4` prefix:

* Service params: `0xa4 | ServiceID -> ProtocolBuffer(ServiceParams)`

## Messages

### MsgCreteService

The `MsgCreateService` can be sent by anyone to create a new service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L58-L93
```

The message will fail under the following conditions:

* The service data are not valid
* The user creating the service has not enough funds to pay for the creation fee set inside the module's params

This message returns a `MsgCreateServiceResponse` that contains the ID of the newly created service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L95-L100
```

### MsgUpdateService

The `MsgUpdateService` can be sent by the service admin to update the service's data.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L102-L133
```

The message will fail under the following conditions:

* The service data are not valid
* The user updating the service is not the service admin

### MsgActivateService

The `MsgActivateService` can be sent by the service admin to activate the service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L138-L148
```

The message will fail under the following conditions:

* The service is already active
* The user activating the service is not the service admin

### MsgDeactivateService

The `MsgDeactivateService` can be sent by the service admin to deactivate the service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L153-L165
```

The message will fail under the following conditions:

* The service is already inactive
* The user deactivating the service is not the service admin

### MsgDeleteService

The `MsgDeleteService` can be sent by the service admin to delete the service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L170-L182
```

The message will fail under the following conditions:

* The user deleting the service is not the service admin

### MsgTransferServiceOwnership

The `MsgTransferServiceOwnership` can be sent by the service admin to transfer the ownership of the service to another
address.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L187-L202
```

The message will fail under the following conditions:

* The user transferring the service's ownership is not the service admin

### MsgSetServiceParams

The `MsgSetServiceParams` can be sent by the service admin to set the service's params.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L208-L223
```

The message will fail under the following conditions:

* The params are not valid
* The user setting the service's params is not the service admin

### MsgAccreditService

The `MsgAccreditService` can be sent by the gov module to recognize a service as "accredited".

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L250-L265
```

The message will fail under the following conditions:

* The sender is not the god module account

### MsgRevokeServiceAccreditation

The `MsgRevokeServiceAccreditation` can be sent by the gov module to revoke the accreditation of a service.

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/messages.proto#L270-L287
```

The message will fail under the following conditions:

* The sender is not the god module account

## Events

### Handlers

#### MsgCreateService

| Type             | Attribute Key | Attribute Value   |
|------------------|---------------|-------------------|
| `create_service` | `service_id`  | `{serviceID}`     |
| `create_service` | `sender`      | `{senderAddress}` |

#### MsgUpdateService

| Type             | Attribute Key | Attribute Value   |
|------------------|---------------|-------------------|
| `update_service` | `service_id`  | `{serviceID}`     |
| `update_service` | `sender`      | `{senderAddress}` |

#### MsgActivateService

| Type               | Attribute Key | Attribute Value   |
|--------------------|---------------|-------------------|
| `activate_service` | `service_id`  | `{serviceID}`     |
| `activate_service` | `sender`      | `{senderAddress}` |

#### MsgDeactivateService

| Type                 | Attribute Key | Attribute Value   |
|----------------------|---------------|-------------------|
| `deactivate_service` | `service_id`  | `{serviceID}`     |
| `deactivate_service` | `sender`      | `{senderAddress}` |

#### MsgDeleteService

| Type             | Attribute Key | Attribute Value   |
|------------------|---------------|-------------------|
| `delete_service` | `service_id`  | `{serviceID}`     |
| `delete_service` | `sender`      | `{senderAddress}` |

#### MsgTransferServiceOwnership

| Type                         | Attribute Key | Attribute Value     |
|------------------------------|---------------|---------------------|
| `transfer_service_ownership` | `service_id`  | `{serviceID}`       |
| `transfer_service_ownership` | `new_admin`   | `{newAdminAddress}` |
| `transfer_service_ownership` | `sender`      | `{senderAddress}`   |

#### MsgSetServiceParams

| Type                 | Attribute Key | Attribute Value   |
|----------------------|---------------|-------------------|
| `set_service_params` | `service_id`  | `{serviceID}`     |
| `set_service_params` | `sender`      | `{senderAddress}` |

#### MsgAccreditService

| Type               | Attribute Key | Attribute Value |
|--------------------|---------------|-----------------|
| `accredit_service` | `service_id`  | `{serviceID}`   |

#### MsgRevokeServiceAccreditation

| Type                           | Attribute Key | Attribute Value |
|--------------------------------|---------------|-----------------|
| `revoke_service_accreditation` | `service_id`  | `{serviceID}`   |

## Params

The services module contains the following parameters:

```protobuf reference
https://github.com/milkyway-labs/milkyway/blob/v8.1.0/proto/milkyway/services/v1/params.proto#L9-L19
``` 
