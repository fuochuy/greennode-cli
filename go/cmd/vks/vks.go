package vks

import (
	"github.com/spf13/cobra"
)

// VksCmd is the parent command for all VKS subcommands.
var VksCmd = &cobra.Command{
	Use:   "vks",
	Short: "VNG Kubernetes Service (VKS) commands",
	Long:  "Manage VKS clusters, node groups, and related resources.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Cluster commands
	VksCmd.AddCommand(listClustersCmd)
	VksCmd.AddCommand(getClusterCmd)
	VksCmd.AddCommand(createClusterCmd)
	VksCmd.AddCommand(updateClusterCmd)
	VksCmd.AddCommand(deleteClusterCmd)

	// Nodegroup commands
	VksCmd.AddCommand(listNodegroupsCmd)
	VksCmd.AddCommand(getNodegroupCmd)
	VksCmd.AddCommand(createNodegroupCmd)
	VksCmd.AddCommand(updateNodegroupCmd)
	VksCmd.AddCommand(deleteNodegroupCmd)

	// Wait commands
	VksCmd.AddCommand(waitClusterActiveCmd)

	// Auto-upgrade commands
	VksCmd.AddCommand(setAutoUpgradeConfigCmd)
	VksCmd.AddCommand(deleteAutoUpgradeConfigCmd)

	// Quota commands
	VksCmd.AddCommand(getQuotaCmd)
}
