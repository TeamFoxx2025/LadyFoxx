package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/TeamFoxx2025/LadyFoxx/command/backup"
	"github.com/TeamFoxx2025/LadyFoxx/command/bridge"
	"github.com/TeamFoxx2025/LadyFoxx/command/genesis"
	"github.com/TeamFoxx2025/LadyFoxx/command/helper"
	"github.com/TeamFoxx2025/LadyFoxx/command/ibft"
	"github.com/TeamFoxx2025/LadyFoxx/command/license"
	"github.com/TeamFoxx2025/LadyFoxx/command/monitor"
	"github.com/TeamFoxx2025/LadyFoxx/command/peers"
	"github.com/TeamFoxx2025/LadyFoxx/command/polybft"
	"github.com/TeamFoxx2025/LadyFoxx/command/polybftsecrets"
	"github.com/TeamFoxx2025/LadyFoxx/command/regenesis"
	"github.com/TeamFoxx2025/LadyFoxx/command/rootchain"
	"github.com/TeamFoxx2025/LadyFoxx/command/secrets"
	"github.com/TeamFoxx2025/LadyFoxx/command/server"
	"github.com/TeamFoxx2025/LadyFoxx/command/status"
	"github.com/TeamFoxx2025/LadyFoxx/command/txpool"
	"github.com/TeamFoxx2025/LadyFoxx/command/version"
)

type RootCommand struct {
	baseCmd *cobra.Command
}

func NewRootCommand() *RootCommand {
	rootCommand := &RootCommand{
		baseCmd: &cobra.Command{
			Short: "Polygon Edge is a framework for building Ethereum-compatible Blockchain networks",
		},
	}

	helper.RegisterJSONOutputFlag(rootCommand.baseCmd)

	rootCommand.registerSubCommands()

	return rootCommand
}

func (rc *RootCommand) registerSubCommands() {
	rc.baseCmd.AddCommand(
		version.GetCommand(),
		txpool.GetCommand(),
		status.GetCommand(),
		secrets.GetCommand(),
		peers.GetCommand(),
		rootchain.GetCommand(),
		monitor.GetCommand(),
		ibft.GetCommand(),
		backup.GetCommand(),
		genesis.GetCommand(),
		server.GetCommand(),
		license.GetCommand(),
		polybftsecrets.GetCommand(),
		polybft.GetCommand(),
		bridge.GetCommand(),
		regenesis.GetCommand(),
	)
}

func (rc *RootCommand) Execute() {
	if err := rc.baseCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}
}
