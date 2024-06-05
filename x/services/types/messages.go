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
func NewMsgCreateService(name, description, website, pictureURL, sender string) *MsgCreateService {
	return &MsgCreateService{
		Name:        name,
		Description: description,
		Website:     website,
		PictureURL:  pictureURL,
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

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgCreateService) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(msg)
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
	return AminoCdc.MustMarshalJSON(msg)
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateService) GetSigners() []sdk.AccAddress {
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
	return AminoCdc.MustMarshalJSON(msg)
}

// GetSigners implements sdk.Msg
func (msg *MsgDeactivateService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}
