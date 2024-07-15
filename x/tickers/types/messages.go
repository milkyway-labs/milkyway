package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterTicker{}
	_ sdk.Msg = &MsgDeregisterTicker{}
	_ sdk.Msg = &MsgUpdateParams{}
)

func NewMsgRegisterTicker(authority, denom, ticker string) *MsgRegisterTicker {
	return &MsgRegisterTicker{
		Authority: authority,
		Denom:     denom,
		Ticker:    ticker,
	}
}

func (msg *MsgRegisterTicker) Validate() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return err
	}
	// TODO: validate ticker
	return nil
}

func NewMsgDeregisterTicker(authority, denom string) *MsgDeregisterTicker {
	return &MsgDeregisterTicker{
		Authority: authority,
		Denom:     denom,
	}
}

func (msg *MsgDeregisterTicker) Validate() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return err
	}
	return nil
}

func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

func (msg *MsgUpdateParams) Validate() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}
	if err := msg.Params.Validate(); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid params: %s", err.Error())
	}
	return nil
}
