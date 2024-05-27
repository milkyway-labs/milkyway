package keeper

import (
	"github.com/milkyway-labs/milkyway/x/icacallbacks/types"
)

var _ types.QueryServer = Keeper{}
