package types

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ParseServiceID parses a string into a uint32
func ParseServiceID(value string) (uint32, error) {
	id, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewService creates a new Service instance
func NewService(
	id uint32,
	status AVSStatus,
	name string,
	description string,
	website string,
	pictureURL string,
	admin string,
) Service {
	return Service{
		ID:          id,
		Status:      status,
		Name:        name,
		Description: description,
		Website:     website,
		PictureURL:  pictureURL,
		Admin:       admin,
	}
}

// Validate checks that the Service has valid values.
func (a *Service) Validate() error {
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

// --------------------------------------------------------------------------------------------------------------------

// ServiceUpdate defines the fields that can be updated in an Service.
type ServiceUpdate struct {
	Name        string
	Description string
	Website     string
	PictureURL  string
}

// NewServiceUpdate returns a new ServiceUpdate instance.
func NewServiceUpdate(
	name string,
	description string,
	website string,
	pictureURL string,
) ServiceUpdate {
	return ServiceUpdate{
		Name:        name,
		Description: description,
		Website:     website,
		PictureURL:  pictureURL,
	}
}

// Update returns a new Service with updated fields.
func (a *Service) Update(update ServiceUpdate) Service {
	if update.Name == DoNotModify {
		update.Name = a.Name
	}

	if update.Description == DoNotModify {
		update.Description = a.Description
	}

	if update.Website == DoNotModify {
		update.Website = a.Website
	}

	if update.PictureURL == DoNotModify {
		update.PictureURL = a.PictureURL
	}

	return NewService(
		a.ID,
		a.Status,
		update.Name,
		update.Description,
		update.Website,
		update.PictureURL,
		a.Admin,
	)
}
