package types

// DONTCOVER

const (
	EventTypeRegisterOperator             = "register_operator"
	EventTypeUpdateOperator               = "update_operator"
	EventTypeStartOperatorInactivation    = "start_operator_inactivation"
	EventTypeCompleteOperatorInactivation = "complete_operator_inactivation"
	EventTypeTransferOperatorOwnership    = "transfer_operator_ownership"
	EventTypeSetOperatorParams            = "set_operator_params"

	AttributeKeyOperatorID        = "operator_id"
	AttributeKeyNewAdmin          = "new_admin"
	AttributeKeyNewCommissionRate = "new_commission_rate"
)
