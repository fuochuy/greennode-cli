package volume

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new volume",
	RunE:  runCreate,
}

func init() {
	f := createCmd.Flags()

	// Required — registered first, then marked required
	f.String("name", "", "Volume name (required)")
	f.String("volume-type-id", "", "Volume type ID (required)")
	f.String("zone-id", "", "Availability zone ID, e.g. HCM03-1A (required)")
	f.Int("size", 0, "Volume size in GiB (required)")

	for _, name := range []string{"name", "volume-type-id", "size"} {
		if err := createCmd.MarkFlagRequired(name); err != nil {
			panic(fmt.Sprintf("BUG: MarkFlagRequired(%q): %v", name, err))
		}
	}

	// Optional
	f.String("description", "", "Volume description")
	f.String("encryption-type", "", "Encryption type")
	f.Bool("multiattach", false, "Allow the volume to be attached to multiple servers")
	f.Bool("is-poc", false, "Mark as PoC (proof-of-concept) volume")
	f.Bool("dry-run", false, "Validate parameters without creating the volume")
}

func runCreate(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	volumeTypeID, _ := cmd.Flags().GetString("volume-type-id")
	zoneID, _ := cmd.Flags().GetString("zone-id")
	size, _ := cmd.Flags().GetInt("size")
	description, _ := cmd.Flags().GetString("description")
	encryptionType, _ := cmd.Flags().GetString("encryption-type")
	multiattach, _ := cmd.Flags().GetBool("multiattach")
	isPoc, _ := cmd.Flags().GetBool("is-poc")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if dryRun {
		return validateCreate(name, size)
	}

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

	body := map[string]interface{}{
		"name":           name,
		"volumeTypeId":   volumeTypeID,
		"zoneId":         zoneID,
		"size":           size,
		"description":    nilIfEmpty(description),
		"encryptionType": nilIfEmpty(encryptionType),
		"multiattach":    multiattach,
		"isPoc":          isPoc,
	}

	result, err := apiClient.Post(fmt.Sprintf("/v2/%s/volumes", projectID), body)
	if err != nil {
		return fmt.Errorf("failed to create volume: %w", err)
	}

	return outputResult(cmd, cfg, result)
}

func validateCreate(name string, size int) error {
	var errs []string

	if len(name) < 1 {
		errs = append(errs, "volume name cannot be empty")
	}
	if size < 1 {
		errs = append(errs, fmt.Sprintf("volume size %d GiB is invalid (minimum 1 GiB)", size))
	}

	fmt.Println("=== DRY RUN: Validation results ===")
	fmt.Println()
	if len(errs) > 0 {
		fmt.Printf("Found %d error(s):\n", len(errs))
		for _, e := range errs {
			fmt.Printf("  - %s\n", e)
		}
		return fmt.Errorf("dry-run validation failed with %d error(s)", len(errs))
	}

	fmt.Println("All parameters are valid. Run without --dry-run to create the volume.")
	return nil
}
