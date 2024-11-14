package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

// OnRecvPacket processes the packet received from the IBC handler
func (k *Keeper) OnRecvPacket(
	ctx sdk.Context,
	data transfertypes.FungibleTokenPacketData,
	msgDepositInsurance types.MsgDepositInsurance,
) error {

	// Ensure the receiver is the x/liquidvesting module account
	if data.Receiver != k.ModuleAddress {
		return fmt.Errorf("the receiver should be the module address, got: %s, expected: %s", data.Receiver, k.ModuleAddress)
	}

	// Ensure that the message is valid
	if err := msgDepositInsurance.ValidateBasic(); err != nil {
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
		accountAddress, err := sdk.AccAddressFromBech32(deposit.Depositor)
		if err != nil {
			return err
		}
		err = k.AddToUserInsuranceFund(ctx, accountAddress, sdk.NewCoins(deposit.Amount))
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
