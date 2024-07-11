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
