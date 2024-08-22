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
