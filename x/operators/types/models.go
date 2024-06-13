package types

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// GetOperatorAddress generates an operator address from its id
func GetOperatorAddress(operatorID uint32) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("operator-%d", operatorID))
}

// ParseOperatorID tries parsing the given value as an operator id
func ParseOperatorID(value string) (uint32, error) {
	operatorID, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid operator ID: %s", value)
	}
	return uint32(operatorID), nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewOperator creates a new Operator object
func NewOperator(
	id uint32,
	status OperatorStatus,
	moniker string,
	website string,
	pictureURL string,
	admin string,
) Operator {
	return Operator{
		ID:         id,
		Status:     status,
		Moniker:    moniker,
		Website:    website,
		PictureURL: pictureURL,
		Admin:      admin,
		Address:    GetOperatorAddress(id).String(),
	}
}

// Validate checks that the Operator has valid values.
func (o *Operator) Validate() error {
	if o.ID == 0 {
		return fmt.Errorf("invalid id: %d", o.ID)
	}

	if o.Status == OPERATOR_STATUS_UNSPECIFIED {
		return fmt.Errorf("invalid status: %s", o.Status)
	}

	if strings.TrimSpace(o.Moniker) == "" {
		return fmt.Errorf("invalid moniker: %s", o.Moniker)
	}

	_, err := sdk.AccAddressFromBech32(o.Admin)
	if err != nil {
		return fmt.Errorf("invalid admin address: %s", o.Admin)
	}

	_, err = sdk.AccAddressFromBech32(o.Address)
	if err != nil {
		return fmt.Errorf("invalid address: %s", o.Address)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// OperatorUpdate defines the fields that can be updated in an Operator.
type OperatorUpdate struct {
	Moniker    string
	Website    string
	PictureURL string
}

func NewOperatorUpdate(
	moniker string,
	website string,
	pictureURL string,
) OperatorUpdate {
	return OperatorUpdate{
		Moniker:    moniker,
		Website:    website,
		PictureURL: pictureURL,
	}
}

// Update returns a new Operator with updated fields.
func (o *Operator) Update(update OperatorUpdate) Operator {
	if update.Moniker == DoNotModify {
		update.Moniker = o.Moniker
	}

	if update.Website == DoNotModify {
		update.Website = o.Website
	}

	if update.PictureURL == DoNotModify {
		update.PictureURL = o.PictureURL
	}

	return NewOperator(
		o.ID,
		o.Status,
		update.Moniker,
		update.Website,
		update.PictureURL,
		o.Admin,
	)
}
