package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgCreateService{}
	_ sdk.Msg = &MsgUpdateService{}
)

// NewMsgCreateService creates a new MsgCreateService instance
func NewMsgCreateService(
	name string,
	description string,
	website string,
	pictureURL string,
	feeAmount sdk.Coins,
	sender string,
) *MsgCreateService {
	return &MsgCreateService{
		Name:        name,
		Description: description,
		Website:     website,
		PictureURL:  pictureURL,
		FeeAmount:   feeAmount,
		Sender:      sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgCreateService) ValidateBasic() error {
	if strings.TrimSpace(msg.Name) == "" || msg.Name == DoNotModify {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid name: %s", msg.Name)
	}

	if msg.Description == DoNotModify {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid description")
	}

	if msg.Website == DoNotModify {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid website")
	}

	if msg.PictureURL == DoNotModify {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid picture URL")
	}

	if err := msg.FeeAmount.Validate(); err != nil {
		return err
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgCreateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgCreateService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgUpdateService creates a new MsgUpdateService instance
func NewMsgUpdateService(serviceID uint32, name, description, website, pictureURL, sender string) *MsgUpdateService {
	return &MsgUpdateService{
		Sender:      sender,
		ServiceID:   serviceID,
		Name:        name,
		Description: description,
		Website:     website,
		PictureURL:  pictureURL,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateService) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service id: %d", msg.ServiceID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgActivateService creates a new MsgActivateService instance
func NewMsgActivateService(serviceID uint32, sender string) *MsgActivateService {
	return &MsgActivateService{
		ServiceID: serviceID,
		Sender:    sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgActivateService) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service id: %d", msg.ServiceID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgActivateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgActivateService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgDeactivateService creates a new MsgDeactivateService instance
func NewMsgDeactivateService(serviceID uint32, sender string) *MsgDeactivateService {
	return &MsgDeactivateService{
		ServiceID: serviceID,
		Sender:    sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgDeactivateService) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service id: %d", msg.ServiceID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgDeactivateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgDeactivateService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgTransferServiceOwnership creates a new MsgTransferServiceOwnership instance
func NewMsgTransferServiceOwnership(serviceID uint32, newAdmin, sender string) *MsgTransferServiceOwnership {
	return &MsgTransferServiceOwnership{
		ServiceID: serviceID,
		NewAdmin:  newAdmin,
		Sender:    sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgTransferServiceOwnership) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	_, err := sdk.AccAddressFromBech32(msg.NewAdmin)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid new admin address")
	}

	_, err = sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgTransferServiceOwnership) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgTransferServiceOwnership) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgDeleteService creates a new MsgDeleteService instance.
func NewMsgDeleteService(serviceID uint32, sender string) *MsgDeleteService {
	return &MsgDeleteService{
		ServiceID: serviceID,
		Sender:    sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgDeleteService) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service id: %d", msg.ServiceID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgDeleteService) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgDeleteService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgUpdateParams(params Params, authority string) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateParams) ValidateBasic() error {
	err := msg.Params.Validate()
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid params: %s", err.Error())
	}

	_, err = sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgAccreditService creates a new MsgAccreditService instance
func NewMsgAccreditService(serviceID uint32, authority string) *MsgAccreditService {
	return &MsgAccreditService{
		ServiceID: serviceID,
		Authority: authority,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgAccreditService) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgAccreditService) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgAccreditService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgRevokeServiceAccreditation creates a new MsgRevokeServiceAccreditation instance
func NewMsgRevokeServiceAccreditation(serviceID uint32, authority string) *MsgRevokeServiceAccreditation {
	return &MsgRevokeServiceAccreditation{
		ServiceID: serviceID,
		Authority: authority,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgRevokeServiceAccreditation) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgRevokeServiceAccreditation) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgRevokeServiceAccreditation) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgSetServiceParams creates a new MsgSetServiceParams instance.
func NewMsgSetServiceParams(serviceID uint32, params ServiceParams, sender string) *MsgSetServiceParams {
	return &MsgSetServiceParams{
		ServiceID:     serviceID,
		ServiceParams: params,
		Sender:        sender,
	}
}

func (msg *MsgSetServiceParams) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service id: %d", msg.ServiceID)
	}

	err := msg.ServiceParams.Validate()
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid params: %s", err.Error())
	}

	_, err = sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}
