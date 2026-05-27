package secgroup

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a security group",
	RunE:  runDelete,
}

func init() {
	f := deleteCmd.Flags()
	f.String("secgroup-id", "", "Security group ID (required)")
	f.Bool("force", false, "Skip confirmation prompt")
	deleteCmd.MarkFlagRequired("secgroup-id")
}

func runDelete(cmd *cobra.Command, args []string) error {
	secgroupID, _ := cmd.Flags().GetString("secgroup-id")
	force, _ := cmd.Flags().GetBool("force")

	if err := validator.ValidateID(secgroupID, "secgroup-id"); err != nil {
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

	response, err := apiClient.Get(fmt.Sprintf("/v2/%s/secgroups/%s", projectID, secgroupID), nil)
	if err != nil {
		return fmt.Errorf("failed to fetch security group %s: %w", secgroupID, err)
	}

	sgData, ok := response.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected response type from API: %T", response)
	}

	printSecgroupDeletePreview(sgData)

	if !force {
		fmt.Print("\nAre you sure you want to delete this security group? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	result, err := apiClient.Delete(fmt.Sprintf("/v2/%s/secgroups/%s", projectID, secgroupID), nil)
	if err != nil {
		return fmt.Errorf("failed to delete security group %s: %w", secgroupID, err)
	}

	return outputResult(cmd, cfg, result)
}

func printSecgroupDeletePreview(sg map[string]interface{}) {
	fmt.Println("The following security group will be deleted:")
	fmt.Println()
	fmt.Printf("  ID:          %v\n", sg["id"])
	fmt.Printf("  Name:        %v\n", sg["name"])
	fmt.Printf("  Description: %v\n", sg["description"])
	fmt.Println()
	fmt.Println("This action is irreversible.")
}
