package types

import (
	"bytes"
	"fmt"

	"github.com/milkyway-labs/milkyway/utils"
	operatorstypes "github.com/milkyway-labs/milkyway/x/operators/types"
	servicestypes "github.com/milkyway-labs/milkyway/x/services/types"
)

var (
	// OperatorParamsPrefix is the prefix used to store the operator params
	// This has been replaced by OperatorServicesPrefix that is
	// used to store the services secured by an operator, the operator params
	// instead have been moved to the x/operators module.
	OperatorParamsPrefix = []byte{0x11}

	// ServiceParamsPrefix is the prefix used to store the service params
	ServiceParamsPrefix = []byte{0x12}
)

// OperatorParamsStoreKey returns the key used to store the operator params
// The operator params are stored in the x/operator module, now
// in this module we only keep the list of services secured by a operator.
func OperatorParamsStoreKey(operatorID uint32) []byte {
	return utils.CompositeKey(OperatorParamsPrefix, operatorstypes.GetOperatorIDBytes(operatorID))
}

// ParseOperatorParamsKey parses the operator ID from the given key
// The operator params are stored in the x/operator module, now
// in this module we only keep the list of services secured by a operator.
func ParseOperatorParamsKey(bz []byte) (operatorID uint32, err error) {
	bz = bytes.TrimPrefix(bz, OperatorParamsPrefix)
	if len(bz) != 4 {
		return 0, fmt.Errorf("invalid key length; expected: 4, got: %d", len(bz))
	}

	return operatorstypes.GetOperatorIDFromBytes(bz), nil
}

// ServiceParamsStoreKey returns the key used to store the service params
func ServiceParamsStoreKey(serviceID uint32) []byte {
	return utils.CompositeKey(ServiceParamsPrefix, servicestypes.GetServiceIDBytes(serviceID))
}

// ParseServiceParamsKey parses the service ID from the given key
func ParseServiceParamsKey(bz []byte) (serviceID uint32, err error) {
	bz = bytes.TrimPrefix(bz, ServiceParamsPrefix)
	if len(bz) != 4 {
		return 0, fmt.Errorf("invalid key length; expected: 4, got: %d", len(bz))
	}

	return servicestypes.GetServiceIDFromBytes(bz), nil
}
