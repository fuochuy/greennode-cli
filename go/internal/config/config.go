package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

// REGIONS maps region names to service endpoints.
var REGIONS = map[string]map[string]string{
	"HCM-3": {
		"vks_endpoint":     "https://vks.api.vngcloud.vn",
		"vserver_endpoint": "https://hcm-3.api.vngcloud.vn/vserver/vserver-gateway",
	},
	"HAN": {
		"vks_endpoint":     "https://vks-han-1.api.vngcloud.vn",
		"vserver_endpoint": "https://han-1.api.vngcloud.vn/vserver/vserver-gateway",
	},
}

// Config holds the resolved CLI configuration.
type Config struct {
	ClientID     string
	ClientSecret string
	Region       string
	Output       string
	Profile      string
	ProjectID    string
	Regions      map[string]map[string]string
}

// DefaultConfigDir returns ~/.greenode
func DefaultConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".greenode")
}

// LoadConfig loads configuration for the given profile.
// Resolution order: env vars -> config files.
func LoadConfig(profile string) (*Config, error) {
	if profile == "" {
		profile = os.Getenv("GRN_PROFILE")
	}
	if profile == "" {
		profile = "default"
	}

	configDir := DefaultConfigDir()
	cfg := &Config{
		Profile: profile,
		Regions: REGIONS,
	}

	// Load credentials — env vars override file
	if v := os.Getenv("GRN_ACCESS_KEY_ID"); v != "" {
		cfg.ClientID = v
	}
	if v := os.Getenv("GRN_SECRET_ACCESS_KEY"); v != "" {
		cfg.ClientSecret = v
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		credsFile := filepath.Join(configDir, "credentials")
		if _, err := os.Stat(credsFile); err == nil {
			iniCreds, err := ini.Load(credsFile)
			if err != nil {
				return nil, fmt.Errorf("failed to parse credentials file: %w", err)
			}
			section, err := iniCreds.GetSection(profile)
			if err != nil {
				return nil, fmt.Errorf("profile '%s' does not exist in %s", profile, credsFile)
			}
			if cfg.ClientID == "" {
				cfg.ClientID = section.Key("client_id").String()
			}
			if cfg.ClientSecret == "" {
				cfg.ClientSecret = section.Key("client_secret").String()
			}
		}
	}

	// Load config file
	configFile := filepath.Join(configDir, "config")
	if _, err := os.Stat(configFile); err == nil {
		iniCfg, err := ini.Load(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}

		sectionName := profile
		if profile != "default" {
			sectionName = "profile " + profile
		}

		section, err := iniCfg.GetSection(sectionName)
		if err != nil && profile == "default" {
			// Try DEFAULT section for default profile
			section = iniCfg.Section("")
		}
		if section != nil {
			if v := section.Key("region").String(); v != "" {
				cfg.Region = v
			}
			if v := section.Key("output").String(); v != "" {
				cfg.Output = v
			}
			if v := section.Key("project_id").String(); v != "" {
				cfg.ProjectID = v
			}
		}
	}

	// Env var overrides for region
	if v := os.Getenv("GRN_DEFAULT_REGION"); v != "" {
		cfg.Region = v
	}

	// Env var override for project_id
	if v := os.Getenv("GRN_DEFAULT_PROJECT_ID"); v != "" {
		cfg.ProjectID = v
	}

	// Default output
	if cfg.Output == "" {
		cfg.Output = "json"
	}

	return cfg, nil
}

// GetEndpoint returns the service endpoint for the configured region.
func (c *Config) GetEndpoint(serviceName string) (string, error) {
	if c.Region == "" {
		return "", fmt.Errorf("region is not configured. Use 'grn configure' or the --region flag")
	}
	regionConfig, ok := c.Regions[c.Region]
	if !ok {
		return "", fmt.Errorf("invalid region: %s", c.Region)
	}
	endpointKey := serviceName + "_endpoint"
	endpoint, ok := regionConfig[endpointKey]
	if !ok {
		return "", fmt.Errorf("endpoint not found for service '%s' in region '%s'", serviceName, c.Region)
	}
	return endpoint, nil
}

// MaskCredential masks a credential string showing only last 4 chars.
func MaskCredential(value string) string {
	if len(value) <= 4 {
		return value
	}
	return "****************" + value[len(value)-4:]
}
