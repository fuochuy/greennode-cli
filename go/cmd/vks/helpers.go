package vks

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/auth"
	"github.com/vngcloud/greennode-cli/internal/client"
	"github.com/vngcloud/greennode-cli/internal/config"
	"github.com/vngcloud/greennode-cli/internal/formatter"
)

// createClient builds a GreenodeClient from the current command flags.
func createClient(cmd *cobra.Command) (*client.GreenodeClient, error) {
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")
	endpointURL, _ := cmd.Flags().GetString("endpoint-url")
	noVerifySSL, _ := cmd.Flags().GetBool("no-verify-ssl")
	debug, _ := cmd.Flags().GetBool("debug")
	readTimeout, _ := cmd.Flags().GetInt("cli-read-timeout")

	cfg, err := config.LoadConfig(profile)
	if err != nil {
		return nil, err
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("credentials not configured. Run 'grn configure' to set up credentials")
	}

	if region != "" {
		cfg.Region = region
	}

	var baseURL string
	if endpointURL != "" {
		baseURL = endpointURL
	} else {
		baseURL, err = cfg.GetEndpoint("vks")
		if err != nil {
			return nil, err
		}
	}

	if noVerifySSL {
		fmt.Fprintln(os.Stderr, "Warning: SSL certificate verification is disabled. This is not recommended for production use.")
	}

	tokenManager := auth.NewTokenManager(cfg.ClientID, cfg.ClientSecret)
	timeout := time.Duration(readTimeout) * time.Second

	return client.NewGreenodeClient(baseURL, tokenManager, timeout, !noVerifySSL, debug), nil
}

// outputResult formats and prints the API response.
func outputResult(cmd *cobra.Command, data interface{}) error {
	output, _ := cmd.Flags().GetString("output")
	query, _ := cmd.Flags().GetString("query")

	if output == "" {
		// Load from config
		profile, _ := cmd.Flags().GetString("profile")
		cfg, _ := config.LoadConfig(profile)
		if cfg != nil {
			output = cfg.Output
		}
	}
	if output == "" {
		output = "json"
	}

	return formatter.Format(data, output, query, os.Stdout)
}

// parseLabels parses "key1=val1,key2=val2" into a map.
func parseLabels(labelsStr string) map[string]string {
	result := make(map[string]string)
	if labelsStr == "" {
		return result
	}
	for _, pair := range strings.Split(labelsStr, ",") {
		pair = strings.TrimSpace(pair)
		if idx := strings.Index(pair, "="); idx > 0 {
			result[strings.TrimSpace(pair[:idx])] = strings.TrimSpace(pair[idx+1:])
		}
	}
	return result
}

// Taint represents a Kubernetes taint.
type Taint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"`
}

// parseTaints parses "key=value:effect,..." into a slice of Taints.
func parseTaints(taintsStr string) []Taint {
	var result []Taint
	if taintsStr == "" {
		return result
	}
	for _, t := range strings.Split(taintsStr, ",") {
		t = strings.TrimSpace(t)
		if colonIdx := strings.LastIndex(t, ":"); colonIdx > 0 {
			kv := t[:colonIdx]
			effect := strings.TrimSpace(t[colonIdx+1:])
			key, value := kv, ""
			if eqIdx := strings.Index(kv, "="); eqIdx > 0 {
				key = strings.TrimSpace(kv[:eqIdx])
				value = strings.TrimSpace(kv[eqIdx+1:])
			}
			result = append(result, Taint{Key: key, Value: value, Effect: effect})
		}
	}
	return result
}

// parseCommaSeparated splits a comma-separated string into a trimmed slice.
func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
