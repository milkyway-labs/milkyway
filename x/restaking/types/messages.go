package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgJoinRestakingPool creates a new MsgJoinRestakePool instance
func NewMsgJoinRestakingPool(amount sdk.Coin, delegator string) *MsgJoinRestakingPool {
	return &MsgJoinRestakingPool{
		Amount:    amount,
		Delegator: delegator,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgJoinRestakingPool) ValidateBasic() error {
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
func (msg *MsgJoinRestakingPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgJoinRestakingPool) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Delegator)
	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgDelegateOperator creates a new MsgDelegateOperator instance
func NewMsgDelegateOperator(operatorID uint32, amount sdk.Coin, delegator string) *MsgDelegateOperator {
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

func NewMsgDelegateService(serviceID uint32, amount sdk.Coin, delegator string) *MsgDelegateService {
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
