package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/v3/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/v3/x/pools/types"
	"github.com/milkyway-labs/milkyway/v3/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v3/x/services/types"
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

	ir.RegisterRoute(types.ModuleName, "allowed-operators-exist",
		AllowedOperatorsExistInvariant(k))
	ir.RegisterRoute(types.ModuleName, "operators-joined-services-exist",
		OperatorsJoinedServicesExistInvariant(k))

	ir.RegisterRoute(types.ModuleName, "total-restaked-assets",
		TotalRestakedAssetsInvariant(k))
}

// AccountsBalancesInvariants checks that the pools, operators and services accounts have the correct balance
// based on the delegations that are stored in the store
func AccountsBalancesInvariants(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// Get all the pools balances and tokens
		poolsBalances := sdk.NewCoins()
		poolsTokens := sdk.NewCoins()
		err := k.poolsKeeper.IteratePools(ctx, func(pool poolstypes.Pool) (stop bool, err error) {
			poolAddress, err := sdk.AccAddressFromBech32(pool.GetAddress())
			if err != nil {
				panic(err)
			}

			poolBalance := k.bankKeeper.GetBalance(ctx, poolAddress, pool.GetDenom())
			poolsBalances = poolsBalances.Add(poolBalance)

			poolsTokens = poolsTokens.Add(sdk.NewCoin(pool.GetDenom(), pool.Tokens))

			return false, nil
		})
		if err != nil {
			panic(err)
		}

		// Get all the operators balances and tokens
		operatorsBalances := sdk.NewCoins()
		operatorsTokens := sdk.NewCoins()
		err = k.operatorsKeeper.IterateOperators(ctx, func(operator operatorstypes.Operator) (stop bool, err error) {
			operatorAddress, err := sdk.AccAddressFromBech32(operator.GetAddress())
			if err != nil {
				return true, err
			}

			operatorBalance := k.bankKeeper.GetAllBalances(ctx, operatorAddress)
			operatorsBalances = operatorsBalances.Add(operatorBalance...)

			operatorsTokens = operatorsTokens.Add(operator.Tokens...)

			return false, nil
		})
		if err != nil {
			panic(err)
		}

		// Get all the services balances and tokens
		servicesBalances := sdk.NewCoins()
		servicesTokens := sdk.NewCoins()
		err = k.servicesKeeper.IterateServices(ctx, func(service servicestypes.Service) (stop bool, err error) {
			serviceAddress, err := sdk.AccAddressFromBech32(service.GetAddress())
			if err != nil {
				panic(err)
			}

			serviceBalance := k.bankKeeper.GetAllBalances(ctx, serviceAddress)
			servicesBalances = servicesBalances.Add(serviceBalance...)

			servicesTokens = servicesTokens.Add(service.Tokens...)

			return false, nil
		})
		if err != nil {
			panic(err)
		}

		// We use IsAllGTE to check that the balances are greater or equal to the tokens
		// This is used because users might have sent tokens to the accounts and if we check using Equals
		// the invariant would be broken
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

		poolDelegations, err := k.GetAllPoolDelegations(ctx)
		if err != nil {
			panic(err)
		}

		for _, delegation := range poolDelegations {
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
		pools, err := k.poolsKeeper.GetPools(ctx)
		if err != nil {
			panic(err)
		}

		poolsDelegatorsShares := map[uint32]sdk.DecCoins{}
		for _, pool := range pools {
			poolsDelegatorsShares[pool.ID] = sdk.NewDecCoins()
		}

		// Iterate through all the pool delegations to calculate the total delegators shares for each pool
		poolDelegations, err := k.GetAllPoolDelegations(ctx)
		if err != nil {
			panic(err)
		}

		for _, delegation := range poolDelegations {
			if _, ok := poolsDelegatorsShares[delegation.TargetID]; !ok {
				poolsDelegatorsShares[delegation.TargetID] = sdk.NewDecCoins()
			}

			poolsDelegatorsShares[delegation.TargetID] = poolsDelegatorsShares[delegation.TargetID].Add(delegation.Shares...)
		}

		for poolID, delegatorsShares := range poolsDelegatorsShares {
			pool, err := k.poolsKeeper.GetPool(ctx, poolID)
			if err != nil {
				if errors.Is(err, collections.ErrNotFound) {
					panic(fmt.Errorf("pool with id %d not found", poolID))
				}
				panic(err)
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

		operatorDelegations, err := k.GetAllOperatorDelegations(ctx)
		if err != nil {
			panic(err)
		}

		for _, delegation := range operatorDelegations {
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
		operators, err := k.operatorsKeeper.GetOperators(ctx)
		if err != nil {
			panic(err)
		}

		operatorsDelegatorsShares := map[uint32]sdk.DecCoins{}
		for _, operator := range operators {
			operatorsDelegatorsShares[operator.ID] = sdk.NewDecCoins()
		}

		// Iterate through all the operator delegations to calculate the total delegators shares for each operator
		operatorDelegations, err := k.GetAllOperatorDelegations(ctx)
		if err != nil {
			panic(err)
		}

		for _, delegation := range operatorDelegations {
			if _, ok := operatorsDelegatorsShares[delegation.TargetID]; !ok {
				operatorsDelegatorsShares[delegation.TargetID] = sdk.NewDecCoins()
			}

			operatorsDelegatorsShares[delegation.TargetID] = operatorsDelegatorsShares[delegation.TargetID].Add(delegation.Shares...)
		}

		for operatorID, delegatorsShares := range operatorsDelegatorsShares {
			operator, err := k.operatorsKeeper.GetOperator(ctx, operatorID)
			if err != nil {
				if errors.Is(err, collections.ErrNotFound) {
					panic(fmt.Errorf("operator with id %d not found", operatorID))
				}
				panic(err)
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

		serviceDelegations, err := k.GetAllServiceDelegations(ctx)
		if err != nil {
			panic(err)
		}

		for _, delegation := range serviceDelegations {
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
		services, err := k.servicesKeeper.GetServices(ctx)
		if err != nil {
			panic(err)
		}
		for _, service := range services {
			servicesDelegatorsShares[service.ID] = sdk.NewDecCoins()
		}

		// Iterate through all the service delegations to calculate the total delegators shares for each service
		serviceDelegations, err := k.GetAllServiceDelegations(ctx)
		if err != nil {
			panic(err)
		}

		for _, delegation := range serviceDelegations {
			if _, ok := servicesDelegatorsShares[delegation.TargetID]; !ok {
				servicesDelegatorsShares[delegation.TargetID] = sdk.NewDecCoins()
			}

			servicesDelegatorsShares[delegation.TargetID] = servicesDelegatorsShares[delegation.TargetID].Add(delegation.Shares...)
		}

		for serviceID, delegatorsShares := range servicesDelegatorsShares {
			service, err := k.servicesKeeper.GetService(ctx, serviceID)
			if err != nil {
				if errors.Is(err, collections.ErrNotFound) {
					panic(fmt.Errorf("service with id %d not found", serviceID))
				}
				panic(err)
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

// AllowedOperatorsExistInvariant checks that all the operators that are allowed to join a service exist
func AllowedOperatorsExistInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// Iterate over all the services joined by operators
		var notFoundOperatorsIDs []uint32
		err := k.IterateAllServicesAllowedOperators(ctx, func(serviceID uint32, operatorID uint32) (stop bool, err error) {
			_, err = k.operatorsKeeper.GetOperator(ctx, operatorID)
			if err != nil {
				if errors.Is(err, collections.ErrNotFound) {
					notFoundOperatorsIDs = append(notFoundOperatorsIDs, operatorID)
				}
				return true, err
			}
			return false, nil
		})
		if err != nil {
			panic(err)
		}

		var msg string
		if len(notFoundOperatorsIDs) > 0 {
			msg = fmt.Sprintf("operators not found: %v\n", notFoundOperatorsIDs)
		}

		return sdk.FormatInvariant(types.ModuleName, "allowed operators exist", fmt.Sprintf(
			"allowed operators exist invariant broken\n%s", msg)), len(notFoundOperatorsIDs) > 0
	}
}

// OperatorsJoinedServicesExistInvariant checks that all the services that are joined by operators exist
func OperatorsJoinedServicesExistInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// Iterate over all the operators joined services
		var notFoundServicesIDs []uint32
		err := k.IterateAllOperatorsJoinedServices(ctx, func(operatorID uint32, serviceID uint32) (stop bool, err error) {
			_, err = k.servicesKeeper.GetService(ctx, serviceID)
			if err != nil {
				if errors.Is(err, collections.ErrNotFound) {
					notFoundServicesIDs = append(notFoundServicesIDs, serviceID)
				}
				return false, err
			}
			return false, nil
		})
		if err != nil {
			panic(err)
		}

		var msg string
		if len(notFoundServicesIDs) > 0 {
			msg = fmt.Sprintf("services not found: %v\n", notFoundServicesIDs)
		}

		return sdk.FormatInvariant(types.ModuleName, "joined services exist", fmt.Sprintf(
			"joined services exist invariant broken\n%s", msg)), len(notFoundServicesIDs) > 0
	}
}

func TotalRestakedAssetsInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// Get all the restaked assets
		totalRestakedAssets, err := k.GetTotalRestakedAssets(ctx)
		if err != nil {
			panic(err)
		}

		// Helper function to get restaked assets from specific type of delegations
		getRestakedAssets := func(iterFn func(context.Context, func(types.Delegation) (bool, error)) error) sdk.Coins {
			restakedAssets := sdk.NewCoins()
			targets := map[uint32]types.DelegationTarget{}
			err := iterFn(ctx, func(del types.Delegation) (stop bool, err error) {
				target, ok := targets[del.TargetID]
				if !ok {
					target, err = k.GetDelegationTargetFromDelegation(ctx, del)
					if err != nil {
						return true, err
					}
					targets[del.TargetID] = target
				}

				tokens := target.TokensFromSharesTruncated(del.Shares)
				tokensTruncated, _ := tokens.TruncateDecimal()
				restakedAssets = restakedAssets.Add(tokensTruncated...)
				return false, nil
			})
			if err != nil {
				panic(err)
			}
			return restakedAssets
		}

		totalRestakedAssetsFromDels := sdk.NewCoins()
		totalRestakedAssetsFromDels = totalRestakedAssetsFromDels.Add(
			getRestakedAssets(k.IterateAllPoolDelegations)...)
		totalRestakedAssetsFromDels = totalRestakedAssetsFromDels.Add(
			getRestakedAssets(k.IterateAllOperatorDelegations)...)
		totalRestakedAssetsFromDels = totalRestakedAssetsFromDels.Add(
			getRestakedAssets(k.IterateAllServiceDelegations)...)

		// Check if any of the total restaked assets calculated from delegations is
		// greater than the total restaked assets stored. Note that because of the
		// truncation, the total restaked assets calculated from delegations should be
		// less than or equal to the total restaked assets stored.
		broken := totalRestakedAssetsFromDels.IsAnyGT(totalRestakedAssets)

		return sdk.FormatInvariant(types.ModuleName, "total restaked assets", fmt.Sprintf(
			"total restaked assets invariant broken\n"+
				"total restaked assets stored: %v\n"+
				"total restaked assets calculated from delegations: %v\n",
			totalRestakedAssets, totalRestakedAssetsFromDels)), broken
	}
}
