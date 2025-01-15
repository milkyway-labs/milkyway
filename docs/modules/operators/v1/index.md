# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [milkyway/operators/v1/models.proto](#milkyway_operators_v1_models-proto)
    - [Operator](#milkyway-operators-v1-Operator)
    - [OperatorParams](#milkyway-operators-v1-OperatorParams)
  
    - [OperatorStatus](#milkyway-operators-v1-OperatorStatus)
  
- [milkyway/operators/v1/params.proto](#milkyway_operators_v1_params-proto)
    - [Params](#milkyway-operators-v1-Params)
  
- [milkyway/operators/v1/genesis.proto](#milkyway_operators_v1_genesis-proto)
    - [GenesisState](#milkyway-operators-v1-GenesisState)
    - [OperatorParamsRecord](#milkyway-operators-v1-OperatorParamsRecord)
    - [UnbondingOperator](#milkyway-operators-v1-UnbondingOperator)
  
- [milkyway/operators/v1/messages.proto](#milkyway_operators_v1_messages-proto)
    - [MsgDeactivateOperator](#milkyway-operators-v1-MsgDeactivateOperator)
    - [MsgDeactivateOperatorResponse](#milkyway-operators-v1-MsgDeactivateOperatorResponse)
    - [MsgDeleteOperator](#milkyway-operators-v1-MsgDeleteOperator)
    - [MsgDeleteOperatorResponse](#milkyway-operators-v1-MsgDeleteOperatorResponse)
    - [MsgReactivateOperator](#milkyway-operators-v1-MsgReactivateOperator)
    - [MsgReactivateOperatorResponse](#milkyway-operators-v1-MsgReactivateOperatorResponse)
    - [MsgRegisterOperator](#milkyway-operators-v1-MsgRegisterOperator)
    - [MsgRegisterOperatorResponse](#milkyway-operators-v1-MsgRegisterOperatorResponse)
    - [MsgSetOperatorParams](#milkyway-operators-v1-MsgSetOperatorParams)
    - [MsgSetOperatorParamsResponse](#milkyway-operators-v1-MsgSetOperatorParamsResponse)
    - [MsgTransferOperatorOwnership](#milkyway-operators-v1-MsgTransferOperatorOwnership)
    - [MsgTransferOperatorOwnershipResponse](#milkyway-operators-v1-MsgTransferOperatorOwnershipResponse)
    - [MsgUpdateOperator](#milkyway-operators-v1-MsgUpdateOperator)
    - [MsgUpdateOperatorResponse](#milkyway-operators-v1-MsgUpdateOperatorResponse)
    - [MsgUpdateParams](#milkyway-operators-v1-MsgUpdateParams)
    - [MsgUpdateParamsResponse](#milkyway-operators-v1-MsgUpdateParamsResponse)
  
    - [Msg](#milkyway-operators-v1-Msg)
  
- [milkyway/operators/v1/query.proto](#milkyway_operators_v1_query-proto)
    - [QueryOperatorParamsRequest](#milkyway-operators-v1-QueryOperatorParamsRequest)
    - [QueryOperatorParamsResponse](#milkyway-operators-v1-QueryOperatorParamsResponse)
    - [QueryOperatorRequest](#milkyway-operators-v1-QueryOperatorRequest)
    - [QueryOperatorResponse](#milkyway-operators-v1-QueryOperatorResponse)
    - [QueryOperatorsRequest](#milkyway-operators-v1-QueryOperatorsRequest)
    - [QueryOperatorsResponse](#milkyway-operators-v1-QueryOperatorsResponse)
    - [QueryParamsRequest](#milkyway-operators-v1-QueryParamsRequest)
    - [QueryParamsResponse](#milkyway-operators-v1-QueryParamsResponse)
  
    - [Query](#milkyway-operators-v1-Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="milkyway_operators_v1_models-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/operators/v1/models.proto



<a name="milkyway-operators-v1-Operator"></a>

### Operator
Operator defines the fields of an operator


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint32](#uint32) |  | ID is the auto-generated unique identifier for the operator |
| status | [OperatorStatus](#milkyway-operators-v1-OperatorStatus) |  | Status is the status of the operator |
| admin | [string](#string) |  | Admin is the address of the user that can manage the operator |
| moniker | [string](#string) |  | Moniker is the identifier of the operator |
| website | [string](#string) |  | Website is the website of the operator |
| picture_url | [string](#string) |  | PictureURL is the URL of the picture of the operator |
| address | [string](#string) |  | Address is the address of the account associated to the operator. This will be used to store tokens that are delegated to this operator. |
| tokens | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Tokens define the delegated tokens. |
| delegator_shares | [cosmos.base.v1beta1.DecCoin](#cosmos-base-v1beta1-DecCoin) | repeated | DelegatorShares define the total shares issued to an operator&#39;s delegators. |






<a name="milkyway-operators-v1-OperatorParams"></a>

### OperatorParams
OperatorParams represent the params that have been set for an individual
operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| commission_rate | [string](#string) |  | CommissionRate defines the commission rate charged to delegators, as a fraction. |





 


<a name="milkyway-operators-v1-OperatorStatus"></a>

### OperatorStatus
OperatorStatus defines the possible statuses of an operator

| Name | Number | Description |
| ---- | ------ | ----------- |
| OPERATOR_STATUS_UNSPECIFIED | 0 | OPERATOR_STATUS_UNSPECIFIED defines an unspecified status |
| OPERATOR_STATUS_ACTIVE | 1 | OPERATOR_STATUS_ACTIVE identifies an active operator which is providing services |
| OPERATOR_STATUS_INACTIVATING | 2 | OPERATOR_STATUS_INACTIVATING identifies an operator that is in the process of becoming inactive |
| OPERATOR_STATUS_INACTIVE | 3 | OPERATOR_STATUS_INACTIVE defines an inactive operator that is not providing services |


 

 

 



<a name="milkyway_operators_v1_params-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/operators/v1/params.proto



<a name="milkyway-operators-v1-Params"></a>

### Params
Params defines the parameters for the operators module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_registration_fee | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | OperatorRegistrationFee represents the fee that an operator must pay in order to register itself with the network. The fee is drawn from the MsgRegisterOperator sender&#39;s account and transferred to the community pool. |
| deactivation_time | [int64](#int64) |  | DeactivationTime represents the amount of time that will pass between the time that an operator signals its willingness to deactivate and the time that it actually becomes inactive. |





 

 

 

 



<a name="milkyway_operators_v1_genesis-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/operators/v1/genesis.proto



<a name="milkyway-operators-v1-GenesisState"></a>

### GenesisState
GenesisState defines the operators module&#39;s genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-operators-v1-Params) |  | Params defines the parameters of the module. |
| next_operator_id | [uint32](#uint32) |  | NextOperatorID defines the ID that will be assigned to the next operator that gets created. |
| operators | [Operator](#milkyway-operators-v1-Operator) | repeated | Operators defines the list of operators. |
| unbonding_operators | [UnbondingOperator](#milkyway-operators-v1-UnbondingOperator) | repeated | UnbondingOperators defines the list of operators that are currently being unbonded. |
| operators_params | [OperatorParamsRecord](#milkyway-operators-v1-OperatorParamsRecord) | repeated | OperatorsParams defines the list of operators params. |






<a name="milkyway-operators-v1-OperatorParamsRecord"></a>

### OperatorParamsRecord
OperatorParamsRecord represents the params that have been set for an
individual operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorID is the ID of the operator. |
| params | [OperatorParams](#milkyway-operators-v1-OperatorParams) |  | Params defines the parameters for the operators module. |






<a name="milkyway-operators-v1-UnbondingOperator"></a>

### UnbondingOperator
UnbondingOperator contains the data about an operator that is currently being
unbonded.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorID is the ID of the operator that is being unbonded. |
| unbonding_completion_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | UnbondingCompletionTime is the time at which the unbonding of the operator will be completed |





 

 

 

 



<a name="milkyway_operators_v1_messages-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/operators/v1/messages.proto



<a name="milkyway-operators-v1-MsgDeactivateOperator"></a>

### MsgDeactivateOperator
MsgDeactivateOperator defines the message structure for the
DeactivateOperator gRPC service method. It allows the operator owner to
signal that the operator will become inactive. This should be used to signal
users that the operator is going to stop performing services and they should
switch to another operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user deactivating the operator |
| operator_id | [uint32](#uint32) |  | OperatorID represents the ID of the operator to be deregistered |






<a name="milkyway-operators-v1-MsgDeactivateOperatorResponse"></a>

### MsgDeactivateOperatorResponse
MsgDeactivateOperatorResponse is the return value of MsgDeactivateOperator.






<a name="milkyway-operators-v1-MsgDeleteOperator"></a>

### MsgDeleteOperator
MsgDeleteOperator defines the message structure for the
DeleteOperator gRPC service method. It allows the operator owner to
delete a deactivated operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user deleting the operator |
| operator_id | [uint32](#uint32) |  | OperatorID represents the ID of the operator to be deleted |






<a name="milkyway-operators-v1-MsgDeleteOperatorResponse"></a>

### MsgDeleteOperatorResponse
MsgDeleteOperatorResponse is the return value of MsgDeleteOperator.






<a name="milkyway-operators-v1-MsgReactivateOperator"></a>

### MsgReactivateOperator
MsgReactivateOperator defines the message structure for the
ReactivateOperator gRPC service method. It allows the operator owner to
reactivate an inactive operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user reactivating the operator |
| operator_id | [uint32](#uint32) |  | OperatorID represents the ID of the operator to be reactivated |






<a name="milkyway-operators-v1-MsgReactivateOperatorResponse"></a>

### MsgReactivateOperatorResponse
MsgReactivateOperatorResponse is the return value of MsgReactivateOperator.






<a name="milkyway-operators-v1-MsgRegisterOperator"></a>

### MsgRegisterOperator
MsgRegisterOperator defines the message structure for the RegisterOperator
gRPC service method. It allows an account to register a new operator that can
opt-in to validate various services. It requires a sender address as well as
the details of the operator to be registered.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user registering the operator |
| moniker | [string](#string) |  | Moniker is the moniker of the operator |
| website | [string](#string) |  | Website is the website of the operator (optional) |
| picture_url | [string](#string) |  | PictureURL is the URL of operator picture (optional) |
| fee_amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | FeeAmount represents the fees that are going to be paid to create the operator. These should always be greater or equals of any of the coins specified inside the OperatorRegistrationFee field of the modules params. If no fees are specified inside the module parameters, this field can be omitted. |






<a name="milkyway-operators-v1-MsgRegisterOperatorResponse"></a>

### MsgRegisterOperatorResponse
MsgRegisterOperatorResponse is the return value of MsgRegisterOperator.
It returns the newly created operator ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| new_operator_id | [uint32](#uint32) |  | NewOperatorID is the ID of the newly registered operator |






<a name="milkyway-operators-v1-MsgSetOperatorParams"></a>

### MsgSetOperatorParams
MsgSetOperatorParams defines the message structure for the
SetOperatorParams gRPC service method. It allows the operator admin to
update the operator&#39;s parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  |  |
| operator_id | [uint32](#uint32) |  |  |
| params | [OperatorParams](#milkyway-operators-v1-OperatorParams) |  |  |






<a name="milkyway-operators-v1-MsgSetOperatorParamsResponse"></a>

### MsgSetOperatorParamsResponse
MsgSetOperatorParamsResponse is the return value of
MsgSetOperatorParams.






<a name="milkyway-operators-v1-MsgTransferOperatorOwnership"></a>

### MsgTransferOperatorOwnership
MsgTransferOperatorOwnership defines the message structure for the
TransferOperatorOwnership gRPC service method. It allows an operator admin to
transfer the ownership of the operator to another account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user transferring the ownership |
| operator_id | [uint32](#uint32) |  | OperatorID represents the ID of the operator to transfer ownership |
| new_admin | [string](#string) |  | NewAdmin is the address of the new admin of the operator |






<a name="milkyway-operators-v1-MsgTransferOperatorOwnershipResponse"></a>

### MsgTransferOperatorOwnershipResponse
MsgTransferOperatorOwnershipResponse is the return value of
MsgTransferOperatorOwnership.






<a name="milkyway-operators-v1-MsgUpdateOperator"></a>

### MsgUpdateOperator
MsgUpdateOperator defines the message structure for the UpdateOperator gRPC
service method. It allows the operator owner to update the details of an
existing operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user updating the operator |
| operator_id | [uint32](#uint32) |  | OperatorID represents the ID of the operator to be updated |
| moniker | [string](#string) |  | Moniker is the new moniker of the operator. If it shouldn&#39;t be changed, use [do-not-modify] instead. |
| website | [string](#string) |  | Website is the new website of the operator. If it shouldn&#39;t be changed, use [do-not-modify] instead. |
| picture_url | [string](#string) |  | PictureURL is the new URL of the operator picture. If it shouldn&#39;t be changed, use [do-not-modify] instead. |






<a name="milkyway-operators-v1-MsgUpdateOperatorResponse"></a>

### MsgUpdateOperatorResponse
MsgUpdateOperatorResponse is the return value of MsgUpdateOperator.






<a name="milkyway-operators-v1-MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams defines the message structure for the UpdateParams gRPC
service method. It allows the authority to update the module parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authority | [string](#string) |  | Authority is the address that controls the module (defaults to x/gov unless overwritten). |
| params | [Params](#milkyway-operators-v1-Params) |  | Params define the parameters to update.

NOTE: All parameters must be supplied. |






<a name="milkyway-operators-v1-MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse is the return value of MsgUpdateParams.





 

 

 


<a name="milkyway-operators-v1-Msg"></a>

### Msg
Msg defines the avs module&#39;s gRPC message service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RegisterOperator | [MsgRegisterOperator](#milkyway-operators-v1-MsgRegisterOperator) | [MsgRegisterOperatorResponse](#milkyway-operators-v1-MsgRegisterOperatorResponse) | RegisterOperator defines the operation for registering a new operator. |
| UpdateOperator | [MsgUpdateOperator](#milkyway-operators-v1-MsgUpdateOperator) | [MsgUpdateOperatorResponse](#milkyway-operators-v1-MsgUpdateOperatorResponse) | UpdateOperator defines the operation for updating an operator&#39;s details. The operator owner can update the moniker, website, and picture URL. |
| DeactivateOperator | [MsgDeactivateOperator](#milkyway-operators-v1-MsgDeactivateOperator) | [MsgDeactivateOperatorResponse](#milkyway-operators-v1-MsgDeactivateOperatorResponse) | DeactivateOperator defines the operation for deactivating an operator. Operators will require some time in order to be deactivated. This time is defined by the governance parameters. |
| ReactivateOperator | [MsgReactivateOperator](#milkyway-operators-v1-MsgReactivateOperator) | [MsgReactivateOperatorResponse](#milkyway-operators-v1-MsgReactivateOperatorResponse) | ReactivateOperator defines the operation for reactivating an inactive operator. |
| DeleteOperator | [MsgDeleteOperator](#milkyway-operators-v1-MsgDeleteOperator) | [MsgDeleteOperatorResponse](#milkyway-operators-v1-MsgDeleteOperatorResponse) | DeleteOperator defines the operation for deleting a deactivated operator. |
| TransferOperatorOwnership | [MsgTransferOperatorOwnership](#milkyway-operators-v1-MsgTransferOperatorOwnership) | [MsgTransferOperatorOwnershipResponse](#milkyway-operators-v1-MsgTransferOperatorOwnershipResponse) | TransferOperatorOwnership defines the operation for transferring the ownership of an operator to another account. |
| SetOperatorParams | [MsgSetOperatorParams](#milkyway-operators-v1-MsgSetOperatorParams) | [MsgSetOperatorParamsResponse](#milkyway-operators-v1-MsgSetOperatorParamsResponse) | SetOperatorParams defines the operation for setting a operator&#39;s parameters. |
| UpdateParams | [MsgUpdateParams](#milkyway-operators-v1-MsgUpdateParams) | [MsgUpdateParamsResponse](#milkyway-operators-v1-MsgUpdateParamsResponse) | UpdateParams defines a governance operation for updating the module parameters. The authority defaults to the x/gov module account. |

 



<a name="milkyway_operators_v1_query-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/operators/v1/query.proto



<a name="milkyway-operators-v1-QueryOperatorParamsRequest"></a>

### QueryOperatorParamsRequest
QueryOperatorParamsRequest is the request type for the Query/OperatorParams
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorID is the ID of the operator for which to query the params |






<a name="milkyway-operators-v1-QueryOperatorParamsResponse"></a>

### QueryOperatorParamsResponse
QueryOperatorParamsResponse is the response type for the Query/OperatorParams
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_params | [OperatorParams](#milkyway-operators-v1-OperatorParams) |  |  |






<a name="milkyway-operators-v1-QueryOperatorRequest"></a>

### QueryOperatorRequest
QueryOperatorRequest is the request type for the Query/Operator RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorId is the ID of the operator to query |






<a name="milkyway-operators-v1-QueryOperatorResponse"></a>

### QueryOperatorResponse
QueryOperatorResponse is the response type for the Query/Operator RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator | [Operator](#milkyway-operators-v1-Operator) |  |  |






<a name="milkyway-operators-v1-QueryOperatorsRequest"></a>

### QueryOperatorsRequest
QueryOperatorsRequest is the request type for the Query/Operators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  |  |






<a name="milkyway-operators-v1-QueryOperatorsResponse"></a>

### QueryOperatorsResponse
QueryOperatorsResponse is the response type for the Query/Operators RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operators | [Operator](#milkyway-operators-v1-Operator) | repeated | Operators is the list of operators |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination response |






<a name="milkyway-operators-v1-QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="milkyway-operators-v1-QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-operators-v1-Params) |  |  |





 

 

 


<a name="milkyway-operators-v1-Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Operator | [QueryOperatorRequest](#milkyway-operators-v1-QueryOperatorRequest) | [QueryOperatorResponse](#milkyway-operators-v1-QueryOperatorResponse) | Operator defines a gRPC query method that returns the operator by the given operator id. |
| OperatorParams | [QueryOperatorParamsRequest](#milkyway-operators-v1-QueryOperatorParamsRequest) | [QueryOperatorParamsResponse](#milkyway-operators-v1-QueryOperatorParamsResponse) | OperatorParams defines a gRPC query method that returns the operator&#39;s params by the given operator id. |
| Operators | [QueryOperatorsRequest](#milkyway-operators-v1-QueryOperatorsRequest) | [QueryOperatorsResponse](#milkyway-operators-v1-QueryOperatorsResponse) | Operators defines a gRPC query method that returns the list of operators. |
| Params | [QueryParamsRequest](#milkyway-operators-v1-QueryParamsRequest) | [QueryParamsResponse](#milkyway-operators-v1-QueryParamsResponse) | Params defines a gRPC query method that returns the parameters of the module. |

 



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

