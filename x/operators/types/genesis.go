package types

import (
	"fmt"

	"github.com/milkyway-labs/milkyway/utils"
)

func NewGenesisState(nextOperatorID uint32, operators []Operator, params Params) *GenesisState {
	return &GenesisState{
		NextOperatorID: nextOperatorID,
		Operators:      operators,
		Params:         params,
	}
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return NewGenesisState(1, nil, DefaultParams())
}

// Validate checks that the genesis state is valid.
func (data *GenesisState) Validate() error {
	if data.NextOperatorID == 0 {
		return fmt.Errorf("invalid next operator ID: %d", data.NextOperatorID)
	}

	// Check for duplicate operators
	if duplicate := findDuplicateOperators(data.Operators); duplicate != nil {
		return fmt.Errorf("duplicated operator: %d", duplicate.ID)
	}

	// Validate operators
	for _, operator := range data.Operators {
		err := operator.Validate()
		if err != nil {
			return fmt.Errorf("invalid operator with id %d: %s", operator.ID, err)
		}
	}

	// Validate params
	err := data.Params.Validate()
	if err != nil {
		return fmt.Errorf("invalid params: %s", err)
	}

	return nil
}

// findDuplicateOperators returns the first duplicated operator in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicateOperators(operators []Operator) *Operator {
	return utils.FindDuplicate(operators, func(a, b Operator) bool {
		return a.ID == b.ID
	})
}
