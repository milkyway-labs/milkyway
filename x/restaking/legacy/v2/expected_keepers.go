package v2

import (
	"context"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

type OperatorsKeeper interface {
	GetOperator(ctx context.Context, operatorID uint32) (operatorstypes.Operator, bool, error)
	SaveOperatorParams(ctx context.Context, operatorID uint32, params operatorstypes.OperatorParams) error
}

type ServicesKeeper interface {
	GetService(ctx context.Context, serviceID uint32) (servicestypes.Service, bool, error)
}

type RestakingKeeper interface {
	AddServiceToOperatorJoinedServices(ctx context.Context, operatorID uint32, serviceID uint32) error
	AddOperatorToServiceAllowList(ctx context.Context, serviceID uint32, operatorID uint32) error
	AddPoolToServiceSecuringPools(ctx context.Context, serviceID uint32, poolID uint32) error
}
