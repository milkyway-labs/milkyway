package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type msgServer struct {
	types.MsgServer
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the distribution MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		MsgServer: distrkeeper.NewMsgServerImpl(keeper.Keeper),
		Keeper:    keeper,
	}
}

func (k msgServer) SetWithdrawAddress(ctx context.Context, msg *types.MsgSetWithdrawAddress) (*types.MsgSetWithdrawAddressResponse, error) {
	delegatorAddress, err := k.authKeeper.AddressCodec().StringToBytes(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}

	withdrawAddress, err := k.authKeeper.AddressCodec().StringToBytes(msg.WithdrawAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid withdraw address: %s", err)
	}

	err = k.SetWithdrawAddr(ctx, delegatorAddress, withdrawAddress)
	if err != nil {
		return nil, err
	}

	err = k.AfterSetWithdrawAddress(ctx, delegatorAddress, withdrawAddress)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetWithdrawAddressResponse{}, nil
}
