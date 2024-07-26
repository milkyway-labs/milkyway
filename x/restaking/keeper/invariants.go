package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k *Keeper) {
	ir.RegisterRoute(types.ModuleName, "accounts-balances",
		AccountsBalancesInvariants(k))

	ir.RegisterRoute(types.ModuleName, "positive-pools-delegations",
		PositivePoolsDelegationsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "pools-delegators-shares",
		PoolsDelegatorsSharesInvariant(k))

	ir.RegisterRoute(types.ModuleName, "positive-operators-delegations",
		PositiveOperatorsDelegationsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "operators-delegators-shares",
		OperatorsDelegatorsSharesInvariant(k))

	ir.RegisterRoute(types.ModuleName, "positive-services-delegations",
		PositiveServicesDelegationsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "services-delegators-shares",
		ServicesDelegatorsSharesInvariant(k))
}

// AccountsBalancesInvariants checks that the pools, operators and services accounts have the correct balance
// based on the delegations that are stored in the store
func AccountsBalancesInvariants(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// Get all the pools balances and tokens
		poolsBalances := sdk.NewCoins()
		poolsTokens := sdk.NewCoins()
		k.poolsKeeper.IteratePools(ctx, func(pool poolstypes.Pool) (stop bool) {
			poolAddress, err := sdk.AccAddressFromBech32(pool.GetAddress())
			if err != nil {
				panic(err)
			}

			poolBalance := k.bankKeeper.GetBalance(ctx, poolAddress, pool.GetDenom())
			poolsBalances = poolsBalances.Add(poolBalance)

			poolsTokens = poolsTokens.Add(sdk.NewCoin(pool.GetDenom(), pool.Tokens))

			return false
		})

		// Get all the operators balances and tokens
		operatorsBalances := sdk.NewCoins()
		operatorsTokens := sdk.NewCoins()
		k.operatorsKeeper.IterateOperators(ctx, func(operator operatorstypes.Operator) (stop bool) {
			operatorAddress, err := sdk.AccAddressFromBech32(operator.GetAddress())
			if err != nil {
				panic(err)
			}

			operatorBalance := k.bankKeeper.GetAllBalances(ctx, operatorAddress)
			operatorsBalances = operatorsBalances.Add(operatorBalance...)

			operatorsTokens = operatorsTokens.Add(operator.Tokens...)

			return false
		})

		// Get all the services balances and tokens
		servicesBalances := sdk.NewCoins()
		servicesTokens := sdk.NewCoins()
		k.servicesKeeper.IterateServices(ctx, func(service servicestypes.Service) (stop bool) {
			serviceAddress, err := sdk.AccAddressFromBech32(service.GetAddress())
			if err != nil {
				panic(err)
			}

			serviceBalance := k.bankKeeper.GetAllBalances(ctx, serviceAddress)
			servicesBalances = servicesBalances.Add(serviceBalance...)

			servicesTokens = servicesTokens.Add(service.Tokens...)

			return false
		})

		poolsBroken := !poolsBalances.IsAllGTE(poolsTokens)
		operatorsBroken := !operatorsBalances.IsAllGTE(operatorsTokens)
		servicesBroken := !servicesBalances.IsAllGTE(servicesTokens)
		broken := poolsBroken || operatorsBroken || servicesBroken

		return sdk.FormatInvariant(types.ModuleName, "delegated module account coins", fmt.Sprintf(
			"\tPools' bonded tokens: %v\n"+
				"\tsum of bonded tokens: %v\n"+
				"\tOperators' bonded tokens: %v\n"+
				"\tsum of bonded tokens: %v\n"+
				"\tServices' bonded tokens: %v\n"+
				"\tsum of bonded tokens: %v\n",
			poolsTokens, poolsBalances, operatorsTokens, operatorsBalances, servicesTokens, servicesBalances)), broken
	}
}

// PositivePoolsDelegationsInvariant checks that all stored pools delegations have > 0 shares
func PositivePoolsDelegationsInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int

		for _, delegation := range k.GetAllPoolDelegations(ctx) {
			if delegation.Shares.IsAnyNegative() {
				count++
				msg += fmt.Sprintf("pool delegation with negative shares: %v\n", delegation)
			}

			if delegation.Shares.IsZero() {
				count++
				msg += fmt.Sprintf("pool delegation with zero shares: %v\n", delegation)
			}
		}

		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "positive pool delegations", fmt.Sprintf(
			"%d invalid pool delegations found\n%s", count, msg)), broken
	}
}

// PoolsDelegatorsSharesInvariant checks that the sum of all delegators shares for each pool is equal
// to the total delegators shares stored in the pool object
func PoolsDelegatorsSharesInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var broken bool

		// Initialize a map: pool id -> its delegators shares
		poolsDelegatorsShares := map[uint32]sdk.DecCoins{}
		for _, pool := range k.poolsKeeper.GetPools(ctx) {
			poolsDelegatorsShares[pool.ID] = sdk.NewDecCoins()
		}

		// Iterate through all the pool delegations to calculate the total delegators shares for each pool
		for _, delegation := range k.GetAllPoolDelegations(ctx) {
			if _, ok := poolsDelegatorsShares[delegation.TargetID]; !ok {
				poolsDelegatorsShares[delegation.TargetID] = sdk.NewDecCoins()
			}

			poolsDelegatorsShares[delegation.TargetID] = poolsDelegatorsShares[delegation.TargetID].Add(delegation.Shares...)
		}

		for poolID, delegatorsShares := range poolsDelegatorsShares {
			pool, found := k.poolsKeeper.GetPool(ctx, poolID)
			if !found {
				panic(fmt.Errorf("pool with id %d not found", poolID))
			}

			sharesAmount := delegatorsShares.AmountOf(pool.GetSharesDenom(pool.Denom))
			if !pool.DelegatorShares.Equal(sharesAmount) {
				broken = true
				msg += fmt.Sprintf("pool %d total shares: %v, delegators shares: %v\n", poolID, pool.DelegatorShares, delegatorsShares)
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "pools delegators shares", fmt.Sprintf(
			"pools delegators shares invariant broken\n%s", msg)), broken
	}
}

