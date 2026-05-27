package volumetype

import (
	"github.com/spf13/cobra"
)

// VolumeTypeCmd is the parent command for all volume type subcommands.
var VolumeTypeCmd = &cobra.Command{
	Use:   "volume-type",
	Short: "Manage vServer volume types",
	Long:  "List available volume types for a zone.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	VolumeTypeCmd.AddCommand(listCmd)
}
