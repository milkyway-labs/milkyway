package types

import (
	"fmt"
	"slices"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

// NewParams creates a new Params instance.
func NewParams(
	insurancePercentage math.LegacyDec,
	burners []string,
	minters []string,
	allowedChannels []string,
) Params {
	return Params{
		InsurancePercentage: insurancePercentage,
		Burners:             burners,
		Minters:             minters,
		AllowedChannels:     allowedChannels,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(math.LegacyNewDec(2), nil, nil, nil)
}

// Validate ensure that the Prams structure is correct
func (p *Params) Validate() error {
	if p.InsurancePercentage.LTE(math.LegacyNewDec(0)) || p.InsurancePercentage.GT(math.LegacyNewDec(100)) {
		return ErrInvalidInsurancePercentage
	}
	for _, address := range p.Minters {
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return err
		}
	}
	for _, address := range p.Burners {
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return err
		}
	}
	for _, channel := range p.AllowedChannels {
		if !channeltypes.IsValidChannelID(channel) {
			return fmt.Errorf("invalid channel id: %s", channel)
		}
	}
	return nil
}

// IsAllowedChannel checks if is allowed to receive
// deposits to the insurance fund from the provided channel.
func (p *Params) IsAllowedChannel(channelID string) bool {
	return slices.Contains(p.AllowedChannels, channelID)
}
