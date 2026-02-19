// Package commands provides CLI commands.
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Initialize a new Terraship policy",
	Long: `Initialize a new Terraship policy file in the specified directory.

This command creates a sample policy file with security and compliance rules
that you can customize for your infrastructure validation needs.

Examples:
  terraship init                    # Create policy in current directory
  terraship init ./my-project       # Create policy in my-project/policies
  terraship init --policy ./custom  # Create policy file named custom.yml`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

var (
	policyFileName string
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&policyFileName, "policy", "terraship-policy.yml", "Name of the policy file to create")
}

func runInit(cmd *cobra.Command, args []string) error {
	directory := "."
	if len(args) > 0 {
		directory = args[0]
	}

	// Create policies subdirectory if it doesn't exist
	policiesDir := filepath.Join(directory, "policies")
	if err := os.MkdirAll(policiesDir, 0755); err != nil {
		return fmt.Errorf("failed to create policies directory: %w", err)
	}

	// Write policy file with default content
	policyFile := filepath.Join(policiesDir, policyFileName)
	policyContent := getDefaultPolicy()
	if err := os.WriteFile(policyFile, []byte(policyContent), 0644); err != nil {
		return fmt.Errorf("failed to write policy file: %w", err)
	}

	ruleCount := countRulesInPolicy(policyContent)

	fmt.Printf("✓ Terraship policy initialized successfully!\n\n")
	fmt.Printf("Policy file created at: %s\n", policyFile)
	fmt.Printf("Policy contains %d comprehensive rules covering:\n", ruleCount)
	fmt.Println("  • Security best practices")
	fmt.Println("  • Compliance requirements")
	fmt.Println("  • Multi-cloud support (AWS, Azure, GCP)")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. Review and customize the policy rules:\n")
	fmt.Printf("     %s\n\n", policyFile)
	fmt.Printf("  2. Validate your infrastructure:\n")
	fmt.Printf("     terraship validate ./ --policy %s\n\n", filepath.Join("policies", policyFileName))
	fmt.Println("  3. For required environment variables:")
	fmt.Println("     terraship validate --help")
	fmt.Println()

	return nil
}

func countRulesInPolicy(content string) int {
	count := strings.Count(content, "\n  - name:")
	return count
}

func getDefaultPolicy() string {
	return `version: "1.0"
name: "Multi-Cloud Security and Compliance Policy"
description: "Comprehensive policy for AWS, Azure, and GCP resources covering security, compliance, and best practices"

rules:
  # Tagging and Governance
  - name: "required-tags"
    description: "Ensure all resources have required tags"
    severity: "error"
    category: "compliance"
    enabled: true
    resource_types:
      - "aws_*"
      - "azurerm_*"
      - "google_*"
    conditions:
      tags.required:
        - "Environment"
        - "Owner"
        - "Project"
    message: "Resources must have Environment, Owner, and Project tags"
    remediation: "Add the required tags to your resource configuration"

  # Encryption
  - name: "encryption-at-rest"
    description: "Ensure encryption at rest is enabled"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_s3_bucket"
      - "aws_ebs_volume"
      - "aws_rds_*"
      - "azurerm_storage_account"
      - "azurerm_managed_disk"
      - "google_storage_bucket"
      - "google_compute_disk"
    conditions:
      encryption.enabled: true
    message: "Encryption at rest must be enabled"
    remediation: "Enable server-side encryption for your resource"

  # Public Access
  - name: "block-public-access"
    description: "Block public access to sensitive resources"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_s3_bucket"
      - "aws_db_instance"
      - "azurerm_storage_account"
      - "azurerm_sql_server"
      - "google_storage_bucket"
      - "google_sql_database_instance"
    conditions:
      public_access.blocked: true
    message: "Public access should be blocked"
    remediation: "Configure the resource to block public access"

  # Versioning
  - name: "enable-versioning"
    description: "Enable versioning for storage resources"
    severity: "warning"
    category: "compliance"
    enabled: true
    resource_types:
      - "aws_s3_bucket"
      - "google_storage_bucket"
    conditions:
      versioning.enabled: true
    message: "Versioning should be enabled for data protection"
    remediation: "Enable versioning in your bucket configuration"

  # Logging
  - name: "enable-logging"
    description: "Enable logging for audit trail"
    severity: "warning"
    category: "compliance"
    enabled: true
    resource_types:
      - "aws_s3_bucket"
      - "aws_cloudtrail"
      - "azurerm_storage_account"
      - "google_storage_bucket"
    conditions:
      logging.enabled: true
    message: "Logging should be enabled for audit purposes"
    remediation: "Configure access logging or diagnostic settings"

  # IAM Best Practices
  - name: "iam-least-privilege"
    description: "Ensure IAM policies follow least privilege principle"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_iam_*"
      - "azurerm_role_*"
      - "google_project_iam_*"
    conditions:
      iam.least_privilege: true
    message: "IAM policies should not use wildcard permissions"
    remediation: "Specify explicit permissions instead of using wildcards"

  # Network Security
  - name: "use-private-subnet"
    description: "Deploy resources in private subnets"
    severity: "warning"
    category: "security"
    enabled: true
    resource_types:
      - "aws_instance"
      - "aws_db_instance"
      - "azurerm_virtual_machine"
      - "google_compute_instance"
    conditions:
      network.private_subnet: true
    message: "Resources should be deployed in private subnets"
    remediation: "Configure the resource to use a private subnet"

  # Backup Configuration
  - name: "backup-enabled"
    description: "Ensure backup is configured"
    severity: "warning"
    category: "reliability"
    enabled: true
    resource_types:
      - "aws_db_instance"
      - "aws_rds_cluster"
      - "azurerm_sql_database"
      - "google_sql_database_instance"
    conditions:
      backup.enabled: true
    message: "Backup should be configured for data protection"
    remediation: "Enable automated backups with appropriate retention period"
`
}
