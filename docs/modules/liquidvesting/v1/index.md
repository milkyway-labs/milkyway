# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [milkyway/liquidvesting/v1/models.proto](#milkyway_liquidvesting_v1_models-proto)
    - [BurnCoins](#milkyway-liquidvesting-v1-BurnCoins)
    - [BurnCoinsList](#milkyway-liquidvesting-v1-BurnCoinsList)
    - [UserInsuranceFund](#milkyway-liquidvesting-v1-UserInsuranceFund)
    - [UserInsuranceFundEntry](#milkyway-liquidvesting-v1-UserInsuranceFundEntry)
  
- [milkyway/liquidvesting/v1/params.proto](#milkyway_liquidvesting_v1_params-proto)
    - [Params](#milkyway-liquidvesting-v1-Params)
  
- [milkyway/liquidvesting/v1/genesis.proto](#milkyway_liquidvesting_v1_genesis-proto)
    - [GenesisState](#milkyway-liquidvesting-v1-GenesisState)
  
- [milkyway/liquidvesting/v1/messages.proto](#milkyway_liquidvesting_v1_messages-proto)
    - [MsgBurnLockedRepresentation](#milkyway-liquidvesting-v1-MsgBurnLockedRepresentation)
    - [MsgBurnLockedRepresentationResponse](#milkyway-liquidvesting-v1-MsgBurnLockedRepresentationResponse)
    - [MsgMintLockedRepresentation](#milkyway-liquidvesting-v1-MsgMintLockedRepresentation)
    - [MsgMintLockedRepresentationResponse](#milkyway-liquidvesting-v1-MsgMintLockedRepresentationResponse)
    - [MsgUpdateParams](#milkyway-liquidvesting-v1-MsgUpdateParams)
    - [MsgUpdateParamsResponse](#milkyway-liquidvesting-v1-MsgUpdateParamsResponse)
    - [MsgWithdrawInsuranceFund](#milkyway-liquidvesting-v1-MsgWithdrawInsuranceFund)
    - [MsgWithdrawInsuranceFundResponse](#milkyway-liquidvesting-v1-MsgWithdrawInsuranceFundResponse)
  
    - [Msg](#milkyway-liquidvesting-v1-Msg)
  
- [milkyway/liquidvesting/v1/query.proto](#milkyway_liquidvesting_v1_query-proto)
    - [QueryInsuranceFundRequest](#milkyway-liquidvesting-v1-QueryInsuranceFundRequest)
    - [QueryInsuranceFundResponse](#milkyway-liquidvesting-v1-QueryInsuranceFundResponse)
    - [QueryParamsRequest](#milkyway-liquidvesting-v1-QueryParamsRequest)
    - [QueryParamsResponse](#milkyway-liquidvesting-v1-QueryParamsResponse)
    - [QueryUserInsuranceFundRequest](#milkyway-liquidvesting-v1-QueryUserInsuranceFundRequest)
    - [QueryUserInsuranceFundResponse](#milkyway-liquidvesting-v1-QueryUserInsuranceFundResponse)
    - [QueryUserInsuranceFundsRequest](#milkyway-liquidvesting-v1-QueryUserInsuranceFundsRequest)
    - [QueryUserInsuranceFundsResponse](#milkyway-liquidvesting-v1-QueryUserInsuranceFundsResponse)
    - [QueryUserRestakableAssetsRequest](#milkyway-liquidvesting-v1-QueryUserRestakableAssetsRequest)
    - [QueryUserRestakableAssetsResponse](#milkyway-liquidvesting-v1-QueryUserRestakableAssetsResponse)
    - [UserInsuranceFundData](#milkyway-liquidvesting-v1-UserInsuranceFundData)
  
    - [Query](#milkyway-liquidvesting-v1-Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="milkyway_liquidvesting_v1_models-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/liquidvesting/v1/models.proto



<a name="milkyway-liquidvesting-v1-BurnCoins"></a>

### BurnCoins
BurnCoins is a struct that contains the information about the coins to burn
once the unbonding period of the locked representation tokens ends.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | Address of who has delegated the coins. |
| completion_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | CompletionTime is the unix time for unbonding completion. |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Amount of coins to be burned from the delegator address. |






<a name="milkyway-liquidvesting-v1-BurnCoinsList"></a>

### BurnCoinsList
BurnCoinsList represents a list of BurnCoins.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [BurnCoins](#milkyway-liquidvesting-v1-BurnCoins) | repeated |  |






<a name="milkyway-liquidvesting-v1-UserInsuranceFund"></a>

### UserInsuranceFund
UserInsuranceFund defines a user&#39;s insurance fund.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| balance | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Amount of coins deposited into the user&#39;s insurance fund. |






<a name="milkyway-liquidvesting-v1-UserInsuranceFundEntry"></a>

### UserInsuranceFundEntry
UserInsuranceFundEntry represents an entry containing the data of a user
insurance fund.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_address | [string](#string) |  | Address of who owns the insurance fund. |
| balance | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Amount of coins deposited into the user&#39;s insurance fund. |





 

 

 

 



<a name="milkyway_liquidvesting_v1_params-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/liquidvesting/v1/params.proto



<a name="milkyway-liquidvesting-v1-Params"></a>

### Params
Params defines the parameters for the module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| insurance_percentage | [string](#string) |  | This value represents the percentage that needs to be sent to the insurance fund in order to restake a certain amount of locked tokens. For example, if this value is 2%, a user must send 2 tokens to the insurance fund to restake 100 locked tokens |
| burners | [string](#string) | repeated | This value represents the list of users who are authorized to execute the MsgBurnLockedRepresentation. |
| minters | [string](#string) | repeated | This value represents the list of users who are authorized to execute the MsgMintLockedRepresentation. |
| allowed_channels | [string](#string) | repeated | List of channels from which is allowed to receive deposits to the insurance fund. |





 

 

 

 



<a name="milkyway_liquidvesting_v1_genesis-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/liquidvesting/v1/genesis.proto



<a name="milkyway-liquidvesting-v1-GenesisState"></a>

### GenesisState
GenesisState defines the liquidvesting module&#39;s genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-liquidvesting-v1-Params) |  | Params defines the parameters of the module. |
| burn_coins | [BurnCoins](#milkyway-liquidvesting-v1-BurnCoins) | repeated | BurnCoins represents the list of coins that should be burned from the users&#39; balances |
| user_insurance_funds | [UserInsuranceFundEntry](#milkyway-liquidvesting-v1-UserInsuranceFundEntry) | repeated | UserInsuranceFunds represents the users&#39; insurance fund. |





 

 

 

 



<a name="milkyway_liquidvesting_v1_messages-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/liquidvesting/v1/messages.proto



<a name="milkyway-liquidvesting-v1-MsgBurnLockedRepresentation"></a>

### MsgBurnLockedRepresentation
MsgBurnLockedRepresentation defines the message structure for the
BurnLockedRepresentation gRPC service method. It allows an authorized
account to burn a user&#39;s staked locked tokens representation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | User that want to trigger the tokens burn. |
| user | [string](#string) |  | User from which we want to burn the tokens. |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | The amount of tokens that will be burned from the user&#39;s balance. |






<a name="milkyway-liquidvesting-v1-MsgBurnLockedRepresentationResponse"></a>

### MsgBurnLockedRepresentationResponse
MsgBurnLockedRepresentationResponse is the return value of
MsgBurnLockedRepresentation.






<a name="milkyway-liquidvesting-v1-MsgMintLockedRepresentation"></a>

### MsgMintLockedRepresentation
MsgMintLockedRepresentation defines the message structure for the
MintLockedRepresentation gRPC service method. It allows an authorized
account to mint a user&#39;s staked locked tokens representation that can be
used in the liquid vesting module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | User that want to trigger the tokens mint. |
| receiver | [string](#string) |  | User that will receive the minted tokens. |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | The amount of tokens that will be minted |






<a name="milkyway-liquidvesting-v1-MsgMintLockedRepresentationResponse"></a>

### MsgMintLockedRepresentationResponse
MsgMintLockedRepresentationResponse is the return value of
MsgMintLockedRepresentation.






<a name="milkyway-liquidvesting-v1-MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams defines the message structure for the UpdateParams gRPC
service method. It allows the authority to update the module parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authority | [string](#string) |  | Authority is the address that controls the module (defaults to x/gov unless overwritten). |
| params | [Params](#milkyway-liquidvesting-v1-Params) |  | Params define the parameters to update.

NOTE: All parameters must be supplied. |






<a name="milkyway-liquidvesting-v1-MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse is the return value of MsgUpdateParams.






<a name="milkyway-liquidvesting-v1-MsgWithdrawInsuranceFund"></a>

### MsgWithdrawInsuranceFund
MsgWithdrawInsuranceFund defines the message structure for the
WithdrawInsuranceFund gRPC service method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  | User that want to withdraw the tokens. |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | The amount of tokens that will be withdrawn from the user&#39;s insurance fund. |






<a name="milkyway-liquidvesting-v1-MsgWithdrawInsuranceFundResponse"></a>

### MsgWithdrawInsuranceFundResponse
MsgWithdrawInsuranceFundResponse is the return value of MsgWithdrawInsuranceFund.





 

 

 


<a name="milkyway-liquidvesting-v1-Msg"></a>

### Msg
Msg defines the services module&#39;s gRPC message service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| MintLockedRepresentation | [MsgMintLockedRepresentation](#milkyway-liquidvesting-v1-MsgMintLockedRepresentation) | [MsgMintLockedRepresentationResponse](#milkyway-liquidvesting-v1-MsgMintLockedRepresentationResponse) | MintLockedRepresentation defines the operation to mint a user&#39;s staked locked tokens representation that can be used in the liquid vesting module. |
| BurnLockedRepresentation | [MsgBurnLockedRepresentation](#milkyway-liquidvesting-v1-MsgBurnLockedRepresentation) | [MsgBurnLockedRepresentationResponse](#milkyway-liquidvesting-v1-MsgBurnLockedRepresentationResponse) | BurnLockedRepresentation defines the operation to burn a user&#39;s staked locked tokens representation. |
| WithdrawInsuranceFund | [MsgWithdrawInsuranceFund](#milkyway-liquidvesting-v1-MsgWithdrawInsuranceFund) | [MsgWithdrawInsuranceFundResponse](#milkyway-liquidvesting-v1-MsgWithdrawInsuranceFundResponse) | WithdrawInsuranceFund defines the operation to withdraw an amount of tokens from the user&#39;s insurance fund. This can be used from the user to withdraw their funds after some of their staking representations have been burned or if the balance in the insurance fund is more than the required to cover all their staking representations. |
| UpdateParams | [MsgUpdateParams](#milkyway-liquidvesting-v1-MsgUpdateParams) | [MsgUpdateParamsResponse](#milkyway-liquidvesting-v1-MsgUpdateParamsResponse) | UpdateParams defines a (governance) operation for updating the module parameters. The authority defaults to the x/gov module account. |

 



<a name="milkyway_liquidvesting_v1_query-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/liquidvesting/v1/query.proto



<a name="milkyway-liquidvesting-v1-QueryInsuranceFundRequest"></a>

### QueryInsuranceFundRequest
QueryInsuranceFundRequest is the request type for the
Query/InsuranceFund RPC method.






<a name="milkyway-liquidvesting-v1-QueryInsuranceFundResponse"></a>

### QueryInsuranceFundResponse
QueryInsuranceFundResponse is the response type for the
Query/InsuranceFund RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | amount is the amount of tokens that are in the insurance fund. |






<a name="milkyway-liquidvesting-v1-QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is the request type for the Query/Params RPC method.






<a name="milkyway-liquidvesting-v1-QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is the response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-liquidvesting-v1-Params) |  |  |






<a name="milkyway-liquidvesting-v1-QueryUserInsuranceFundRequest"></a>

### QueryUserInsuranceFundRequest
QueryUserInsuranceFundRequest is the request type for the
Query/UserInsuranceFund RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_address | [string](#string) |  | user_address is the address of the user to query. |






<a name="milkyway-liquidvesting-v1-QueryUserInsuranceFundResponse"></a>

### QueryUserInsuranceFundResponse
QueryUserInsuranceFundResponse is the response type for the
Query/UserInsuranceFund RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| balance | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | balance is the amount of tokens that is in the user&#39;s insurance fund. |
| used | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | used is the amount of tokens is being used to cover the user&#39;s restaked assets. |






<a name="milkyway-liquidvesting-v1-QueryUserInsuranceFundsRequest"></a>

### QueryUserInsuranceFundsRequest
QueryUserInsuranceFundsRequest is the request type for the
Query/UserInsuranceFunds RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  |  |






<a name="milkyway-liquidvesting-v1-QueryUserInsuranceFundsResponse"></a>

### QueryUserInsuranceFundsResponse
QueryUserInsuranceFundsResponse is the response type for the
Query/UserInsuranceFunds RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| insurance_funds | [UserInsuranceFundData](#milkyway-liquidvesting-v1-UserInsuranceFundData) | repeated | insurance_funds is the list of users insurance funds. |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination response |






<a name="milkyway-liquidvesting-v1-QueryUserRestakableAssetsRequest"></a>

### QueryUserRestakableAssetsRequest
QueryUserRestakableAssetsRequest is the request type for the
Query/UserRestakableAssets RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_address | [string](#string) |  | user_address is the address of the user to query. |






<a name="milkyway-liquidvesting-v1-QueryUserRestakableAssetsResponse"></a>

### QueryUserRestakableAssetsResponse
QueryUserRestakableAssetsResponse is the response type for the
Query/UserRestakableAssets RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | amount is the amount of tokens that the user can restake. |






<a name="milkyway-liquidvesting-v1-UserInsuranceFundData"></a>

### UserInsuranceFundData
UserInsuranceFundData is the structure that contains the information about
a user&#39;s insurance fund.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_address | [string](#string) |  | user_address is the address of who owns the insurance fund. |
| balance | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | balance is the amount of tokens that is in the user&#39;s insurance fund. |
| used | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | used is the amount of tokens that is to cover the user&#39;s restaked assets. |





 

 

 


<a name="milkyway-liquidvesting-v1-Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| UserInsuranceFund | [QueryUserInsuranceFundRequest](#milkyway-liquidvesting-v1-QueryUserInsuranceFundRequest) | [QueryUserInsuranceFundResponse](#milkyway-liquidvesting-v1-QueryUserInsuranceFundResponse) | UserInsuranceFund defines a gRPC query method that returns the user&#39;s insurance fund balance given their address. |
| UserInsuranceFunds | [QueryUserInsuranceFundsRequest](#milkyway-liquidvesting-v1-QueryUserInsuranceFundsRequest) | [QueryUserInsuranceFundsResponse](#milkyway-liquidvesting-v1-QueryUserInsuranceFundsResponse) | UserInsuranceFunds defines a gRPC query method that returns all user&#39;s insurance fund balance. |
| UserRestakableAssets | [QueryUserRestakableAssetsRequest](#milkyway-liquidvesting-v1-QueryUserRestakableAssetsRequest) | [QueryUserRestakableAssetsResponse](#milkyway-liquidvesting-v1-QueryUserRestakableAssetsResponse) | UserRestakableAssets defines a gRPC query method that returns the amount of assets that can be restaked from the one minted by this module. |
| InsuranceFund | [QueryInsuranceFundRequest](#milkyway-liquidvesting-v1-QueryInsuranceFundRequest) | [QueryInsuranceFundResponse](#milkyway-liquidvesting-v1-QueryInsuranceFundResponse) | InsuranceFund defines a gRPC query method that returns the amount of tokens that are in the insurance fund. |
| Params | [QueryParamsRequest](#milkyway-liquidvesting-v1-QueryParamsRequest) | [QueryParamsResponse](#milkyway-liquidvesting-v1-QueryParamsResponse) | Params defines a gRPC query method that returns the parameters of the module. |

 



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

