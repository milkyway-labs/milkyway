package types

const (
	EventTypeUpdateOperatorParams = "update_operator_params"
	EventTypeUpdateServiceParams  = "update_service_params"
	EventTypeDelegatePool         = "delegate_pool"
	EventTypeDelegateOperator     = "delegate_operator"
	EventTypeDelegateService      = "delegate_service"

	AttributeKeyCommissionRate         = "commission_rate"
	AttributeKeyJoinedServiceIDs       = "joined_service_ids"
	AttributeKeySlashFraction          = "slash_fraction"
	AttributeKeyWhitelistedPoolIDs     = "whitelisted_pool_ids"
	AttributeKeyWhitelistedOperatorIDs = "whitelisted_operator_ids"
	AttributeKeyDelegator              = "delegator"
	AttributeKeyOperatorID             = "operator_id"
	AttributeKeyServiceID              = "service_id"
	AttributeKeyNewShares              = "new_shares"
)
