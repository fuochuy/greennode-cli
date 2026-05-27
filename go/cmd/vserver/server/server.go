package server

import (
	"github.com/spf13/cobra"
)

// ServerCmd is the parent command for all server subcommands.
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage vServer instances",
	Long:  "Create, list, get, and manage vServer instances.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ServerCmd.AddCommand(listCmd)
	ServerCmd.AddCommand(getCmd)
	ServerCmd.AddCommand(createCmd)
	ServerCmd.AddCommand(startCmd)
	ServerCmd.AddCommand(stopCmd)
	ServerCmd.AddCommand(rebootCmd)
	ServerCmd.AddCommand(resizeCmd)
	ServerCmd.AddCommand(deleteCmd)
}
