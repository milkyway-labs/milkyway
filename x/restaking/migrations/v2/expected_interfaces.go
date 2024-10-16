package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

type OperatorsKeeper interface {
	GetOperator(ctx sdk.Context, operatorID uint32) (operatorstypes.Operator, bool)
	SaveOperatorParams(ctx sdk.Context, operatorID uint32, params operatorstypes.OperatorParams) error
}

type RestakingKeeper interface {
	GetOperatorSecuredServices(ctx sdk.Context, operatorID uint32) (types.OperatorSecuredServices, error)
	SetOperatorSecuredServices(ctx sdk.Context, operatorID uint32, securedServices types.OperatorSecuredServices) error
}
