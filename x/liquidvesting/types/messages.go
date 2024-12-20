package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgMintLockedRepresentation creates a new MsgMintLockedRepresentation instance.
func NewMsgMintLockedRepresentation(
	sender string,
	receiver string,
	amount sdk.Coins,
) *MsgMintLockedRepresentation {
	return &MsgMintLockedRepresentation{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}
}

func (msg *MsgMintLockedRepresentation) ValidateBasic() error {
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

	if msg.Amount.IsZero() {
		return ErrInvalidAmount
	}

	return nil
}

// NewMsgBurnLockedRepresentation creates a new MsgBurnLockedRepresentation instance.
func NewMsgBurnLockedRepresentation(
	sender string,
	user string,
	amount sdk.Coins,
) *MsgBurnLockedRepresentation {
	return &MsgBurnLockedRepresentation{
		Sender: sender,
		User:   user,
		Amount: amount,
	}
}

func (msg *MsgBurnLockedRepresentation) ValidateBasic() error {
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

	if msg.Amount.IsZero() {
		return ErrInvalidAmount
	}

	return nil
}

// NewMsgWithdrawInsuranceFund creates a new MsgWithdrawInsuranceFund instance.
func NewMsgWithdrawInsuranceFund(
	sender string,
	amount sdk.Coins,
) *MsgWithdrawInsuranceFund {
	return &MsgWithdrawInsuranceFund{
		Sender: sender,
		Amount: amount,
	}
}

func (msg *MsgWithdrawInsuranceFund) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}
	return msg.Amount.Validate()
}

// NewMsgUpdateParams creates a new MsgUpdateParams instance
func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	return msg.Params.Validate()
}
