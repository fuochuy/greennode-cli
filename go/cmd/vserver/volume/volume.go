package volume

import (
	"github.com/spf13/cobra"
)

// VolumeCmd is the parent command for all volume subcommands.
var VolumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Manage vServer volumes",
	Long:  "Create, list, get, and delete vServer block storage volumes.\n\nTo see available volume types for a zone, run: grn vserver volume-type list",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	VolumeCmd.AddCommand(listCmd)
	VolumeCmd.AddCommand(getCmd)
	VolumeCmd.AddCommand(createCmd)
	VolumeCmd.AddCommand(resizeCmd)
	VolumeCmd.AddCommand(deleteCmd)
}