// PositiveOperatorsDelegationsInvariant checks that all stored operator delegations have > 0 shares.
func PositiveOperatorsDelegationsInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int

		for _, delegation := range k.GetAllOperatorDelegations(ctx) {
			if delegation.Shares.IsAnyNegative() {
				count++
				msg += fmt.Sprintf("operator delegation with negative shares: %v\n", delegation)
			}

			if delegation.Shares.IsZero() {
				count++
				msg += fmt.Sprintf("operator delegation with zero shares: %v\n", delegation)
			}
		}

		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "positive operator delegations", fmt.Sprintf(
			"%d invalid operator delegations found\n%s", count, msg)), broken
	}
}

// OperatorsDelegatorsSharesInvariant checks that the sum of all delegators shares for each operator is equal
// to the total delegators shares stored in the operator object
func OperatorsDelegatorsSharesInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var broken bool

		// Initialize a map: operator id -> its delegators shares
		operatorsDelegatorsShares := map[uint32]sdk.DecCoins{}
		for _, operator := range k.operatorsKeeper.GetOperators(ctx) {
			operatorsDelegatorsShares[operator.ID] = sdk.NewDecCoins()
		}

		// Iterate through all the operator delegations to calculate the total delegators shares for each operator
		for _, delegation := range k.GetAllOperatorDelegations(ctx) {
			if _, ok := operatorsDelegatorsShares[delegation.TargetID]; !ok {
				operatorsDelegatorsShares[delegation.TargetID] = sdk.NewDecCoins()
			}

			operatorsDelegatorsShares[delegation.TargetID] = operatorsDelegatorsShares[delegation.TargetID].Add(delegation.Shares...)
		}

		for operatorID, delegatorsShares := range operatorsDelegatorsShares {
			operator, found := k.operatorsKeeper.GetOperator(ctx, operatorID)
			if !found {
				panic(fmt.Errorf("operator with id %d not found", operatorID))
			}

			if !operator.DelegatorShares.Equal(delegatorsShares) {
				broken = true
				msg += fmt.Sprintf("operator %d total shares: %v, delegators shares: %v\n", operatorID, operator.DelegatorShares, delegatorsShares)
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "operators delegators shares", fmt.Sprintf(
			"operators delegators shares invariant broken\n%s", msg)), broken
	}
}

// PositiveServicesDelegationsInvariant checks that all stored service delegations have > 0 shares.
func PositiveServicesDelegationsInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int

		for _, delegation := range k.GetAllServiceDelegations(ctx) {
			if delegation.Shares.IsAnyNegative() {
				count++
				msg += fmt.Sprintf("service delegation with negative shares: %v\n", delegation)
			}

			if delegation.Shares.IsZero() {
				count++
				msg += fmt.Sprintf("service delegation with zero shares: %v\n", delegation)
			}
		}

		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "positive service delegations", fmt.Sprintf(
			"%d invalid service delegations found\n%s", count, msg)), broken
	}
}

// ServicesDelegatorsSharesInvariant checks that the sum of all delegators shares for each service is equal
// to the total delegators shares stored in the service object
func ServicesDelegatorsSharesInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var broken bool

		// Initialize a map: service id -> its delegators shares
		servicesDelegatorsShares := map[uint32]sdk.DecCoins{}
		for _, service := range k.servicesKeeper.GetServices(ctx) {
			servicesDelegatorsShares[service.ID] = sdk.NewDecCoins()
		}

		// Iterate through all the service delegations to calculate the total delegators shares for each service
		for _, delegation := range k.GetAllServiceDelegations(ctx) {
			if _, ok := servicesDelegatorsShares[delegation.TargetID]; !ok {
				servicesDelegatorsShares[delegation.TargetID] = sdk.NewDecCoins()
			}

			servicesDelegatorsShares[delegation.TargetID] = servicesDelegatorsShares[delegation.TargetID].Add(delegation.Shares...)
		}

		for serviceID, delegatorsShares := range servicesDelegatorsShares {
			service, found := k.servicesKeeper.GetService(ctx, serviceID)
			if !found {
				panic(fmt.Errorf("service with id %d not found", serviceID))
			}

			if !service.DelegatorShares.Equal(delegatorsShares) {
				broken = true
				msg += fmt.Sprintf("service %d total shares: %v, delegators shares: %v\n", serviceID, service.DelegatorShares, delegatorsShares)
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "services delegators shares", fmt.Sprintf(
			"services delegators shares invariant broken\n%s", msg)), broken
	}
}
