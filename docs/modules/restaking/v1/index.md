# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [milkyway/restaking/v1/models.proto](#milkyway_restaking_v1_models-proto)
    - [DTData](#milkyway-restaking-v1-DTData)
    - [DTDataList](#milkyway-restaking-v1-DTDataList)
    - [Delegation](#milkyway-restaking-v1-Delegation)
    - [DelegationResponse](#milkyway-restaking-v1-DelegationResponse)
    - [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation)
    - [UnbondingDelegationEntry](#milkyway-restaking-v1-UnbondingDelegationEntry)
    - [UserPreferences](#milkyway-restaking-v1-UserPreferences)
  
    - [DelegationType](#milkyway-restaking-v1-DelegationType)
  
- [milkyway/restaking/v1/params.proto](#milkyway_restaking_v1_params-proto)
    - [Params](#milkyway-restaking-v1-Params)
  
- [milkyway/restaking/v1/genesis.proto](#milkyway_restaking_v1_genesis-proto)
    - [GenesisState](#milkyway-restaking-v1-GenesisState)
    - [OperatorJoinedServices](#milkyway-restaking-v1-OperatorJoinedServices)
    - [ServiceAllowedOperators](#milkyway-restaking-v1-ServiceAllowedOperators)
    - [ServiceSecuringPools](#milkyway-restaking-v1-ServiceSecuringPools)
    - [UserPreferencesEntry](#milkyway-restaking-v1-UserPreferencesEntry)
  
- [milkyway/restaking/v1/messages.proto](#milkyway_restaking_v1_messages-proto)
    - [MsgAddOperatorToAllowList](#milkyway-restaking-v1-MsgAddOperatorToAllowList)
    - [MsgAddOperatorToAllowListResponse](#milkyway-restaking-v1-MsgAddOperatorToAllowListResponse)
    - [MsgBorrowPoolSecurity](#milkyway-restaking-v1-MsgBorrowPoolSecurity)
    - [MsgBorrowPoolSecurityResponse](#milkyway-restaking-v1-MsgBorrowPoolSecurityResponse)
    - [MsgCeasePoolSecurityBorrow](#milkyway-restaking-v1-MsgCeasePoolSecurityBorrow)
    - [MsgCeasePoolSecurityBorrowResponse](#milkyway-restaking-v1-MsgCeasePoolSecurityBorrowResponse)
    - [MsgDelegateOperator](#milkyway-restaking-v1-MsgDelegateOperator)
    - [MsgDelegateOperatorResponse](#milkyway-restaking-v1-MsgDelegateOperatorResponse)
    - [MsgDelegatePool](#milkyway-restaking-v1-MsgDelegatePool)
    - [MsgDelegatePoolResponse](#milkyway-restaking-v1-MsgDelegatePoolResponse)
    - [MsgDelegateService](#milkyway-restaking-v1-MsgDelegateService)
    - [MsgDelegateServiceResponse](#milkyway-restaking-v1-MsgDelegateServiceResponse)
    - [MsgJoinService](#milkyway-restaking-v1-MsgJoinService)
    - [MsgJoinServiceResponse](#milkyway-restaking-v1-MsgJoinServiceResponse)
    - [MsgLeaveService](#milkyway-restaking-v1-MsgLeaveService)
    - [MsgLeaveServiceResponse](#milkyway-restaking-v1-MsgLeaveServiceResponse)
    - [MsgRemoveOperatorFromAllowlist](#milkyway-restaking-v1-MsgRemoveOperatorFromAllowlist)
    - [MsgRemoveOperatorFromAllowlistResponse](#milkyway-restaking-v1-MsgRemoveOperatorFromAllowlistResponse)
    - [MsgSetUserPreferences](#milkyway-restaking-v1-MsgSetUserPreferences)
    - [MsgSetUserPreferencesResponse](#milkyway-restaking-v1-MsgSetUserPreferencesResponse)
    - [MsgUndelegateOperator](#milkyway-restaking-v1-MsgUndelegateOperator)
    - [MsgUndelegatePool](#milkyway-restaking-v1-MsgUndelegatePool)
    - [MsgUndelegateResponse](#milkyway-restaking-v1-MsgUndelegateResponse)
    - [MsgUndelegateService](#milkyway-restaking-v1-MsgUndelegateService)
    - [MsgUpdateParams](#milkyway-restaking-v1-MsgUpdateParams)
    - [MsgUpdateParamsResponse](#milkyway-restaking-v1-MsgUpdateParamsResponse)
  
    - [Msg](#milkyway-restaking-v1-Msg)
  
- [milkyway/restaking/v1/query.proto](#milkyway_restaking_v1_query-proto)
    - [QueryDelegatorOperatorDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorOperatorDelegationsRequest)
    - [QueryDelegatorOperatorDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorOperatorDelegationsResponse)
    - [QueryDelegatorOperatorRequest](#milkyway-restaking-v1-QueryDelegatorOperatorRequest)
    - [QueryDelegatorOperatorResponse](#milkyway-restaking-v1-QueryDelegatorOperatorResponse)
    - [QueryDelegatorOperatorUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorOperatorUnbondingDelegationsRequest)
    - [QueryDelegatorOperatorUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorOperatorUnbondingDelegationsResponse)
    - [QueryDelegatorOperatorsRequest](#milkyway-restaking-v1-QueryDelegatorOperatorsRequest)
    - [QueryDelegatorOperatorsResponse](#milkyway-restaking-v1-QueryDelegatorOperatorsResponse)
    - [QueryDelegatorPoolDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorPoolDelegationsRequest)
    - [QueryDelegatorPoolDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorPoolDelegationsResponse)
    - [QueryDelegatorPoolRequest](#milkyway-restaking-v1-QueryDelegatorPoolRequest)
    - [QueryDelegatorPoolResponse](#milkyway-restaking-v1-QueryDelegatorPoolResponse)
    - [QueryDelegatorPoolUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorPoolUnbondingDelegationsRequest)
    - [QueryDelegatorPoolUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorPoolUnbondingDelegationsResponse)
    - [QueryDelegatorPoolsRequest](#milkyway-restaking-v1-QueryDelegatorPoolsRequest)
    - [QueryDelegatorPoolsResponse](#milkyway-restaking-v1-QueryDelegatorPoolsResponse)
    - [QueryDelegatorServiceDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorServiceDelegationsRequest)
    - [QueryDelegatorServiceDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorServiceDelegationsResponse)
    - [QueryDelegatorServiceRequest](#milkyway-restaking-v1-QueryDelegatorServiceRequest)
    - [QueryDelegatorServiceResponse](#milkyway-restaking-v1-QueryDelegatorServiceResponse)
    - [QueryDelegatorServiceUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorServiceUnbondingDelegationsRequest)
    - [QueryDelegatorServiceUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorServiceUnbondingDelegationsResponse)
    - [QueryDelegatorServicesRequest](#milkyway-restaking-v1-QueryDelegatorServicesRequest)
    - [QueryDelegatorServicesResponse](#milkyway-restaking-v1-QueryDelegatorServicesResponse)
    - [QueryOperatorDelegationRequest](#milkyway-restaking-v1-QueryOperatorDelegationRequest)
    - [QueryOperatorDelegationResponse](#milkyway-restaking-v1-QueryOperatorDelegationResponse)
    - [QueryOperatorDelegationsRequest](#milkyway-restaking-v1-QueryOperatorDelegationsRequest)
    - [QueryOperatorDelegationsResponse](#milkyway-restaking-v1-QueryOperatorDelegationsResponse)
    - [QueryOperatorJoinedServicesRequest](#milkyway-restaking-v1-QueryOperatorJoinedServicesRequest)
    - [QueryOperatorJoinedServicesResponse](#milkyway-restaking-v1-QueryOperatorJoinedServicesResponse)
    - [QueryOperatorUnbondingDelegationRequest](#milkyway-restaking-v1-QueryOperatorUnbondingDelegationRequest)
    - [QueryOperatorUnbondingDelegationResponse](#milkyway-restaking-v1-QueryOperatorUnbondingDelegationResponse)
    - [QueryOperatorUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryOperatorUnbondingDelegationsRequest)
    - [QueryOperatorUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryOperatorUnbondingDelegationsResponse)
    - [QueryParamsRequest](#milkyway-restaking-v1-QueryParamsRequest)
    - [QueryParamsResponse](#milkyway-restaking-v1-QueryParamsResponse)
    - [QueryPoolDelegationRequest](#milkyway-restaking-v1-QueryPoolDelegationRequest)
    - [QueryPoolDelegationResponse](#milkyway-restaking-v1-QueryPoolDelegationResponse)
    - [QueryPoolDelegationsRequest](#milkyway-restaking-v1-QueryPoolDelegationsRequest)
    - [QueryPoolDelegationsResponse](#milkyway-restaking-v1-QueryPoolDelegationsResponse)
    - [QueryPoolUnbondingDelegationRequest](#milkyway-restaking-v1-QueryPoolUnbondingDelegationRequest)
    - [QueryPoolUnbondingDelegationResponse](#milkyway-restaking-v1-QueryPoolUnbondingDelegationResponse)
    - [QueryPoolUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryPoolUnbondingDelegationsRequest)
    - [QueryPoolUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryPoolUnbondingDelegationsResponse)
    - [QueryServiceAllowedOperatorsRequest](#milkyway-restaking-v1-QueryServiceAllowedOperatorsRequest)
    - [QueryServiceAllowedOperatorsResponse](#milkyway-restaking-v1-QueryServiceAllowedOperatorsResponse)
    - [QueryServiceDelegationRequest](#milkyway-restaking-v1-QueryServiceDelegationRequest)
    - [QueryServiceDelegationResponse](#milkyway-restaking-v1-QueryServiceDelegationResponse)
    - [QueryServiceDelegationsRequest](#milkyway-restaking-v1-QueryServiceDelegationsRequest)
    - [QueryServiceDelegationsResponse](#milkyway-restaking-v1-QueryServiceDelegationsResponse)
    - [QueryServiceOperatorsRequest](#milkyway-restaking-v1-QueryServiceOperatorsRequest)
    - [QueryServiceOperatorsResponse](#milkyway-restaking-v1-QueryServiceOperatorsResponse)
    - [QueryServiceSecuringPoolsRequest](#milkyway-restaking-v1-QueryServiceSecuringPoolsRequest)
    - [QueryServiceSecuringPoolsResponse](#milkyway-restaking-v1-QueryServiceSecuringPoolsResponse)
    - [QueryServiceUnbondingDelegationRequest](#milkyway-restaking-v1-QueryServiceUnbondingDelegationRequest)
    - [QueryServiceUnbondingDelegationResponse](#milkyway-restaking-v1-QueryServiceUnbondingDelegationResponse)
    - [QueryServiceUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryServiceUnbondingDelegationsRequest)
    - [QueryServiceUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryServiceUnbondingDelegationsResponse)
    - [QueryUserPreferencesRequest](#milkyway-restaking-v1-QueryUserPreferencesRequest)
    - [QueryUserPreferencesResponse](#milkyway-restaking-v1-QueryUserPreferencesResponse)
  
    - [Query](#milkyway-restaking-v1-Query)
  
- [Scalar Value Types](#scalar-value-types)



<a name="milkyway_restaking_v1_models-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/restaking/v1/models.proto



<a name="milkyway-restaking-v1-DTData"></a>

### DTData
DTData is a struct that contains the basic information about an unbonding
delegation. It is intended to be used as a marshalable pointer. For example,
a DTData can be used to construct the key to getting an UnbondingDelegation
from state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegation_type | [DelegationType](#milkyway-restaking-v1-DelegationType) |  |  |
| delegator_address | [string](#string) |  |  |
| target_id | [uint32](#uint32) |  |  |






<a name="milkyway-restaking-v1-DTDataList"></a>

### DTDataList
DTDataList defines an array of DTData objects.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [DTData](#milkyway-restaking-v1-DTData) | repeated |  |






<a name="milkyway-restaking-v1-Delegation"></a>

### Delegation
Delegation represents the bond with tokens held by an account with a
given target.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [DelegationType](#milkyway-restaking-v1-DelegationType) |  | Type is the type of delegation. |
| user_address | [string](#string) |  | UserAddress is the encoded address of the user. |
| target_id | [uint32](#uint32) |  | TargetID is the id of the target to which the delegation is associated (pool, operator, service). |
| shares | [cosmos.base.v1beta1.DecCoin](#cosmos-base-v1beta1-DecCoin) | repeated | Shares define the delegation shares received. |






<a name="milkyway-restaking-v1-DelegationResponse"></a>

### DelegationResponse
DelegationResponse is equivalent to Delegation except that it
contains a balance in addition to shares which is more suitable for client
responses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation | [Delegation](#milkyway-restaking-v1-Delegation) |  |  |
| balance | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated |  |






<a name="milkyway-restaking-v1-UnbondingDelegation"></a>

### UnbondingDelegation
UnbondingDelegation stores all of a single delegator&#39;s unbonding bonds
for a single target in an time-ordered list.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [DelegationType](#milkyway-restaking-v1-DelegationType) |  | Type is the type of the unbonding delegation. |
| delegator_address | [string](#string) |  | DelegatorAddress is the encoded address of the delegator. |
| target_id | [uint32](#uint32) |  | TargetID is the ID of the target from which the tokens will be undelegated (pool, service, operator) |
| entries | [UnbondingDelegationEntry](#milkyway-restaking-v1-UnbondingDelegationEntry) | repeated | Entries are the unbonding delegation entries.

unbonding delegation entries |






<a name="milkyway-restaking-v1-UnbondingDelegationEntry"></a>

### UnbondingDelegationEntry
UnbondingDelegationEntry defines an unbonding object with relevant metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| creation_height | [int64](#int64) |  | CreationHeight is the height which the unbonding took place. |
| completion_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | CompletionTime is the unix time for unbonding completion. |
| initial_balance | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | InitialBalance defines the tokens initially scheduled to receive at completion. |
| balance | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Balance defines the tokens to receive at completion. |
| unbonding_id | [uint64](#uint64) |  | Incrementing id that uniquely identifies this entry |






<a name="milkyway-restaking-v1-UserPreferences"></a>

### UserPreferences
UserPreferences is a struct that contains a user&#39;s preferences for
restaking.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trust_non_accredited_services | [bool](#bool) |  | TrustNonAccreditedServices tells whether the user trusts all non-accredited services present on the platform. |
| trust_accredited_services | [bool](#bool) |  | TrustAccreditedServices tells whether the user trusts all accredited services present on the platform. |
| trusted_services_ids | [uint32](#uint32) | repeated | TrustedServicesIDs is a list of service IDs that the user trusts (both accredited and non-accredited). |





 


<a name="milkyway-restaking-v1-DelegationType"></a>

### DelegationType
DelegationType defines the type of delegation.

| Name | Number | Description |
| ---- | ------ | ----------- |
| DELEGATION_TYPE_UNSPECIFIED | 0 | DELEGATION_TYPE_UNSPECIFIED defines an unspecified delegation type. |
| DELEGATION_TYPE_POOL | 1 | DELEGATION_TYPE_POOL defines a delegation to a pool. |
| DELEGATION_TYPE_OPERATOR | 2 | DELEGATION_TYPE_OPERATOR defines a delegation to an operator. |
| DELEGATION_TYPE_SERVICE | 3 | DELEGATION_TYPE_SERVICE defines a delegation to a service. |


 

 

 



<a name="milkyway_restaking_v1_params-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/restaking/v1/params.proto



<a name="milkyway-restaking-v1-Params"></a>

### Params
Params defines the parameters for the module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_time | [int64](#int64) |  | UnbondingTime represents the time that will take for assets to be unbonded after the user initiates an unbonding request. This will be applied to all types of restaking: pool, operator and service restaking. |
| allowed_denoms | [string](#string) | repeated | AllowedDenoms represents the list of denoms allowed for restaking and that will be considered when computing rewards. If no denoms are set, all denoms will be considered as restakable. |
| restaking_cap | [string](#string) |  | RestakingCap represents the maximum USD value of overall restaked assets inside the chain. If set to 0, it indicates no limit, allowing any amount of assets to be restaked. |





 

 

 

 



<a name="milkyway_restaking_v1_genesis-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/restaking/v1/genesis.proto



<a name="milkyway-restaking-v1-GenesisState"></a>

### GenesisState
GenesisState defines the restaking module&#39;s genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-restaking-v1-Params) |  | Params defines the parameters of the module. |
| operators_joined_services | [OperatorJoinedServices](#milkyway-restaking-v1-OperatorJoinedServices) | repeated | OperatorsJoinedServices defines the list of the services that each operator has joined. |
| services_allowed_operators | [ServiceAllowedOperators](#milkyway-restaking-v1-ServiceAllowedOperators) | repeated | ServiceAllowedOperators defines the operators allowed to secure each service. |
| services_securing_pools | [ServiceSecuringPools](#milkyway-restaking-v1-ServiceSecuringPools) | repeated | ServicesSecuringPools defines the whitelisted pools for each service. |
| delegations | [Delegation](#milkyway-restaking-v1-Delegation) | repeated | Delegations represents the delegations. |
| unbonding_delegations | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) | repeated | UnbondingDelegations represents the unbonding delegations. |
| users_preferences | [UserPreferencesEntry](#milkyway-restaking-v1-UserPreferencesEntry) | repeated | UserPreferences represents the user preferences. |






<a name="milkyway-restaking-v1-OperatorJoinedServices"></a>

### OperatorJoinedServices
OperatorJoinedServicesRecord represents the services joined by a
individual operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorID is the ID of the operator. |
| service_ids | [uint32](#uint32) | repeated | ServiceIDs represents the list of services joined by the operator. |






<a name="milkyway-restaking-v1-ServiceAllowedOperators"></a>

### ServiceAllowedOperators
ServiceAllowedOperators represents the operators allowed to secure a
a service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceID is the ID of the service. |
| operator_ids | [uint32](#uint32) | repeated | OperatorIDs defines the allowed operator IDs. |






<a name="milkyway-restaking-v1-ServiceSecuringPools"></a>

### ServiceSecuringPools
ServiceSecuringPools represents the list pools from which a service can
borrow security


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceID is the ID of the service. |
| pool_ids | [uint32](#uint32) | repeated | PoolIDs defines the IDs of the pools from which the service can borrow security. |






<a name="milkyway-restaking-v1-UserPreferencesEntry"></a>

### UserPreferencesEntry
UserPreferencesEntry represents the user preferences.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_address | [string](#string) |  | UserAddress is the encoded address of the user. |
| preferences | [UserPreferences](#milkyway-restaking-v1-UserPreferences) |  | Preferences is the user preferences. |





 

 

 

 



<a name="milkyway_restaking_v1_messages-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/restaking/v1/messages.proto



<a name="milkyway-restaking-v1-MsgAddOperatorToAllowList"></a>

### MsgAddOperatorToAllowList
MsgAddOperatorToAllowList defines the message structure for the
AddOperatorToAllowList gRPC service method. It allows the service admin
to add an operator to the list of allowed operator to secure the service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  |  |
| service_id | [uint32](#uint32) |  |  |
| operator_id | [uint32](#uint32) |  |  |






<a name="milkyway-restaking-v1-MsgAddOperatorToAllowListResponse"></a>

### MsgAddOperatorToAllowListResponse
MsgAddOperatorToAllowListResponse is the return value of
MsgAddOperatorToAllowList.






<a name="milkyway-restaking-v1-MsgBorrowPoolSecurity"></a>

### MsgBorrowPoolSecurity
MsgBorrowPoolSecurity defines the message structure for the
BorrowPoolSecurity gRPC service method. It allows the service admin
to add a pool to the list of pools from which the service has chosen
to borrow security.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  |  |
| service_id | [uint32](#uint32) |  |  |
| pool_id | [uint32](#uint32) |  |  |






<a name="milkyway-restaking-v1-MsgBorrowPoolSecurityResponse"></a>

### MsgBorrowPoolSecurityResponse
MsgBorrowPoolSecurityResponse is the return value of MsgBorrowPoolSecurity.






<a name="milkyway-restaking-v1-MsgCeasePoolSecurityBorrow"></a>

### MsgCeasePoolSecurityBorrow
MsgCeasePoolSecurityBorrow defines the message structure for the
CeaseBorrowPoolSecurity gRPC service method. It allows the service admin
to remove a pool from the list of pools from which the service has chosen
to borrow security.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  |  |
| service_id | [uint32](#uint32) |  |  |
| pool_id | [uint32](#uint32) |  |  |






<a name="milkyway-restaking-v1-MsgCeasePoolSecurityBorrowResponse"></a>

### MsgCeasePoolSecurityBorrowResponse
MsgCeasePoolSecurityBorrowResponse is the return value of
MsgCeasePoolSecurityBorrow.






<a name="milkyway-restaking-v1-MsgDelegateOperator"></a>

### MsgDelegateOperator
MsgDelegateOperator defines the message structure for the DelegateOperator
gRPC service method. It allows a user to delegate their assets to an
operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator | [string](#string) |  | Delegator is the address of the user delegating to the operator |
| operator_id | [uint32](#uint32) |  | OperatorID is the ID of the operator to delegate to |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Amount is the amount of coins to be delegated |






<a name="milkyway-restaking-v1-MsgDelegateOperatorResponse"></a>

### MsgDelegateOperatorResponse
MsgDelegateOperatorResponse is the return value of MsgDelegateOperator.






<a name="milkyway-restaking-v1-MsgDelegatePool"></a>

### MsgDelegatePool
MsgDelegatePool defines the message structure for the DelegatePool gRPC
service method. It allows a user to put their assets into a restaking pool
that will later be used to provide cryptoeconomic security to services that
choose it.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator | [string](#string) |  | Delegator is the address of the user joining the pool |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) |  | Amount is the amount of coins to be staked |






<a name="milkyway-restaking-v1-MsgDelegatePoolResponse"></a>

### MsgDelegatePoolResponse
MsgDelegatePoolResponse defines the return value of MsgDelegatePool.






<a name="milkyway-restaking-v1-MsgDelegateService"></a>

### MsgDelegateService
MsgDelegateService defines the message structure for the DelegateService gRPC
service method. It allows a user to delegate their assets to a service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator | [string](#string) |  | Delegator is the address of the user delegating to the service |
| service_id | [uint32](#uint32) |  | ServiceID is the ID of the service to delegate to |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Amount is the amount of coins to be delegated |






<a name="milkyway-restaking-v1-MsgDelegateServiceResponse"></a>

### MsgDelegateServiceResponse
MsgDelegateServiceResponse is the return value of MsgDelegateService.






<a name="milkyway-restaking-v1-MsgJoinService"></a>

### MsgJoinService
MsgJoinService defines the message structure for the
JoinService gRPC service method. It allows the operator admin to
start securing a AVS.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  |  |
| operator_id | [uint32](#uint32) |  |  |
| service_id | [uint32](#uint32) |  |  |






<a name="milkyway-restaking-v1-MsgJoinServiceResponse"></a>

### MsgJoinServiceResponse
MsgJoinServiceResponse is the return value of MsgJoinService.






<a name="milkyway-restaking-v1-MsgLeaveService"></a>

### MsgLeaveService
MsgLeaveService defines the message structure for the
LeaveService gRPC service method. It allows the operator admin to
stop securing a AVS.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  |  |
| operator_id | [uint32](#uint32) |  |  |
| service_id | [uint32](#uint32) |  |  |






<a name="milkyway-restaking-v1-MsgLeaveServiceResponse"></a>

### MsgLeaveServiceResponse
MsgLeaveServiceResponse is the return value of MsgLeaveService.






<a name="milkyway-restaking-v1-MsgRemoveOperatorFromAllowlist"></a>

### MsgRemoveOperatorFromAllowlist
MsgRemoveOperatorFromAllowlist defines the message structure for the
RemoveOperatorFromAllowlist gRPC service method. It allows the service admin
to remove a previously added operator from the list of allowed operators
to secure the service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sender | [string](#string) |  |  |
| service_id | [uint32](#uint32) |  |  |
| operator_id | [uint32](#uint32) |  |  |






<a name="milkyway-restaking-v1-MsgRemoveOperatorFromAllowlistResponse"></a>

### MsgRemoveOperatorFromAllowlistResponse
MsgRemoveOperatorFromAllowlistResponse is the return value of
MsgRemoveOperatorFromAllowlist.






<a name="milkyway-restaking-v1-MsgSetUserPreferences"></a>

### MsgSetUserPreferences
MsgSetUserPreferences is the message structure for the SetUserPreferences
gRPC service method. It allows a user to set their preferences for the
restaking module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [string](#string) |  | User is the address of the user setting their preferences |
| preferences | [UserPreferences](#milkyway-restaking-v1-UserPreferences) |  | Preferences is the user&#39;s preferences |






<a name="milkyway-restaking-v1-MsgSetUserPreferencesResponse"></a>

### MsgSetUserPreferencesResponse
MsgSetUserPreferencesResponse is the return value of MsgSetUserPreferences.






<a name="milkyway-restaking-v1-MsgUndelegateOperator"></a>

### MsgUndelegateOperator
MsgUndelegateOperator the message structure for the UndelegateOperator gRPC
service method. It allows a user to undelegate their assets from a restaking
operator.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator | [string](#string) |  | Delegator is the address of the user undelegating from the operator. |
| operator_id | [uint32](#uint32) |  | OperatorID is the ID of the operator to undelegate from. |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Amount is the amount of coins to be undelegated. |






<a name="milkyway-restaking-v1-MsgUndelegatePool"></a>

### MsgUndelegatePool
MsgUndelegatePool the message structure for the UndelegatePool gRPC service
method. It allows a user to undelegate their assets from a restaking pool.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator | [string](#string) |  | Delegator is the address of the user undelegating from the pool. |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) |  | Amount is the amount of coins to be undelegated. |






<a name="milkyway-restaking-v1-MsgUndelegateResponse"></a>

### MsgUndelegateResponse
MsgUndelegateResponse defines the response type for the undelegation methods.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| completion_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | CompletionTime represents the time at which the undelegation will be complete |






<a name="milkyway-restaking-v1-MsgUndelegateService"></a>

### MsgUndelegateService
MsgUndelegateService the message structure for the UndelegateService gRPC
service method. It allows a user to undelegate their assets from a restaking
service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator | [string](#string) |  | Delegator is the address of the user undelegating from the service. |
| service_id | [uint32](#uint32) |  | ServiceID is the ID of the service to undelegate from. |
| amount | [cosmos.base.v1beta1.Coin](#cosmos-base-v1beta1-Coin) | repeated | Amount is the amount of coins to be undelegated. |






<a name="milkyway-restaking-v1-MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams defines the message structure for the UpdateParams gRPC
service method. It allows the authority to update the module parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authority | [string](#string) |  | Authority is the address that controls the module (defaults to x/gov unless overwritten). |
| params | [Params](#milkyway-restaking-v1-Params) |  | Params define the parameters to update.

NOTE: All parameters must be supplied. |






<a name="milkyway-restaking-v1-MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse is the return value of MsgUpdateParams.





 

 

 


<a name="milkyway-restaking-v1-Msg"></a>

### Msg
Msg defines the restaking module&#39;s gRPC message service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| JoinService | [MsgJoinService](#milkyway-restaking-v1-MsgJoinService) | [MsgJoinServiceResponse](#milkyway-restaking-v1-MsgJoinServiceResponse) | JoinService defines the operation that allows the operator admin to start securing an AVS |
| LeaveService | [MsgLeaveService](#milkyway-restaking-v1-MsgLeaveService) | [MsgLeaveServiceResponse](#milkyway-restaking-v1-MsgLeaveServiceResponse) | LeaveService defines the operation that allows the operator admin to stop securing an AVS |
| AddOperatorToAllowList | [MsgAddOperatorToAllowList](#milkyway-restaking-v1-MsgAddOperatorToAllowList) | [MsgAddOperatorToAllowListResponse](#milkyway-restaking-v1-MsgAddOperatorToAllowListResponse) | AddOperatorToAllowList defines the operation that allows the service admin to add an operator to the list of allowed operator to secure the service |
| RemoveOperatorFromAllowlist | [MsgRemoveOperatorFromAllowlist](#milkyway-restaking-v1-MsgRemoveOperatorFromAllowlist) | [MsgRemoveOperatorFromAllowlistResponse](#milkyway-restaking-v1-MsgRemoveOperatorFromAllowlistResponse) | RemoveOperatorFromAllowlist defines the operation that allows the service admin to remove a previously added operator from the list of allowed operators to secure the service |
| BorrowPoolSecurity | [MsgBorrowPoolSecurity](#milkyway-restaking-v1-MsgBorrowPoolSecurity) | [MsgBorrowPoolSecurityResponse](#milkyway-restaking-v1-MsgBorrowPoolSecurityResponse) | BorrowPoolSecurity defines the operation that allows the service admin to add a pool to the list of pools from which the service has chosen to borrow security. |
| CeasePoolSecurityBorrow | [MsgCeasePoolSecurityBorrow](#milkyway-restaking-v1-MsgCeasePoolSecurityBorrow) | [MsgCeasePoolSecurityBorrowResponse](#milkyway-restaking-v1-MsgCeasePoolSecurityBorrowResponse) | CeasePoolSecurityBorrow defines the operation that allows the service admin to remove a pool from the list of pools from which the service has chosen to borrow security. |
| DelegatePool | [MsgDelegatePool](#milkyway-restaking-v1-MsgDelegatePool) | [MsgDelegatePoolResponse](#milkyway-restaking-v1-MsgDelegatePoolResponse) | DelegatePool defines the operation that allows users to delegate any amount of an asset to a pool that can then be used to provide services with cryptoeconomic security. |
| DelegateOperator | [MsgDelegateOperator](#milkyway-restaking-v1-MsgDelegateOperator) | [MsgDelegateOperatorResponse](#milkyway-restaking-v1-MsgDelegateOperatorResponse) | DelegateOperator defines the operation that allows users to delegate their assets to a specific operator. |
| DelegateService | [MsgDelegateService](#milkyway-restaking-v1-MsgDelegateService) | [MsgDelegateServiceResponse](#milkyway-restaking-v1-MsgDelegateServiceResponse) | DelegateService defines the operation that allows users to delegate their assets to a specific service. |
| UpdateParams | [MsgUpdateParams](#milkyway-restaking-v1-MsgUpdateParams) | [MsgUpdateParamsResponse](#milkyway-restaking-v1-MsgUpdateParamsResponse) | UpdateParams defines a (governance) operation for updating the module parameters. The authority defaults to the x/gov module account. |
| UndelegatePool | [MsgUndelegatePool](#milkyway-restaking-v1-MsgUndelegatePool) | [MsgUndelegateResponse](#milkyway-restaking-v1-MsgUndelegateResponse) | UndelegatePool defines the operation that allows users to undelegate their assets from a pool. |
| UndelegateOperator | [MsgUndelegateOperator](#milkyway-restaking-v1-MsgUndelegateOperator) | [MsgUndelegateResponse](#milkyway-restaking-v1-MsgUndelegateResponse) | UndelegateOperator defines the operation that allows users to undelegate their assets from a specific operator. |
| UndelegateService | [MsgUndelegateService](#milkyway-restaking-v1-MsgUndelegateService) | [MsgUndelegateResponse](#milkyway-restaking-v1-MsgUndelegateResponse) | UndelegateService defines the operation that allows users to undelegate their assets from a specific service. |
| SetUserPreferences | [MsgSetUserPreferences](#milkyway-restaking-v1-MsgSetUserPreferences) | [MsgSetUserPreferencesResponse](#milkyway-restaking-v1-MsgSetUserPreferencesResponse) | SetUserPreferences defines the operation that allows users to set their preferences for the restaking module. |

 



<a name="milkyway_restaking_v1_query-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## milkyway/restaking/v1/query.proto



<a name="milkyway-restaking-v1-QueryDelegatorOperatorDelegationsRequest"></a>

### QueryDelegatorOperatorDelegationsRequest
QueryDelegatorOperatorDelegationsRequest is request type for the
Query/DelegatorOperatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryDelegatorOperatorDelegationsResponse"></a>

### QueryDelegatorOperatorDelegationsResponse
QueryDelegatorOperatorDelegationsResponse is response type for the
Query/DelegatorOperatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegations | [DelegationResponse](#milkyway-restaking-v1-DelegationResponse) | repeated | Delegations is the list of delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryDelegatorOperatorRequest"></a>

### QueryDelegatorOperatorRequest
QueryDelegatorOperatorRequest is request type for the Query/DelegatorOperator
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| operator_id | [uint32](#uint32) |  | OperatorId is the ID of the operator to query |






<a name="milkyway-restaking-v1-QueryDelegatorOperatorResponse"></a>

### QueryDelegatorOperatorResponse
QueryDelegatorOperatorResponse is response type for the
Query/DelegatorOperator RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator | [milkyway.operators.v1.Operator](#milkyway-operators-v1-Operator) |  | Operator is the operator |






<a name="milkyway-restaking-v1-QueryDelegatorOperatorUnbondingDelegationsRequest"></a>

### QueryDelegatorOperatorUnbondingDelegationsRequest
QueryDelegatorOperatorUnbondingDelegationsRequest is request type for the
Query/DelegatorOperatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryDelegatorOperatorUnbondingDelegationsResponse"></a>

### QueryDelegatorOperatorUnbondingDelegationsResponse
QueryDelegatorOperatorUnbondingDelegationsResponse is response type for the
Query/DelegatorOperatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegations | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) | repeated | UnbondingDelegations is the list of unbonding delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryDelegatorOperatorsRequest"></a>

### QueryDelegatorOperatorsRequest
QueryDelegatorOperatorsRequest is request type for the
Query/DelegatorOperators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryDelegatorOperatorsResponse"></a>

### QueryDelegatorOperatorsResponse
QueryDelegatorOperatorsResponse is response type for the
Query/DelegatorOperators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operators | [milkyway.operators.v1.Operator](#milkyway-operators-v1-Operator) | repeated | Operators is the list of operators |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryDelegatorPoolDelegationsRequest"></a>

### QueryDelegatorPoolDelegationsRequest
QueryDelegatorPoolDelegationsRequest is request type for the
Query/DelegatorPoolDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryDelegatorPoolDelegationsResponse"></a>

### QueryDelegatorPoolDelegationsResponse
QueryDelegatorPoolDelegationsResponse is response type for the
Query/DelegatorPoolDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegations | [DelegationResponse](#milkyway-restaking-v1-DelegationResponse) | repeated | Delegations is the list of delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryDelegatorPoolRequest"></a>

### QueryDelegatorPoolRequest
QueryDelegatorPoolRequest is request type for the Query/DelegatorPool RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pool_id | [uint32](#uint32) |  | PoolId is the ID of the pool to query |






<a name="milkyway-restaking-v1-QueryDelegatorPoolResponse"></a>

### QueryDelegatorPoolResponse
QueryDelegatorPoolResponse is response type for the Query/DelegatorPool RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool | [milkyway.pools.v1.Pool](#milkyway-pools-v1-Pool) |  | Pool is the pool |






<a name="milkyway-restaking-v1-QueryDelegatorPoolUnbondingDelegationsRequest"></a>

### QueryDelegatorPoolUnbondingDelegationsRequest
QueryDelegatorPoolUnbondingDelegationsRequest is request type for the
Query/DelegatorPoolUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryDelegatorPoolUnbondingDelegationsResponse"></a>

### QueryDelegatorPoolUnbondingDelegationsResponse
QueryDelegatorPoolUnbondingDelegationsResponse is response type for the
Query/DelegatorPoolUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegations | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) | repeated | UnbondingDelegations is the list of unbonding delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryDelegatorPoolsRequest"></a>

### QueryDelegatorPoolsRequest
QueryDelegatorPoolsRequest is request type for the Query/DelegatorPools RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryDelegatorPoolsResponse"></a>

### QueryDelegatorPoolsResponse
QueryDelegatorPoolsResponse is response type for the Query/DelegatorPools RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pools | [milkyway.pools.v1.Pool](#milkyway-pools-v1-Pool) | repeated | Pools is the list of pools |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryDelegatorServiceDelegationsRequest"></a>

### QueryDelegatorServiceDelegationsRequest
QueryDelegatorServiceDelegationsRequest is request type for the
Query/DelegatorServiceDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryDelegatorServiceDelegationsResponse"></a>

### QueryDelegatorServiceDelegationsResponse
QueryDelegatorServiceDelegationsResponse is response type for the
Query/DelegatorServiceDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegations | [DelegationResponse](#milkyway-restaking-v1-DelegationResponse) | repeated | Delegations is the list of delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryDelegatorServiceRequest"></a>

### QueryDelegatorServiceRequest
QueryDelegatorServiceRequest is request type for the Query/DelegatorService
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| service_id | [uint32](#uint32) |  | ServiceId is the ID of the service to query |






<a name="milkyway-restaking-v1-QueryDelegatorServiceResponse"></a>

### QueryDelegatorServiceResponse
QueryDelegatorServiceResponse is response type for the Query/DelegatorService
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [milkyway.services.v1.Service](#milkyway-services-v1-Service) |  | Service is the service |






<a name="milkyway-restaking-v1-QueryDelegatorServiceUnbondingDelegationsRequest"></a>

### QueryDelegatorServiceUnbondingDelegationsRequest
QueryDelegatorServiceUnbondingDelegationsRequest is request type for the
Query/DelegatorServiceUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryDelegatorServiceUnbondingDelegationsResponse"></a>

### QueryDelegatorServiceUnbondingDelegationsResponse
QueryDelegatorServiceUnbondingDelegationsResponse is response type for the
Query/DelegatorServiceUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegations | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) | repeated | UnbondingDelegations is the list of unbonding delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryDelegatorServicesRequest"></a>

### QueryDelegatorServicesRequest
QueryDelegatorServicesRequest is request type for the Query/DelegatorServices
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryDelegatorServicesResponse"></a>

### QueryDelegatorServicesResponse
QueryDelegatorServicesResponse is response type for the
Query/DelegatorServices RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| services | [milkyway.services.v1.Service](#milkyway-services-v1-Service) | repeated | Services is the list of services |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryOperatorDelegationRequest"></a>

### QueryOperatorDelegationRequest
QueryOperatorDelegationRequest is request type for the
Query/OperatorDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorId is the ID of the operator to query |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |






<a name="milkyway-restaking-v1-QueryOperatorDelegationResponse"></a>

### QueryOperatorDelegationResponse
QueryOperatorDelegationResponse is response type for the
Query/OperatorDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation | [DelegationResponse](#milkyway-restaking-v1-DelegationResponse) |  | Delegation is the delegation |






<a name="milkyway-restaking-v1-QueryOperatorDelegationsRequest"></a>

### QueryOperatorDelegationsRequest
QueryOperatorDelegationsRequest is request type for the
Query/OperatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorId is the ID of the operator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryOperatorDelegationsResponse"></a>

### QueryOperatorDelegationsResponse
QueryOperatorDelegationsResponse is response type for the
Query/OperatorDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegations | [DelegationResponse](#milkyway-restaking-v1-DelegationResponse) | repeated | Delegations is the list of delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryOperatorJoinedServicesRequest"></a>

### QueryOperatorJoinedServicesRequest
QueryOperatorJoinedServicesRequest is request type for the
Query/OperatorJoinedServices RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorId is the ID of the operator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryOperatorJoinedServicesResponse"></a>

### QueryOperatorJoinedServicesResponse
QueryOperatorJoinedServicesResponse is response type for the
Query/OperatorJoinedServices RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_ids | [uint32](#uint32) | repeated | ServiceIds is the list of services joined by the operator. |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryOperatorUnbondingDelegationRequest"></a>

### QueryOperatorUnbondingDelegationRequest
QueryOperatorUnbondingDelegationRequest is request type for the
Query/OperatorUnbondingDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorId is the ID of the operator to query |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |






<a name="milkyway-restaking-v1-QueryOperatorUnbondingDelegationResponse"></a>

### QueryOperatorUnbondingDelegationResponse
QueryOperatorUnbondingDelegationResponse is response type for the
Query/OperatorUnbondingDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegation | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) |  | UnbondingDelegation is the unbonding delegation |






<a name="milkyway-restaking-v1-QueryOperatorUnbondingDelegationsRequest"></a>

### QueryOperatorUnbondingDelegationsRequest
QueryOperatorUnbondingDelegationsRequest is request type for the
Query/OperatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_id | [uint32](#uint32) |  | OperatorId is the ID of the operator to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryOperatorUnbondingDelegationsResponse"></a>

### QueryOperatorUnbondingDelegationsResponse
QueryOperatorUnbondingDelegationsResponse is response type for the
Query/OperatorUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegations | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) | repeated | UnbondingDelegations is the list of unbonding delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryParamsRequest"></a>

### QueryParamsRequest
QueryParamsRequest is request type for the Query/Params RPC method.






<a name="milkyway-restaking-v1-QueryParamsResponse"></a>

### QueryParamsResponse
QueryParamsResponse is response type for the Query/Params RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| params | [Params](#milkyway-restaking-v1-Params) |  | params holds all the parameters of this module. |






<a name="milkyway-restaking-v1-QueryPoolDelegationRequest"></a>

### QueryPoolDelegationRequest
QueryPoolDelegationRequest is request type for the Query/PoolDelegation RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_id | [uint32](#uint32) |  | PoolId is the ID of the pool to query |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |






<a name="milkyway-restaking-v1-QueryPoolDelegationResponse"></a>

### QueryPoolDelegationResponse
QueryPoolDelegationResponse is response type for the Query/PoolDelegation RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation | [DelegationResponse](#milkyway-restaking-v1-DelegationResponse) |  | Delegation is the delegation |






<a name="milkyway-restaking-v1-QueryPoolDelegationsRequest"></a>

### QueryPoolDelegationsRequest
QueryPoolDelegationsRequest is request type for the Query/PoolDelegations RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_id | [uint32](#uint32) |  | PoolId is the ID of the pool to query. |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryPoolDelegationsResponse"></a>

### QueryPoolDelegationsResponse
QueryPoolDelegationsResponse is response type for the Query/PoolDelegations
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegations | [DelegationResponse](#milkyway-restaking-v1-DelegationResponse) | repeated | Delegations is the list of delegations. |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryPoolUnbondingDelegationRequest"></a>

### QueryPoolUnbondingDelegationRequest
QueryPoolUnbondingDelegationRequest is request type for the
Query/PoolUnbondingDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_id | [uint32](#uint32) |  |  |
| delegator_address | [string](#string) |  |  |






<a name="milkyway-restaking-v1-QueryPoolUnbondingDelegationResponse"></a>

### QueryPoolUnbondingDelegationResponse
QueryPoolUnbondingDelegationResponse is response type for the
Query/PoolUnbondingDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegation | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) |  |  |






<a name="milkyway-restaking-v1-QueryPoolUnbondingDelegationsRequest"></a>

### QueryPoolUnbondingDelegationsRequest
QueryPoolUnbondingDelegationsRequest is request type for the
Query/PoolUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_id | [uint32](#uint32) |  |  |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  |  |






<a name="milkyway-restaking-v1-QueryPoolUnbondingDelegationsResponse"></a>

### QueryPoolUnbondingDelegationsResponse
QueryPoolUnbondingDelegationsResponse is response type for the
Query/PoolUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegations | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) | repeated |  |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  |  |






<a name="milkyway-restaking-v1-QueryServiceAllowedOperatorsRequest"></a>

### QueryServiceAllowedOperatorsRequest
QueryServiceAllowedOperatorsRequest is request type for the
Query/ServiceAllowedOperators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceId is the ID of the service to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryServiceAllowedOperatorsResponse"></a>

### QueryServiceAllowedOperatorsResponse
QueryServiceAllowedOperatorsResponse is response type for the
Query/ServiceAllowedOperators RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operator_ids | [uint32](#uint32) | repeated | OperatorIds is the list of operators allowed to validate the service |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryServiceDelegationRequest"></a>

### QueryServiceDelegationRequest
QueryServiceDelegationRequest is request type for the Query/ServiceDelegation
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceId is the ID of the service to query |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |






<a name="milkyway-restaking-v1-QueryServiceDelegationResponse"></a>

### QueryServiceDelegationResponse
QueryServiceDelegationResponse is response type for the
Query/ServiceDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegation | [DelegationResponse](#milkyway-restaking-v1-DelegationResponse) |  | Delegation is the delegation |






<a name="milkyway-restaking-v1-QueryServiceDelegationsRequest"></a>

### QueryServiceDelegationsRequest
QueryServiceDelegationsRequest is request type for the
Query/ServiceDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceId is the ID of the service to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryServiceDelegationsResponse"></a>

### QueryServiceDelegationsResponse
QueryServiceDelegationsResponse is response type for the
Query/ServiceDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| delegations | [DelegationResponse](#milkyway-restaking-v1-DelegationResponse) | repeated | Delegations is the list of delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryServiceOperatorsRequest"></a>

### QueryServiceOperatorsRequest
QueryServiceOperatorsRequest is request type for the Query/ServiceOperators
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceId is the ID of the service to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryServiceOperatorsResponse"></a>

### QueryServiceOperatorsResponse
QueryServiceOperatorsResponse is response type for the Query/ServiceOperators
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operators | [milkyway.operators.v1.Operator](#milkyway-operators-v1-Operator) | repeated | Operators is the list of operators |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryServiceSecuringPoolsRequest"></a>

### QueryServiceSecuringPoolsRequest
QueryServiceSecuringPoolsRequest is request type for the
Query/ServiceSecuringPools RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceId is the ID of the service to query. |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryServiceSecuringPoolsResponse"></a>

### QueryServiceSecuringPoolsResponse
QueryServiceSecuringPoolsResponse is response type for the
Query/ServiceSecuringPools RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pool_ids | [uint32](#uint32) | repeated | PoolIds is the list of pools from which the service is allowed to borrow security. |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryServiceUnbondingDelegationRequest"></a>

### QueryServiceUnbondingDelegationRequest
QueryServiceUnbondingDelegationRequest is request type for the
Query/ServiceUnbondingDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceId is the ID of the service to query |
| delegator_address | [string](#string) |  | DelegatorAddress is the address of the delegator to query |






<a name="milkyway-restaking-v1-QueryServiceUnbondingDelegationResponse"></a>

### QueryServiceUnbondingDelegationResponse
QueryServiceUnbondingDelegationResponse is response type for the
Query/ServiceUnbondingDelegation RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegation | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) |  | UnbondingDelegation is the unbonding delegation |






<a name="milkyway-restaking-v1-QueryServiceUnbondingDelegationsRequest"></a>

### QueryServiceUnbondingDelegationsRequest
QueryServiceUnbondingDelegationsRequest is request type for the
Query/ServiceUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [uint32](#uint32) |  | ServiceId is the ID of the service to query |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos-base-query-v1beta1-PageRequest) |  | Pagination defines an optional pagination for the request. |






<a name="milkyway-restaking-v1-QueryServiceUnbondingDelegationsResponse"></a>

### QueryServiceUnbondingDelegationsResponse
QueryServiceUnbondingDelegationsResponse is response type for the
Query/ServiceUnbondingDelegations RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| unbonding_delegations | [UnbondingDelegation](#milkyway-restaking-v1-UnbondingDelegation) | repeated | UnbondingDelegations is the list of unbonding delegations |
| pagination | [cosmos.base.query.v1beta1.PageResponse](#cosmos-base-query-v1beta1-PageResponse) |  | Pagination defines the pagination in the response. |






<a name="milkyway-restaking-v1-QueryUserPreferencesRequest"></a>

### QueryUserPreferencesRequest
QueryUserPreferences is request type for the Query/UserPreferences RPC
method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_address | [string](#string) |  | UserAddress is the address of the user to query the preferences for |






<a name="milkyway-restaking-v1-QueryUserPreferencesResponse"></a>

### QueryUserPreferencesResponse
QueryUserPreferencesResponse is response type for the Query/UserPreferences
RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| preferences | [UserPreferences](#milkyway-restaking-v1-UserPreferences) |  | Preferences is the user preferences |





 

 

 


<a name="milkyway-restaking-v1-Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| OperatorJoinedServices | [QueryOperatorJoinedServicesRequest](#milkyway-restaking-v1-QueryOperatorJoinedServicesRequest) | [QueryOperatorJoinedServicesResponse](#milkyway-restaking-v1-QueryOperatorJoinedServicesResponse) | OperatorJoinedServices queries the services that an operator has joined. |
| ServiceAllowedOperators | [QueryServiceAllowedOperatorsRequest](#milkyway-restaking-v1-QueryServiceAllowedOperatorsRequest) | [QueryServiceAllowedOperatorsResponse](#milkyway-restaking-v1-QueryServiceAllowedOperatorsResponse) | ServiceAllowedOperators queries the allowed operators for a given service. |
| ServiceSecuringPools | [QueryServiceSecuringPoolsRequest](#milkyway-restaking-v1-QueryServiceSecuringPoolsRequest) | [QueryServiceSecuringPoolsResponse](#milkyway-restaking-v1-QueryServiceSecuringPoolsResponse) | ServiceSecuringPools queries the pools that are securing a given service. |
| ServiceOperators | [QueryServiceOperatorsRequest](#milkyway-restaking-v1-QueryServiceOperatorsRequest) | [QueryServiceOperatorsResponse](#milkyway-restaking-v1-QueryServiceOperatorsResponse) | ServiceOperators queries the operators for a given service. |
| PoolDelegations | [QueryPoolDelegationsRequest](#milkyway-restaking-v1-QueryPoolDelegationsRequest) | [QueryPoolDelegationsResponse](#milkyway-restaking-v1-QueryPoolDelegationsResponse) | PoolDelegations queries the delegations info for the given pool. |
| PoolDelegation | [QueryPoolDelegationRequest](#milkyway-restaking-v1-QueryPoolDelegationRequest) | [QueryPoolDelegationResponse](#milkyway-restaking-v1-QueryPoolDelegationResponse) | PoolDelegation queries the delegation info for the given pool and delegator. |
| PoolUnbondingDelegations | [QueryPoolUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryPoolUnbondingDelegationsRequest) | [QueryPoolUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryPoolUnbondingDelegationsResponse) | PoolUnbondingDelegations queries the unbonding delegations info for the given pool. |
| PoolUnbondingDelegation | [QueryPoolUnbondingDelegationRequest](#milkyway-restaking-v1-QueryPoolUnbondingDelegationRequest) | [QueryPoolUnbondingDelegationResponse](#milkyway-restaking-v1-QueryPoolUnbondingDelegationResponse) | PoolUnbondingDelegation queries the unbonding delegation info for the given pool and delegator. |
| OperatorDelegations | [QueryOperatorDelegationsRequest](#milkyway-restaking-v1-QueryOperatorDelegationsRequest) | [QueryOperatorDelegationsResponse](#milkyway-restaking-v1-QueryOperatorDelegationsResponse) | OperatorDelegations queries the delegations info for the given operator. |
| OperatorDelegation | [QueryOperatorDelegationRequest](#milkyway-restaking-v1-QueryOperatorDelegationRequest) | [QueryOperatorDelegationResponse](#milkyway-restaking-v1-QueryOperatorDelegationResponse) | OperatorDelegation queries the delegation info for the given operator and delegator. |
| OperatorUnbondingDelegations | [QueryOperatorUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryOperatorUnbondingDelegationsRequest) | [QueryOperatorUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryOperatorUnbondingDelegationsResponse) | OperatorUnbondingDelegations queries the unbonding delegations info for the given operator. |
| OperatorUnbondingDelegation | [QueryOperatorUnbondingDelegationRequest](#milkyway-restaking-v1-QueryOperatorUnbondingDelegationRequest) | [QueryOperatorUnbondingDelegationResponse](#milkyway-restaking-v1-QueryOperatorUnbondingDelegationResponse) | OperatorUnbondingDelegation queries the unbonding delegation info for the given operator and delegator. |
| ServiceDelegations | [QueryServiceDelegationsRequest](#milkyway-restaking-v1-QueryServiceDelegationsRequest) | [QueryServiceDelegationsResponse](#milkyway-restaking-v1-QueryServiceDelegationsResponse) | ServiceDelegations queries the delegations info for the given service. |
| ServiceDelegation | [QueryServiceDelegationRequest](#milkyway-restaking-v1-QueryServiceDelegationRequest) | [QueryServiceDelegationResponse](#milkyway-restaking-v1-QueryServiceDelegationResponse) | ServiceDelegation queries the delegation info for the given service and delegator. |
| ServiceUnbondingDelegations | [QueryServiceUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryServiceUnbondingDelegationsRequest) | [QueryServiceUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryServiceUnbondingDelegationsResponse) | ServiceUnbondingDelegations queries the unbonding delegations info for the given service. |
| ServiceUnbondingDelegation | [QueryServiceUnbondingDelegationRequest](#milkyway-restaking-v1-QueryServiceUnbondingDelegationRequest) | [QueryServiceUnbondingDelegationResponse](#milkyway-restaking-v1-QueryServiceUnbondingDelegationResponse) | ServiceUnbondingDelegation queries the unbonding delegation info for the given service and delegator. |
| DelegatorPoolDelegations | [QueryDelegatorPoolDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorPoolDelegationsRequest) | [QueryDelegatorPoolDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorPoolDelegationsResponse) | DelegatorPoolDelegations queries all the pool delegations of a given delegator address. |
| DelegatorPoolUnbondingDelegations | [QueryDelegatorPoolUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorPoolUnbondingDelegationsRequest) | [QueryDelegatorPoolUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorPoolUnbondingDelegationsResponse) | DelegatorPoolUnbondingDelegations queries all the pool unbonding delegations of a given delegator address. |
| DelegatorOperatorDelegations | [QueryDelegatorOperatorDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorOperatorDelegationsRequest) | [QueryDelegatorOperatorDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorOperatorDelegationsResponse) | DelegatorOperatorDelegations queries all the operator delegations of a given delegator address. |
| DelegatorOperatorUnbondingDelegations | [QueryDelegatorOperatorUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorOperatorUnbondingDelegationsRequest) | [QueryDelegatorOperatorUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorOperatorUnbondingDelegationsResponse) | DelegatorOperatorUnbondingDelegations queries all the operator unbonding delegations of a given delegator address. |
| DelegatorServiceDelegations | [QueryDelegatorServiceDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorServiceDelegationsRequest) | [QueryDelegatorServiceDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorServiceDelegationsResponse) | DelegatorServiceDelegations queries all the service delegations of a given delegator address. |
| DelegatorServiceUnbondingDelegations | [QueryDelegatorServiceUnbondingDelegationsRequest](#milkyway-restaking-v1-QueryDelegatorServiceUnbondingDelegationsRequest) | [QueryDelegatorServiceUnbondingDelegationsResponse](#milkyway-restaking-v1-QueryDelegatorServiceUnbondingDelegationsResponse) | DelegatorServiceUnbondingDelegations queries all the service unbonding delegations of a given delegator address. |
| DelegatorPools | [QueryDelegatorPoolsRequest](#milkyway-restaking-v1-QueryDelegatorPoolsRequest) | [QueryDelegatorPoolsResponse](#milkyway-restaking-v1-QueryDelegatorPoolsResponse) | DelegatorPools queries all pools info for given delegator address. |
| DelegatorPool | [QueryDelegatorPoolRequest](#milkyway-restaking-v1-QueryDelegatorPoolRequest) | [QueryDelegatorPoolResponse](#milkyway-restaking-v1-QueryDelegatorPoolResponse) | DelegatorPool queries the pool info for given delegator and pool id. |
| DelegatorOperators | [QueryDelegatorOperatorsRequest](#milkyway-restaking-v1-QueryDelegatorOperatorsRequest) | [QueryDelegatorOperatorsResponse](#milkyway-restaking-v1-QueryDelegatorOperatorsResponse) | DelegatorOperators queries all operators info for given delegator |
| DelegatorOperator | [QueryDelegatorOperatorRequest](#milkyway-restaking-v1-QueryDelegatorOperatorRequest) | [QueryDelegatorOperatorResponse](#milkyway-restaking-v1-QueryDelegatorOperatorResponse) | DelegatorOperator queries the operator info for given delegator and operator id. |
| DelegatorServices | [QueryDelegatorServicesRequest](#milkyway-restaking-v1-QueryDelegatorServicesRequest) | [QueryDelegatorServicesResponse](#milkyway-restaking-v1-QueryDelegatorServicesResponse) | DelegatorServices queries all services info for given delegator |
| DelegatorService | [QueryDelegatorServiceRequest](#milkyway-restaking-v1-QueryDelegatorServiceRequest) | [QueryDelegatorServiceResponse](#milkyway-restaking-v1-QueryDelegatorServiceResponse) | DelegatorService queries the service info for given delegator and service id. |
| UserPreferences | [QueryUserPreferencesRequest](#milkyway-restaking-v1-QueryUserPreferencesRequest) | [QueryUserPreferencesResponse](#milkyway-restaking-v1-QueryUserPreferencesResponse) | UserPreferences queries the user preferences. |
| Params | [QueryParamsRequest](#milkyway-restaking-v1-QueryParamsRequest) | [QueryParamsResponse](#milkyway-restaking-v1-QueryParamsResponse) | Params queries the restaking parameters. |

 



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

