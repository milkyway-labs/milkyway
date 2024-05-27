package keeper

import (
	"github.com/milkyway-labs/milkyway/x/records/types"
)

var _ types.QueryServer = Keeper{}
