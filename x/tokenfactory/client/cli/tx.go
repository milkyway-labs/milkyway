package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/v10/x/tokenfactory/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		NewCreateDenomCmd(),
		NewMintCmd(),
		NewBurnCmd(),
		// NewForceTransferCmd(),
		NewChangeAdminCmd(),
		NewSetBeforeSendHookCmd(),
		NewMsgSetDenomMetadata(),
	)

	return cmd
}

func NewMsgSetDenomMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-denom-metadata [metadata-file] [flags]",
		Short: "Sets a bank metadata for a token denom. Must have admin authority to do so.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			metadataBytes, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			var metadata banktypes.Metadata
			if err = json.Unmarshal(metadataBytes, &metadata); err != nil {
				return err
			}

			msg := types.NewMsgSetDenomMetadata(
				clientCtx.GetFromAddress().String(),
				metadata,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewCreateDenomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-denom [subdenom] [flags]",
		Short: "create a new denom from an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateDenom(
				clientCtx.FromAddress.String(),
				args[0],
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewMintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [amount] [mint-to-address] [flags]",
		Short: "Mint a denom to an address. Must have admin authority to do so.",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			mintToAddress := ""
			if len(args) == 2 {
				mintToAddress = args[1]
			}

			msg := types.NewMsgMintTo(
				clientCtx.GetFromAddress().String(),
				amount,
				mintToAddress,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewBurnCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [amount] [burn-from-address] [flags]",
		Short: "Burn tokens from an address. Must have admin authority to do so.",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			burnFromAddress := ""
			if len(args) == 2 {
				burnFromAddress = args[1]
			}

			msg := types.NewMsgBurnFrom(
				clientCtx.GetFromAddress().String(),
				amount,
				burnFromAddress,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewChangeAdminCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "change-admin [denom] [new-admin] [flags]",
		Short: "Changes the admin address for a factory-created denom. Must have admin authority to do so.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgChangeAdmin(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewChangeAdminCmd broadcast MsgChangeAdmin
func NewSetBeforeSendHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-beforesend-hook [denom] [cosmwasm-address] [flags]",
		Short: "Set a cosmwasm contract to be the beforesend hook for a factory-created denom. Must have admin authority to do so.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			txf = txf.WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			msg := types.NewMsgSetBeforeSendHook(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
			)

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
