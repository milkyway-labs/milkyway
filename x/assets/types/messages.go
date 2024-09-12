package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterAsset{}
	_ sdk.Msg = &MsgDeregisterAsset{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// NewMsgRegisterAsset creates a new MsgRegisterAsset instance
func NewMsgRegisterAsset(authority string, asset Asset) *MsgRegisterAsset {
	return &MsgRegisterAsset{
		Authority: authority,
		Asset:     asset,
	}
}

// Validate validates the MsgRegisterAsset instance
func (msg *MsgRegisterAsset) Validate() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}

	err = msg.Asset.Validate()
	if err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// NewMsgDeregisterAsset creates a new MsgDeregisterAsset instance
func NewMsgDeregisterAsset(authority, denom string) *MsgDeregisterAsset {
	return &MsgDeregisterAsset{
		Authority: authority,
		Denom:     denom,
	}
}

// Validate validates the MsgDeregisterAsset instance
func (msg *MsgDeregisterAsset) Validate() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}
	err = sdk.ValidateDenom(msg.Denom)
	if err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return nil
}

// NewMsgUpdateParams creates a new MsgUpdateParams instance
func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

// Validate validates the MsgUpdateParams instance
func (msg *MsgUpdateParams) Validate() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address")
	}
	err = msg.Params.Validate()
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid params: %s", err)
	}
	return nil
}
