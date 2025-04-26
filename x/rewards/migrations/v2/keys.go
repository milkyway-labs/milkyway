package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v12/x/rewards/types"
)

func PlanStoreKey(id uint64) []byte {
	return append(types.RewardsPlanKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
