package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

func RandomServiceStatus(r *rand.Rand) types.ServiceStatus {
	value := (int32(r.Uint64())) % 3
	// Here we add 1 since 0 is SERVICE_STATUS_UNSPECIFIED.
	return types.ServiceStatus(value + 1)
}

func RandomService(r *rand.Rand, id uint32, admin string) types.Service {
	return types.NewService(id,
		RandomServiceStatus(r),
		simulation.RandStringOfLength(r, 24),
		simulation.RandStringOfLength(r, 24),
		simulation.RandStringOfLength(r, 24),
		simulation.RandStringOfLength(r, 24),
		admin,
		(r.Uint64()%2) == 0,
	)
}
