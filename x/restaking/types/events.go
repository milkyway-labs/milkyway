package types

const (
	EventTypeCompleteUnbonding       = "complete_unbonding"
	EventTypeJoinService             = "join_service"
	EventTypeLeaveService            = "leave_service"
	EventTypeAllowOperator           = "allow_operator"
	EventTypeRemoveAllowedOperator   = "remove_allowed_operator"
	EventTypeBorrowPoolSecurity      = "borrow_pool_security"
	EventTypeCeasePoolSecurityBorrow = "cease_pool_security_borrow"
	EventTypeDelegatePool            = "delegate_pool"
	EventTypeDelegateOperator        = "delegate_operator"
	EventTypeDelegateService         = "delegate_service"
	EventTypeUnbondPool              = "unbond_pool"
	EventTypeUnbondOperator          = "unbond_operator"
	EventTypeUnbondService           = "unbond_service"
	EventTypeSetUserPreferences      = "set_user_preferences"

	AttributeKeyDelegator            = "delegator"
	AttributeKeyNewShares            = "new_shares"
	AttributeKeyCompletionTime       = "completion_time"
	AttributeUnbondingDelegationType = "unbonding_delegation"
	AttributeTargetID                = "target_id"
	AttributeKeyUser                 = "user"
)
