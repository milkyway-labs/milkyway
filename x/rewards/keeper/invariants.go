package keeper

import (
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
		k.restakingKeeper.IterateAllPoolDelegations(ctx, func(del restakingtypes.Delegation) (stop bool) {
			delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
			if err != nil {
				panic(err)
			}
			poolID := del.TargetID
			poolDelegationAddrs[poolID] = append(poolDelegationAddrs[poolID], delAddr)
			return false
		})
		operatorDelegationAddrs := make(map[uint32][]sdk.AccAddress)
		k.restakingKeeper.IterateAllOperatorDelegations(ctx, func(del restakingtypes.Delegation) (stop bool) {
			delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
			if err != nil {
				panic(err)
			}
			operatorID := del.TargetID
			operatorDelegationAddrs[operatorID] = append(operatorDelegationAddrs[operatorID], delAddr)
			return false
		})
		serviceDelegationAddrs := make(map[uint32][]sdk.AccAddress)
		k.restakingKeeper.IterateAllServiceDelegations(ctx, func(del restakingtypes.Delegation) (stop bool) {
			delAddr, err := k.accountKeeper.AddressCodec().StringToBytes(del.UserAddress)
			if err != nil {
				panic(err)
			}
			serviceID := del.TargetID
			serviceDelegationAddrs[serviceID] = append(serviceDelegationAddrs[serviceID], delAddr)
			return false
		})

		// iterate over all pools
		k.poolsKeeper.IteratePools(ctx, func(pool poolstypes.Pool) (stop bool) {
			target, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_POOL, pool.ID)
			if err != nil {
				panic(err)
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
				panic(err)
			}
			if remaining.IsAnyNegative() {
				return true
			}
			return false
		})
		broken := remaining.IsAnyNegative()
		if broken {
			return sdk.FormatInvariant(types.ModuleName, "can withdraw",
				fmt.Sprintf("pools remaining coins: %v\n", remaining)), true
		}

		// iterate over all operators
		k.operatorsKeeper.IterateOperators(ctx, func(operator operatorstypes.Operator) (stop bool) {
			target, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_OPERATOR, operator.ID)
			if err != nil {
				panic(err)
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
				panic(err)
			}
			if remaining.IsAnyNegative() {
				return true
			}
			return false
		})
		broken = remaining.IsAnyNegative()
		if broken {
			return sdk.FormatInvariant(types.ModuleName, "can withdraw",
				fmt.Sprintf("operators remaining coins: %v\n", remaining)), true
		}

		// iterate over all services
		k.servicesKeeper.IterateServices(ctx, func(service servicestypes.Service) (stop bool) {
			target, err := k.GetDelegationTarget(ctx, restakingtypes.DELEGATION_TYPE_SERVICE, service.ID)
			if err != nil {
				panic(err)
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
				panic(err)
			}
			if remaining.IsAnyNegative() {
				return true
			}
			return false
		})
		broken = remaining.IsAnyNegative()
		return sdk.FormatInvariant(types.ModuleName, "can withdraw",
			fmt.Sprintf("services remaining coins: %v\n", remaining)), broken
	}
}

// ReferenceCountInvariant checks that the number of historical rewards records is correct
func ReferenceCountInvariant(k *Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		targetCount := uint64(0)
		k.poolsKeeper.IteratePools(ctx, func(_ poolstypes.Pool) (stop bool) {
			targetCount++
			return false
		})
		k.operatorsKeeper.IterateOperators(ctx, func(_ operatorstypes.Operator) (stop bool) {
			targetCount++
			return false
		})
		k.servicesKeeper.IterateServices(ctx, func(_ servicestypes.Service) (stop bool) {
			targetCount++
			return false
		})

		delCount := uint64(0)
		k.restakingKeeper.IterateAllPoolDelegations(ctx, func(_ restakingtypes.Delegation) (stop bool) {
			delCount++
			return false
		})
		k.restakingKeeper.IterateAllOperatorDelegations(ctx, func(_ restakingtypes.Delegation) (stop bool) {
			delCount++
			return false
		})
		k.restakingKeeper.IterateAllServiceDelegations(ctx, func(_ restakingtypes.Delegation) (stop bool) {
			delCount++
			return false
		})

		// one record per delegation target (last tracked period), one record per
		// delegation (previous period)
		// TODO: handle slash events
		expected := targetCount + delCount
		count := uint64(0)
		err := k.PoolHistoricalRewards.Walk(
			ctx, nil, func(key collections.Pair[uint32, uint64], rewards types.HistoricalRewards) (stop bool, err error) {
				count += uint64(rewards.ReferenceCount)
				return false, nil
			},
		)
		if err != nil {
			panic(err)
		}

		broken := count != expected

		return sdk.FormatInvariant(types.ModuleName, "reference count",
			fmt.Sprintf("expected historical reference count: %d = %v delegation targets + %v delegations\n"+
				"total validator historical reference count: %d\n",
				expected, targetCount, delCount, count)), broken
	}
}

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

		balances := k.bankKeeper.GetAllBalances(ctx, types.RewardsPoolAddress)
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
