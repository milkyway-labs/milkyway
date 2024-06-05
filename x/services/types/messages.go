package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterService{}
	_ sdk.Msg = &MsgUpdateService{}
)

func NewMsgRegisterService(name, description, website, pictureURL, sender string) *MsgRegisterService {
	return &MsgRegisterService{
		Name:        name,
		Description: description,
		Website:     website,
		PictureURL:  pictureURL,
		Sender:      sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgRegisterService) ValidateBasic() error {
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
func (msg *MsgRegisterService) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(msg)
}

// GetSigners implements sdk.Msg
func (msg *MsgRegisterService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

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
