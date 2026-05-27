package vserver

import (
	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/cmd/vserver/flavor"
	"github.com/vngcloud/greennode-cli/cmd/vserver/image"
	"github.com/vngcloud/greennode-cli/cmd/vserver/secgroup"
	"github.com/vngcloud/greennode-cli/cmd/vserver/server"
	"github.com/vngcloud/greennode-cli/cmd/vserver/subnet"
	"github.com/vngcloud/greennode-cli/cmd/vserver/volume"
	"github.com/vngcloud/greennode-cli/cmd/vserver/volumetype"
	"github.com/vngcloud/greennode-cli/cmd/vserver/vpc"
)

// VServerCmd is the parent command for all vServer subcommands.
var VServerCmd = &cobra.Command{
	Use:   "vserver",
	Short: "VNG Virtual Server (vServer) commands",
	Long:  "Manage vServer instances and related resources.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	VServerCmd.AddCommand(server.ServerCmd)
	VServerCmd.AddCommand(volume.VolumeCmd)
	VServerCmd.AddCommand(vpc.VpcCmd)
	VServerCmd.AddCommand(subnet.SubnetCmd)
	VServerCmd.AddCommand(secgroup.SecgroupCmd)
	VServerCmd.AddCommand(flavor.FlavorCmd)
	VServerCmd.AddCommand(volumetype.VolumeTypeCmd)
	VServerCmd.AddCommand(image.ImageCmd)
}
