package server

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a vServer instance",
	RunE:  runGet,
}

func init() {
	getCmd.Flags().String("server-id", "", "Server ID (required)")
	getCmd.MarkFlagRequired("server-id")
}

func runGet(cmd *cobra.Command, args []string) error {
	serverID, _ := cmd.Flags().GetString("server-id")
	if err := validator.ValidateID(serverID, "server-id"); err != nil {
		return err
	}

	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	result, err := apiClient.Get(fmt.Sprintf("/v2/%s/servers/%s", projectID, serverID), nil)
	if err != nil {
		return fmt.Errorf("failed to get server %s: %w", serverID, err)
	}

	return outputServerDetail(cmd, cfg, transformServerResult(result))
}
