package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	"github.com/milkyway-labs/milkyway/x/rewards/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// RegisterInvariants registers all module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k *Keeper) {
	ir.RegisterRoute(types.ModuleName, "nonnegative-outstanding",
		NonNegativeOutstandingInvariant(k))
	ir.RegisterRoute(types.ModuleName, "can-withdraw",
		CanWithdrawInvariant(k))
	ir.RegisterRoute(types.ModuleName, "reference-count",
		ReferenceCountInvariant(k))
	ir.RegisterRoute(types.ModuleName, "module-account",
		ModuleAccountInvariant(k))
}

// AllInvariants runs all invariants of the module
func AllInvariants(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := CanWithdrawInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		res, stop = NonNegativeOutstandingInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		res, stop = ReferenceCountInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		return ModuleAccountInvariant(k)(ctx)
	}
}

// NonNegativeOutstandingInvariant checks that outstanding unwithdrawn fees are never negative
func NonNegativeOutstandingInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int

		err := k.PoolOutstandingRewards.Walk(ctx, nil, func(poolID uint32, rewards types.OutstandingRewards) (stop bool, err error) {
			outstanding := rewards.Rewards
			if outstanding.IsAnyNegative() {
				count++
				msg += fmt.Sprintf("\tpool %d has negative outstanding coins: %v\n", poolID, outstanding)
			}
			return false, nil
		})
		if err != nil {
			return sdk.FormatInvariant(types.ModuleName, "nonnegative outstanding", err.Error()), true
		}
		err = k.OperatorOutstandingRewards.Walk(ctx, nil, func(operatorID uint32, rewards types.OutstandingRewards) (stop bool, err error) {
			outstanding := rewards.Rewards
			if outstanding.IsAnyNegative() {
				count++
				msg += fmt.Sprintf("\toperator %d has negative outstanding coins: %v\n", operatorID, outstanding)
			}
			return false, nil
		})
		if err != nil {
			return sdk.FormatInvariant(types.ModuleName, "nonnegative outstanding", err.Error()), true
		}
		err = k.ServiceOutstandingRewards.Walk(ctx, nil, func(serviceID uint32, rewards types.OutstandingRewards) (stop bool, err error) {
			outstanding := rewards.Rewards
			if outstanding.IsAnyNegative() {
				count++
				msg += fmt.Sprintf("\tservice %d has negative outstanding coins: %v\n", serviceID, outstanding)
			}
			return false, nil
		})
		if err != nil {
			return sdk.FormatInvariant(types.ModuleName, "nonnegative outstanding", err.Error()), true
		}
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "nonnegative outstanding",
			fmt.Sprintf("found %d delegation targets with negative outstanding rewards\n%s", count, msg)), broken
	}
}

// CanWithdrawInvariant checks that current rewards can be completely withdrawn
func CanWithdrawInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

		var remaining types.DecPools

		poolDelegationAddrs := make(map[uint32][]sdk.AccAddress)
		err := k.restakingKeeper.IterateAllPoolDelegations(ctx, func(del restakingtypes.Delegation) (stop bool, err error) {
			delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
			if err != nil {
				return true, err
			}

			poolDelegationAddrs[del.TargetID] = append(poolDelegationAddrs[del.TargetID], delAddr)
			return false, nil
		})
		if err != nil {
			panic(err)
		}

		operatorDelegationAddrs := make(map[uint32][]sdk.AccAddress)
		err = k.restakingKeeper.IterateAllOperatorDelegations(ctx, func(del restakingtypes.Delegation) (stop bool, err error) {
			delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
			if err != nil {
				return true, err
			}

			operatorDelegationAddrs[del.TargetID] = append(operatorDelegationAddrs[del.TargetID], delAddr)
			return false, nil
		})
		if err != nil {
			panic(err)
		}

		serviceDelegationAddrs := make(map[uint32][]sdk.AccAddress)
		err = k.restakingKeeper.IterateAllServiceDelegations(ctx, func(del restakingtypes.Delegation) (stop bool, err error) {
			delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
			if err != nil {
				return true, err
			}

			serviceDelegationAddrs[del.TargetID] = append(serviceDelegationAddrs[del.TargetID], delAddr)
			return false, nil
		})
		if err != nil {
			panic(err)
		}

		// iterate over all pools
		err = k.poolsKeeper.IteratePools(ctx, func(pool poolstypes.Pool) (stop bool, err error) {
			target, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_POOL, pool.ID)
			if err != nil {
				return true, err
			}

			delegationAddrs, ok := poolDelegationAddrs[pool.ID]
			if ok {
				for _, delAddr := range delegationAddrs {
					if _, err := k.WithdrawDelegationRewards(ctx, delAddr, target); err != nil {
						panic(err)
					}
				}
			}
			remaining, err = k.GetOutstandingRewardsCoins(ctx, target)
			if err != nil {
				return true, err
			}

			if remaining.IsAnyNegative() {
				return true, nil
			}

			return false, nil
		})
		if err != nil {
			panic(err)
		}

		broken := remaining.IsAnyNegative()
		if broken {
			return sdk.FormatInvariant(types.ModuleName, "can withdraw",
				fmt.Sprintf("pools remaining coins: %v\n", remaining)), true
		}

		// iterate over all operators
		err = k.operatorsKeeper.IterateOperators(ctx, func(operator operatorstypes.Operator) (stop bool, err error) {
			target, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operator.ID)
			if err != nil {
				return true, err
			}

			delegationAddrs, ok := operatorDelegationAddrs[operator.ID]
			if ok {
				for _, delAddr := range delegationAddrs {
					if _, err := k.WithdrawDelegationRewards(ctx, delAddr, target); err != nil {
						panic(err)
					}
				}
			}

			remaining, err = k.GetOutstandingRewardsCoins(ctx, target)
			if err != nil {
				return true, err
			}

			if remaining.IsAnyNegative() {
				return true, nil
			}

			return false, nil
		})
		if err != nil {
			panic(err)
		}

		broken = remaining.IsAnyNegative()
		if broken {
			return sdk.FormatInvariant(types.ModuleName, "can withdraw",
				fmt.Sprintf("operators remaining coins: %v\n", remaining)), true
		}

		// iterate over all services
		err = k.servicesKeeper.IterateServices(ctx, func(service servicestypes.Service) (stop bool, err error) {
			target, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, service.ID)
			if err != nil {
				return true, err
			}

			delegationAddrs, ok := serviceDelegationAddrs[service.ID]
			if ok {
				for _, delAddr := range delegationAddrs {
					if _, err := k.WithdrawDelegationRewards(ctx, delAddr, target); err != nil {
						panic(err)
					}
				}
			}
			remaining, err = k.GetOutstandingRewardsCoins(ctx, target)
			if err != nil {
				return true, err
			}

			if remaining.IsAnyNegative() {
				return true, nil
			}

			return false, nil
		})
		if err != nil {
			panic(err)
		}

		broken = remaining.IsAnyNegative()
		return sdk.FormatInvariant(types.ModuleName, "can withdraw",
			fmt.Sprintf("services remaining coins: %v\n", remaining)), broken
	}
}

