package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/v10/utils"
)

func TestGetTokenDenomFromSharesDenom_IBCDenom(t *testing.T) {
	tokenDenom := "ibc/37A3FB4FED4CA04ED6D9E5DA36C6D27248645F0E22F585576A1488B8A89C5A50"
	sharesDenom := utils.GetSharesDenomFromTokenDenom("operator", 1, tokenDenom)
	require.Equal(t, "operator/1/ibc/37A3FB4FED4CA04ED6D9E5DA36C6D27248645F0E22F585576A1488B8A89C5A50", sharesDenom)
	res := utils.GetTokenDenomFromSharesDenom(sharesDenom)
	require.Equal(t, tokenDenom, res)
}
