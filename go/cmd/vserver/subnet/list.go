package subnet

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all subnets",
	RunE:  runList,
}

func init() {
	listCmd.Flags().Int("page", 1, "Page number (1-based)")
	listCmd.Flags().Int("page-size", 50, "Number of items per page")
	listCmd.Flags().String("vpc-id", "", "VPC (network) ID (required)")

	if err := listCmd.MarkFlagRequired("vpc-id"); err != nil {
		panic(fmt.Sprintf("BUG: MarkFlagRequired(\"vpc-id\"): %v", err))
	}
}

func runList(cmd *cobra.Command, args []string) error {
	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("page-size")
	vpcID, _ := cmd.Flags().GetString("vpc-id")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	params := map[string]string{
		"page": fmt.Sprintf("%d", page),
		"size": fmt.Sprintf("%d", pageSize),
	}

	result, err := apiClient.Get(fmt.Sprintf("/v2/%s/networks/%s/subnets", projectID, vpcID), params)
	if err != nil {
		return fmt.Errorf("failed to list subnets for VPC %s: %w", vpcID, err)
	}

	return outputResult(cmd, cfg, result)
}
