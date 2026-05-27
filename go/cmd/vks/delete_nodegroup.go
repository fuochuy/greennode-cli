package vks

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var deleteNodegroupCmd = &cobra.Command{
	Use:   "delete-nodegroup",
	Short: "Delete a node group",
	RunE:  runDeleteNodegroup,
}

func init() {
	f := deleteNodegroupCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("nodegroup-id", "", "Node group ID (required)")
	f.Bool("force-delete", false, "Force delete on API side")
	f.Bool("dry-run", false, "Preview what will be deleted without executing")
	f.Bool("force", false, "Skip confirmation prompt")

	deleteNodegroupCmd.MarkFlagRequired("cluster-id")
	deleteNodegroupCmd.MarkFlagRequired("nodegroup-id")
}

func runDeleteNodegroup(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	nodegroupID, _ := cmd.Flags().GetString("nodegroup-id")
	forceDelete, _ := cmd.Flags().GetBool("force-delete")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	force, _ := cmd.Flags().GetBool("force")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}
	if err := validator.ValidateID(nodegroupID, "nodegroup-id"); err != nil {
		return err
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	// Fetch nodegroup info for preview
	ng, err := apiClient.Get(
		fmt.Sprintf("/v1/clusters/%s/node-groups/%s", clusterID, nodegroupID), nil,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch node group: %w", err)
	}

	ngMap, _ := ng.(map[string]interface{})
	fmt.Println("The following node group will be deleted:")
	fmt.Println()
	fmt.Printf("  ID:      %v\n", ngMap["id"])
	fmt.Printf("  Name:    %v\n", ngMap["name"])
	fmt.Printf("  Status:  %v\n", ngMap["status"])
	fmt.Printf("  Nodes:   %v\n", ngMap["numNodes"])
	fmt.Println()
	fmt.Println("This action is irreversible.")

	if dryRun {
		fmt.Println("Run without --dry-run to delete.")
		return nil
	}

	if !force {
		fmt.Print("\nAre you sure you want to delete this node group? (yes/no): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(response)) != "yes" {
			fmt.Println("Delete cancelled.")
			return nil
		}
	}

	params := map[string]string{}
	if forceDelete {
		params["forceDelete"] = "true"
	}

	var paramsArg map[string]string
	if len(params) > 0 {
		paramsArg = params
	}

	result, err := apiClient.Delete(
		fmt.Sprintf("/v1/clusters/%s/node-groups/%s", clusterID, nodegroupID), paramsArg,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
