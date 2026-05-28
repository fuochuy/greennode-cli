package server

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running vServer instance",
	RunE:  runStop,
}

func init() {
	stopCmd.Flags().String("server-id", "", "Server ID (required)")
	stopCmd.MarkFlagRequired("server-id")
}

func runStop(cmd *cobra.Command, args []string) error {
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

	result, err := apiClient.Put(fmt.Sprintf("/v2/%s/servers/%s/stop", projectID, serverID), nil)
	if err != nil {
		return fmt.Errorf("failed to stop server %s: %w", serverID, err)
	}

	return outputResult(cmd, cfg, transformServerResult(result))
}
