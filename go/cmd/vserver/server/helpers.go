package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/client"
	"github.com/vngcloud/greennode-cli/internal/config"
	"github.com/vngcloud/greennode-cli/internal/vserverclient"
)

func createClient(cmd *cobra.Command) (*client.GreenodeClient, *config.Config, error) {
	return vserverclient.BuildClient(cmd)
}

func getProjectID(cfg *config.Config) (string, error) {
	return vserverclient.ProjectID(cfg)
}

func outputResult(cmd *cobra.Command, cfg *config.Config, data interface{}) error {
	return vserverclient.Output(cmd, cfg, data)
}

func parseCommaSeparated(s string) []string {
	result := []string{}
	if s == "" {
		return result
	}
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func suggestZones(apiClient *client.GreenodeClient, projectID string) error {
	return vserverclient.SuggestZoneIDs(apiClient, projectID)
}

func suggestVPCs(apiClient *client.GreenodeClient, projectID string) error {
	result, err := apiClient.Get(fmt.Sprintf("/v2/%s/networks", projectID), map[string]string{"page": "1", "size": "50"})
	if err != nil {
		return fmt.Errorf("--network-id is required (also failed to fetch VPCs: %w)", err)
	}
	fmt.Fprintln(os.Stderr, "Flag --network-id is required. Available VPCs:")
	printItems(result, []string{"listData"}, func(obj map[string]interface{}) {
		fmt.Fprintf(os.Stderr, "  - %-40s  name: %v  cidr: %v\n", obj["id"], obj["displayName"], obj["cidr"])
	})
	return fmt.Errorf("flag --network-id is required")
}

func suggestSubnets(apiClient *client.GreenodeClient, projectID, networkID string) error {
	result, err := apiClient.Get(fmt.Sprintf("/v2/%s/networks/%s/subnets", projectID, networkID), map[string]string{"page": "1", "size": "50"})
	if err != nil {
		return fmt.Errorf("--subnet-id is required (also failed to fetch subnets: %w)", err)
	}
	fmt.Fprintln(os.Stderr, "Flag --subnet-id is required. Available subnets for VPC "+networkID+":")
	printItems(result, []string{"data"}, func(obj map[string]interface{}) {
		id := obj["uuid"]
		if id == nil {
			id = obj["id"]
		}
		fmt.Fprintf(os.Stderr, "  - %-40v  name: %v  cidr: %v\n", id, obj["name"], obj["cidr"])
	})
	return fmt.Errorf("flag --subnet-id is required")
}

func suggestImages() error {
	fmt.Fprintln(os.Stderr, "Flag --image-id is required. To see available images, run:")
	fmt.Fprintln(os.Stderr, "  grn vserver image list --type os")
	fmt.Fprintln(os.Stderr, "  grn vserver image list --type gpu")
	return fmt.Errorf("flag --image-id is required")
}

func suggestFlavors() error {
	fmt.Fprintln(os.Stderr, "Flag --flavor-id is required. To see available flavors, run:")
	fmt.Fprintln(os.Stderr, "  grn vserver flavor list-families          # see instance families")
	fmt.Fprintln(os.Stderr, "  grn vserver flavor list-codes             # see CPU platform codes")
	fmt.Fprintln(os.Stderr, "  grn vserver flavor list --family <family> --code <code>")
	return fmt.Errorf("flag --flavor-id is required")
}

func suggestRootDiskTypes(zoneID string) error {
	fmt.Fprintln(os.Stderr, "Flag --root-disk-type-id is required. To see available volume types, run:")
	fmt.Fprintf(os.Stderr, "  grn vserver volume-type list --zone-id %s\n", zoneID)
	return fmt.Errorf("flag --root-disk-type-id is required")
}

// printItems iterates the items array from a response envelope and calls fn for each object.
func printItems(result interface{}, keys []string, fn func(map[string]interface{})) {
	var items []interface{}
	switch v := result.(type) {
	case []interface{}:
		items = v
	case map[string]interface{}:
		for _, key := range keys {
			if d, ok := v[key].([]interface{}); ok {
				items = d
				break
			}
		}
	}
	for _, item := range items {
		if obj, ok := item.(map[string]interface{}); ok {
			fn(obj)
		}
	}
}
