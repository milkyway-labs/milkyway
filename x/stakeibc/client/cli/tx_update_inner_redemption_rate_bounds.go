package cli

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milk/x/stakeibc/types"
)

func CmdUpdateInnerRedemptionRateBounds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-redemption-rate-bounds [chainid] [min-bound] [max-bound]",
		Short: "Broadcast message set-redemption-rate-bounds",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChainId := args[0]
			minInnerRedemptionRate := sdkmath.LegacyMustNewDecFromStr(args[1])
			maxInnerRedemptionRate := sdkmath.LegacyMustNewDecFromStr(args[2])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateInnerRedemptionRateBounds(
				clientCtx.GetFromAddress().String(),
				argChainId,
				minInnerRedemptionRate,
				maxInnerRedemptionRate,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
