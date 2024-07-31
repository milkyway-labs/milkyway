package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgUpdateOperatorParams creates a new MsgUpdateOperatorParams instance
func NewMsgUpdateOperatorParams(
	operatorID uint32, params OperatorParams, sender string) *MsgUpdateOperatorParams {
	return &MsgUpdateOperatorParams{
		OperatorID:     operatorID,
		OperatorParams: params,
		Sender:         sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateOperatorParams) ValidateBasic() error {
	if msg.OperatorID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator ID: %d", msg.OperatorID)
	}

	err := msg.OperatorParams.Validate()
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator params: %s", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateOperatorParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateOperatorParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgUpdateServiceParams creates a new MsgUpdateServiceParams instance
func NewMsgUpdateServiceParams(
	operatorID uint32, params ServiceParams, sender string) *MsgUpdateServiceParams {
	return &MsgUpdateServiceParams{
		ServiceID:     operatorID,
		ServiceParams: params,
		Sender:        sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateServiceParams) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	err := msg.ServiceParams.Validate()
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service params: %s", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateServiceParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgUpdateServiceParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgDelegatePool creates a new MsgDelegatePool instance
func NewMsgDelegatePool(amount sdk.Coin, delegator string) *MsgDelegatePool {
	return &MsgDelegatePool{
		Amount:    amount,
		Delegator: delegator,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgDelegatePool) ValidateBasic() error {
	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid amount")
	}

	_, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgDelegatePool) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgDelegatePool) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Delegator)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgDelegateOperator creates a new MsgDelegateOperator instance
func NewMsgDelegateOperator(operatorID uint32, amount sdk.Coins, delegator string) *MsgDelegateOperator {
	return &MsgDelegateOperator{
		OperatorID: operatorID,
		Amount:     amount,
		Delegator:  delegator,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgDelegateOperator) ValidateBasic() error {
	if msg.OperatorID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator id")
	}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid amount")
	}

	_, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgDelegateOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgDelegateOperator) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Delegator)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgDelegateService(serviceID uint32, amount sdk.Coins, delegator string) *MsgDelegateService {
	return &MsgDelegateService{
		ServiceID: serviceID,
		Amount:    amount,
		Delegator: delegator,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgDelegateService) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service id")
	}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "invalid amount")
	}

	_, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgDelegateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgDelegateService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Delegator)
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
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid params: %s", err)
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
