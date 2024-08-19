package types

import (
	"cosmossdk.io/errors"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	poolstypes "github.com/milkyway-labs/milkyway/x/pools/types"
	restakingtypes "github.com/milkyway-labs/milkyway/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

// GetDelegationTargetType returns the delegation target's type.
func GetDelegationTargetType(target restakingtypes.DelegationTarget) (restakingtypes.DelegationType, error) {
	switch target.(type) {
	case *poolstypes.Pool:
		return restakingtypes.DELEGATION_TYPE_POOL, nil
	case *operatorstypes.Operator:
		return restakingtypes.DELEGATION_TYPE_OPERATOR, nil
	case *servicestypes.Service:
		return restakingtypes.DELEGATION_TYPE_SERVICE, nil
	default:
		return restakingtypes.DELEGATION_TYPE_UNSPECIFIED, errors.Wrapf(
			restakingtypes.ErrInvalidDelegationType, "invalid delegation target type: %T", target,
		)
	}
}
