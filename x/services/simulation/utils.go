package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/testutils/simtesting"
	"github.com/milkyway-labs/milkyway/utils"
	"github.com/milkyway-labs/milkyway/x/services/keeper"
	"github.com/milkyway-labs/milkyway/x/services/types"
)

// RandomService returns a random service
func RandomService(r *rand.Rand, accs []simtypes.Account) types.Service {
	adminAccount, _ := simtypes.RandomAcc(r, accs)
	return types.NewService(
		r.Uint32(),
		randomServiceStatus(r),
		simtypes.RandStringOfLength(r, 24),
		simtypes.RandStringOfLength(r, 24),
		simtypes.RandStringOfLength(r, 24),
		simtypes.RandStringOfLength(r, 24),
		adminAccount.Address.String(),
		(r.Uint64()%2) == 0,
	)
}

func randomServiceStatus(r *rand.Rand) types.ServiceStatus {
	statusesSize := len(types.ServiceStatus_name)
	return types.ServiceStatus(r.Intn(statusesSize-1) + 1)
}

// RandomServiceParams returns a random service params
func RandomServiceParams(r *rand.Rand) types.ServiceParams {
	var allowedDenom []string
	if r.Intn(2) == 0 {
		// 50% chance of having an empty list of allowed denoms
		for i := 0; i < r.Intn(10); i++ {
			generatedDenom := simtypes.RandStringOfLength(r, 5)
			if sdk.ValidateDenom(generatedDenom) != nil {
				continue
			}

			allowedDenom = append(allowedDenom, generatedDenom)
		}
	}

	return types.NewServiceParams(allowedDenom)
}

// RandomParams returns a random params
func RandomParams(r *rand.Rand, bondDenom string) types.Params {
	return types.NewParams(
		sdk.NewCoins(simtesting.RandomCoin(r, bondDenom, 10)),
	)
}

// GetRandomExistingService returns a random existing service
func GetRandomExistingService(r *rand.Rand, ctx sdk.Context, k *keeper.Keeper, filter func(s types.Service) bool) (types.Service, bool) {
	services, err := k.GetServices(ctx)
	if err != nil {
		panic(err)
	}

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
