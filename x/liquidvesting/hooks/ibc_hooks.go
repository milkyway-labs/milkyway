package hooks

import (
	"encoding/json"
	"fmt"
	"slices"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	chan4types "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibchooks "github.com/initia-labs/initia/x/ibc-hooks"

	milkywaytypes "github.com/milkyway-labs/milkyway/types"
	"github.com/milkyway-labs/milkyway/utils"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/keeper"
	"github.com/milkyway-labs/milkyway/x/liquidvesting/types"
)

var _ ibchooks.OnRecvPacketOverrideHooks = IBCHooks{}

// IBCHooks represents the structure that implements the
// ibc_hooks.OnRecvPacketOverrideHooks interface to execute
// custom logic when an IBC token transfer packet is received.
type IBCHooks struct {
	*keeper.Keeper
}

func NewIBCHooks(k *keeper.Keeper) IBCHooks {
	return IBCHooks{k}
}

func (h IBCHooks) onRecvIcs20Packet(
	ctx sdk.Context,
	im ibchooks.IBCMiddleware,
	packet chan4types.Packet,
	relayer sdk.AccAddress,
	ics20Packet transfertypes.FungibleTokenPacketData,
) exported.Acknowledgement {
	objFound, object := utils.JSONStringHasKey(ics20Packet.GetMemo(), types.ModuleName)
	if !objFound {
		// Module payload not found, pass the packet to next middleware
		return im.App.OnRecvPacket(ctx, packet, relayer)
	}

	// Ensure the receiver is the x/liquidvesting module account
	if ics20Packet.Receiver != h.ModuleAddress {
		return milkywaytypes.NewEmitErrorAcknowledgement(
			fmt.Errorf("the receiver should be the module address, got: %s, expected: %s", ics20Packet.Receiver, h.ModuleAddress),
		)
	}

	// Ensure that the sender is allowed to deposit
	canDeposit, err := h.isAllowedDepositor(ctx, ics20Packet.Sender)
	if err != nil {
		return milkywaytypes.NewEmitErrorAcknowledgement(err)
	}
	if !canDeposit {
		return milkywaytypes.NewEmitErrorAcknowledgement(
			fmt.Errorf("the sender %s is not allowed to deposit", ics20Packet.Sender))
	}

	// Parse the message from the memo
	bytes, err := json.Marshal(object[types.ModuleName])
	if err != nil {
		return milkywaytypes.NewEmitErrorAcknowledgement(err)
	}
	var depositMsg types.MsgDepositInsurance
	if err := json.Unmarshal(bytes, &depositMsg); err != nil {
		return milkywaytypes.NewEmitErrorAcknowledgement(err)
	}

	// Ensure that the message is valid
	if err := depositMsg.ValidateBasic(); err != nil {
		return milkywaytypes.NewEmitErrorAcknowledgement(err)
	}

	// Get the total deposit amount from the message
	totalDeposit, err := depositMsg.GetTotalDepositAmount()
	if err != nil {
		return milkywaytypes.NewEmitErrorAcknowledgement(err)
	}

	// Parse the amount from the ics20Packet
	amount, ok := math.NewIntFromString(ics20Packet.GetAmount())
	if !ok {
		return milkywaytypes.NewEmitErrorAcknowledgement(fmt.Errorf("invalid ics20 amount"))
	}
	receivedAmount := sdk.NewCoin(ics20Packet.Denom, amount)

	// Ensure that we have received the same amount of tokens
	// as the ones that needs to be added to the users' insurance fund
	if !receivedAmount.Equal(totalDeposit) {
		return milkywaytypes.NewEmitErrorAcknowledgement(
			fmt.Errorf("amount received is not equal to the amounts to deposit in the users' insurance fund"),
		)
	}

	// Deposit the amounts into the users' insurance fund
	for _, deposit := range depositMsg.Amounts {
		accountAddress, err := sdk.AccAddressFromBech32(deposit.Depositor)
		if err != nil {
			return milkywaytypes.NewEmitErrorAcknowledgement(err)
		}
		err = h.AddToUserInsuranceFund(ctx, accountAddress, sdk.NewCoins(deposit.Amount))
		if err != nil {
			return milkywaytypes.NewEmitErrorAcknowledgement(err)
		}

		// Dispatch the deposit event.
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeDepositToUserInsuranceFund,
				sdk.NewAttribute(types.AttributeKeyUser, deposit.Depositor),
				sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
			),
		)
	}

	return im.App.OnRecvPacket(ctx, packet, relayer)
}

// OnRecvPacketOverride implements ibc_hooks.OnRecvPacketOverrideHooks.
func (h IBCHooks) OnRecvPacketOverride(
	im ibchooks.IBCMiddleware,
	ctx sdk.Context,
	packet chan4types.Packet,
	relayer sdk.AccAddress,
) exported.Acknowledgement {
	if isIcs20, ics20Packet := utils.IsIcs20Packet(packet.GetData()); isIcs20 {
		return h.onRecvIcs20Packet(ctx, im, packet, relayer, ics20Packet)
	}

	return im.App.OnRecvPacket(ctx, packet, relayer)
}

// IsAllowedDepositor checks if the provided address is allowed to deposit funds
// to the insurance fund.
func (h IBCHooks) isAllowedDepositor(ctx sdk.Context, address string) (bool, error) {
	params, err := h.GetParams(ctx)
	if err != nil {
		return false, err
	}

	return slices.Contains(params.TrustedDelegates, address), nil
}