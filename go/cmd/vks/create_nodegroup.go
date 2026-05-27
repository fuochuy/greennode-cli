package vks

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var createNodegroupCmd = &cobra.Command{
	Use:   "create-nodegroup",
	Short: "Create a new node group",
	RunE:  runCreateNodegroup,
}

func init() {
	f := createNodegroupCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("name", "", "Node group name (required)")
	f.String("image-id", "", "Image ID (required)")
	f.String("flavor-id", "", "Flavor ID (required)")
	f.String("disk-type", "", "Disk type ID (required)")
	f.String("ssh-key-id", "", "SSH key ID (required)")
	f.Bool("enable-private-nodes", false, "Enable private nodes")
	f.Int("num-nodes", 1, "Number of nodes (0-10)")
	f.Int("disk-size", 100, "Disk size in GiB (20-5000)")
	f.String("security-groups", "", "Security group IDs (comma-separated)")
	f.String("subnet-id", "", "Subnet ID for node group")
	f.String("labels", "", "Node labels as key=value pairs (comma-separated)")
	f.String("taints", "", "Node taints as key=value:effect (comma-separated)")
	f.Bool("enable-encryption-volume", false, "Enable volume encryption")
	f.Bool("dry-run", false, "Validate parameters without creating")

	for _, name := range []string{"cluster-id", "name", "image-id", "flavor-id", "disk-type", "ssh-key-id"} {
		createNodegroupCmd.MarkFlagRequired(name)
	}
}

func runCreateNodegroup(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	name, _ := cmd.Flags().GetString("name")
	imageID, _ := cmd.Flags().GetString("image-id")
	flavorID, _ := cmd.Flags().GetString("flavor-id")
	diskType, _ := cmd.Flags().GetString("disk-type")
	sshKeyID, _ := cmd.Flags().GetString("ssh-key-id")
	enablePrivateNodes, _ := cmd.Flags().GetBool("enable-private-nodes")
	numNodes, _ := cmd.Flags().GetInt("num-nodes")
	diskSize, _ := cmd.Flags().GetInt("disk-size")
	securityGroups, _ := cmd.Flags().GetString("security-groups")
	subnetID, _ := cmd.Flags().GetString("subnet-id")
	labelsStr, _ := cmd.Flags().GetString("labels")
	taintsStr, _ := cmd.Flags().GetString("taints")
	enableEncryption, _ := cmd.Flags().GetBool("enable-encryption-volume")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	body := map[string]interface{}{
		"name":                    name,
		"numNodes":                numNodes,
		"imageId":                 imageID,
		"flavorId":                flavorID,
		"diskSize":                diskSize,
		"diskType":                diskType,
		"enablePrivateNodes":      enablePrivateNodes,
		"sshKeyId":                sshKeyID,
		"enabledEncryptionVolume": enableEncryption,
		"securityGroups":          []string{},
		"upgradeConfig": map[string]interface{}{
			"maxSurge":       1,
			"maxUnavailable": 0,
			"strategy":       "SURGE",
		},
	}

	if securityGroups != "" {
		body["securityGroups"] = parseCommaSeparated(securityGroups)
	}
	if subnetID != "" {
		body["subnetId"] = subnetID
	}
	if labelsStr != "" {
		body["labels"] = parseLabels(labelsStr)
	}
	if taintsStr != "" {
		body["taints"] = parseTaints(taintsStr)
	}

	if dryRun {
		return validateCreateNodegroup(name, diskSize, numNodes)
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Post(
		fmt.Sprintf("/v1/clusters/%s/node-groups", clusterID), body,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}

func validateCreateNodegroup(name string, diskSize, numNodes int) error {
	ngNameRE := regexp.MustCompile(`^[a-z0-9][a-z0-9-]{3,13}[a-z0-9]$`)
	var errors []string

	if !ngNameRE.MatchString(name) {
		errors = append(errors, fmt.Sprintf(
			"Node group name '%s' is invalid. Must be 5-15 chars, lowercase alphanumeric and hyphens.", name))
	}
	if diskSize < 20 || diskSize > 5000 {
		errors = append(errors, fmt.Sprintf("Disk size %d out of range (20-5000 GiB)", diskSize))
	}
	if numNodes < 0 || numNodes > 10 {
		errors = append(errors, fmt.Sprintf("Number of nodes %d out of range (0-10)", numNodes))
	}

	fmt.Println("=== DRY RUN: Validation results ===")
	fmt.Println()
	if len(errors) > 0 {
		fmt.Printf("Found %d error(s):\n", len(errors))
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
		os.Exit(1)
	}

	fmt.Println("All parameters are valid. Run without --dry-run to create.")
	return nil
}
