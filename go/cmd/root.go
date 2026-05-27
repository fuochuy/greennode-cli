package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/cmd/configure"
	"github.com/vngcloud/greennode-cli/cmd/mcp"
	"github.com/vngcloud/greennode-cli/cmd/vks"
	"github.com/vngcloud/greennode-cli/cmd/vserver"
)

const cliVersion = "1.3.1"

// Global flags
var (
	Profile           string
	Region            string
	Output            string
	Query             string
	EndpointURL       string
	NoVerifySSL       bool
	Debug             bool
	CLIReadTimeout    int
	CLIConnectTimeout int
	Color             string
)

var rootCmd = &cobra.Command{
	Use:     "grn",
	Short:   "GreenNode CLI - unified command-line tool for GreenNode (VNG Cloud) services",
	Version: fmt.Sprintf("%s Go/%s %s/%s", cliVersion, runtime.Version()[2:], runtime.GOOS, runtime.GOARCH),
	Long: `GreenNode CLI (grn) is a unified command-line tool for managing
GreenNode (VNG Cloud) services including VKS (VNG Kubernetes Service).

To get started, run:
  grn configure

For help on any command:
  grn <command> --help`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&Profile, "profile", "", "Use a specific profile from credentials file")
	rootCmd.PersistentFlags().StringVar(&Region, "region", "", "The region to use (e.g. HCM-3, HAN)")
	rootCmd.PersistentFlags().StringVar(&Output, "output", "", "The output format (json, text, table)")
	rootCmd.PersistentFlags().StringVar(&Query, "query", "", "JMESPath query to filter output")
	rootCmd.PersistentFlags().StringVar(&EndpointURL, "endpoint-url", "", "Override the service endpoint URL")
	rootCmd.PersistentFlags().BoolVar(&NoVerifySSL, "no-verify-ssl", false, "Disable SSL certificate verification")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().IntVar(&CLIReadTimeout, "cli-read-timeout", 30, "HTTP read timeout in seconds")
	rootCmd.PersistentFlags().IntVar(&CLIConnectTimeout, "cli-connect-timeout", 30, "HTTP connect timeout in seconds")
	rootCmd.PersistentFlags().StringVar(&Color, "color", "auto", "Color output (on, off, auto)")

	rootCmd.SetVersionTemplate("grn-cli/{{.Version}}\n")
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	rootCmd.AddCommand(configure.ConfigureCmd)
	rootCmd.AddCommand(vks.VksCmd)
	rootCmd.AddCommand(vserver.VServerCmd)
	rootCmd.AddCommand(mcp.McpCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
