package simulation

import (
	"context"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/milkyway-labs/milkyway/v12/testutils/simtesting"
	"github.com/milkyway-labs/milkyway/v12/utils"
	"github.com/milkyway-labs/milkyway/v12/x/operators/keeper"
	"github.com/milkyway-labs/milkyway/v12/x/operators/types"
)

// RandomOperator returns a random operator
func RandomOperator(r *rand.Rand, accounts []simtypes.Account) types.Operator {
	adminAccount, _ := simtypes.RandomAcc(r, accounts)

	return types.NewOperator(
		r.Uint32(),
		randomOperatorStatus(r),
		simtypes.RandStringOfLength(r, 10),
		simtypes.RandStringOfLength(r, 20),
		simtypes.RandStringOfLength(r, 20),
		adminAccount.Address.String(),
	)
}

// randomOperatorStatus returns a random operator status
func randomOperatorStatus(r *rand.Rand) types.OperatorStatus {
	statusesSize := len(types.OperatorStatus_name)
	return types.OperatorStatus(r.Intn(statusesSize-1) + 1)
}

// RandomOperatorParams returns random operator params
func RandomOperatorParams(r *rand.Rand) types.OperatorParams {
	return types.NewOperatorParams(
		sdkmath.LegacyNewDecWithPrec(int64(r.Intn(100)), 2),
	)
}

// RandomParams returns random params
func RandomParams(r *rand.Rand, stakeDenom string) types.Params {
	return types.NewParams(
		sdk.NewCoins(simtesting.RandomCoin(r, stakeDenom, 10)),
		simtesting.RandomDuration(r, 5*time.Minute, 3*24*time.Hour),
	)
}

// GetRandomExistingOperator returns a random existing operator
func GetRandomExistingOperator(r *rand.Rand, ctx context.Context, k *keeper.Keeper, filter func(operator types.Operator) bool) (types.Operator, bool) {
	operators, err := k.GetOperators(ctx)
	if err != nil {
		panic(err)
	}

	if len(operators) == 0 {
		return types.Operator{}, false
	}

	if filter != nil {
		operators = utils.Filter(operators, filter)
		if len(operators) == 0 {
			return types.Operator{}, false
		}
	}

	randomOperatorIndex := r.Intn(len(operators))
	return operators[randomOperatorIndex], true
}
