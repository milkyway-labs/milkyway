package types

// DONTCOVER

const (
	EventTypeRegisterOperator             = "register_operator"
	EventTypeUpdateOperator               = "update_operator"
	EventTypeStartOperatorInactivation    = "start_operator_inactivation"
	EventTypeCompleteOperatorInactivation = "complete_operator_inactivation"
	EventTypeTransferOperatorOwnership    = "transfer_operator_ownership"

	AttributeKeyOperatorID = "operator_id"
	AttributeKeyNewAdmin   = "new_admin"
)
