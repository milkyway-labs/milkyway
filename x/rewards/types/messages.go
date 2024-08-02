package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
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

func NewMsgSetWithdrawAddress(delAddr, withdrawAddr string) *MsgSetWithdrawAddress {
	return &MsgSetWithdrawAddress{
		DelegatorAddress: delAddr,
		WithdrawAddress:  withdrawAddr,
	}
}

// NewMsgWithdrawDelegationReward creates a new MsgWithdrawDelegationReward instance
func NewMsgWithdrawDelegationReward(delAddr string, delType restakingtypes.DelegationType, poolID uint32) *MsgWithdrawDelegationReward {
	return &MsgWithdrawDelegationReward{
		DelegatorAddress: delAddr,
		DelegationType:   delType,
		TargetID:         poolID,
	}
}

// NewMsgWithdrawOperatorCommission creates a new MsgWithdrawOperatorCommission instance
func NewMsgWithdrawOperatorCommission(operatorID uint32) *MsgWithdrawOperatorCommission {
	return &MsgWithdrawOperatorCommission{
		OperatorID: operatorID,
	}
}
