package configure

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/config"
)

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run:   runSet,
}

func runSet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]
	profile := cmd.Flag("profile").Value.String()
	if profile == "" {
		profile = os.Getenv("GRN_PROFILE")
	}
	if profile == "" {
		profile = "default"
	}

	writer := config.NewConfigFileWriter()

	switch key {
	case "client_id":
		cfg, _ := config.LoadConfig(profile)
		if err := writer.WriteCredentials(profile, value, cfg.ClientSecret); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "client_secret":
		cfg, _ := config.LoadConfig(profile)
		if err := writer.WriteCredentials(profile, cfg.ClientID, value); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "region":
		cfg, _ := config.LoadConfig(profile)
		if err := writer.WriteConfig(profile, value, cfg.Output, cfg.ProjectID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "output":
		cfg, _ := config.LoadConfig(profile)
		if err := writer.WriteConfig(profile, cfg.Region, value, cfg.ProjectID); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "project_id":
		cfg, _ := config.LoadConfig(profile)
		if err := writer.WriteConfig(profile, cfg.Region, cfg.Output, value); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown configuration key: %s\n", key)
		os.Exit(1)
	}

	fmt.Printf("Set '%s' to '%s' for profile '%s'.\n", key, value, profile)
}
