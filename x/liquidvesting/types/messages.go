package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgMintVestedRepresentation creates a new MsgMintVestedRepresentation instance.
func NewMsgMintVestedRepresentation(
	sender string,
	receiver string,
	amount sdk.Coins,
) *MsgMintVestedRepresentation {
	return &MsgMintVestedRepresentation{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}
}

func (msg *MsgMintVestedRepresentation) ValidateBasic() error {
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

// NewMsgBurnVestedRepresentation creates a new MsgBurnVestedRepresentation instance.
func NewMsgBurnVestedRepresentation(
	sender string,
	user string,
	amount sdk.Coins,
) *MsgBurnVestedRepresentation {
	return &MsgBurnVestedRepresentation{
		Sender: sender,
		User:   user,
		Amount: amount,
	}
}

func (msg *MsgBurnVestedRepresentation) ValidateBasic() error {
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
