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

// NewMsgLeaveService creates a new MsgLeaveService instance
func NewMsgLeaveService(operatorID uint32, serviceID uint32, sender string) *MsgLeaveService {
	return &MsgLeaveService{
		OperatorID: operatorID,
		ServiceID:  serviceID,
		Sender:     sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgLeaveService) ValidateBasic() error {
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

// NewMsgAddOperatorToAllowList creates a new MsgAddOperatorToAllowList instance
func NewMsgAddOperatorToAllowList(serviceID uint32, operatorID uint32, sender string) *MsgAddOperatorToAllowList {
	return &MsgAddOperatorToAllowList{
		ServiceID:  serviceID,
		OperatorID: operatorID,
		Sender:     sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgAddOperatorToAllowList) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	if msg.OperatorID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator ID: %d", msg.OperatorID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgRemoveOperatorFromAllowList creates a new MsgRemoveOperatorFromAllowlist instance
func NewMsgRemoveOperatorFromAllowList(serviceID uint32, operatorID uint32, sender string) *MsgRemoveOperatorFromAllowlist {
	return &MsgRemoveOperatorFromAllowlist{
		ServiceID:  serviceID,
		OperatorID: operatorID,
		Sender:     sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgRemoveOperatorFromAllowlist) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	if msg.OperatorID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid operator ID: %d", msg.OperatorID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgBorrowPoolSecurity creates a new MsgBorrowPoolSecurity instance
func NewMsgBorrowPoolSecurity(serviceID uint32, poolID uint32, sender string) *MsgBorrowPoolSecurity {
	return &MsgBorrowPoolSecurity{
		ServiceID: serviceID,
		PoolID:    poolID,
		Sender:    sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgBorrowPoolSecurity) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	if msg.PoolID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid pool ID: %d", msg.PoolID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgCeasePoolSecurityBorrow creates a new MsgCeasePoolSecurityBorrow instance
func NewMsgCeasePoolSecurityBorrow(serviceID uint32, poolID uint32, sender string) *MsgCeasePoolSecurityBorrow {
	return &MsgCeasePoolSecurityBorrow{
		ServiceID: serviceID,
		PoolID:    poolID,
		Sender:    sender,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgCeasePoolSecurityBorrow) ValidateBasic() error {
	if msg.ServiceID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service ID: %d", msg.ServiceID)
	}

	if msg.PoolID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid pool ID: %d", msg.PoolID)
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return nil
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

// --------------------------------------------------------------------------------------------------------------------

// NewMsgSetUserPreferences creates a new MsgSetUserPreferences instance
func NewMsgSetUserPreferences(preferences UserPreferences, userAddress string) *MsgSetUserPreferences {
	return &MsgSetUserPreferences{
		Preferences: preferences,
		User:        userAddress,
	}
}

// ValidateBasic implements sdk.Msg
func (msg *MsgSetUserPreferences) ValidateBasic() error {
	err := msg.Preferences.Validate()
	if err != nil {
		return err
	}

	_, err = sdk.AccAddressFromBech32(msg.User)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid user address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (msg *MsgSetUserPreferences) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg *MsgSetUserPreferences) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.User)
	return []sdk.AccAddress{addr}
}
