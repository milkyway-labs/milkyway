package types

import (
	"fmt"
	"time"

	"github.com/milkyway-labs/milkyway/utils"
)

func NewGenesisState(
	nextOperatorID uint32,
	operators []Operator,
	unbondingOperators []UnbondingOperator,
	params Params,
) *GenesisState {
	return &GenesisState{
		NextOperatorID:     nextOperatorID,
		Operators:          operators,
		UnbondingOperators: unbondingOperators,
		Params:             params,
	}
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return NewGenesisState(1, nil, nil, DefaultParams())
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

	// Check for duplicated unbonding operators
	if duplicate := findDuplicateUnbondingOperators(data.UnbondingOperators); duplicate != nil {
		return fmt.Errorf("duplicated unbonding operator: %d", duplicate.OperatorID)
	}

	// Validate unbonding operators
	for _, operator := range data.UnbondingOperators {
		err := operator.Validate()
		if err != nil {
			return fmt.Errorf("invalid unbonding operator with id %d: %s", operator.OperatorID, err)
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

// findDuplicateUnbondingOperators returns the first duplicated unbonding operator in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicateUnbondingOperators(operators []UnbondingOperator) *UnbondingOperator {
	return utils.FindDuplicate(operators, func(a, b UnbondingOperator) bool {
		return a.OperatorID == b.OperatorID
	})
}

// --------------------------------------------------------------------------------------------------------------------

// NewUnbondingOperator creates a new UnbondingOperator instance.
func NewUnbondingOperator(operatorID uint32, completionTime time.Time) UnbondingOperator {
	return UnbondingOperator{
		OperatorID:              operatorID,
		UnbondingCompletionTime: completionTime,
	}
}

// Validate checks that the UnbondingOperator has valid values.
func (o *UnbondingOperator) Validate() error {
	if o.OperatorID == 0 {
		return fmt.Errorf("invalid operator ID: %d", o.OperatorID)
	}

	if o.UnbondingCompletionTime.IsZero() {
		return fmt.Errorf("invalid unbond completion time: %s", o.UnbondingCompletionTime)
	}

	return nil
}
