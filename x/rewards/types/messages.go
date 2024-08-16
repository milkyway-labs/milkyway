package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
)

// NewMsgCreateRewardsPlan creates a new MsgCreateRewardsPlan instance
func NewMsgCreateRewardsPlan(
	sender string, description string, serviceID uint32, amt sdk.Coins, startTime, endTime time.Time,
	poolsDistribution Distribution, operatorsDistribution Distribution,
	usersDistribution UsersDistribution) *MsgCreateRewardsPlan {
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
func NewMsgSetWithdrawAddress(senderAddr, withdrawAddr string) *MsgSetWithdrawAddress {
	return &MsgSetWithdrawAddress{
		Sender:          senderAddr,
		WithdrawAddress: withdrawAddr,
	}
}

// NewMsgWithdrawDelegatorReward creates a new MsgWithdrawDelegatorReward instance
func NewMsgWithdrawDelegatorReward(
	delAddr string, delType restakingtypes.DelegationType, targetID uint32,
) *MsgWithdrawDelegatorReward {
	return &MsgWithdrawDelegatorReward{
		DelegatorAddress:   delAddr,
		DelegationType:     delType,
		DelegationTargetID: targetID,
	}
}

// NewMsgWithdrawOperatorCommission creates a new MsgWithdrawOperatorCommission instance
func NewMsgWithdrawOperatorCommission(sender string, operatorID uint32) *MsgWithdrawOperatorCommission {
	return &MsgWithdrawOperatorCommission{
		Sender:     sender,
		OperatorID: operatorID,
	}
}
