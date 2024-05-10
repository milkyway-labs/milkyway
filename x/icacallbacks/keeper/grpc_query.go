package keeper

import (
	"github.com/milkyway-labs/milk/x/icacallbacks/types"
)

var _ types.QueryServer = Keeper{}
