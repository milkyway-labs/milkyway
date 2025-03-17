package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

	"github.com/milkyway-labs/milkyway/v10/x/liquidvesting/types"
)

// OnRecvPacket processes the packet received from the IBC handler
func (k *Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data transfertypes.FungibleTokenPacketData,
	msgDepositInsurance types.MsgDepositInsurance,
) error {
	// Get the params
	params, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	// Check if is allowed to receive deposits to the insurance fund
	// from the channel
	if !params.IsAllowedChannel(packet.DestinationChannel) {
		return fmt.Errorf("deposit not allowed using channel %s", packet.DestinationChannel)
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
	totalDeposit := msgDepositInsurance.GetTotalDepositAmount()

	// Parse the amount from the ics20Packet
	amount, ok := math.NewIntFromString(data.GetAmount())
	if !ok {
		return fmt.Errorf("invalid ics20 amount")
	}
	receivedAmount := sdk.NewCoin(data.Denom, amount)

	// Ensure that we have received the same amount of tokens
	// as the ones that needs to be added to the users' insurance fund
	if !receivedAmount.Amount.Equal(totalDeposit) {
		return fmt.Errorf("amount received is not equal to the amounts to deposit in the users' insurance fund")
	}

	// Deposit the amounts into the users' insurance fund
	for _, deposit := range msgDepositInsurance.Amounts {
		// Convert the coin denom to its IBC representation
		sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
		// NOTE: sourcePrefix contains the trailing "/"
		prefixedDenom := sourcePrefix + data.Denom
		// construct the denomination trace from the full raw denomination
		denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
		// Convert the insurance deposit to its IBC representation
		ibcCoin := sdk.NewCoin(denomTrace.IBCDenom(), deposit.Amount)

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