// --------------------------------------------------------------------------------------------------------------------

// ReferenceCountInvariant checks that the number of historical rewards records is correct
func ReferenceCountInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// Check the reference count for pools
		msg, broken := checkReferencesCount(
			ctx,
			restakingtypes.DELEGATION_TYPE_POOL,
			k.poolsKeeper.IteratePools,
			k.restakingKeeper.IterateAllPoolDelegations,
			k.PoolHistoricalRewards,
		)
		if broken {
			return sdk.FormatInvariant(types.ModuleName, "reference count", msg), broken
		}

		// Check the reference count for operators
		msg, broken = checkReferencesCount(
			ctx,
			restakingtypes.DELEGATION_TYPE_OPERATOR,
			k.operatorsKeeper.IterateOperators,
			k.restakingKeeper.IterateAllOperatorDelegations,
			k.OperatorHistoricalRewards,
		)
		if broken {
			return sdk.FormatInvariant(types.ModuleName, "reference count", msg), broken
		}

		// Check the reference count for services
		msg, broken = checkReferencesCount(
			ctx,
			restakingtypes.DELEGATION_TYPE_SERVICE,
			k.servicesKeeper.IterateServices,
			k.restakingKeeper.IterateAllServiceDelegations,
			k.ServiceHistoricalRewards,
		)
		if broken {
			return sdk.FormatInvariant(types.ModuleName, "reference count", msg), broken
		}

		return "", false
	}
}

// checkReferencesCount checks the reference count for a given delegation target type
func checkReferencesCount[T any](
	ctx sdk.Context,
	delegationTargetType restakingtypes.DelegationType,
	targetsIterator func(ctx context.Context, fn func(T) (bool, error)) error,
	delegationsIterator func(ctx context.Context, fn func(restakingtypes.Delegation) (bool, error)) error,
	historicalRewardsCollection collections.Map[collections.Pair[uint32, uint64], types.HistoricalRewards],
) (msg string, broken bool) {

	targetCount := uint64(0)
	err := targetsIterator(ctx, func(_ T) (bool, error) {
		targetCount++
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	delegationsCount := uint64(0)
	err = delegationsIterator(ctx, func(_ restakingtypes.Delegation) (bool, error) {
		delegationsCount++
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	referencesCount := uint64(0)
	err = historicalRewardsCollection.Walk(ctx, nil, func(key collections.Pair[uint32, uint64], value types.HistoricalRewards) (stop bool, err error) {
		referencesCount += uint64(value.ReferenceCount)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	// Make sure we have one record per delegation target (last tracked period) and
	// one record per delegation (previous period)
	expected := targetCount + delegationsCount

	broken = referencesCount != expected

	return fmt.Sprintf("expected historical reference count: %d = %v delegation targets + %v delegations\n"+
		"total %s historical reference count: %d\n",
		expected, targetCount, delegationsCount, delegationTargetType, referencesCount,
	), broken
}

// --------------------------------------------------------------------------------------------------------------------

// ModuleAccountInvariant checks that the coins held by the global rewards pool
// is consistent with the sum of outstanding rewards
func ModuleAccountInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var expectedCoins sdk.DecCoins
		err := k.PoolOutstandingRewards.Walk(ctx, nil, func(_ uint32, rewards types.OutstandingRewards) (stop bool, err error) {
			expectedCoins = expectedCoins.Add(rewards.Rewards.Sum()...)
			return false, nil
		})
		if err != nil {
			return sdk.FormatInvariant(types.ModuleName, "module account coins", err.Error()), true
		}

		expectedInt, _ := expectedCoins.TruncateDecimal()

		rewardsPoolAddr := k.accountKeeper.GetModuleAddress(types.RewardsPoolName)
		balances := k.bankKeeper.GetAllBalances(ctx, rewardsPoolAddr)
		broken := !balances.IsAllGTE(expectedInt)
		return sdk.FormatInvariant(
			types.ModuleName, "ModuleAccount coins",
			fmt.Sprintf("\texpected ModuleAccount coins:     %s\n"+
				"\tdistribution ModuleAccount coins: %s\n",
				expectedInt, balances,
			),
		), broken
	}
}
