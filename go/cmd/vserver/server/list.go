package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all vServer instances",
	RunE:  runList,
}

func init() {
	listCmd.Flags().Int("page", 1, "Page number (1-based)")
	listCmd.Flags().Int("page-size", 50, "Number of items per page")
	listCmd.Flags().Bool("no-paginate", false, "Disable auto-pagination")
	listCmd.Flags().String("name", "", "Filter by server name (substring match)")
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
	noPaginate, _ := cmd.Flags().GetBool("no-paginate")
	filterName, _ := cmd.Flags().GetString("name")

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
	if filterName != "" {
		params["name"] = filterName
	}
	if noPaginate {
		delete(params, "page")
		delete(params, "size")
	}

	result, err := apiClient.Get(fmt.Sprintf("/v2/%s/servers", projectID), params)
	if err != nil {
		return fmt.Errorf("failed to list servers: %w", err)
	}

	return outputServerList(cmd, cfg, transformServerResult(result))
}
