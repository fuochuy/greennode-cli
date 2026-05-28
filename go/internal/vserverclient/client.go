package vserverclient

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/auth"
	"github.com/vngcloud/greennode-cli/internal/client"
	"github.com/vngcloud/greennode-cli/internal/config"
	"github.com/vngcloud/greennode-cli/internal/formatter"
)

// BuildClient creates a GreenodeClient from cobra command flags.
func BuildClient(cmd *cobra.Command) (*client.GreenodeClient, *config.Config, error) {
	profile, _ := cmd.Flags().GetString("profile")
	region, _ := cmd.Flags().GetString("region")
	endpointURL, _ := cmd.Flags().GetString("endpoint-url")
	noVerifySSL, _ := cmd.Flags().GetBool("no-verify-ssl")
	debug, _ := cmd.Flags().GetBool("debug")
	readTimeout, _ := cmd.Flags().GetInt("cli-read-timeout")

	cfg, err := config.LoadConfig(profile)
	if err != nil {
		return nil, nil, err
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, nil, fmt.Errorf("credentials not configured. Run 'grn configure' to set up credentials")
	}

	if region != "" {
		cfg.Region = region
	}

	var baseURL string
	if endpointURL != "" {
		baseURL = endpointURL
	} else {
		baseURL, err = cfg.GetEndpoint("vserver")
		if err != nil {
			return nil, nil, err
		}
	}

	if noVerifySSL {
		fmt.Fprintln(os.Stderr, "Warning: SSL certificate verification is disabled. This is not recommended for production use.")
	}

	tokenManager := auth.NewTokenManager(cfg.ClientID, cfg.ClientSecret)
	timeout := time.Duration(readTimeout) * time.Second

	return client.NewGreenodeClient(baseURL, tokenManager, timeout, !noVerifySSL, debug), cfg, nil
}

// ProjectID extracts and validates the project ID from config.
func ProjectID(cfg *config.Config) (string, error) {
	if cfg.ProjectID == "" {
		return "", fmt.Errorf("project_id is not configured. Run 'grn configure' or set GRN_DEFAULT_PROJECT_ID")
	}
	return cfg.ProjectID, nil
}

// Output formats and writes the API result to stdout.
func Output(cmd *cobra.Command, cfg *config.Config, data interface{}) error {
	output, _ := cmd.Flags().GetString("output")
	query, _ := cmd.Flags().GetString("query")

	if output == "" && cfg != nil {
		output = cfg.Output
	}
	if output == "" {
		output = "json"
	}

	return formatter.Format(data, output, query, os.Stdout)
}

// OutputWithColumns formats and writes the API result to stdout.
// When the output format is "table", only the specified columns are shown in the given order.
func OutputWithColumns(cmd *cobra.Command, cfg *config.Config, data interface{}, columns []string) error {
	output, _ := cmd.Flags().GetString("output")
	query, _ := cmd.Flags().GetString("query")

	if output == "" && cfg != nil {
		output = cfg.Output
	}
	if output == "" {
		output = "json"
	}

	if output == "table" && len(columns) > 0 {
		return formatter.FormatTableWithColumns(data, columns, query, os.Stdout)
	}
	return formatter.Format(data, output, query, os.Stdout)
}
