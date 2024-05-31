package types

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Default init params
var (
	// these are default intervals _in epochs_ NOT in blocks
	DefaultDepositInterval        uint64 = 1
	DefaultDelegateInterval       uint64 = 1
	DefaultReinvestInterval       uint64 = 1
	DefaultRewardsInterval        uint64 = 1
	DefaultRedemptionRateInterval uint64 = 1
	// you apparently cannot safely encode floats, so we make commission / 100
	DefaultStrideCommission             uint64 = 10
	DefaultICATimeoutNanos              uint64 = 600000000000
	DefaultBufferSize                   uint64 = 5             // 1/5=20% of the epoch
	DefaultIbcTimeoutBlocks             uint64 = 300           // 300 blocks ~= 30 minutes
	DefaultFeeTransferTimeoutNanos      uint64 = 1800000000000 // 30 minutes
	DefaultMinRedemptionRateThreshold   uint64 = 90            // divide by 100, so 90 = 0.9
	DefaultMaxRedemptionRateThreshold   uint64 = 150           // divide by 100, so 150 = 1.5
	DefaultMaxStakeICACallsPerEpoch     uint64 = 100
	DefaultIBCTransferTimeoutNanos      uint64 = 1800000000000 // 30 minutes
	DefaultValidatorSlashQueryThreshold uint64 = 1             // denominated in percentage of TVL (1 => 1%)
	DefaultValidatorWeightCap           uint64 = 10            // percentage (10 => 10%)
)

// NewParams creates a new Params instance
func NewParams(
	depositInterval uint64,
	delegateInterval uint64,
	rewardsInterval uint64,
	redemptionRateInterval uint64,
	strideCommission uint64,
	reinvestInterval uint64,
	icaTimeoutNanos uint64,
	bufferSize uint64,
	ibcTimeoutBlocks uint64,
	feeTransferTimeoutNanos uint64,
	maxStakeIcaCallsPerEpoch uint64,
	defaultMinRedemptionRateThreshold uint64,
	defaultMaxRedemptionRateThreshold uint64,
	ibcTransferTimeoutNanos uint64,
	validatorSlashQueryInterval uint64,
	validatorWeightCap uint64,
) Params {
	return Params{
		DepositInterval:                   depositInterval,
		DelegateInterval:                  delegateInterval,
		RewardsInterval:                   rewardsInterval,
		RedemptionRateInterval:            redemptionRateInterval,
		StrideCommission:                  strideCommission,
		ReinvestInterval:                  reinvestInterval,
		IcaTimeoutNanos:                   icaTimeoutNanos,
		BufferSize:                        bufferSize,
		IbcTimeoutBlocks:                  ibcTimeoutBlocks,
		FeeTransferTimeoutNanos:           feeTransferTimeoutNanos,
		MaxStakeIcaCallsPerEpoch:          maxStakeIcaCallsPerEpoch,
		DefaultMinRedemptionRateThreshold: defaultMinRedemptionRateThreshold,
		DefaultMaxRedemptionRateThreshold: defaultMaxRedemptionRateThreshold,
		IbcTransferTimeoutNanos:           ibcTransferTimeoutNanos,
		ValidatorSlashQueryThreshold:      validatorSlashQueryInterval,
		ValidatorWeightCap:                validatorWeightCap,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultDepositInterval,
		DefaultDelegateInterval,
		DefaultRewardsInterval,
		DefaultRedemptionRateInterval,
		DefaultStrideCommission,
		DefaultReinvestInterval,
		DefaultICATimeoutNanos,
		DefaultBufferSize,
		DefaultIbcTimeoutBlocks,
		DefaultFeeTransferTimeoutNanos,
		DefaultMaxStakeICACallsPerEpoch,
		DefaultMinRedemptionRateThreshold,
		DefaultMaxRedemptionRateThreshold,
		DefaultIBCTransferTimeoutNanos,
		DefaultValidatorSlashQueryThreshold,
		DefaultValidatorWeightCap,
	)
}

func validTimeoutNanos(i interface{}) error {
	ival, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("parameter not accepted: %T", i)
	}

	tenMin := uint64(600000000000)
	oneHour := uint64(600000000000 * 6)

	if ival < tenMin {
		return fmt.Errorf("parameter must be g.t. 600000000000ns: %d", ival)
	}
	if ival > oneHour {
		return fmt.Errorf("parameter must be less than %dns: %d", oneHour, ival)
	}
	return nil
}

func validMaxRedemptionRateThreshold(i interface{}) error {
	ival, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("parameter not accepted: %T", i)
	}

	maxVal := uint64(1000) // divide by 100, so 1000 => 10

	if ival > maxVal {
		return fmt.Errorf("parameter must be l.t. 1000: %d", ival)
	}

	return nil
}

func validMinRedemptionRateThreshold(i interface{}) error {
	ival, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("parameter not accepted: %T", i)
	}

	minVal := uint64(75) // divide by 100, so 75 => 0.75

	if ival < minVal {
		return fmt.Errorf("parameter must be g.t. 75: %d", ival)
	}

	return nil
}

func validValidatorWeightCap(i interface{}) error {
	if err := isPositive(i); err != nil {
		return err
	}
	if err := isPercentage(i); err != nil {
		return err
	}
	return nil
}

func isPositive(i interface{}) error {
	ival, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("parameter not accepted: %T", i)
	}

	if ival <= 0 {
		return fmt.Errorf("parameter must be positive: %d", ival)
	}
	return nil
}

func isPercentage(i interface{}) error {
	ival, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("parameter not accepted: %T", i)
	}

	if ival > 100 {
		return fmt.Errorf("parameter must be between 0 and 100: %d", ival)
	}
	return nil
}

func isCommission(i interface{}) error {
	ival, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("commission not accepted: %T", i)
	}

	if ival > 100 {
		return fmt.Errorf("commission must be less than 100: %d", ival)
	}
	return nil
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := isPositive(p.DepositInterval); err != nil {
		return err
	}
	if err := isPositive(p.DelegateInterval); err != nil {
		return err
	}
	if err := isPositive(p.RewardsInterval); err != nil {
		return err
	}
	if err := isPositive(p.RedemptionRateInterval); err != nil {
		return err
	}
	if err := isCommission(p.StrideCommission); err != nil {
		return err
	}
	if err := isPositive(p.ReinvestInterval); err != nil {
		return err
	}
	if err := isPositive(p.IcaTimeoutNanos); err != nil {
		return err
	}
	if err := isPositive(p.BufferSize); err != nil {
		return err
	}
	if err := isPositive(p.IbcTimeoutBlocks); err != nil {
		return err
	}
	if err := validTimeoutNanos(p.FeeTransferTimeoutNanos); err != nil {
		return err
	}
	if err := isPositive(p.MaxStakeIcaCallsPerEpoch); err != nil {
		return err
	}
	if err := validMinRedemptionRateThreshold(p.DefaultMinRedemptionRateThreshold); err != nil {
		return err
	}
	if err := validMaxRedemptionRateThreshold(p.DefaultMaxRedemptionRateThreshold); err != nil {
		return err
	}
	if err := validTimeoutNanos(p.IbcTransferTimeoutNanos); err != nil {
		return err
	}
	if err := validValidatorWeightCap(p.ValidatorWeightCap); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
