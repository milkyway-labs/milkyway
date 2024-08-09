package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/milkyway-labs/milkyway/utils"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

func TestDelegationTarget_TokensFromShares(t *testing.T) {
	pool := poolstypes.NewPool(1, "umilk")
	pool.Tokens = math.NewInt(10_000000)
	pool.DelegatorShares = math.LegacyNewDec(9_999990)

	target := types.NewDelegationTarget(&pool)

	tokens := target.TokensFromShares(utils.MustParseDecCoins("9_999990pool/1/umilk"))
	require.Equal(t, utils.MustParseDecCoins("10_000000umilk"), tokens)

	tokens = target.TokensFromShares(utils.MustParseDecCoins("1_000000pool/1/umilk"))
	require.Equal(t, utils.MustParseDecCoins("1_000001.000001000001000001umilk"), tokens)
}
