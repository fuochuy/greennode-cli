package secgroup

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new security group",
	RunE:  runCreate,
}

func init() {
	f := createCmd.Flags()
	f.String("name", "", "Security group name (required)")
	f.String("description", "", "Security group description")
	createCmd.MarkFlagRequired("name")
}

func runCreate(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")

	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	body := map[string]interface{}{
		"name":        name,
		"description": nilIfEmpty(description),
	}

	result, err := apiClient.Post(fmt.Sprintf("/v2/%s/secgroups", projectID), body)
	if err != nil {
		return fmt.Errorf("failed to create security group: %w", err)
	}

	return outputResult(cmd, cfg, result)
}
