package keeper

import (
	"github.com/milkyway-labs/milk/x/records/types"
)

var _ types.QueryServer = Keeper{}
