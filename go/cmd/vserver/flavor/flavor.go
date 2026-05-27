package flavor

import (
	"github.com/spf13/cobra"
)

// FlavorCmd is the parent command for all flavor subcommands.
var FlavorCmd = &cobra.Command{
	Use:   "flavor",
	Short: "Manage vServer flavors",
	Long:  "List and inspect available vServer flavors.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	FlavorCmd.AddCommand(listCmd)
	FlavorCmd.AddCommand(listFamiliesCmd)
	FlavorCmd.AddCommand(listCodesCmd)
}
