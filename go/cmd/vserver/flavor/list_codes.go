package flavor

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCodesCmd = &cobra.Command{
	Use:   "list-codes",
	Short: "List available CPU platform codes",
	RunE:  runListCodes,
}

func runListCodes(cmd *cobra.Command, args []string) error {
	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	result, err := apiClient.Get(fmt.Sprintf("/v1/%s/flavor_zones/codes", projectID), nil)
	if err != nil {
		return fmt.Errorf("failed to list CPU platform codes: %w", err)
	}

	return outputResult(cmd, cfg, transformCodes(result))
}

// completeCodes is used by RegisterFlagCompletionFunc for --cpu-platform.
func completeCodes(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	result, err := apiClient.Get(fmt.Sprintf("/v1/%s/flavor_zones/codes", projectID), nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return extractStringSlice(result), cobra.ShellCompDirectiveNoFileComp
}
