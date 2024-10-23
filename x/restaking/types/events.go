package types

const (
	EventTypeCompleteUnbonding   = "complete_unbonding"
	EventTypeJoinService         = "join_service"
	EventTypeLeaveService        = "leave_service"
	EventTypeUpdateServiceParams = "update_service_params"
	EventTypeDelegatePool        = "delegate_pool"
	EventTypeDelegateOperator    = "delegate_operator"
	EventTypeDelegateService     = "delegate_service"
	EventTypeUnbondPool          = "unbond_pool"
	EventTypeUnbondOperator      = "unbond_operator"
	EventTypeUnbondService       = "unbond_service"

	AttributeKeyJoinedServiceID      = "joined_services_id"
	AttributeKeyDelegator            = "delegator"
	AttributeKeyNewShares            = "new_shares"
	AttributeKeyCompletionTime       = "completion_time"
	AttributeUnbondingDelegationType = "unbonding_delegation"
	AttributeTargetID                = "target_id"
)
