# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [milkyway/pools/v1/models.proto](#milkyway_pools_v1_models-proto)
    - [Pool](#milkyway-pools-v1-Pool)
  
- [milkyway/pools/v1/genesis.proto](#milkyway_pools_v1_genesis-proto)
    - [GenesisState](#milkyway-pools-v1-GenesisState)
  
- [milkyway/pools/v1/query.proto](#milkyway_pools_v1_query-proto)
    - [QueryPoolByDenomRequest](#milkyway-pools-v1-QueryPoolByDenomRequest)
    - [QueryPoolByIdRequest](#milkyway-pools-v1-QueryPoolByIdRequest)
    - [QueryPoolResponse](#milkyway-pools-v1-QueryPoolResponse)
    - [QueryPoolsRequest](#milkyway-pools-v1-QueryPoolsRequest)
    - [QueryPoolsResponse](#milkyway-pools-v1-QueryPoolsResponse)
  
    - [Query](#milkyway-pools-v1-Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="milkyway_pools_v1_models-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/pools/v1/models.proto



<a name="milkyway-pools-v1-Pool"></a>

### Pool
Pool defines the structure of a restaking pool


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [uint32](#uint32) |  | ID is the auto-generated unique identifier for the pool |
| denom | [string](#string) |  | Denom represents the denomination of the tokens that are staked in the pool |
| address | [string](#string) |  | Address represents the address of the account that is associated with this pool. This will be used to store tokens that users delegate to this pool. |
| tokens | [string](#string) |  | Tokens define the delegated tokens. |
| delegator_shares | [string](#string) |  | DelegatorShares defines total shares issued to a pool&#39;s delegators. |





 

 

 

 



<a name="milkyway_pools_v1_genesis-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/pools/v1/genesis.proto



<a name="milkyway-pools-v1-GenesisState"></a>

### GenesisState
GenesisState defines the pools module&#39;s genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| next_pool_id | [uint32](#uint32) |  | NextPoolID represents the id to be used when creating the next pool. |
| pools | [Pool](#milkyway-pools-v1-Pool) | repeated | Pools defines the list of pools. |





 

 

 

 



<a name="milkyway_pools_v1_query-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/pools/v1/query.proto



<a name="milkyway-pools-v1-QueryPoolByDenomRequest"></a>

### QueryPoolByDenomRequest
QueryPoolByDenomRequest is the request type for the Query/PollByDenom RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| denom | [string](#string) |  | Denom is the denom for which the pool is to be queried |






<a name="milkyway-pools-v1-QueryPoolByIdRequest"></a>

### QueryPoolByIdRequest
QueryPoolByIdRequest is the request type for the Query/PoolById RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_id | [uint32](#uint32) |  | PoolID is the ID of the pool to query |






<a name="milkyway-pools-v1-QueryPoolResponse"></a>

### QueryPoolResponse
QueryPoolResponse is the response type for the Query/PoolById and
Query/PoolByDenom RPC methods.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool | [Pool](#milkyway-pools-v1-Pool) |  | Pool is the queried pool |






<a name="milkyway-pools-v1-QueryPoolsRequest"></a>

### QueryPoolsRequest
QueryPoolsRequest is the request type for the Query/Pools RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  |  |






<a name="milkyway-pools-v1-QueryPoolsResponse"></a>

### QueryPoolsResponse
QueryPoolsResponse is the response type for the Query/Pools RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pools | [Pool](#milkyway-pools-v1-Pool) | repeated | Pools is the list of pool |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination response |





 

 

 


<a name="milkyway-pools-v1-Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| PoolByID | [QueryPoolByIdRequest](#milkyway-pools-v1-QueryPoolByIdRequest) | [QueryPoolResponse](#milkyway-pools-v1-QueryPoolResponse) | PoolByID defines a gRPC query method that returns the pool by the given ID. |
| PoolByDenom | [QueryPoolByDenomRequest](#milkyway-pools-v1-QueryPoolByDenomRequest) | [QueryPoolResponse](#milkyway-pools-v1-QueryPoolResponse) | PoolByDenom defines a gRPC query method that returns the pool by the given denom. |
| Pools | [QueryPoolsRequest](#milkyway-pools-v1-QueryPoolsRequest) | [QueryPoolsResponse](#milkyway-pools-v1-QueryPoolsResponse) | Pools defines a gRPC query method that returns all pools. |

 



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

