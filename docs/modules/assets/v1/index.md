# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [milkyway/assets/v1/models.proto](#milkyway_assets_v1_models-proto)
    - [Asset](#milkyway-assets-v1-Asset)
  
- [milkyway/assets/v1/genesis.proto](#milkyway_assets_v1_genesis-proto)
    - [GenesisState](#milkyway-assets-v1-GenesisState)
  
- [milkyway/assets/v1/messages.proto](#milkyway_assets_v1_messages-proto)
    - [MsgDeregisterAsset](#milkyway-assets-v1-MsgDeregisterAsset)
    - [MsgDeregisterAssetResponse](#milkyway-assets-v1-MsgDeregisterAssetResponse)
    - [MsgRegisterAsset](#milkyway-assets-v1-MsgRegisterAsset)
    - [MsgRegisterAssetResponse](#milkyway-assets-v1-MsgRegisterAssetResponse)
  
    - [Msg](#milkyway-assets-v1-Msg)
  
- [milkyway/assets/v1/query.proto](#milkyway_assets_v1_query-proto)
    - [QueryAssetRequest](#milkyway-assets-v1-QueryAssetRequest)
    - [QueryAssetResponse](#milkyway-assets-v1-QueryAssetResponse)
    - [QueryAssetsRequest](#milkyway-assets-v1-QueryAssetsRequest)
    - [QueryAssetsResponse](#milkyway-assets-v1-QueryAssetsResponse)
  
    - [Query](#milkyway-assets-v1-Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="milkyway_assets_v1_models-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/assets/v1/models.proto



<a name="milkyway-assets-v1-Asset"></a>

### Asset
Asset represents an asset that can be registered on the chain.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| denom | [string](#string) |  | Denom is the denomination of the asset. |
| ticker | [string](#string) |  | Ticker is the ticker of the asset. |
| exponent | [uint32](#uint32) |  | Exponent represents power of 10 exponent that one must raise the denom to in order to equal the given ticker. 1 ticker = 10^exponent denom |





 

 

 

 



<a name="milkyway_assets_v1_genesis-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/assets/v1/genesis.proto



<a name="milkyway-assets-v1-GenesisState"></a>

### GenesisState
GenesisState defines the module&#39;s genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assets | [Asset](#milkyway-assets-v1-Asset) | repeated | Assets defines the registered assets. |





 

 

 

 



<a name="milkyway_assets_v1_messages-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/assets/v1/messages.proto



<a name="milkyway-assets-v1-MsgDeregisterAsset"></a>

### MsgDeregisterAsset
MsgDeregisterAsset defines the message structure for the DeregisterAsset
gRPC service method. It allows the authority to de-register an asset with
the token denomination.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authority | [string](#string) |  | Authority is the address that controls the module (defaults to x/gov unless overwritten). |
| denom | [string](#string) |  | Denom represents the denomination of the token associated with the asset. |






<a name="milkyway-assets-v1-MsgDeregisterAssetResponse"></a>

### MsgDeregisterAssetResponse
MsgRegisterAssetResponse is the return value of MsgDeregisterAsset.






<a name="milkyway-assets-v1-MsgRegisterAsset"></a>

### MsgRegisterAsset
MsgRegisterAsset defines the message structure for the RegisterAsset
gRPC service method. It allows the authority to register an asset.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authority | [string](#string) |  | Authority is the address that controls the module (defaults to x/gov unless overwritten). |
| asset | [Asset](#milkyway-assets-v1-Asset) |  | Asset represents the asset to be registered. |






<a name="milkyway-assets-v1-MsgRegisterAssetResponse"></a>

### MsgRegisterAssetResponse
MsgRegisterAssetResponse is the return value of MsgRegisterAsset.





 

 

 


<a name="milkyway-assets-v1-Msg"></a>

### Msg
Msg defines the assets module&#39;s gRPC message service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| RegisterAsset | [MsgRegisterAsset](#milkyway-assets-v1-MsgRegisterAsset) | [MsgRegisterAssetResponse](#milkyway-assets-v1-MsgRegisterAssetResponse) | RegisterAsset defines the operation for registering an asset. |
| DeregisterAsset | [MsgDeregisterAsset](#milkyway-assets-v1-MsgDeregisterAsset) | [MsgDeregisterAssetResponse](#milkyway-assets-v1-MsgDeregisterAssetResponse) | DeregisterAsset defines the operation for de-registering an asset with its denomination. |

 



<a name="milkyway_assets_v1_query-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/assets/v1/query.proto



<a name="milkyway-assets-v1-QueryAssetRequest"></a>

### QueryAssetRequest
QueryAssetRequest is the request type for the Query/Asset RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| denom | [string](#string) |  | Denom is the token denomination for which the ticker is to be queried. |






<a name="milkyway-assets-v1-QueryAssetResponse"></a>

### QueryAssetResponse
QueryAssetResponse is the response type for the Query/Asset RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| asset | [Asset](#milkyway-assets-v1-Asset) |  | Asset is the asset associated with the token denomination. |






<a name="milkyway-assets-v1-QueryAssetsRequest"></a>

### QueryAssetsRequest
QueryAssetsRequest is the request type for the Query/Assets RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ticker | [string](#string) |  | Ticker defines an optional filter parameter to query assets with the given ticker. |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-assets-v1-QueryAssetsResponse"></a>

### QueryAssetsResponse
QueryAssetsResponse is the response type for the Query/Assets RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assets | [Asset](#milkyway-assets-v1-Asset) | repeated | Assets represents all the assets registered. |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |





 

 

 


<a name="milkyway-assets-v1-Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Assets | [QueryAssetsRequest](#milkyway-assets-v1-QueryAssetsRequest) | [QueryAssetsResponse](#milkyway-assets-v1-QueryAssetsResponse) | Assets defined a gRPC query method that returns all assets registered. |
| Asset | [QueryAssetRequest](#milkyway-assets-v1-QueryAssetRequest) | [QueryAssetResponse](#milkyway-assets-v1-QueryAssetResponse) | Asset defines a gRPC query method that returns the asset associated with the given token denomination. |

 



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

