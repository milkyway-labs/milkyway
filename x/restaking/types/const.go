package types

const (
	// BaseDelegationGasCost is the gas cost for a delegation transaction that will be paid per each delegation
	// already existing for the same target. This is used to implement delegations count-based scaling costs.
	BaseDelegationGasCost uint64 = 20_000

	// BaseDelegationDenomCost is the gas cost for storing or deleting a coin denom for each delegation.
	// Examples:
	// * if a user wants to create a new delegation with 3 denoms, they will be charged 3 * BaseDelegationDenomCost
	// * if a user wants to undelegate from a delegation with 2 denoms, they will be charged 2 * BaseDelegationDenomCost
	BaseDelegationDenomCost uint64 = 20_000
)
