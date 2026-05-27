package image

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var supportedTypes = []string{"os", "gpu"}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available images",
	RunE:  runList,
}

func init() {
	f := listCmd.Flags()
	f.String("type", "", "Image type: os or gpu")
	f.Int("page", 1, "Page number (1-based)")
	f.Int("page-size", 50, "Number of items per page")
	f.String("image-version", "", "Filter by imageVersion (client-side substring match)")

	listCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		return supportedTypes, cobra.ShellCompDirectiveNoFileComp
	})
}

func runList(cmd *cobra.Command, args []string) error {
	imageType, _ := cmd.Flags().GetString("type")
	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("page-size")
	imageVersion, _ := cmd.Flags().GetString("image-version")

	if imageType == "" {
		fmt.Fprintln(os.Stderr, "Flag --type is required. Available image types:")
		for _, t := range supportedTypes {
			fmt.Fprintf(os.Stderr, "  - %s\n", t)
		}
		return fmt.Errorf("flag --type is required")
	}

	if !validImageType(imageType) {
		fmt.Fprintln(os.Stderr, "Invalid --type value. Available image types:")
		for _, t := range supportedTypes {
			fmt.Fprintf(os.Stderr, "  - %s\n", t)
		}
		return fmt.Errorf("invalid image type %q (must be one of: os, gpu)", imageType)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}

	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	params := map[string]string{
		"page": fmt.Sprintf("%d", page),
		"size": fmt.Sprintf("%d", pageSize),
	}

	result, err := apiClient.Get(fmt.Sprintf("/v1/%s/images/%s", projectID, imageType), params)
	if err != nil {
		return fmt.Errorf("failed to list %s images: %w", imageType, err)
	}

	return outputResult(cmd, cfg, filterByImageVersion(dropImageFields(result), imageVersion))
}

func filterByImageVersion(result interface{}, version string) interface{} {
	if version == "" {
		return result
	}
	filter := func(items []interface{}) []interface{} {
		out := make([]interface{}, 0, len(items))
		for _, item := range items {
			obj, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			v, _ := obj["imageVersion"].(string)
			if strings.Contains(strings.ToLower(v), strings.ToLower(version)) {
				out = append(out, item)
			}
		}
		return out
	}

	switch v := result.(type) {
	case []interface{}:
		return filter(v)
	case map[string]interface{}:
		for _, key := range []string{"images", "data"} {
			if items, ok := v[key].([]interface{}); ok {
				v[key] = filter(items)
				return v
			}
		}
	}
	return result
}

func dropImageFields(result interface{}) interface{} {
	drop := func(items []interface{}) {
		for _, item := range items {
			if obj, ok := item.(map[string]interface{}); ok {
				delete(obj, "flavorZoneIds")
				delete(obj, "zoneId")
			}
		}
	}

	switch v := result.(type) {
	case []interface{}:
		drop(v)
		return result
	case map[string]interface{}:
		// Unwrap envelope — keep only the images list
		for _, key := range []string{"images", "data"} {
			if items, ok := v[key].([]interface{}); ok {
				drop(items)
				return map[string]interface{}{key: items}
			}
		}
	}
	return result
}

func validImageType(t string) bool {
	for _, s := range supportedTypes {
		if s == t {
			return true
		}
	}
	return false
}
