package vpc

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VPC",
	RunE:  runCreate,
}

func init() {
	f := createCmd.Flags()

	f.String("name", "", "VPC name (required)")
	f.String("cidr", "", "CIDR block for the VPC, e.g. 10.0.0.0/16 (required)")

	for _, name := range []string{"name", "cidr"} {
		if err := createCmd.MarkFlagRequired(name); err != nil {
			panic(fmt.Sprintf("BUG: MarkFlagRequired(%q): %v", name, err))
		}
	}

	f.String("description", "", "VPC description")
	f.Bool("is-default", false, "Mark as the default VPC")
	f.Bool("dry-run", false, "Validate parameters without creating the VPC")
}

func runCreate(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	cidr, _ := cmd.Flags().GetString("cidr")
	description, _ := cmd.Flags().GetString("description")
	isDefault, _ := cmd.Flags().GetBool("is-default")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if dryRun {
		return validateCreate(name, cidr)
	}

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
		"cidr":        cidr,
		"description": nilIfEmpty(description),
		"isDefault":   isDefault,
	}

	result, err := apiClient.Post(fmt.Sprintf("/v2/%s/networks", projectID), body)
	if err != nil {
		return fmt.Errorf("failed to create VPC: %w", err)
	}

	return outputResult(cmd, cfg, result)
}

func validateCreate(name, cidr string) error {
	var errs []string

	if len(name) < 1 {
		errs = append(errs, "VPC name cannot be empty")
	}
	if _, _, err := net.ParseCIDR(cidr); err != nil {
		errs = append(errs, fmt.Sprintf("CIDR %q is invalid: %v", cidr, err))
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

	fmt.Println("All parameters are valid. Run without --dry-run to create the VPC.")
	return nil
}
