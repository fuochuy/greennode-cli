package flavor

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listFamiliesCmd = &cobra.Command{
	Use:   "list-families",
	Short: "List available instance families",
	RunE:  runListFamilies,
}

func runListFamilies(cmd *cobra.Command, args []string) error {
	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	result, err := apiClient.Get(fmt.Sprintf("/v1/%s/flavor_zones/families", projectID), nil)
	if err != nil {
		return fmt.Errorf("failed to list instance families: %w", err)
	}

	return outputResult(cmd, cfg, transformFamilies(result))
}

// completeFamilies is used by RegisterFlagCompletionFunc for --instance-family.
func completeFamilies(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	result, err := apiClient.Get(fmt.Sprintf("/v1/%s/flavor_zones/families", projectID), nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return extractStringSlice(result), cobra.ShellCompDirectiveNoFileComp
}
