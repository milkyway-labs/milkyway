package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	LockedRepresentationDenomPrefix = "locked"
)

// GetLockedRepresentationDenom returns the denom used to
// represent the locked version of the provided denom.
func GetLockedRepresentationDenom(denom string) (string, error) {
	if IsLockedRepresentationDenom(denom) {
		return "", ErrInvalidDenom
	}

	// Create the locked representation of the provided denom
	lockedDenom := strings.Join([]string{LockedRepresentationDenomPrefix, denom}, "/")
	return lockedDenom, sdk.ValidateDenom(lockedDenom)
}

// IsLockedRepresentationDenom tells if the provided denom is
// a representation of a locked denom.
func IsLockedRepresentationDenom(denom string) bool {
	return strings.HasPrefix(denom, LockedRepresentationDenomPrefix+"/")
}

// LockedDenomToNative converts the denom of a locked token representation
// to its native denom.
func LockedDenomToNative(denom string) (string, error) {
	if !strings.HasPrefix(denom, LockedRepresentationDenomPrefix+"/") {
		return "", ErrInvalidDenom
	}

	return strings.TrimPrefix(denom, LockedRepresentationDenomPrefix+"/"), nil
}
