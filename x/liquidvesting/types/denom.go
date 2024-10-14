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
	if IsVestedRepresentationDenom(denom) {
		return "", ErrInvalidDenom
	}

	// Create the vested representation of the provided denom
	vestedDenom := strings.Join([]string{VestedRepresentationDenomPrefix, denom}, "/")
	return vestedDenom, sdk.ValidateDenom(vestedDenom)
}

// IsVestedRepresentationDenom tells if the provided denom is
// a representation of a vested denom.
func IsVestedRepresentationDenom(denom string) bool {
	return strings.HasPrefix(denom, VestedRepresentationDenomPrefix+"/")
}

// VestedDenomToNative converts the denom of a vested token representation
// to its native denom.
func VestedDenomToNative(denom string) (string, error) {
	if !strings.HasPrefix(denom, VestedRepresentationDenomPrefix+"/") {
		return "", ErrInvalidDenom
	}

	return strings.TrimPrefix(denom, VestedRepresentationDenomPrefix+"/"), nil
}
