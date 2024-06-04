package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterAVS{}
	_ sdk.Msg = &MsgDeregisterAVS{}
)

// ValidateBasic implements sdk.Msg
func (msg *MsgRegisterAVS) ValidateBasic() error {
	if strings.TrimSpace(msg.Name) == "" {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid AVS name: %s", msg.Name)
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
func (msg *MsgDeregisterAVS) ValidateBasic() error {
	if msg.AVSID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid id: %d", msg.AVSID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgDeregisterAVS) GetSignBytes() []byte {
	return AminoCdc.MustMarshalJSON(msg)
}

// GetSigners implements sdk.Msg
func (msg *MsgDeregisterAVS) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}
