package types

// DONTCOVER

const (
	EventTypeCreateService              = "create_service"
	EventTypeUpdateService              = "update_service"
	EventTypeActivateService            = "activate_service"
	EventTypeDeactivateService          = "deactivate_service"
	EventTypeDeleteService              = "delete_service"
	EventTypeTransferServiceOwnership   = "transfer_service_ownership"
	EventTypeAccreditService            = "accredit_service"
	EventTypeRevokeServiceAccreditation = "revoke_service_accreditation"
	EventTypeSetServiceParams           = "set_service_params"

	AttributeKeyServiceID = "service_id"
	AttributeKeyNewAdmin  = "new_admin"
)
