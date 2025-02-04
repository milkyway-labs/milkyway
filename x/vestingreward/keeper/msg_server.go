package keeper

import (
	"context"

	"cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/milkyway-labs/milkyway/v7/x/vestingreward/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(k *Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

// UpdateVestingAccountsRewardRatio defines the rpc method for Msg/UpdateVestingAccountsRewardRatio
func (k msgServer) UpdateVestingAccountsRewardRatio(ctx context.Context, msg *types.MsgUpdateVestingAccountsRewardRatio) (*types.MsgUpdateVestingAccountsRewardRatioResponse, error) {
	// Check the authority
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	err := k.Keeper.UpdateVestingAccountsRewardRatio(ctx, msg.VestingAccountsRewardRatio)
	if err != nil {
		return nil, err
	}
	return &types.MsgUpdateVestingAccountsRewardRatioResponse{}, nil
}
