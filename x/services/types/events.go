package types

// DONTCOVER

const (
	EventTypeCreateService            = "create_service"
	EventTypeUpdateService            = "update_service"
	EventTypeActivateService          = "activate_service"
	EventTypeDeactivateService        = "deactivate_service"
	EventTypeTransferServiceOwnership = "transfer_service_ownership"
	EventTypeSetServiceParams         = "set_service_params"

	AttributeKeyServiceID = "service_id"
	AttributeKeyNewAdmin  = "new_admin"
)
