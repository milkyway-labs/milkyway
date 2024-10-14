package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgRegisterOperator creates a new MsgRegisterOperator instance
func NewMsgRegisterOperator(moniker string, website string, pictureURL string, sender string) *MsgRegisterOperator {
	return &MsgRegisterOperator{
		Moniker:    moniker,
		Website:    website,
		PictureURL: pictureURL,
		Sender:     sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgRegisterOperator) ValidateBasic() error {
	if strings.TrimSpace(msg.Moniker) == "" || msg.Moniker == DoNotModify {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid moniker: %s", msg.Moniker)
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
func (msg *MsgRegisterOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgRegisterOperator) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgUpdateOperator creates a new MsgUpdateOperator instance
func NewMsgUpdateOperator(operatorID uint32, moniker string, website string, pictureURL string, sender string) *MsgUpdateOperator {
	return &MsgUpdateOperator{
		OperatorID: operatorID,
		Moniker:    moniker,
		Website:    website,
		PictureURL: pictureURL,
		Sender:     sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateOperator) ValidateBasic() error {
	if msg.OperatorID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator ID: %d", msg.OperatorID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateOperator) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgDeactivateOperator creates a new MsgDeactivateOperator instance
func NewMsgDeactivateOperator(operatorID uint32, sender string) *MsgDeactivateOperator {
	return &MsgDeactivateOperator{
		OperatorID: operatorID,
		Sender:     sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgDeactivateOperator) ValidateBasic() error {
	if msg.OperatorID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator ID: %d", msg.OperatorID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgDeactivateOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgDeactivateOperator) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgTransferOperatorOwnership creates a new MsgTransferOperatorOwnership instance
func NewMsgTransferOperatorOwnership(operatorID uint32, newAdmin, sender string) *MsgTransferOperatorOwnership {
	return &MsgTransferOperatorOwnership{
		OperatorID: operatorID,
		NewAdmin:   newAdmin,
		Sender:     sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgTransferOperatorOwnership) ValidateBasic() error {
	if msg.OperatorID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator ID: %d", msg.OperatorID)
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
func (msg *MsgTransferOperatorOwnership) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgTransferOperatorOwnership) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgUpdateParams(params Params, authority string) *MsgUpdateParams {
	return &MsgUpdateParams{
		Params:    params,
		Authority: authority,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateParams) ValidateBasic() error {
	err := msg.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %w", err)
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

func NewMsgSetOperatorParams(sender string, operatorId uint32, params OperatorParams) *MsgSetOperatorParams {
	return &MsgSetOperatorParams{
		Sender:     sender,
		OperatorID: operatorId,
		Params:     params,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgSetOperatorParams) ValidateBasic() error {
	if msg.OperatorID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator ID: %d", msg.OperatorID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	err = msg.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid operator params: %w", err)
	}

	return nil
}
