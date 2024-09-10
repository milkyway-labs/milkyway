package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterAsset{}
	_ sdk.Msg = &MsgDeregisterAsset{}
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
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
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
