package vks

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

var createClusterCmd = &cobra.Command{
	Use:   "create-cluster",
	Short: "Create a new VKS cluster",
	RunE:  runCreateCluster,
}

func init() {
	f := createClusterCmd.Flags()
	// Cluster settings (required)
	f.String("name", "", "Cluster name (required)")
	f.String("k8s-version", "", "Kubernetes version (required)")
	f.String("network-type", "", "Network type: CALICO, CILIUM_OVERLAY, CILIUM_NATIVE_ROUTING (required)")
	f.String("vpc-id", "", "VPC ID (required)")
	f.String("subnet-id", "", "Subnet ID (required)")
	// Node group settings (required)
	f.String("node-group-name", "", "Default node group name (required)")
	f.String("flavor-id", "", "Flavor ID for node group (required)")
	f.String("image-id", "", "Image ID for node group (required)")
	f.String("disk-type", "", "Disk type ID (required)")
	f.String("ssh-key-id", "", "SSH key ID for node group (required)")

	for _, name := range []string{"name", "k8s-version", "network-type", "vpc-id", "subnet-id", "node-group-name", "flavor-id", "image-id", "disk-type", "ssh-key-id"} {
		createClusterCmd.MarkFlagRequired(name)
	}

	// Cluster settings (optional)
	f.String("cidr", "", "CIDR block (required for CALICO and CILIUM_OVERLAY)")
	f.String("description", "", "Cluster description")
	f.Bool("enable-private-cluster", false, "Enable private cluster")
	f.String("release-channel", "STABLE", "Release channel (RAPID, STABLE)")
	f.Bool("no-load-balancer-plugin", false, "Disable load balancer plugin")
	f.Bool("no-block-store-csi-plugin", false, "Disable block store CSI plugin")

	// Node group settings (optional)
	f.Int("disk-size", 100, "Disk size in GiB (20-5000)")
	f.Int("num-nodes", 1, "Number of nodes (0-10)")
	f.Bool("enable-private-nodes", false, "Enable private nodes")
	f.String("security-groups", "", "Security group IDs (comma-separated)")
	f.String("labels", "", "Node labels as key=value pairs (comma-separated)")
	f.String("taints", "", "Node taints as key=value:effect (comma-separated)")
	f.Bool("dry-run", false, "Validate parameters without creating the cluster")
}

func runCreateCluster(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	k8sVersion, _ := cmd.Flags().GetString("k8s-version")
	networkType, _ := cmd.Flags().GetString("network-type")
	vpcID, _ := cmd.Flags().GetString("vpc-id")
	subnetID, _ := cmd.Flags().GetString("subnet-id")
	cidr, _ := cmd.Flags().GetString("cidr")
	description, _ := cmd.Flags().GetString("description")
	enablePrivateCluster, _ := cmd.Flags().GetBool("enable-private-cluster")
	releaseChannel, _ := cmd.Flags().GetString("release-channel")
	noLBPlugin, _ := cmd.Flags().GetBool("no-load-balancer-plugin")
	noCSIPlugin, _ := cmd.Flags().GetBool("no-block-store-csi-plugin")

	ngName, _ := cmd.Flags().GetString("node-group-name")
	flavorID, _ := cmd.Flags().GetString("flavor-id")
	imageID, _ := cmd.Flags().GetString("image-id")
	diskType, _ := cmd.Flags().GetString("disk-type")
	sshKeyID, _ := cmd.Flags().GetString("ssh-key-id")
	diskSize, _ := cmd.Flags().GetInt("disk-size")
	numNodes, _ := cmd.Flags().GetInt("num-nodes")
	enablePrivateNodes, _ := cmd.Flags().GetBool("enable-private-nodes")
	securityGroups, _ := cmd.Flags().GetString("security-groups")
	labels, _ := cmd.Flags().GetString("labels")
	taints, _ := cmd.Flags().GetString("taints")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	// Build node group
	nodeGroup := map[string]interface{}{
		"name":               ngName,
		"flavorId":           flavorID,
		"imageId":            imageID,
		"diskSize":           diskSize,
		"diskType":           diskType,
		"numNodes":           numNodes,
		"enablePrivateNodes": enablePrivateNodes,
		"sshKeyId":           sshKeyID,
		"upgradeConfig": map[string]interface{}{
			"maxSurge":       1,
			"maxUnavailable": 0,
			"strategy":       "SURGE",
		},
		"subnetId":       subnetID,
		"securityGroups": []string{},
	}

	if securityGroups != "" {
		nodeGroup["securityGroups"] = parseCommaSeparated(securityGroups)
	}
	if labels != "" {
		nodeGroup["labels"] = parseLabels(labels)
	}
	if taints != "" {
		nodeGroup["taints"] = parseTaints(taints)
	}

	// Build cluster body
	body := map[string]interface{}{
		"name":                       name,
		"version":                    k8sVersion,
		"networkType":                networkType,
		"vpcId":                      vpcID,
		"subnetId":                   subnetID,
		"enablePrivateCluster":       enablePrivateCluster,
		"releaseChannel":             releaseChannel,
		"enabledBlockStoreCsiPlugin": !noCSIPlugin,
		"enabledLoadBalancerPlugin":  !noLBPlugin,
		"enabledServiceEndpoint":     false,
		"azStrategy":                 "SINGLE",
		"nodeGroups":                 []interface{}{nodeGroup},
	}

	if cidr != "" {
		body["cidr"] = cidr
	}
	if description != "" {
		body["description"] = description
	}

	if dryRun {
		return validateCreateCluster(name, ngName, networkType, cidr, diskSize, numNodes)
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Post("/v1/clusters", body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}

func validateCreateCluster(name, ngName, networkType, cidr string, diskSize, numNodes int) error {
	clusterNameRE := regexp.MustCompile(`^[a-z0-9][a-z0-9\-]{3,18}[a-z0-9]$`)
	ngNameRE := regexp.MustCompile(`^[a-z0-9][a-z0-9-]{3,13}[a-z0-9]$`)

	var errors []string

	if !clusterNameRE.MatchString(name) {
		errors = append(errors, fmt.Sprintf(
			"Cluster name '%s' is invalid. Must be 5-20 chars, lowercase alphanumeric and hyphens, start/end with alphanumeric.", name))
	}

	if (networkType == "CALICO" || networkType == "CILIUM_OVERLAY") && cidr == "" {
		errors = append(errors, fmt.Sprintf("--cidr is required when network-type is %s", networkType))
	}

	if !ngNameRE.MatchString(ngName) {
		errors = append(errors, fmt.Sprintf(
			"Node group name '%s' is invalid. Must be 5-15 chars, lowercase alphanumeric and hyphens, start/end with alphanumeric.", ngName))
	}

	if diskSize < 20 || diskSize > 5000 {
		errors = append(errors, fmt.Sprintf("Disk size %s out of range (20-5000 GiB)", strconv.Itoa(diskSize)))
	}

	if numNodes < 0 || numNodes > 10 {
		errors = append(errors, fmt.Sprintf("Number of nodes %s out of range (0-10)", strconv.Itoa(numNodes)))
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

	fmt.Println("All parameters are valid. Run without --dry-run to create the cluster.")
	return nil
}
