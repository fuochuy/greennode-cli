package secgroup

import (
	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/cmd/vserver/secgroup/rule"
)

// SecgroupCmd is the parent command for all security group subcommands.
var SecgroupCmd = &cobra.Command{
	Use:   "secgroup",
	Short: "Manage security groups",
	Long:  "Create, list, and delete security groups and their rules.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	SecgroupCmd.AddCommand(listCmd)
	SecgroupCmd.AddCommand(getCmd)
	SecgroupCmd.AddCommand(createCmd)
	SecgroupCmd.AddCommand(deleteCmd)
	SecgroupCmd.AddCommand(rule.RuleCmd)
}
