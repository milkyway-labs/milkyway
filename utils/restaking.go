package utils

import (
	"fmt"
	"strings"
)

// GetSharesDenomFromTokenDenom returns the shares denom from the token denom.
// The returned shares denom will be in the format "{prefix}/{id}/{tokenDenom}".
func GetSharesDenomFromTokenDenom(prefix string, id uint32, tokenDenom string) string {
	return fmt.Sprintf("%s/%d/%s", prefix, id, tokenDenom)
}

// GetTokenDenomFromSharesDenom returns the token denom from the shares denom.
// It expects the shares denom to be in the format "{xxxxxx}/{xxxxxx}/{tokenDenom}".
func GetTokenDenomFromSharesDenom(sharesDenom string) string {
	parts := strings.Split(sharesDenom, "/")
	if len(parts) != 3 {
		return ""
	}
	return parts[2]
}
