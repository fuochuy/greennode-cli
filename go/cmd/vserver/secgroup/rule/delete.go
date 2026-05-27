package rule

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a rule from a security group",
	RunE:  runDelete,
}

func init() {
	deleteCmd.Flags().String("secgroup-id", "", "Security group ID (required)")
	deleteCmd.Flags().String("rule-id", "", "Security group rule ID (required)")
	deleteCmd.MarkFlagRequired("secgroup-id")
	deleteCmd.MarkFlagRequired("rule-id")
}

func runDelete(cmd *cobra.Command, args []string) error {
	secgroupID, _ := cmd.Flags().GetString("secgroup-id")
	ruleID, _ := cmd.Flags().GetString("rule-id")

	if err := validator.ValidateID(secgroupID, "secgroup-id"); err != nil {
		return err
	}
	if err := validator.ValidateID(ruleID, "rule-id"); err != nil {
		return err
	}

	apiClient, cfg, err := createClient(cmd)
	if err != nil {
		return err
	}

	projectID, err := getProjectID(cfg)
	if err != nil {
		return err
	}

	result, err := apiClient.Delete(
		fmt.Sprintf("/v2/%s/secgroups/%s/secgroupRules/%s", projectID, secgroupID, ruleID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete rule %s from security group %s: %w", ruleID, secgroupID, err)
	}

	return outputResult(cmd, cfg, result)
}
