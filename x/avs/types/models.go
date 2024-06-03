package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAVS creates a new AVS instance
func NewAVS(id uint32, name string, admin string) AVS {
	return AVS{
		ID:    id,
		Name:  name,
		Admin: admin,
	}
}

// Validate checks that the AVS has valid values.
func (a *AVS) Validate() error {
	if a.ID == 0 {
		return fmt.Errorf("invalid id: %d", a.ID)
	}

	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("invalid name: %s", a.Name)
	}

	_, err := sdk.AccAddressFromBech32(a.Admin)
	if err != nil {
		return fmt.Errorf("invalid admin address")
	}

	return nil
}
