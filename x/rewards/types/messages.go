package types

import (
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	restakingtypes "github.com/milkyway-labs/milkyway/v6/x/restaking/types"
)

// NewMsgCreateRewardsPlan creates a new MsgCreateRewardsPlan instance
func NewMsgCreateRewardsPlan(
	serviceID uint32,
	description string,
	amt sdk.Coins,
	startTime,
	endTime time.Time,
	poolsDistribution Distribution,
	operatorsDistribution Distribution,
	usersDistribution UsersDistribution,
	feeAmount sdk.Coins,
	sender string,
) *MsgCreateRewardsPlan {
	return &MsgCreateRewardsPlan{
		Sender:                sender,
		Description:           description,
		ServiceID:             serviceID,
		Amount:                amt,
		StartTime:             startTime,
		EndTime:               endTime,
		PoolsDistribution:     poolsDistribution,
		OperatorsDistribution: operatorsDistribution,
		UsersDistribution:     usersDistribution,
		FeeAmount:             feeAmount,
	}
}

// ValidateBasic implements sdk.Msg
func (m *MsgCreateRewardsPlan) ValidateBasic() error {
	if len(m.Description) > MaxRewardsPlanDescriptionLength {
		return fmt.Errorf("too long description")
	}

	if m.ServiceID == 0 {
		return fmt.Errorf("invalid service ID: %d", m.ServiceID)
	}

	err := m.Amount.Validate()
	if err != nil {
		return fmt.Errorf("invalid amount per day: %w", err)
	}

	if !m.EndTime.After(m.StartTime) {
		return fmt.Errorf(
			"end time must be after start time: %s <= %s",
			m.EndTime.Format(time.RFC3339),
			m.StartTime.Format(time.RFC3339),
		)
	}

	if m.PoolsDistribution.DelegationType != restakingtypes.DELEGATION_TYPE_POOL {
		return fmt.Errorf("pools distribution has invalid delegation type: %v", m.PoolsDistribution.DelegationType)
	}

	if m.OperatorsDistribution.DelegationType != restakingtypes.DELEGATION_TYPE_OPERATOR {
		return fmt.Errorf("operators distribution has invalid delegation type: %v", m.OperatorsDistribution.DelegationType)
	}

	err = m.FeeAmount.Validate()
	if err != nil {
		return fmt.Errorf("invalid fee amount: %w", err)
	}

	_, err = sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return fmt.Errorf("invalid sender address: %s", m.Sender)
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgCreateRewardsPlan) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgCreateRewardsPlan) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (m *MsgCreateRewardsPlan) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	err := m.PoolsDistribution.UnpackInterfaces(unpacker)
	if err != nil {
		return nil
	}

	err = m.OperatorsDistribution.UnpackInterfaces(unpacker)
	if err != nil {
		return nil
	}

	err = m.UsersDistribution.UnpackInterfaces(unpacker)
	if err != nil {
		return nil
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgEditRewardsPlan creates a new MsgEditRewardsPlan instance.
func NewMsgEditRewardsPlan(
	id uint64,
	description string,
	amount sdk.Coins,
	startTime,
	endTime time.Time,
	poolsDistribution Distribution,
	operatorsDistribution Distribution,
	usersDistribution UsersDistribution,
	sender string,
) *MsgEditRewardsPlan {
	return &MsgEditRewardsPlan{
		ID:                    id,
		Sender:                sender,
		Description:           description,
		Amount:                amount,
		StartTime:             startTime,
		EndTime:               endTime,
		PoolsDistribution:     poolsDistribution,
		OperatorsDistribution: operatorsDistribution,
		UsersDistribution:     usersDistribution,
	}
}

// ValidateBasic implements sdk.Msg
func (m *MsgEditRewardsPlan) ValidateBasic() error {
	if m.ID == 0 {
		return fmt.Errorf("invalid ID: %d", m.ID)
	}

	err := m.Amount.Validate()
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	if !m.EndTime.After(m.StartTime) {
		return fmt.Errorf(
			"end time must be after start time: %s <= %s",
			m.EndTime.Format(time.RFC3339),
			m.StartTime.Format(time.RFC3339),
		)
	}

	if m.PoolsDistribution.DelegationType != restakingtypes.DELEGATION_TYPE_POOL {
		return fmt.Errorf("pools distribution has invalid delegation type: %v", m.PoolsDistribution.DelegationType)
	}

	if m.OperatorsDistribution.DelegationType != restakingtypes.DELEGATION_TYPE_OPERATOR {
		return fmt.Errorf("operators distribution has invalid delegation type: %v", m.OperatorsDistribution.DelegationType)
	}

	_, err = sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return fmt.Errorf("invalid sender address: %s, %w", m.Sender, err)
	}

	// We need a codec to properly validate the rewards plan, we do that
	// when handling the message.

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgEditRewardsPlan) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgEditRewardsPlan) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (m *MsgEditRewardsPlan) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	err := m.PoolsDistribution.UnpackInterfaces(unpacker)
	if err != nil {
		return nil
	}

	err = m.OperatorsDistribution.UnpackInterfaces(unpacker)
	if err != nil {
		return nil
	}

	err = m.UsersDistribution.UnpackInterfaces(unpacker)
	if err != nil {
		return nil
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgSetWithdrawAddress creates a new NewMsgSetWithdrawAddress instance
func NewMsgSetWithdrawAddress(withdrawAddress string, sender string) *MsgSetWithdrawAddress {
	return &MsgSetWithdrawAddress{
		Sender:          sender,
		WithdrawAddress: withdrawAddress,
	}
}

// ValidateBasic implements sdk.Msg
func (m *MsgSetWithdrawAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return fmt.Errorf("invalid sender address: %s", m.Sender)
	}

	_, err = sdk.AccAddressFromBech32(m.WithdrawAddress)
	if err != nil {
		return fmt.Errorf("invalid withdraw address: %s", m.WithdrawAddress)
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgSetWithdrawAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgWithdrawDelegatorReward creates a new MsgWithdrawDelegatorReward instance
func NewMsgWithdrawDelegatorReward(
	delegationType restakingtypes.DelegationType,
	targetID uint32,
	delegatorAddress string,
) *MsgWithdrawDelegatorReward {
	return &MsgWithdrawDelegatorReward{
		DelegatorAddress:   delegatorAddress,
		DelegationType:     delegationType,
		DelegationTargetID: targetID,
	}
}

// ValidateBasic implements sdk.Msg
func (m *MsgWithdrawDelegatorReward) ValidateBasic() error {
	if m.DelegationType == restakingtypes.DELEGATION_TYPE_UNSPECIFIED {
		return fmt.Errorf("invalid delegation type: %v", m.DelegationType)
	}

	if m.DelegationTargetID == 0 {
		return fmt.Errorf("invalid delegation target ID: %d", m.DelegationTargetID)
	}

	_, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		return fmt.Errorf("invalid delegator address: %s", m.DelegatorAddress)
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgWithdrawDelegatorReward) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgWithdrawOperatorCommission creates a new MsgWithdrawOperatorCommission instance
func NewMsgWithdrawOperatorCommission(operatorID uint32, senderAddress string) *MsgWithdrawOperatorCommission {
	return &MsgWithdrawOperatorCommission{
		Sender:     senderAddress,
		OperatorID: operatorID,
	}
}

// ValidateBasic implements sdk.Msg
func (m *MsgWithdrawOperatorCommission) ValidateBasic() error {
	if m.OperatorID == 0 {
		return fmt.Errorf("invalid operator ID: %d", m.OperatorID)
	}

	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return fmt.Errorf("invalid sender address: %s", m.Sender)
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgWithdrawOperatorCommission) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgWithdrawOperatorCommission) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgUpdateParams creates a new MsgUpdateParams instance
func NewMsgUpdateParams(params Params, authority string) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}
