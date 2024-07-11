package types

// NewParams creates a new Params object
func NewParams() Params {
	return Params{}
}

// DefaultParams returns default Params
func DefaultParams() Params {
	return Params{}
}

// Validate checks that the parameters have valid values.
func (p *Params) Validate() error {
	return nil
}
