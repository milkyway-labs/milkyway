package types

import (
	"fmt"
	"time"

	"github.com/milkyway-labs/milkyway/v4/utils"
)

func NewGenesisState(
	nextOperatorID uint32,
	operators []Operator,
	operatorParams []OperatorParamsRecord,
	unbondingOperators []UnbondingOperator,
	params Params,
) *GenesisState {
	return &GenesisState{
		NextOperatorID:     nextOperatorID,
		Operators:          operators,
		OperatorsParams:    operatorParams,
		UnbondingOperators: unbondingOperators,
		Params:             params,
	}
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return NewGenesisState(1, nil, nil, nil, DefaultParams())
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

	// Check for duplicate operator params
	if duplicate := findDuplicateOperatorsParams(data.OperatorsParams); duplicate != nil {
		return fmt.Errorf("duplicated operator params for operator: %d", duplicate.OperatorID)
	}

	// Validate the operator params
	for _, operatorParams := range data.OperatorsParams {
		err := operatorParams.Validate()
		if err != nil {
			return fmt.Errorf("invalid operator params for operator %d: %s", operatorParams.OperatorID, err)
		}
	}

	// Check for duplicated unbonding operators
	if duplicate := findDuplicateUnbondingOperators(data.UnbondingOperators); duplicate != nil {
		return fmt.Errorf("duplicated unbonding operator: %d", duplicate.OperatorID)
	}

	// Validate unbonding operators
	for _, unbondingOperator := range data.UnbondingOperators {
		err := unbondingOperator.Validate()
		if err != nil {
			return fmt.Errorf("invalid unbonding operator with id %d: %s", unbondingOperator.OperatorID, err)
		}

		// Make sure the operator status is inactivating
		operator, found := utils.Find(data.Operators, func(operator Operator) bool {
			return operator.ID == unbondingOperator.OperatorID
		})

		if !found {
			return fmt.Errorf("unbonding operator with id %d not found", unbondingOperator.OperatorID)
		}

		if operator.Status != OPERATOR_STATUS_INACTIVATING {
			return fmt.Errorf("operator with id %d is not inactivating", unbondingOperator.OperatorID)
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
	return utils.FindDuplicateFunc(operators, func(a, b Operator) bool {
		return a.ID == b.ID
	})
}

// findDuplicateOperatorsParams returns the first duplicated operator params in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicateOperatorsParams(operators []OperatorParamsRecord) *OperatorParamsRecord {
	return utils.FindDuplicateFunc(operators, func(a, b OperatorParamsRecord) bool {
		return a.OperatorID == b.OperatorID
	})
}

// findDuplicateUnbondingOperators returns the first duplicated unbonding operator in the slice.
// If no duplicates are found, it returns nil instead.
func findDuplicateUnbondingOperators(operators []UnbondingOperator) *UnbondingOperator {
	return utils.FindDuplicateFunc(operators, func(a, b UnbondingOperator) bool {
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

// --------------------------------------------------------------------------------------------------------------------

// NewOperatorParamsRecord creates a new OperatorParamsRecord instance.
func NewOperatorParamsRecord(operatorID uint32, operatorParams OperatorParams) OperatorParamsRecord {
	return OperatorParamsRecord{
		OperatorID: operatorID,
		Params:     operatorParams,
	}
}

// Validate checks that the OperatorParamsRecord has valid values.
func (o *OperatorParamsRecord) Validate() error {
	return o.Params.Validate()
}
