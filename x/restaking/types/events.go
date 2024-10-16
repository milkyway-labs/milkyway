package types

const (
	EventTypeCompleteUnbonding    = "complete_unbonding"
	EventTypeUpdateOperatorParams = "update_operator_params"
	EventTypeUpdateServiceParams  = "update_service_params"
	EventTypeDelegatePool         = "delegate_pool"
	EventTypeDelegateOperator     = "delegate_operator"
	EventTypeDelegateService      = "delegate_service"
	EventTypeUnbondPool           = "unbond_pool"
	EventTypeUnbondOperator       = "unbond_operator"
	EventTypeUnbondService        = "unbond_service"

	AttributeKeyCommissionRate         = "commission_rate"
	AttributeKeyJoinedServiceID        = "joined_services_id"
	AttributeKeySlashFraction          = "slash_fraction"
	AttributeKeyWhitelistedPoolIDs     = "whitelisted_pools_ids"
	AttributeKeyWhitelistedOperatorIDs = "whitelisted_operators_ids"
	AttributeKeyDelegator              = "delegator"
	AttributeKeyNewShares              = "new_shares"
	AttributeKeyCompletionTime         = "completion_time"
	AttributeUnbondingDelegationType   = "unbonding_delegation"
	AttributeTargetID                  = "target_id"
)
