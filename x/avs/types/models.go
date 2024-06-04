package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAVS creates a new AVS instance
func NewAVS(
	id uint32,
	status AVSStatus,
	name string,
	description string,
	website string,
	pictureURL string,
	admin string,
) AVS {
	return AVS{
		ID:          id,
		Status:      status,
		Name:        name,
		Description: description,
		Website:     website,
		PictureURL:  pictureURL,
		Admin:       admin,
	}
}

// Validate checks that the AVS has valid values.
func (a *AVS) Validate() error {
	if a.Status == AVS_STATUS_UNSPECIFIED {
		return fmt.Errorf("invalid status: %s", a.Status)
	}

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
