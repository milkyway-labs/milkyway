package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/x/restaking/types"
)

type RestakingKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
}
