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

  # ===== GRANULAR RULES FOR PRODUCTION (25 NEW RULES) =====

  # === ENCRYPTION SPECIFICITY ===
  - name: "kms-encryption-mandatory"
    description: "Enforce customer-managed KMS encryption for sensitive resources"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_s3_bucket"
      - "aws_ebs_volume"
      - "aws_rds_cluster"
      - "azurerm_storage_account"
    conditions:
      encryption_type: "customer-managed"
      key_rotation_enabled: true
    message: "Must use customer-managed KMS encryption with automatic key rotation"
    remediation: "Enable CMK encryption and configure automatic rotation every 365 days"

  - name: "tls-minimum-version-1-2"
    description: "Enforce TLS 1.2 or higher for all communications"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_api_gateway"
      - "azurerm_app_service"
    conditions:
      tls_version: "1.2"
    message: "TLS 1.2 or higher must be enforced"
    remediation: "Set minimum_tls_version to 1.2 or higher"

  - name: "database-encryption-in-transit"
    description: "Enable encryption in transit for databases"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_db_instance"
      - "azurerm_sql_server"
    conditions:
      enforce_ssl: true
      storage_encrypted: true
    message: "Database must enforce SSL/TLS and encrypt storage"
    remediation: "Enable storage_encrypted and enforce_ssl=true"

  # === AUTHENTICATION & ACCESS CONTROL ===
  - name: "mfa-enforced-users"
    description: "MFA must be enabled for user accounts"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_iam_user"
      - "azurerm_user"
    conditions:
      mfa_enabled: true
    message: "Multi-factor authentication (MFA) must be enabled for all users"
    remediation: "Enable MFA for user accounts with authenticator app or hardware token"

  - name: "root-account-hardened"
    description: "Root account must have MFA and no access keys"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_account"
    conditions:
      root_mfa_enabled: true
      root_access_key_disabled: true
    message: "Root account must have MFA and no programmatic access keys"
    remediation: "Enable root MFA, delete any access keys, use roles instead"

  - name: "service-principal-credential-rotation"
    description: "Service principals must be rotated every 90 days"
    severity: "warning"
    category: "security"
    enabled: true
    resource_types:
      - "azurerm_service_principal"
    conditions:
      credential_rotation_days: 90
    message: "Service principal credentials must be rotated every 90 days"
    remediation: "Set credential expiration and rotate before expiry"

  - name: "cross-account-access-restricted"
    description: "Cross-account access must be restricted"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_iam_role"
    conditions:
      cross_account_access_restricted: true
    message: "Cross-account access must be explicitly restricted"
    remediation: "Use principal restrictions and require explicit account IDs"

  # === AUDIT & COMPLIANCE ===
  - name: "cloudtrail-multi-region-enabled"
    description: "CloudTrail must be multi-region with log file validation"
    severity: "error"
    category: "compliance"
    enabled: true
    resource_types:
      - "aws_cloudtrail"
    conditions:
      is_multi_region_trail: true
      enable_log_file_validation: true
      is_logging: true
    message: "CloudTrail must be multi-region with log file validation"
    remediation: "Enable multi_region=true and log_file_validation=true"

  - name: "audit-logs-immutable-storage"
    description: "Audit logs must use immutable storage (Object Lock, WORM)"
    severity: "error"
    category: "compliance"
    enabled: true
    resource_types:
      - "aws_s3_bucket"
      - "azurerm_storage_account"
    conditions:
      object_lock_enabled: true
      versioning_enabled: true
    message: "Audit logs must be stored in immutable/write-once storage"
    remediation: "Enable S3 Object Lock (COMPLIANCE mode) and versioning"

  - name: "log-retention-minimum-90-days"
    description: "Logs must be retained for minimum 90 days"
    severity: "warning"
    category: "compliance"
    enabled: true
    resource_types:
      - "aws_cloudwatch_log_group"
      - "azurerm_log_analytics_workspace"
    conditions:
      retention_in_days: 90
    message: "Logs must be retained for minimum 90 days"
    remediation: "Set log_retention_in_days ≥ 90"

  - name: "vpc-flow-logs-enabled"
    description: "VPC flow logs must be enabled for network monitoring"
    severity: "warning"
    category: "compliance"
    enabled: true
    resource_types:
      - "aws_vpc"
      - "azurerm_virtual_network"
    conditions:
      flow_logs_enabled: true
    message: "Flow logs must be enabled for all VPCs/virtual networks"
    remediation: "Enable VPC/Network Flow Logs to CloudWatch Logs"

  # === NETWORK SECURITY ===
  - name: "security-group-restrict-ssh-rdp"
    description: "Security groups must restrict SSH/RDP to 0.0.0.0/0"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_security_group"
      - "azurerm_network_security_group"
    conditions:
      allow_ssh_world_open: false
      allow_rdp_world_open: false
    message: "SSH (22) and RDP (3389) must NOT be open to 0.0.0.0/0"
    remediation: "Remove rules allowing 0.0.0.0/0 access to ports 22 and 3389"

  - name: "network-acl-deny-by-default"
    description: "Network ACLs must explicitly allow, implicitly deny"
    severity: "warning"
    category: "security"
    enabled: true
    resource_types:
      - "aws_network_acl"
    conditions:
      default_action: "deny"
    message: "Network ACLs should follow explicit allow, implicit deny principle"
    remediation: "Set default network ACL rules to DENY, explicitly allow required traffic"

  - name: "nat-gateway-for-private-subnets"
    description: "Private subnets must route through NAT Gateway"
    severity: "warning"
    category: "security"
    enabled: true
    resource_types:
      - "aws_subnet"
    conditions:
      has_nat_gateway: true
    message: "Private subnets must use NAT Gateway for outbound internet"
    remediation: "Configure NAT Gateway in public subnet for private subnets"

  - name: "waf-enabled-on-apis"
    description: "WAF must be enabled on public APIs and load balancers"
    severity: "warning"
    category: "security"
    enabled: true
    resource_types:
      - "aws_api_gateway"
      - "aws_lb"
    conditions:
      waf_enabled: true
    message: "Web Application Firewall must be enabled"
    remediation: "Enable AWS WAF with OWASP Top 10 rule groups"

  # === DATABASE HARDENING ===
  - name: "database-delete-protection"
    description: "Production databases must have deletion protection"
    severity: "error"
    category: "reliability"
    enabled: true
    resource_types:
      - "aws_db_instance"
      - "azurerm_sql_server"
    conditions:
      deletion_protection: true
    message: "Production databases must be protected against deletion"
    remediation: "Enable deletion_protection=true"

  - name: "database-backup-retention-14-days"
    description: "Database backups must be retained 14+ days"
    severity: "warning"
    category: "reliability"
    enabled: true
    resource_types:
      - "aws_db_instance"
      - "azurerm_sql_server"
    conditions:
      backup_retention_period: 14
    message: "Backups must be retained 14+ days"
    remediation: "Set backup_retention_period ≥ 14 days"

  - name: "database-enhanced-monitoring"
    description: "Enhanced monitoring must be enabled for databases"
    severity: "warning"
    category: "compliance"
    enabled: true
    resource_types:
      - "aws_db_instance"
    conditions:
      enhanced_monitoring_enabled: true
      monitoring_interval: 60
    message: "Enhanced monitoring must be enabled (60 second interval)"
    remediation: "Enable enhanced_monitoring with 60-second granularity"

  - name: "database-not-publicly-accessible"
    description: "Production databases must NOT be publicly accessible"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_db_instance"
      - "azurerm_sql_server"
    conditions:
      publicly_accessible: false
    message: "Production databases must not be publicly accessible"
    remediation: "Set publicly_accessible=false, restrict to VPC/private networks"

  # === TAGGING & GOVERNANCE ===
  - name: "comprehensive-resource-tagging"
    description: "All resources must have comprehensive governance tags"
    severity: "warning"
    category: "governance"
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
        - "CostCenter"
        - "DataClassification"
    message: "All resources must have comprehensive tags"
    remediation: "Add tags: Environment, Owner, Project, CostCenter, DataClassification"

  - name: "temporary-resource-expiration"
    description: "Temporary resources must have expiration dates"
    severity: "info"
    category: "governance"
    enabled: true
    resource_types:
      - "aws_instance"
      - "aws_security_group"
    conditions:
      tags.required: ["ExpirationDate"]
    message: "Temporary resources should have ExpirationDate tag"
    remediation: "Add ExpirationDate in YYYY-MM-DD format for temporary resources"

  # === COST OPTIMIZATION ===
  - name: "use-newest-instance-types"
    description: "Use latest generation instance types (gp3, t4g, m7g)"
    severity: "info"
    category: "cost"
    enabled: true
    resource_types:
      - "aws_instance"
      - "aws_db_instance"
    conditions:
      instance_generation: "latest"
    message: "Use latest generation instances for better cost/performance"
    remediation: "Migrate to gp3 volumes, t4g/m7g instances (20-30% cost savings)"

  - name: "compute-auto-scaling"
    description: "Production compute must have auto-scaling configured"
    severity: "warning"
    category: "reliability"
    enabled: true
    resource_types:
      - "aws_autoscaling_group"
      - "azurerm_virtual_machine_scale_set"
    conditions:
      auto_scaling_enabled: true
      min_size: ">=1"
      max_size: ">=2"
    message: "Production compute must have auto-scaling (min≥1, max≥2)"
    remediation: "Configure auto-scaling with min_size≥1, max_size≥2"

  - name: "cross-region-replication"
    description: "Critical data must be replicated to 2+ regions"
    severity: "warning"
    category: "reliability"
    enabled: true
    resource_types:
      - "aws_s3_bucket"
      - "aws_rds_cluster"
    conditions:
      replication_enabled: true
      replica_regions_minimum: 2
    message: "Critical data must be replicated to ≥2 regions"
    remediation: "Enable cross-region replication for disaster recovery"
`
}
