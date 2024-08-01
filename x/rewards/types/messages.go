package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMsgCreateRewardsPlan creates a new MsgCreateRewardsPlan instance
func NewMsgCreateRewardsPlan(
	sender string, description string, serviceID uint32, amt sdk.Coins, startTime, endTime time.Time,
	poolsDistribution PoolsDistribution, operatorsDistribution OperatorsDistribution,
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

// NewMsgWithdrawPoolDelReward creates a new MsgWithdrawPoolDelReward instance
func NewMsgWithdrawPoolDelReward(delAddr string, poolID uint32) *MsgWithdrawPoolDelReward {
	return &MsgWithdrawPoolDelReward{
		DelegatorAddress: delAddr,
		PoolID:           poolID,
	}
}

// NewMsgWithdrawOperatorDelReward creates a new MsgWithdrawOperatorDelReward instance
func NewMsgWithdrawOperatorDelReward(delAddr string, operatorID uint32) *MsgWithdrawOperatorDelReward {
	return &MsgWithdrawOperatorDelReward{
		DelegatorAddress: delAddr,
		OperatorID:       operatorID,
	}
}

// NewMsgWithdrawServiceDelReward creates a new MsgWithdrawServiceDelReward instance
func NewMsgWithdrawServiceDelReward(delAddr string, serviceID uint32) *MsgWithdrawServiceDelReward {
	return &MsgWithdrawServiceDelReward{
		DelegatorAddress: delAddr,
		ServiceID:        serviceID,
	}
}

// NewMsgWithdrawOperatorCommission creates a new MsgWithdrawOperatorCommission instance
func NewMsgWithdrawOperatorCommission(operatorID uint32) *MsgWithdrawOperatorCommission {
	return &MsgWithdrawOperatorCommission{
		OperatorID: operatorID,
	}
}
