package server

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot a vServer instance",
	RunE:  runReboot,
}

func init() {
	rebootCmd.Flags().String("server-id", "", "Server ID (required)")
	rebootCmd.MarkFlagRequired("server-id")
}

func runReboot(cmd *cobra.Command, args []string) error {
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

	result, err := apiClient.Put(fmt.Sprintf("/v2/%s/servers/%s/reboot", projectID, serverID), nil)
	if err != nil {
		return fmt.Errorf("failed to reboot server %s: %w", serverID, err)
	}

	return outputResult(cmd, cfg, transformServerResult(result))
}
