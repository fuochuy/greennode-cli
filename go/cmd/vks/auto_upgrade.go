package vks

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var setAutoUpgradeConfigCmd = &cobra.Command{
	Use:   "set-auto-upgrade-config",
	Short: "Configure auto-upgrade schedule for a cluster",
	RunE:  runSetAutoUpgradeConfig,
}

var deleteAutoUpgradeConfigCmd = &cobra.Command{
	Use:   "delete-auto-upgrade-config",
	Short: "Delete auto-upgrade config for a cluster",
	RunE:  runDeleteAutoUpgradeConfig,
}

func init() {
	// set-auto-upgrade-config flags
	f := setAutoUpgradeConfigCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("weekdays", "", "Days of the week, e.g. Mon,Wed,Fri (required)")
	f.String("time", "", "Time of day in 24h format HH:mm, e.g. 03:00 (required)")
	setAutoUpgradeConfigCmd.MarkFlagRequired("cluster-id")
	setAutoUpgradeConfigCmd.MarkFlagRequired("weekdays")
	setAutoUpgradeConfigCmd.MarkFlagRequired("time")

	// delete-auto-upgrade-config flags
	g := deleteAutoUpgradeConfigCmd.Flags()
	g.String("cluster-id", "", "Cluster ID (required)")
	g.Bool("force", false, "Skip confirmation prompt")
	deleteAutoUpgradeConfigCmd.MarkFlagRequired("cluster-id")
}

func runSetAutoUpgradeConfig(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	weekdays, _ := cmd.Flags().GetString("weekdays")
	timeVal, _ := cmd.Flags().GetString("time")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	body := map[string]interface{}{
		"weekdays": weekdays,
		"time":     timeVal,
	}

	result, err := apiClient.Put(
		fmt.Sprintf("/v1/clusters/%s/auto-upgrade-config", clusterID), body,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}

func runDeleteAutoUpgradeConfig(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	force, _ := cmd.Flags().GetBool("force")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	if !force {
		fmt.Print("Are you sure you want to delete the auto-upgrade config? (yes/no): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(response)) != "yes" {
			fmt.Println("Delete cancelled.")
			return nil
		}
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Delete(
		fmt.Sprintf("/v1/clusters/%s/auto-upgrade-config", clusterID), nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
