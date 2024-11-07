package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
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
	}
}

// NewMsgSetWithdrawAddress creates a new NewMsgSetWithdrawAddress instance
func NewMsgSetWithdrawAddress(withdrawAddress string, userAddress string) *MsgSetWithdrawAddress {
	return &MsgSetWithdrawAddress{
		Sender:          userAddress,
		WithdrawAddress: withdrawAddress,
	}
}

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

// NewMsgWithdrawOperatorCommission creates a new MsgWithdrawOperatorCommission instance
func NewMsgWithdrawOperatorCommission(operatorID uint32, senderAddress string) *MsgWithdrawOperatorCommission {
	return &MsgWithdrawOperatorCommission{
		Sender:     senderAddress,
		OperatorID: operatorID,
	}
}

// -------------------------------------------------------------------------------

// NewMsgEditRewardsPlan creates a new MsgEditRewardsPlan instance.
func NewMsgEditRewardsPlan(
	id uint64,
	description string,
	amt sdk.Coins,
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
		Amount:                amt,
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

	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return fmt.Errorf("invalid sender address: %s, %w", m.Sender, err)
	}

	// We need a codec to properly validate the rewards plan, we do that
	// when handling the message.

	return nil
}
