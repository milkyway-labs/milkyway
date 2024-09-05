package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgMintStakingRepresentation creates a new MsgMintStakingRepresentation instance.
func NewMsgMintStakingRepresentation(
	sender string,
	receiver string,
	amount types.Coins,
) *MsgMintStakingRepresentation {
	return &MsgMintStakingRepresentation{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}
}

func (msg *MsgMintStakingRepresentation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	_, err = sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receiver address")
	}

	if err := msg.Amount.Validate(); err != nil {
		return err
	}

	return nil
}

// NewMsgBurnStakingRepresentation creates a new MsgBurnStakingRepresentation instance.
func NewMsgBurnStakingRepresentation(
	sender string,
	user string,
	amount types.Coins,
) *MsgBurnStakingRepresentation {
	return &MsgBurnStakingRepresentation{
		Sender: sender,
		User:   user,
		Amount: amount,
	}
}

func (msg *MsgBurnStakingRepresentation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	_, err = sdk.AccAddressFromBech32(msg.User)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid user address")
	}

	if err := msg.Amount.Validate(); err != nil {
		return err
	}

	return nil
}
