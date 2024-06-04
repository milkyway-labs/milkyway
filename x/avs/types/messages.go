package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterAVS{}
	_ sdk.Msg = &MsgUpdateAVS{}
)

// ValidateBasic implements sdk.Msg
func (msg *MsgRegisterAVS) ValidateBasic() error {
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
func (msg *MsgRegisterAVS) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(msg)
}

// GetSigners implements sdk.Msg
func (msg *MsgRegisterAVS) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateAVS) ValidateBasic() error {
	if msg.AVSID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid avs id: %d", msg.AVSID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateAVS) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(msg)
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateAVS) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}
