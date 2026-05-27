package volumetype

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available volume types for a zone",
	RunE:  runList,
}

func init() {
	f := listCmd.Flags()
	f.String("zone-id", "", "Availability zone ID (required)")
	f.String("type", "", "Volume type zone name to filter by (e.g. SSD, NVMe)")
}

func runList(cmd *cobra.Command, args []string) error {
	zoneID, _ := cmd.Flags().GetString("zone-id")
	typeName, _ := cmd.Flags().GetString("type")

	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	if zoneID == "" {
		return suggestZones(apiClient, projectID)
	}

	zoneResult, err := apiClient.Get(
		fmt.Sprintf("/v1/%s/volume_type_zones", projectID),
		map[string]string{"zoneId": zoneID},
	)
	if err != nil {
		return fmt.Errorf("failed to fetch volume type zones for %s: %w", zoneID, err)
	}

	if typeName == "" {
		fmt.Fprintln(os.Stderr, "Flag --type is required. Available volume type zones:")
		for _, name := range extractVolumeTypeZoneNames(zoneResult) {
			fmt.Fprintf(os.Stderr, "  - %s\n", name)
		}
		return fmt.Errorf("flag --type is required")
	}

	volumeTypeZoneID, err := extractVolumeTypeZoneID(zoneResult, typeName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Available volume type zones:")
		for _, name := range extractVolumeTypeZoneNames(zoneResult) {
			fmt.Fprintf(os.Stderr, "  - %s\n", name)
		}
		return fmt.Errorf("volume type zone %q not found in zone %s", typeName, zoneID)
	}

	result, err := apiClient.Get(
		fmt.Sprintf("/v1/%s/%s/volume_types", projectID, volumeTypeZoneID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to list volume types: %w", err)
	}

	return outputResult(cmd, cfg, result)
}
