# Expected Terraship Validation Results
# Generated for: test-all-scenarios.tf

## Summary

When running:
```
terraship validate . --policy ../../policies/sample-policy.yml --provider azure
```

## Expected Output

```
═══════════════════════════════════════════════════════════════════════════
                        TERRASHIP VALIDATION REPORT                          
═══════════════════════════════════════════════════════════════════════════

Validating Azure infrastructure against policy...
Provider: Azure
Policy: ../../policies/sample-policy.yml
Mode: validate-existing

─────────────────────────────────────────────────────────────────────────────
SUMMARY
─────────────────────────────────────────────────────────────────────────────

Total Resources:    26
✓ Passed:           9  (34.6%)
✗ Failed:           17 (65.4%)
⚠ Warnings:         12
⨯ Errors:           14
↔ Drift Detected:   0

─────────────────────────────────────────────────────────────────────────────
VIOLATIONS
─────────────────────────────────────────────────────────────────────────────

❌ ERRORS (Must Fix - Blocks Production)
─────────────────────────────────────────────────────────────────────────────

1. azurerm_resource_group.missing_tags_rg
   File: test-all-scenarios.tf:40
   Rule: required-tags
   Severity: error
   Message: Resources must have Environment, Owner, and Project tags
   Current Tags: { Description }
   Missing Tags: Environment, Owner, Project
   Remediation: Add required tags:
     tags = {
       Environment = "Production"
       Owner       = "Your Team Name"
       Project     = "Your Project Name"
     }

2. azurerm_resource_group.missing_costcenter_rg
   File: test-all-scenarios.tf:49
   Rule: cost-tagging
   Severity: warning
   Message: Resources should have a CostCenter tag for billing
   Current Tags: { Environment, Owner, Project }
   Missing Tags: CostCenter
   Remediation: Add CostCenter tag:
     tags = {
       ...
       CostCenter = "Engineering"
     }

3. azurerm_resource_group.BAD_NAMING_RG
   File: test-all-scenarios.tf:62
   Rule: naming-convention
   Severity: info
   Message: Resource names should follow the naming convention
   Current Name: TerraShip-BAD-NAMING-RG
   Expected Pattern: ^[a-z0-9-]+$ (lowercase, numbers, hyphens only)
   Remediation: Rename to: terraship-bad-naming-rg

4. azurerm_storage_account.insecure_storage
   File: test-all-scenarios.tf:117
   Rule: azure-storage-https-only
   Severity: error
   Message: Storage account should enforce HTTPS only
   Current: enable_https_traffic_only = false
   Required: enable_https_traffic_only = true
   Remediation: Set enable_https_traffic_only = true

5. azurerm_storage_account.insecure_storage
   File: test-all-scenarios.tf:117
   Rule: block-public-access
   Severity: error
   Message: Public access should be blocked
   Current: allow_nested_items_to_be_public = true
   Required: allow_nested_items_to_be_public = false
   Remediation: Block public access:
     allow_nested_items_to_be_public = false
     public_network_access_enabled = false

6. azurerm_storage_container.public_container
   File: test-all-scenarios.tf:153
   Rule: block-public-access
   Severity: error
   Message: Public access should be blocked
   Current: container_access_type = "blob"
   Required: container_access_type = "private"
   Remediation: Change to: container_access_type = "private"

7. azurerm_managed_disk.unencrypted_disk
   File: test-all-scenarios.tf:189
   Rule: encryption-at-rest
   Severity: error
   Message: Encryption at rest must be enabled
   Current: No encryption_settings block
   Remediation: Add encryption:
     encryption_settings {
       enabled = true
     }

8. azurerm_mssql_server.insecure_sql_server
   File: test-all-scenarios.tf:248
   Rule: block-public-access
   Severity: error
   Message: Public access should be blocked
   Current: public_network_access_enabled = true
   Required: public_network_access_enabled = false
   Remediation: Disable public access:
     public_network_access_enabled = false

9. azurerm_mssql_database.insecure_db
   File: test-all-scenarios.tf:269
   Rule: encryption-at-rest
   Severity: error
   Message: Encryption at rest must be enabled
   Current: transparent_data_encryption_enabled = false
   Required: transparent_data_encryption_enabled = true
   Remediation: Enable TDE:
     transparent_data_encryption_enabled = true

10. azurerm_mssql_database.insecure_db
    File: test-all-scenarios.tf:269
    Rule: backup-enabled
    Severity: warning
    Message: Backup should be configured for data protection
    Current: No retention policies configured
    Remediation: Add backup configuration:
      short_term_retention_policy {
        retention_days = 7
      }

11. azurerm_key_vault.insecure_kv
    File: test-all-scenarios.tf:406
    Rule: block-public-access
    Severity: error
    Message: Public access should be blocked
    Current: public_network_access_enabled = true
    Required: public_network_access_enabled = false
    Remediation: Disable public access:
      public_network_access_enabled = false

12. azurerm_role_assignment.overpermissive_owner
    File: test-all-scenarios.tf:431
    Rule: iam-least-privilege
    Severity: error
    Message: IAM policies should not use overly broad permissions
    Current: role_definition_name = "Owner"
    Issue: Owner role grants full control (violates least privilege)
    Remediation: Use specific roles instead:
      - Reader (read-only)
      - Contributor (read/write without role assignments)
      - Custom roles with specific permissions

13. azurerm_linux_web_app.insecure_app
    File: test-all-scenarios.tf:491
    Rule: azure-storage-https-only
    Severity: error
    Message: Storage account should enforce HTTPS only
    Current: https_only = false
    Required: https_only = true
    Remediation: Enforce HTTPS:
      https_only = true
      site_config {
        minimum_tls_version = "1.2"
      }

14. Multiple resources missing required tags
    Affected: insecure_storage, unencrypted_disk, insecure_db, insecure_kv, insecure_app
    Rule: required-tags
    Severity: error
    Message: All resources must have Environment, Owner, Project tags

─────────────────────────────────────────────────────────────────────────────
✓ PASSED RESOURCES
─────────────────────────────────────────────────────────────────────────────

✅ azurerm_resource_group.compliant_rg
   All required tags present ✓
   Proper naming convention ✓
   Cost center tagged ✓

✅ azurerm_storage_account.compliant_storage
   HTTPS enforcement enabled ✓
   Public access blocked ✓
   Encryption enabled ✓
   Proper tagging ✓

✅ azurerm_storage_container.compliant_container
   Private access configured ✓

✅ azurerm_managed_disk.compliant_disk
   Encryption enabled ✓
   All required tags present ✓

✅ azurerm_mssql_server.compliant_sql_server
   Public access disabled ✓
   Secure configuration ✓
   Proper tagging ✓

✅ azurerm_mssql_database.compliant_db
   TDE (Transparent Data Encryption) enabled ✓
   Backup policies configured ✓
   Short-term retention: 7 days ✓
   Long-term retention: Weekly, Monthly, Yearly ✓
   Proper tagging ✓

✅ azurerm_linux_virtual_machine.compliant_vm
   Private subnet deployment ✓
   Secure configuration ✓
   Proper tagging ✓

✅ azurerm_key_vault.compliant_kv
   Public access disabled ✓
   Purge protection enabled ✓
   Soft delete enabled ✓
   RBAC authorization enabled ✓
   Proper tagging ✓

✅ azurerm_linux_web_app.compliant_app
   HTTPS enforcement enabled ✓
   TLS 1.2 minimum ✓
   Proper tagging ✓

─────────────────────────────────────────────────────────────────────────────
POLICY RULES TESTED
─────────────────────────────────────────────────────────────────────────────

✓ required-tags              (14 checks: 9 passed, 5 failed)
✓ cost-tagging               (14 checks: 10 passed, 4 failed)
✓ naming-convention          (26 checks: 25 passed, 1 failed)
✓ encryption-at-rest         (8 checks: 5 passed, 3 failed)
✓ block-public-access        (10 checks: 5 passed, 5 failed)
✓ azure-storage-https-only   (3 checks: 2 passed, 1 failed)
✓ enable-logging             (4 checks: 4 passed)
✓ backup-enabled             (2 checks: 1 passed, 1 failed)
✓ use-private-subnet         (1 check: 1 passed)
✓ iam-least-privilege        (2 checks: 1 passed, 1 failed)

─────────────────────────────────────────────────────────────────────────────
RECOMMENDATIONS
─────────────────────────────────────────────────────────────────────────────

High Priority (Fix Immediately):
  1. Enable encryption on all storage resources
  2. Block public access to databases and storage
  3. Enforce HTTPS on all web-facing services
  4. Add required tags to all resources

Medium Priority (Fix Soon):
  5. Configure backup policies for databases
  6. Review and reduce IAM permissions (least privilege)
  7. Add cost center tags for billing tracking

Low Priority (Nice to Have):
  8. Standardize resource naming conventions
  9. Enable advanced logging and monitoring

─────────────────────────────────────────────────────────────────────────────
COMPLIANCE SCORE
─────────────────────────────────────────────────────────────────────────────

Overall Score:     34.6% 🔴 FAIL
Security Score:    40.0% 🔴 FAIL
Compliance Score:  30.0% 🔴 FAIL
Cost Score:        71.4% 🟡 NEEDS IMPROVEMENT

Threshold: 80% required to pass
Status: ❌ VALIDATION FAILED - 17 resources require fixes

─────────────────────────────────────────────────────────────────────────────
NEXT STEPS
─────────────────────────────────────────────────────────────────────────────

1. Review all ERROR-level violations above
2. Fix issues in test-all-scenarios.tf
3. Re-run validation: terraship validate . --policy ../../policies/sample-policy.yml
4. Repeat until all resources pass or only INFO-level warnings remain

For detailed remediation steps, see the specific violation messages above.

═══════════════════════════════════════════════════════════════════════════
                            END OF REPORT                                    
═══════════════════════════════════════════════════════════════════════════
```

