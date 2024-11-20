package keeper

import (
	"fmt"
	"slices"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

// OnRecvPacket processes the packet received from the IBC handler
func (k *Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data transfertypes.FungibleTokenPacketData,
	msgDepositInsurance types.MsgDepositInsurance,
) error {
	// Ensure that the sender is allowed to deposit
	canDeposit, err := k.isAllowedDepositor(ctx, data.Sender)
	if err != nil {
		return err
	}

	if !canDeposit {
		return fmt.Errorf("the sender %s is not allowed to deposit", data.Sender)
	}

	// Ensure the receiver is the x/liquidvesting module account
	if data.Receiver != k.ModuleAddress {
		return fmt.Errorf("the receiver should be the module address, got: %s, expected: %s", data.Receiver, k.ModuleAddress)
	}

	// Ensure that the message is valid
	err = msgDepositInsurance.ValidateBasic()
	if err != nil {
		return err
	}

	// Get the total deposit amount from the message
	totalDeposit, err := msgDepositInsurance.GetTotalDepositAmount()
	if err != nil {
		return err
	}

	// Parse the amount from the ics20Packet
	amount, ok := math.NewIntFromString(data.GetAmount())
	if !ok {
		return fmt.Errorf("invalid ics20 amount")
	}
	receivedAmount := sdk.NewCoin(data.Denom, amount)

	// Ensure that we have received the same amount of tokens
	// as the ones that needs to be added to the users' insurance fund
	if !receivedAmount.Equal(totalDeposit) {
		return fmt.Errorf("amount received is not equal to the amounts to deposit in the users' insurance fund")
	}

	// Deposit the amounts into the users' insurance fund
	for _, deposit := range msgDepositInsurance.Amounts {
		// Convert the coin denom to its IBC representation
		sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
		// NOTE: sourcePrefix contains the trailing "/"
		prefixedDenom := sourcePrefix + deposit.Amount.Denom
		// construct the denomination trace from the full raw denomination
		denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
		// Convert the insurance deposit to its IBC representation
		ibcCoin := sdk.NewCoin(denomTrace.IBCDenom(), deposit.Amount.Amount)

		err = k.AddToUserInsuranceFund(ctx, deposit.Depositor, sdk.NewCoins(ibcCoin))
		if err != nil {
			return err
		}

		// Dispatch the deposit event.
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeDepositToUserInsuranceFund,
				sdk.NewAttribute(types.AttributeKeyUser, deposit.Depositor),
				sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			),
		)
	}

	return nil
}

// IsAllowedDepositor checks if the provided address is allowed to deposit funds
// to the insurance fund.
func (k *Keeper) isAllowedDepositor(ctx sdk.Context, address string) (bool, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return false, err
	}

	return slices.Contains(params.TrustedDelegates, address), nil
}
