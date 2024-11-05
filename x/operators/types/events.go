package types

// DONTCOVER

const (
	EventTypeRegisterOperator             = "register_operator"
	EventTypeUpdateOperator               = "update_operator"
	EventTypeStartOperatorInactivation    = "start_operator_inactivation"
	EventTypeCompleteOperatorInactivation = "complete_operator_inactivation"
	EventTypeReactivateOperator           = "reactivate_operator"
	EventTypeTransferOperatorOwnership    = "transfer_operator_ownership"
	EventTypeSetOperatorParams            = "set_operator_params"
	EventTypeDeleteOperator               = "delete_operator"

	AttributeKeyOperatorID = "operator_id"
	AttributeKeyNewAdmin   = "new_admin"
)