## How to Actually Test (Once Provider Issue is Resolved)

### Option 1: Test on Non-Corporate Network
```powershell
# Use home internet or mobile hotspot
# Then run:
terraform init
terraship validate . --policy ../../policies/sample-policy.yml
```

### Option 2: Use Terraform Plugin Cache from Another Machine
```powershell
# On machine with internet:
terraform init

# Copy from:
%APPDATA%\terraform.d\plugins\

# To same location on corporate machine
# Then run validation
```

### Option 3: Skip Terraform Init (Syntax Validation Only)
```powershell
# This won't actually validate the resources but will check syntax
terraform validate
```

## Metrics Summary

| Metric | Value |
|--------|-------|
| **Total Resources** | 26 |
| **Compliant Resources** | 9 (34.6%) |
| **Non-Compliant Resources** | 17 (65.4%) |
| **Policy Rules Tested** | 10 |
| **Errors** | 14 |
| **Warnings** | 12 |
| **Pass Threshold** | 80% |
| **Result** | ❌ FAIL |

## Key Takeaways

This test file demonstrates:
1. ✅ What compliant resources look like
2. ❌ Common security misconfigurations
3. 📊 How Terraship validates each rule
4. 💡 Clear remediation guidance for each violation

The intentional mix of compliant and non-compliant resources shows how Terraship catches policy violations before they reach production!
