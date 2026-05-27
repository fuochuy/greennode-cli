package vks

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var waitClusterActiveCmd = &cobra.Command{
	Use:   "wait-cluster-active",
	Short: "Wait until cluster reaches ACTIVE status",
	RunE:  runWaitClusterActive,
}

func init() {
	f := waitClusterActiveCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.Int("delay", 15, "Seconds between polls")
	f.Int("max-attempts", 40, "Maximum poll attempts")

	waitClusterActiveCmd.MarkFlagRequired("cluster-id")
}

func runWaitClusterActive(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	delay, _ := cmd.Flags().GetInt("delay")
	maxAttempts, _ := cmd.Flags().GetInt("max-attempts")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		result, err := apiClient.Get(fmt.Sprintf("/v1/clusters/%s", clusterID), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\rWaiting for cluster %s: error fetching status (attempt %d/%d)",
				clusterID, attempt, maxAttempts)
		} else {
			resultMap, _ := result.(map[string]interface{})
			status, _ := resultMap["status"].(string)

			fmt.Fprintf(os.Stderr, "\rWaiting for cluster %s: %s (attempt %d/%d)",
				clusterID, status, attempt, maxAttempts)

			if status == "ACTIVE" {
				fmt.Fprintln(os.Stderr)
				fmt.Println("Successfully waited for cluster to reach ACTIVE")
				return nil
			}

			if status == "ERROR" || status == "FAILED" {
				fmt.Fprintln(os.Stderr)
				fmt.Fprintf(os.Stderr, "Waiter failed: cluster reached %s\n", status)
				os.Exit(255)
			}
		}

		if attempt < maxAttempts {
			time.Sleep(time.Duration(delay) * time.Second)
		}
	}

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Waiter timed out after %d attempts\n", maxAttempts)
	os.Exit(255)
	return nil
}
