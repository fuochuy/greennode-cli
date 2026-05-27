package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var updateNodegroupCmd = &cobra.Command{
	Use:   "update-nodegroup",
	Short: "Update a node group",
	RunE:  runUpdateNodegroup,
}

func init() {
	f := updateNodegroupCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("nodegroup-id", "", "Node group ID (required)")
	f.String("image-id", "", "Image ID (required)")
	f.String("num-nodes", "", "New number of nodes")
	f.String("security-groups", "", "Security group IDs (comma-separated)")
	f.String("labels", "", "Node labels as key=value pairs (comma-separated)")
	f.String("taints", "", "Node taints as key=value:effect (comma-separated)")
	f.String("auto-scale-min", "", "Auto-scale minimum nodes")
	f.String("auto-scale-max", "", "Auto-scale maximum nodes")
	f.String("upgrade-strategy", "", "Upgrade strategy (SURGE)")
	f.String("upgrade-max-surge", "", "Max surge during upgrade")
	f.String("upgrade-max-unavailable", "", "Max unavailable during upgrade")
	f.Bool("dry-run", false, "Preview update without executing")

	updateNodegroupCmd.MarkFlagRequired("cluster-id")
	updateNodegroupCmd.MarkFlagRequired("nodegroup-id")
	updateNodegroupCmd.MarkFlagRequired("image-id")
}

func runUpdateNodegroup(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	nodegroupID, _ := cmd.Flags().GetString("nodegroup-id")
	imageID, _ := cmd.Flags().GetString("image-id")
	numNodes, _ := cmd.Flags().GetString("num-nodes")
	securityGroups, _ := cmd.Flags().GetString("security-groups")
	labelsStr, _ := cmd.Flags().GetString("labels")
	taintsStr, _ := cmd.Flags().GetString("taints")
	autoScaleMin, _ := cmd.Flags().GetString("auto-scale-min")
	autoScaleMax, _ := cmd.Flags().GetString("auto-scale-max")
	upgradeStrategy, _ := cmd.Flags().GetString("upgrade-strategy")
	upgradeMaxSurge, _ := cmd.Flags().GetString("upgrade-max-surge")
	upgradeMaxUnavail, _ := cmd.Flags().GetString("upgrade-max-unavailable")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}
	if err := validator.ValidateID(nodegroupID, "nodegroup-id"); err != nil {
		return err
	}

	body := map[string]interface{}{
		"imageId": imageID,
	}

	if numNodes != "" {
		body["numNodes"] = toInt(numNodes)
	}
	if securityGroups != "" {
		body["securityGroups"] = parseCommaSeparated(securityGroups)
	}
	if labelsStr != "" {
		body["labels"] = parseLabels(labelsStr)
	}
	if taintsStr != "" {
		body["taints"] = parseTaints(taintsStr)
	}

	if autoScaleMin != "" || autoScaleMax != "" {
		autoScaleConfig := map[string]interface{}{}
		if autoScaleMin != "" {
			autoScaleConfig["minSize"] = toInt(autoScaleMin)
		}
		if autoScaleMax != "" {
			autoScaleConfig["maxSize"] = toInt(autoScaleMax)
		}
		body["autoScaleConfig"] = autoScaleConfig
	}

	if upgradeStrategy != "" || upgradeMaxSurge != "" || upgradeMaxUnavail != "" {
		upgradeConfig := map[string]interface{}{}
		if upgradeStrategy != "" {
			upgradeConfig["strategy"] = upgradeStrategy
		}
		if upgradeMaxSurge != "" {
			upgradeConfig["maxSurge"] = toInt(upgradeMaxSurge)
		}
		if upgradeMaxUnavail != "" {
			upgradeConfig["maxUnavailable"] = toInt(upgradeMaxUnavail)
		}
		body["upgradeConfig"] = upgradeConfig
	}

	if dryRun {
		fmt.Println("=== DRY RUN: Update node group ===")
		fmt.Println()
		fmt.Printf("Cluster ID: %s\n", clusterID)
		fmt.Printf("Node group ID: %s\n", nodegroupID)
		for key, value := range body {
			fmt.Printf("  %s: %v\n", key, value)
		}
		fmt.Println("\nRun without --dry-run to update.")
		return nil
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Put(
		fmt.Sprintf("/v1/clusters/%s/node-groups/%s", clusterID, nodegroupID), body,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}

func toInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
