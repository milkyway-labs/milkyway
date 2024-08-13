package types

const (
	EventTypeUpdateOperatorParams = "update_operator_params"
	EventTypeUpdateServiceParams  = "update_service_params"
	EventTypeDelegatePool         = "delegate_pool"
	EventTypeDelegateOperator     = "delegate_operator"
	EventTypeDelegateService      = "delegate_service"
	EventTypeUnbondPool           = "unbond_pool"

	AttributeKeyCommissionRate         = "commission_rate"
	AttributeKeyJoinedServiceIDs       = "joined_services_ids"
	AttributeKeySlashFraction          = "slash_fraction"
	AttributeKeyWhitelistedPoolIDs     = "whitelisted_pools_ids"
	AttributeKeyWhitelistedOperatorIDs = "whitelisted_operators_ids"
	AttributeKeyDelegator              = "delegator"
	AttributeKeyOperatorID             = "operator_id"
	AttributeKeyServiceID              = "service_id"
	AttributeKeyNewShares              = "new_shares"
	AttributeKeyCompletionTime         = "completion_time"
)
