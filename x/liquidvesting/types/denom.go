package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	VestedRepresentationDenomPrefix = "vested"
)

// GetVestedRepresentationDenom returns the denom used to
// represent the vested version of the provided denom.
func GetVestedRepresentationDenom(denom string) (string, error) {
	// Create the vested representation of the provided denom
	vestedDenom := strings.Join([]string{VestedRepresentationDenomPrefix, denom}, "/")
	return vestedDenom, sdk.ValidateDenom(vestedDenom)
}

// IsVestedRepresentationDenom tells if the provided denom is
// a representation of a vested denom.
func IsVestedRepresentationDenom(denom string) (string, bool) {
	if !strings.HasPrefix(denom, VestedRepresentationDenomPrefix+"/") {
		return "", false
	}

	return strings.TrimPrefix(denom, VestedRepresentationDenomPrefix+"/"), true
}
