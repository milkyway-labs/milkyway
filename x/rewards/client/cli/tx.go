package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/milkyway-labs/milkyway/x/rewards/types"
)

// GetTxCmd returns a new command to perform operators transactions
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Rewards transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdCreateRewardsPlan(),
		// The other commands are generated through the auto CLI
	)

	return txCmd
}

// GetCmdCreateRewardsPlan returns the command allowing to create a rewards plan
// for a service.
func GetCmdCreateRewardsPlan() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-rewards-plan [path/to/rewards_plan.json]",
		Args:    cobra.ExactArgs(1),
		Short:   "Creates a rewards plan for a service",
		Example: fmt.Sprintf(`%s tx %s create-rewards-plan path/to/rewards_plan.json --from alice`, version.AppName, types.ModuleName),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Creates a rewards plan for a service.
The rewards plan should be defined in a JSON file.

Example:
$ %s tx %s create-rewards-plan rewards_plan.json --from alice

Where rewards_plan.json contains:

{
  "service_id": 1,
  "description": "test plan",
  "amount_per_day": "1000uinit",
  "start_time": "2024-01-01T00:00:00Z",
  "end_time": "2024-12-31T23:59:59Z",
  "pools_distribution": {
      "weight": 1,
      "type": {
          "@type":"/milkyway.rewards.v1.DistributionTypeBasic"
      }
  },
  "operators_distribution": {
      "weight": 1,
      "type": {
          "@type": "/milkyway.rewards.v1.DistributionTypeBasic"
      }
  },
  "users_distribution": {
      "weight": 1,
      "type": {
          "@type": "/milkyway.rewards.v1.UsersDistributionTypeBasic"
      }
  }
}
`, version.AppName, types.ModuleName),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			rewardsPlan, err := ParseRewardsPlan(clientCtx.Codec, args[0])
			if err != nil {
				return fmt.Errorf("parsing rewards plan json: %w", err)
			}
			err = rewardsPlan.Validate(clientCtx.Codec)
			if err != nil {
				return fmt.Errorf("invalid rewards plan json: %w", err)
			}

			creator := clientCtx.FromAddress.String()
			msg := types.NewMsgCreateRewardsPlan(
				rewardsPlan.ServiceID,
				rewardsPlan.Description,
				rewardsPlan.AmountPerDay,
				rewardsPlan.StartTime,
				rewardsPlan.EndTime,
				rewardsPlan.PoolsDistribution,
				rewardsPlan.OperatorsDistribution,
				rewardsPlan.UsersDistribution,
				creator,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
