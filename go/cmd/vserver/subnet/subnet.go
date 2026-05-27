package subnet

import (
	"github.com/spf13/cobra"
)

// SubnetCmd is the parent command for all subnet subcommands.
var SubnetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "Manage subnets",
	Long:  "Create, list, get, and delete subnets within a VPC.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	SubnetCmd.AddCommand(listCmd)
	SubnetCmd.AddCommand(getCmd)
	SubnetCmd.AddCommand(createCmd)
	SubnetCmd.AddCommand(deleteCmd)
}
