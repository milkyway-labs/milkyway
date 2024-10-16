package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgJoinService creates a new MsgJoinService instance
func NewMsgJoinService(operatorID uint32, serviceID uint32, sender string) *MsgJoinService {
	return &MsgJoinService{
		OperatorID: operatorID,
		ServiceID:  serviceID,
		Sender:     sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgJoinService) ValidateBasic() error {
	if msg.OperatorID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator ID: %d", msg.OperatorID)
	}

	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgUpdateServiceParams creates a new MsgUpdateServiceParams instance
func NewMsgUpdateServiceParams(operatorID uint32, params ServiceParams, sender string) *MsgUpdateServiceParams {
	return &MsgUpdateServiceParams{
		ServiceID: operatorID,
		Params:    params,
		Sender:    sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateServiceParams) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	err := msg.Params.Validate()
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

// NewMsgDelegateService creates a new MsgDelegateService instance
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

// NewMsgUpdateParams creates a new MsgUpdateParams instance
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

// --------------------------------------------------------------------------------------------------------------------

// NewMsgUndelegatePool creates a new MsgUndelegatePool instance
func NewMsgUndelegatePool(amount sdk.Coin, delegator string) *MsgUndelegatePool {
	return &MsgUndelegatePool{
		Amount:    amount,
		Delegator: delegator,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUndelegatePool) ValidateBasic() error {
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
func (msg *MsgUndelegatePool) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgUndelegatePool) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Delegator)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgUndelegateOperator creates a new MsgUndelegateOperator instance
func NewMsgUndelegateOperator(operatorID uint32, amount sdk.Coins, delegator string) *MsgUndelegateOperator {
	return &MsgUndelegateOperator{
		OperatorID: operatorID,
		Amount:     amount,
		Delegator:  delegator,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUndelegateOperator) ValidateBasic() error {
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
func (msg *MsgUndelegateOperator) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgUndelegateOperator) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Delegator)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgUndelegateService creates a new MsgUndelegateService instance
func NewMsgUndelegateService(serviceID uint32, amount sdk.Coins, delegator string) *MsgUndelegateService {
	return &MsgUndelegateService{
		ServiceID: serviceID,
		Amount:    amount,
		Delegator: delegator,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUndelegateService) ValidateBasic() error {
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
func (msg *MsgUndelegateService) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgUndelegateService) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Delegator)
	return []sdk.AccAddress{addr}
}
