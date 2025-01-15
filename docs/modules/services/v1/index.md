# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [milkyway/services/v1/models.proto](#milkyway_services_v1_models-proto)
    - [Service](#milkyway-services-v1-Service)
    - [ServiceParams](#milkyway-services-v1-ServiceParams)
  
    - [ServiceStatus](#milkyway-services-v1-ServiceStatus)
  
- [milkyway/services/v1/params.proto](#milkyway_services_v1_params-proto)
    - [Params](#milkyway-services-v1-Params)
  
- [milkyway/services/v1/genesis.proto](#milkyway_services_v1_genesis-proto)
    - [GenesisState](#milkyway-services-v1-GenesisState)
    - [ServiceParamsRecord](#milkyway-services-v1-ServiceParamsRecord)
  
- [milkyway/services/v1/messages.proto](#milkyway_services_v1_messages-proto)
    - [MsgAccreditService](#milkyway-services-v1-MsgAccreditService)
    - [MsgAccreditServiceResponse](#milkyway-services-v1-MsgAccreditServiceResponse)
    - [MsgActivateService](#milkyway-services-v1-MsgActivateService)
    - [MsgActivateServiceResponse](#milkyway-services-v1-MsgActivateServiceResponse)
    - [MsgCreateService](#milkyway-services-v1-MsgCreateService)
    - [MsgCreateServiceResponse](#milkyway-services-v1-MsgCreateServiceResponse)
    - [MsgDeactivateService](#milkyway-services-v1-MsgDeactivateService)
    - [MsgDeactivateServiceResponse](#milkyway-services-v1-MsgDeactivateServiceResponse)
    - [MsgDeleteService](#milkyway-services-v1-MsgDeleteService)
    - [MsgDeleteServiceResponse](#milkyway-services-v1-MsgDeleteServiceResponse)
    - [MsgRevokeServiceAccreditation](#milkyway-services-v1-MsgRevokeServiceAccreditation)
    - [MsgRevokeServiceAccreditationResponse](#milkyway-services-v1-MsgRevokeServiceAccreditationResponse)
    - [MsgSetServiceParams](#milkyway-services-v1-MsgSetServiceParams)
    - [MsgSetServiceParamsResponse](#milkyway-services-v1-MsgSetServiceParamsResponse)
    - [MsgTransferServiceOwnership](#milkyway-services-v1-MsgTransferServiceOwnership)
    - [MsgTransferServiceOwnershipResponse](#milkyway-services-v1-MsgTransferServiceOwnershipResponse)
    - [MsgUpdateParams](#milkyway-services-v1-MsgUpdateParams)
    - [MsgUpdateParamsResponse](#milkyway-services-v1-MsgUpdateParamsResponse)
    - [MsgUpdateService](#milkyway-services-v1-MsgUpdateService)
    - [MsgUpdateServiceResponse](#milkyway-services-v1-MsgUpdateServiceResponse)
  
    - [Msg](#milkyway-services-v1-Msg)
  
- [milkyway/services/v1/query.proto](#milkyway_services_v1_query-proto)
    - [QueryParamsRequest](#milkyway-services-v1-QueryParamsRequest)
    - [QueryParamsResponse](#milkyway-services-v1-QueryParamsResponse)
    - [QueryServiceParamsRequest](#milkyway-services-v1-QueryServiceParamsRequest)
    - [QueryServiceParamsResponse](#milkyway-services-v1-QueryServiceParamsResponse)
    - [QueryServiceRequest](#milkyway-services-v1-QueryServiceRequest)
    - [QueryServiceResponse](#milkyway-services-v1-QueryServiceResponse)
    - [QueryServicesRequest](#milkyway-services-v1-QueryServicesRequest)
    - [QueryServicesResponse](#milkyway-services-v1-QueryServicesResponse)
  
    - [Query](#milkyway-services-v1-Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="milkyway_services_v1_models-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/services/v1/models.proto



<a name="milkyway-services-v1-Service"></a>

### Service
Service defines the fields of a service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint32](#uint32) |  | ID is the unique identifier of the service |
| status | [ServiceStatus](#milkyway-services-v1-ServiceStatus) |  | Status is the status of the service |
| admin | [string](#string) |  | Admin is the address of the user that has administrative power over the service |
| name | [string](#string) |  | Name is the name of the service |
| description | [string](#string) |  | Description is the description of the service |
| website | [string](#string) |  | Website is the website of the service |
| picture_url | [string](#string) |  | PictureURL is the URL of the picture of the service |
| address | [string](#string) |  | Address is the address of the account associated with the service. This will be used in order to store all the tokens that are delegated to this service by various users. |
| tokens | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Tokens define the delegated tokens. |
| delegator_shares | [cosmos.base.v1beta1.DecCoin](#cosmos-base-v1beta1-DecCoin) | repeated | DelegatorShares define the total shares issued to a service&#39;s delegators. |
| accredited | [bool](#bool) |  | Accredited defines if the service is accredited. Note: We use this term instead of &#34;trusted&#34; of &#34;verified&#34; in order to represent something more generic. Initially, services will be accredited by the on-chain governance process. In the future, we may add more ways to accredit services (e.g. automatic ones based on the operators that decide to run the service, or the amount of cryptoeconomic security that the service was able to capture). |






<a name="milkyway-services-v1-ServiceParams"></a>

### ServiceParams
ServiceParams defines the parameters of a service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| allowed_denoms | [string](#string) | repeated | AllowedDenoms defines the list of denoms that can be restaked toward the service. If the list is empty, any denom can be used. |





 


<a name="milkyway-services-v1-ServiceStatus"></a>

### ServiceStatus
ServiceStatus defines the status of a service

| Name | Number | Description |
| ---- | ------ | ----------- |
| SERVICE_STATUS_UNSPECIFIED | 0 | SERVICE_STATUS_UNSPECIFIED defines an unspecified status |
| SERVICE_STATUS_CREATED | 1 | SERVICE_STATUS_CREATED identifies a recently created service that is not yet active |
| SERVICE_STATUS_ACTIVE | 2 | SERVICE_STATUS_ACTIVE identifies an active service |
| SERVICE_STATUS_INACTIVE | 3 | SERVICE_STATUS_INACTIVE identifies an inactive service |


 

 

 



<a name="milkyway_services_v1_params-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/services/v1/params.proto



<a name="milkyway-services-v1-Params"></a>

### Params
Params defines the parameters for the module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_registration_fee | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | ServiceRegistrationFee defines the fee to register a new service. The fee is drawn from the MsgRegisterService sender&#39;s account, and transferred to the community pool. |





 

 

 

 



<a name="milkyway_services_v1_genesis-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/services/v1/genesis.proto



<a name="milkyway-services-v1-GenesisState"></a>

### GenesisState
GenesisState defines the services module&#39;s genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-services-v1-Params) |  | Params defines the parameters of the module. |
| services | [Service](#milkyway-services-v1-Service) | repeated | Services defines the list of services. |
| next_service_id | [uint32](#uint32) |  | NextServiceID defines the ID that will be assigned to the next service that gets created. |
| services_params | [ServiceParamsRecord](#milkyway-services-v1-ServiceParamsRecord) | repeated | ServicesParams defines the list of service parameters. |






<a name="milkyway-services-v1-ServiceParamsRecord"></a>

### ServiceParamsRecord
ServiceParamsRecord represents the parameters that have been set for
a specific service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceID represents the ID of the service to which the parameters should be set. |
| params | [ServiceParams](#milkyway-services-v1-ServiceParams) |  | Params represents the parameters that should be set to the service. |





 

 

 

 



<a name="milkyway_services_v1_messages-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/services/v1/messages.proto



<a name="milkyway-services-v1-MsgAccreditService"></a>

### MsgAccreditService
MsgAccreditService defines the message structure for the AccreditService gRPC
service method. It allows the authority to accredit a service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authority | [string](#string) |  | Authority is the address that controls the module (defaults to x/gov unless overwritten). |
| service_id | [uint32](#uint32) |  | ServiceID represents the ID of the service to be accredited |






<a name="milkyway-services-v1-MsgAccreditServiceResponse"></a>

### MsgAccreditServiceResponse
MsgAccreditServiceResponse is the return value of MsgAccreditService.






<a name="milkyway-services-v1-MsgActivateService"></a>

### MsgActivateService
MsgActivateService defines the message structure for the ActivateService gRPC


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user that wants to activate the service |
| service_id | [uint32](#uint32) |  | ServiceID represents the ID of the service to be activated |






<a name="milkyway-services-v1-MsgActivateServiceResponse"></a>

### MsgActivateServiceResponse
MsgActivateServiceResponse is the return value of MsgActivateService.






<a name="milkyway-services-v1-MsgCreateService"></a>

### MsgCreateService
MsgCreateServiceResponse defines the message structure for the
CreateService gRPC service method. It allows an account to register a new
service that can be validated by operators. It requires a sender address
as well as the details of the service to be registered.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user registering the service |
| name | [string](#string) |  | Name is the name of the service |
| description | [string](#string) |  | Description is the description of the service |
| website | [string](#string) |  | Website is the website of the service |
| picture_url | [string](#string) |  | PictureURL is the URL of the service picture |
| fee_amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | FeeAmount represents the fees that are going to be paid to create the service. These should always be greater or equals of any of the coins specified inside the ServiceRegistrationFee field of the modules params. If no fees are specified inside the module parameters, this field can be omitted. |






<a name="milkyway-services-v1-MsgCreateServiceResponse"></a>

### MsgCreateServiceResponse
MsgCreateServiceResponse is the return value of MsgCreateService.
It returns the newly created service ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| new_service_id | [uint32](#uint32) |  | NewServiceID is the ID of the newly registered service |






<a name="milkyway-services-v1-MsgDeactivateService"></a>

### MsgDeactivateService
MsgDeactivateService defines the message structure for the DeactivateService
gRPC service method. It allows the service admin to deactivate an existing
service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user that wants to deactivate the service |
| service_id | [uint32](#uint32) |  | ServiceID represents the ID of the service to be deactivated |






<a name="milkyway-services-v1-MsgDeactivateServiceResponse"></a>

### MsgDeactivateServiceResponse
MsgDeactivateServiceResponse is the return value of MsgDeactivateService.






<a name="milkyway-services-v1-MsgDeleteService"></a>

### MsgDeleteService
MsgDeleteService defines the message structure for the DeleteService
gRPC service method. It allows the service admin to delete a previously
deactivated service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user that wants to delete the service |
| service_id | [uint32](#uint32) |  | ServiceID represents the ID of the service to be deleted |






<a name="milkyway-services-v1-MsgDeleteServiceResponse"></a>

### MsgDeleteServiceResponse
MsgDeleteServiceResponse is the return value of MsgDeleteService.






<a name="milkyway-services-v1-MsgRevokeServiceAccreditation"></a>

### MsgRevokeServiceAccreditation
MsgRevokeServiceAccreditation defines the message structure for the
RevokeServiceAccreditation gRPC service method. It allows the authority to
revoke a service&#39;s accreditation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authority | [string](#string) |  | Authority is the address that controls the module (defaults to x/gov unless overwritten). |
| service_id | [uint32](#uint32) |  | ServiceID represents the ID of the service to have its accreditation revoked |






<a name="milkyway-services-v1-MsgRevokeServiceAccreditationResponse"></a>

### MsgRevokeServiceAccreditationResponse
MsgRevokeServiceAccreditationResponse is the return value of
MsgRevokeServiceAccreditation.






<a name="milkyway-services-v1-MsgSetServiceParams"></a>

### MsgSetServiceParams
MsgSetServiceParams defines the message structure for the
SetServiceParams gRPC service method. It allows a service admin to
update the parameters of a service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user setting the parameters |
| service_id | [uint32](#uint32) |  | ServiceID is the ID of the service whose parameters are being set |
| service_params | [ServiceParams](#milkyway-services-v1-ServiceParams) |  | ServiceParams defines the new parameters of the service |






<a name="milkyway-services-v1-MsgSetServiceParamsResponse"></a>

### MsgSetServiceParamsResponse
MsgSetServiceParamsResponse is the return value of MsgSetServiceParams.






<a name="milkyway-services-v1-MsgTransferServiceOwnership"></a>

### MsgTransferServiceOwnership
MsgTransferServiceOwnership defines the message structure for the
TransferServiceOwnership gRPC service method. It allows a service admin to
transfer the ownership of the service to another account.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user transferring the ownership |
| service_id | [uint32](#uint32) |  | ServiceID represents the ID of the service to transfer ownership |
| new_admin | [string](#string) |  | NewAdmin is the address of the new admin of the service |






<a name="milkyway-services-v1-MsgTransferServiceOwnershipResponse"></a>

### MsgTransferServiceOwnershipResponse
MsgTransferServiceOwnershipResponse is the return value of
MsgTransferServiceOwnership.






<a name="milkyway-services-v1-MsgUpdateParams"></a>

### MsgUpdateParams
MsgDeactivateService defines the message structure for the UpdateParams gRPC
service method. It allows the authority to update the module parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authority | [string](#string) |  | Authority is the address that controls the module (defaults to x/gov unless overwritten). |
| params | [Params](#milkyway-services-v1-Params) |  | Params define the parameters to update.

NOTE: All parameters must be supplied. |






<a name="milkyway-services-v1-MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgDeactivateServiceResponse is the return value of MsgUpdateParams.






<a name="milkyway-services-v1-MsgUpdateService"></a>

### MsgUpdateService
MsgUpdateService defines the message structure for the UpdateService gRPC
service method. It allows the service admin to update the details of
an existing service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | Sender is the address of the user updating the service |
| service_id | [uint32](#uint32) |  | ID represents the ID of the service to be updated |
| name | [string](#string) |  | Name is the new name of the service. If it shouldn&#39;t be changed, use [do-not-modify] instead. |
| description | [string](#string) |  | Description is the new description of the service. If it shouldn&#39;t be changed, use [do-not-modify] instead. |
| website | [string](#string) |  | Website is the new website of the service. If it shouldn&#39;t be changed, use [do-not-modify] instead. |
| picture_url | [string](#string) |  | PictureURL is the new URL of the service picture. If it shouldn&#39;t be changed, use [do-not-modify] instead. |






<a name="milkyway-services-v1-MsgUpdateServiceResponse"></a>

### MsgUpdateServiceResponse
MsgUpdateServiceResponse is the return value of MsgUpdateService.





 

 

 


<a name="milkyway-services-v1-Msg"></a>

### Msg
Msg defines the services module&#39;s gRPC message service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateService | [MsgCreateService](#milkyway-services-v1-MsgCreateService) | [MsgCreateServiceResponse](#milkyway-services-v1-MsgCreateServiceResponse) | CreateService defines the operation for registering a new service. |
| UpdateService | [MsgUpdateService](#milkyway-services-v1-MsgUpdateService) | [MsgUpdateServiceResponse](#milkyway-services-v1-MsgUpdateServiceResponse) | UpdateService defines the operation for updating an existing service. |
| ActivateService | [MsgActivateService](#milkyway-services-v1-MsgActivateService) | [MsgActivateServiceResponse](#milkyway-services-v1-MsgActivateServiceResponse) | ActivateService defines the operation for activating an existing service. |
| DeactivateService | [MsgDeactivateService](#milkyway-services-v1-MsgDeactivateService) | [MsgDeactivateServiceResponse](#milkyway-services-v1-MsgDeactivateServiceResponse) | DeactivateService defines the operation for deactivating an existing service. |
| DeleteService | [MsgDeleteService](#milkyway-services-v1-MsgDeleteService) | [MsgDeleteServiceResponse](#milkyway-services-v1-MsgDeleteServiceResponse) | DeleteService defines the operation for deleting an existing service that has been deactivated. |
| TransferServiceOwnership | [MsgTransferServiceOwnership](#milkyway-services-v1-MsgTransferServiceOwnership) | [MsgTransferServiceOwnershipResponse](#milkyway-services-v1-MsgTransferServiceOwnershipResponse) | TransferServiceOwnership defines the operation for transferring the ownership of a service to another account. |
| SetServiceParams | [MsgSetServiceParams](#milkyway-services-v1-MsgSetServiceParams) | [MsgSetServiceParamsResponse](#milkyway-services-v1-MsgSetServiceParamsResponse) | SetServiceParams defines the operation for setting a service&#39;s parameters. |
| UpdateParams | [MsgUpdateParams](#milkyway-services-v1-MsgUpdateParams) | [MsgUpdateParamsResponse](#milkyway-services-v1-MsgUpdateParamsResponse) | UpdateParams defines a (governance) operation for updating the module parameters. The authority defaults to the x/gov module account. |
| AccreditService | [MsgAccreditService](#milkyway-services-v1-MsgAccreditService) | [MsgAccreditServiceResponse](#milkyway-services-v1-MsgAccreditServiceResponse) | AccreditService defines a (governance) operation for accrediting a service. Since: v1.4.0 |
| RevokeServiceAccreditation | [MsgRevokeServiceAccreditation](#milkyway-services-v1-MsgRevokeServiceAccreditation) | [MsgRevokeServiceAccreditationResponse](#milkyway-services-v1-MsgRevokeServiceAccreditationResponse) | RevokeServiceAccreditation defines a (governance) operation for revoking a service&#39;s accreditation. Since: v1.4.0 |

 



<a name="milkyway_services_v1_query-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/services/v1/query.proto



<a name="milkyway-services-v1-QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="milkyway-services-v1-QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-services-v1-Params) |  |  |






<a name="milkyway-services-v1-QueryServiceParamsRequest"></a>

### QueryServiceParamsRequest
QueryServiceParamsRequest is the request type for the Query/ServiceParams RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  |  |






<a name="milkyway-services-v1-QueryServiceParamsResponse"></a>

### QueryServiceParamsResponse
QueryServiceParamsResponse is the response type for the Query/ServiceParams
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_params | [ServiceParams](#milkyway-services-v1-ServiceParams) |  |  |






<a name="milkyway-services-v1-QueryServiceRequest"></a>

### QueryServiceRequest
QueryServiceRequest is the request type for the Query/Service RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceID is the ID of the service to query |






<a name="milkyway-services-v1-QueryServiceResponse"></a>

### QueryServiceResponse
QueryServiceResponse is the response type for the Query/Service RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [Service](#milkyway-services-v1-Service) |  |  |






<a name="milkyway-services-v1-QueryServicesRequest"></a>

### QueryServicesRequest
QueryServicesRequest is the request type for the Query/Services RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  |  |






<a name="milkyway-services-v1-QueryServicesResponse"></a>

### QueryServicesResponse
QueryServicesResponse is the response type for the Query/Services RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| services | [Service](#milkyway-services-v1-Service) | repeated | Services services defines the list of actively validates services |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination response |





 

 

 


<a name="milkyway-services-v1-Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Service | [QueryServiceRequest](#milkyway-services-v1-QueryServiceRequest) | [QueryServiceResponse](#milkyway-services-v1-QueryServiceResponse) | Service defines a gRPC query method that returns the service by the given service id. |
| Services | [QueryServicesRequest](#milkyway-services-v1-QueryServicesRequest) | [QueryServicesResponse](#milkyway-services-v1-QueryServicesResponse) | Services defines a gRPC query method that returns the actively validates services currently registered in the module. |
| ServiceParams | [QueryServiceParamsRequest](#milkyway-services-v1-QueryServiceParamsRequest) | [QueryServiceParamsResponse](#milkyway-services-v1-QueryServiceParamsResponse) | ServiceParams defines a gRPC query method that returns the parameters of service. |
| Params | [QueryParamsRequest](#milkyway-services-v1-QueryParamsRequest) | [QueryParamsResponse](#milkyway-services-v1-QueryParamsResponse) | Params defines a gRPC query method that returns the parameters of the module. |

 



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

