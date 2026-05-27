package volume

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a volume",
	RunE:  runGet,
}

func init() {
	getCmd.Flags().String("volume-id", "", "Volume ID (required)")
	getCmd.MarkFlagRequired("volume-id")
}

func runGet(cmd *cobra.Command, args []string) error {
	volumeID, _ := cmd.Flags().GetString("volume-id")
	if err := validator.ValidateID(volumeID, "volume-id"); err != nil {
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

	result, err := apiClient.Get(fmt.Sprintf("/v2/%s/volumes/%s", projectID, volumeID), nil)
	if err != nil {
		return fmt.Errorf("failed to get volume %s: %w", volumeID, err)
	}

	return outputResult(cmd, cfg, result)
}
