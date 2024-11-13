package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/utils"
	"github.com/milkyway-labs/milkyway/x/services/keeper"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

func RandomServiceStatus(r *rand.Rand) types.ServiceStatus {
	switch r.Intn(2) {
	case 0:
		return types.SERVICE_STATUS_INACTIVE
	case 1:
		return types.SERVICE_STATUS_CREATED
	default:
		return types.SERVICE_STATUS_ACTIVE
	}
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

func GetRandomExistingService(r *rand.Rand, ctx sdk.Context, k *keeper.Keeper, filter func(s types.Service) bool) (types.Service, bool) {
	services := k.GetServices(ctx)
	if len(services) == 0 {
		return types.Service{}, false
	}
	if filter != nil {
		services = utils.Filter(services, filter)
		if len(services) == 0 {
			return types.Service{}, false
		}
	}

	randomServiceIndex := r.Intn(len(services))
	return services[randomServiceIndex], true
}
