package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/milkyway-labs/milkyway/v9/x/investors/types"
)

var _ types.DistrHooks = Hooks{}

type Hooks struct {
	*Keeper
}

func (k *Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) AfterSetWithdrawAddress(ctx context.Context, delAddr, withdrawAddr sdk.AccAddress) error {
	delegator, err := h.accountKeeper.AddressCodec().BytesToString(delAddr)
	if err != nil {
		return err
	}
	withdrawAddress, err := h.accountKeeper.AddressCodec().BytesToString(withdrawAddr)
	if err != nil {
		return err
	}
	// Associate the delegator address with the withdraw address, which will
	// construct a reverse lookup table.
	err = h.Delegators.Set(ctx, withdrawAddress, delegator)
	if err != nil {
		return err
	}
	return nil
}
