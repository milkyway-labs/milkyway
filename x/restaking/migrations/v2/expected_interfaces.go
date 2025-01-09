package v2

import (
	"context"

	restakingtypes "github.com/milkyway-labs/milkyway/v7/x/restaking/types"
)

type Keeper interface {
	IterateAllOperatorsJoinedServices(ctx context.Context, action func(operatorID uint32, serviceID uint32) (stop bool, err error)) error
	CanOperatorValidateService(ctx context.Context, serviceID uint32, operatorID uint32) (bool, error)
	RemoveServiceFromOperatorJoinedServices(ctx context.Context, operatorID uint32, serviceID uint32) error

	GetAllPoolDelegations(ctx context.Context) ([]restakingtypes.Delegation, error)
	GetAllOperatorDelegations(ctx context.Context) ([]restakingtypes.Delegation, error)
	GetAllServiceDelegations(ctx context.Context) ([]restakingtypes.Delegation, error)
}
