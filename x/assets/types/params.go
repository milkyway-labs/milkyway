package types

// NewParams creates a new Params object
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// Validate checks that the parameters have valid values.
func (p *Params) Validate() error {
	return nil
}