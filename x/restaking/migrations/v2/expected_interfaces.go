package v2

import (
	"context"
)

type Keeper interface {
	IterateAllOperatorsJoinedServices(ctx context.Context, action func(operatorID uint32, serviceID uint32) (stop bool, err error)) error
	CanOperatorValidateService(ctx context.Context, serviceID uint32, operatorID uint32) (bool, error)
	RemoveServiceFromOperatorJoinedServices(ctx context.Context, operatorID uint32, serviceID uint32) error
}
