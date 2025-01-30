package cli

import (
	"strings"

	poolstypes "github.com/milkyway-labs/milkyway/v8/x/pools/types"
	"github.com/milkyway-labs/milkyway/v8/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v8/x/services/types"
)

// ParseTrustedServiceEntry parses a string into a TrustedServiceEntry.
// The value provided must be in the format: "<serviceID>-<poolID>,<poolID>,<poolID>"
func ParseTrustedServiceEntry(value string) (types.TrustedServiceEntry, error) {
	parts := strings.SplitN(value, "-", 2)

	serviceIDString := parts[0]
	serviceID, err := servicestypes.ParseServiceID(serviceIDString)
	if err != nil {
		return types.TrustedServiceEntry{}, err
	}

	var poolsIDs []uint32
	if len(parts) == 2 {
		poolsIDsStrings := parts[1]
		for _, poolID := range strings.Split(poolsIDsStrings, ",") {
			parsedPoolID, err := poolstypes.ParsePoolID(poolID)
			if err != nil {
				return types.TrustedServiceEntry{}, err
			}
			poolsIDs = append(poolsIDs, parsedPoolID)
		}
	}

	return types.NewTrustedServiceEntry(serviceID, poolsIDs), nil
}
