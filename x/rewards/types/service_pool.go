package types

import (
	"cmp"
	"fmt"
)

// ServicePools defines denom and sdk.DecCoins wrapper to represents
// rewards pools for multi-token staking
type ServicePools []ServicePool

// Add adds two sets of ServicePools
func (pools ServicePools) Add(poolsB ...ServicePool) ServicePools {
	return pools.safeAdd(poolsB)
}

// Add will perform addition of two ServicePools sets.
func (pools ServicePools) safeAdd(poolsB ServicePools) ServicePools {
	sum := ([]ServicePool)(nil)
	indexA, indexB := 0, 0
	lenA, lenB := len(pools), len(poolsB)

	for {
		if indexA == lenA {
			if indexB == lenB {
				// return nil pools if both sets are empty
				return sum
			}

			// return set B (excluding zero pools) if set A is empty
			return append(sum, removeZeroServicePools(poolsB[indexB:])...)
		} else if indexB == lenB {
			// return set A (excluding zero pools) if set B is empty
			return append(sum, removeZeroServicePools(pools[indexA:])...)
		}

		poolA, poolB := pools[indexA], poolsB[indexB]

		switch cmp.Compare(poolA.ServiceID, poolB.ServiceID) {
		case -1: // pool A service ID < pool B service ID
			if !poolA.IsEmpty() {
				sum = append(sum, poolA)
			}

			indexA++

		case 0: // pool A service ID == pool B service ID
			res := poolA.Add(poolB)
			if !res.IsEmpty() {
				sum = append(sum, res)
			}

			indexA++
			indexB++

		case 1: // pool A service ID > pool B service ID
			if !poolB.IsEmpty() {
				sum = append(sum, poolB)
			}

			indexB++
		}
	}
}

// Sub subtracts a set of ServicePools from another (adds the inverse).
func (pools ServicePools) Sub(poolsB ServicePools) ServicePools {
	diff, hasNeg := pools.SafeSub(poolsB)
	if hasNeg {
		panic("negative pool coins")
	}

	return diff
}

// SafeSub performs the same arithmetic as Sub but returns a boolean if any
// negative ServicePool coins amount was returned.
func (pools ServicePools) SafeSub(poolsB ServicePools) (ServicePools, bool) {
	diff := pools.safeAdd(poolsB.negative())
	return diff, diff.IsAnyNegative()
}

// IsAnyNegative returns true if there is at least one coin whose amount
// is negative; returns false otherwise. It returns false if the ServicePools set
// is empty too.
func (pools ServicePools) IsAnyNegative() bool {
	for _, pool := range pools {
		if pool.DecPools.IsAnyNegative() {
			return true
		}
	}

	return false
}

// negative returns a set of coins with all amount negative.
func (pools ServicePools) negative() ServicePools {
	res := make([]ServicePool, 0, len(pools))
	for _, pool := range pools {
		res = append(res, ServicePool{
			ServiceID: pool.ServiceID,
			DecPools:  pool.DecPools.negative(),
		})
	}
	return res
}

// String implements the Stringer interface for ServicePools. It returns a
// human-readable representation of ServicePools.
func (pools ServicePools) String() string {
	if len(pools) == 0 {
		return ""
	}

	out := ""
	for _, pool := range pools {
		out += fmt.Sprintf("%v,", pool.String())
	}

	return out[:len(out)-1]
}

func removeZeroServicePools(pools ServicePools) ServicePools {
	result := make([]ServicePool, 0, len(pools))

	for _, pool := range pools {
		if !pool.IsEmpty() {
			result = append(result, pool)
		}
	}

	return result
}

//-----------------------------------------------------------------------------
// ServicePool functions

// NewServicePool return new ServicePool instance
func NewServicePool(serviceID uint32, pools DecPools) ServicePool {
	return ServicePool{serviceID, pools}
}

// IsEmpty returns whether the pool coins are empty or not
func (pool ServicePool) IsEmpty() bool {
	return pool.DecPools.IsEmpty()
}

// Add adds amounts of two ServicePool with same service ID.
func (pool ServicePool) Add(poolB ServicePool) ServicePool {
	if pool.ServiceID != poolB.ServiceID {
		panic(fmt.Sprintf("service ID different: %v %v\n", pool.ServiceID, poolB.ServiceID))
	}
	return ServicePool{pool.ServiceID, pool.DecPools.Add(poolB.DecPools...)}
}
