package flavor

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available flavors",
	RunE:  runList,
}

func init() {
	f := listCmd.Flags()
	f.Int("page", 1, "Page number (1-based)")
	f.Int("page-size", 50, "Number of items per page")
	f.String("zone-id", "", "Filter results by availability zone ID (optional)")
	f.String("family", "", "Filter by instance family (run 'flavor list-families' to see options) (required)")
	f.String("code", "", "Filter by CPU platform code (run 'flavor list-codes' to see options) (required)")

	listCmd.RegisterFlagCompletionFunc("family", completeFamilies) //nolint:errcheck
	listCmd.RegisterFlagCompletionFunc("code", completeCodes)      //nolint:errcheck
}

func runList(cmd *cobra.Command, args []string) error {
	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("page-size")
	zoneID, _ := cmd.Flags().GetString("zone-id")
	instanceFamily, _ := cmd.Flags().GetString("family")
	cpuPlatform, _ := cmd.Flags().GetString("code")

	if instanceFamily == "" {
		return suggestFamilies(apiClient, projectID)
	}
	if cpuPlatform == "" {
		return suggestCodes(apiClient, projectID)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	params := map[string]string{
		"page": fmt.Sprintf("%d", page),
		"size": fmt.Sprintf("%d", pageSize),
	}
	if zoneID != "" {
		params["zoneId"] = zoneID
	}

	result, err := apiClient.Get(fmt.Sprintf("/v1/%s/flavors/families/%s/platforms/%s", projectID, instanceFamily, cpuPlatform), params)
	if err != nil {
		return fmt.Errorf("failed to list flavors: %w", err)
	}

	return outputResult(cmd, cfg, filterFlavors(result))
}

func suggestFamilies(apiClient interface {
	Get(string, map[string]string) (interface{}, error)
}, projectID string) error {
	result, err := apiClient.Get(fmt.Sprintf("/v1/%s/flavor_zones/families", projectID), nil)
	if err != nil {
		return fmt.Errorf("--family is required. Also failed to fetch available families: %w", err)
	}
	fmt.Fprintln(os.Stderr, "Flag --family is required. Available instance families:")
	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	enc.Encode(transformFamilies(result)) //nolint:errcheck
	return fmt.Errorf("flag --family is required")
}

func suggestCodes(apiClient interface {
	Get(string, map[string]string) (interface{}, error)
}, projectID string) error {
	result, err := apiClient.Get(fmt.Sprintf("/v1/%s/flavor_zones/codes", projectID), nil)
	if err != nil {
		return fmt.Errorf("--code is required. Also failed to fetch available codes: %w", err)
	}
	fmt.Fprintln(os.Stderr, "Flag --code is required. Available CPU platform codes:")
	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	enc.Encode(transformCodes(result)) //nolint:errcheck
	return fmt.Errorf("flag --code is required")
}

// filterFlavors keeps only flavors with remainingVms > 1 and removes
// the metadata and remainingVms fields from each item.
func filterFlavors(result interface{}) interface{} {
	// Unwrap envelope: {"data": [...]} or bare [...]
	var items []interface{}
	var envelope map[string]interface{}

	switch v := result.(type) {
	case []interface{}:
		items = v
	case map[string]interface{}:
		envelope = v
		if d, ok := v["data"].([]interface{}); ok {
			items = d
		} else {
			return result
		}
	default:
		return result
	}

	filtered := make([]interface{}, 0, len(items))
	for _, item := range items {
		flavor, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		remaining, _ := flavor["remainingVms"].(float64)
		if remaining <= 1 {
			continue
		}
		delete(flavor, "metaData")
		delete(flavor, "remainingVms")
		delete(flavor, "zoneId")
		delete(flavor, "isSoldOut")
		filtered = append(filtered, flavor)
	}

	if envelope != nil {
		out := make(map[string]interface{}, len(envelope))
		for k, v := range envelope {
			out[k] = v
		}
		out["data"] = filtered
		return out
	}
	return filtered
}
