package txpool

import (
	"github.com/TeamFoxx2025/LadyFoxx/command/helper"
	"github.com/TeamFoxx2025/LadyFoxx/command/txpool/status"
	"github.com/TeamFoxx2025/LadyFoxx/command/txpool/subscribe"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	txPoolCmd := &cobra.Command{
		Use:   "txpool",
		Short: "Top level command for interacting with the transaction pool. Only accepts subcommands.",
	}

	helper.RegisterGRPCAddressFlag(txPoolCmd)

	registerSubcommands(txPoolCmd)

	return txPoolCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// txpool status
		status.GetCommand(),
		// txpool subscribe
		subscribe.GetCommand(),
	)
}
