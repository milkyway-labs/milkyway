package types

import (
	"strings"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	VestedRepresentationDenomPrefix = "vested"
)

// GetVestedRepresentationDenom returns the denom used to
// represent the vested version of the provided denom.
func GetVestedRepresentationDenom(denom string) (string, error) {
	// Ensure that
	strParts := strings.Split(denom, "/")
	if len(strParts) > 0 {
		return "", errors.Wrapf(ErrInvalidDenom, "denom %s is invalid", denom)
	}

	vestedDenom := strings.Join([]string{VestedRepresentationDenomPrefix, denom}, "/")
	return vestedDenom, sdk.ValidateDenom(vestedDenom)
}
