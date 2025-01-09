package cli

import (
	"fmt"
	"strings"

	poolstypes "github.com/milkyway-labs/milkyway/v7/x/pools/types"
	"github.com/milkyway-labs/milkyway/v7/x/restaking/types"
	servicestypes "github.com/milkyway-labs/milkyway/v7/x/services/types"
)

// ParseTrustedServiceEntry parses a string into a TrustedServiceEntry.
// The value provided must be in the format: "<serviceID>-<poolID>,<poolID>,<poolID>"
func ParseTrustedServiceEntry(value string) (types.TrustedServiceEntry, error) {
	var serviceIDString, poolsIDsStrings string
	_, err := fmt.Sscanf(value, "%s-%s", &serviceIDString, &poolsIDsStrings)
	if err != nil {
		return types.TrustedServiceEntry{}, err
	}

	serviceID, err := servicestypes.ParseServiceID(serviceIDString)
	if err != nil {
		return types.TrustedServiceEntry{}, err
	}

	var poolsIDs []uint32
	for _, poolID := range strings.Split(poolsIDsStrings, ",") {
		parsedPoolID, err := poolstypes.ParsePoolID(poolID)
		if err != nil {
			return types.TrustedServiceEntry{}, err
		}
		poolsIDs = append(poolsIDs, parsedPoolID)
	}

	return types.NewTrustedServiceEntry(serviceID, poolsIDs), nil
}
