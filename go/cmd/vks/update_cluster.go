package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var updateClusterCmd = &cobra.Command{
	Use:   "update-cluster",
	Short: "Update a VKS cluster",
	RunE:  runUpdateCluster,
}

func init() {
	f := updateClusterCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("k8s-version", "", "Kubernetes version (required)")
	f.String("whitelist-node-cidrs", "", "Whitelist CIDRs, comma-separated (required)")
	f.Bool("enabled-load-balancer-plugin", false, "Enable load balancer plugin")
	f.Bool("no-load-balancer-plugin", false, "Disable load balancer plugin")
	f.Bool("enabled-block-store-csi-plugin", false, "Enable block store CSI plugin")
	f.Bool("no-block-store-csi-plugin", false, "Disable block store CSI plugin")
	f.Bool("dry-run", false, "Validate parameters without updating")

	updateClusterCmd.MarkFlagRequired("cluster-id")
	updateClusterCmd.MarkFlagRequired("k8s-version")
	updateClusterCmd.MarkFlagRequired("whitelist-node-cidrs")
}

func runUpdateCluster(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	k8sVersion, _ := cmd.Flags().GetString("k8s-version")
	whitelistCIDRs, _ := cmd.Flags().GetString("whitelist-node-cidrs")
	enabledLB, _ := cmd.Flags().GetBool("enabled-load-balancer-plugin")
	noLB, _ := cmd.Flags().GetBool("no-load-balancer-plugin")
	enabledCSI, _ := cmd.Flags().GetBool("enabled-block-store-csi-plugin")
	noCSI, _ := cmd.Flags().GetBool("no-block-store-csi-plugin")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	body := map[string]interface{}{
		"version":            k8sVersion,
		"whitelistNodeCIDRs": parseCommaSeparated(whitelistCIDRs),
	}

	if enabledLB {
		body["enabledLoadBalancerPlugin"] = true
	} else if noLB {
		body["enabledLoadBalancerPlugin"] = false
	}

	if enabledCSI {
		body["enabledBlockStoreCsiPlugin"] = true
	} else if noCSI {
		body["enabledBlockStoreCsiPlugin"] = false
	}

	if dryRun {
		fmt.Println("=== DRY RUN: Update cluster ===")
		fmt.Println()
		fmt.Printf("Cluster ID: %s\n", clusterID)
		fmt.Printf("New version: %s\n", k8sVersion)
		fmt.Printf("Whitelist CIDRs: %s\n", whitelistCIDRs)
		if v, ok := body["enabledLoadBalancerPlugin"]; ok {
			fmt.Printf("Load balancer plugin: %v\n", v)
		}
		if v, ok := body["enabledBlockStoreCsiPlugin"]; ok {
			fmt.Printf("Block store CSI plugin: %v\n", v)
		}
		fmt.Println("\nRun without --dry-run to update.")
		return nil
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Put(fmt.Sprintf("/v1/clusters/%s", clusterID), body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
