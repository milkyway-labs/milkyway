package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v11/utils"
	"github.com/milkyway-labs/milkyway/v11/x/pools/keeper"
	"github.com/milkyway-labs/milkyway/v11/x/pools/types"
)

// GetRandomExistingPool returns a random existing pool
func GetRandomExistingPool(r *rand.Rand, ctx sdk.Context, k *keeper.Keeper, filter func(s types.Pool) bool) (types.Pool, bool) {
	pools, err := k.GetPools(ctx)
	if err != nil {
		panic(err)
	}

	if filter != nil {
		pools = utils.Filter(pools, filter)
	}

	if len(pools) == 0 {
		return types.Pool{}, false
	}

	randomServiceIndex := r.Intn(len(pools))
	return pools[randomServiceIndex], true
}
