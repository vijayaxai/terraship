# Azure Test Scenarios for Terraship

This directory contains comprehensive test scenarios for validating Azure infrastructure with Terraship.

## Files

### `test-all-scenarios.tf`
Comprehensive test file covering all Terraship policy rules for Azure, including:

- ✅ **Compliant Resources** - Examples of properly configured resources
- ❌ **Non-Compliant Resources** - Examples demonstrating policy violations

## Scenarios Covered

### 1. Tagging Compliance
- Required tags (Environment, Owner, Project)
- Cost center tagging
- Naming conventions

### 2. Storage Security
- Encryption at rest
- HTTPS enforcement
- Public access control
- Network security

### 3. Database Security
- SQL Server/Database encryption (TDE)
- Backup configuration
- Public access blocking
- Network isolation

### 4. Compute Security
- VM disk encryption
- Private subnet deployment
- Network interface configuration

### 5. Key Management
- Key Vault security
- Soft delete and purge protection
- Network isolation
- RBAC authorization

### 6. IAM & Access Control
- Least privilege principle
- Role assignments
- Permission scoping

### 7. Application Services
- App Service HTTPS enforcement
- TLS version requirements
- Logging configuration

## How to Use

### Prerequisites

1. **Install Terraform**
   ```powershell
   choco install terraform -y
   ```

2. **Install Azure CLI** (optional, for authentication)
   ```powershell
   choco install azure-cli -y
   ```

3. **Build Terraship**
   ```powershell
   cd ../../
   go build -o bin\terraship.exe .\cmd\terraship
   ```

### Testing Scenarios

#### Option 1: Dry Run Validation (No Azure Credentials Needed)

This validates the Terraform syntax and policy rules without connecting to Azure:

```powershell
# Navigate to Azure examples
cd examples\azure

# Validate with Terraship (dry run)
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --provider azure
```

**Expected Output:**
```
═══════════════════════════════════════════════════════════════
                TERRASHIP VALIDATION REPORT                  
═══════════════════════════════════════════════════════════════

SUMMARY:
  Total Resources:    ~30
  ✓ Passed:           ~15 (compliant resources)
  ✗ Failed:           ~15 (non-compliant resources)
  ⚠ Warnings:         Multiple
  ⨯ Errors:           Multiple

VIOLATIONS FOUND:
  - azurerm_resource_group.missing_tags_rg: Missing required tags
  - azurerm_storage_account.insecure_storage: HTTPS not enforced
  - azurerm_managed_disk.unencrypted_disk: Encryption not enabled
  - azurerm_mssql_database.insecure_db: TDE not enabled
  ... (more violations)
```

#### Option 2: Validate Existing Infrastructure (Requires Azure Credentials)

This connects to your Azure subscription and validates actual deployed resources:

```powershell
# Login to Azure
az login

# Set subscription
az account set --subscription "your-subscription-id"

# Validate existing resources
..\..\bin\terraship.exe validate . --mode validate-existing --policy ..\..\policies\sample-policy.yml --provider azure
```

#### Option 3: Ephemeral Sandbox Testing (Creates & Destroys Resources)

⚠️ **Warning**: This will create actual Azure resources and may incur costs!

```powershell
# Set Azure credentials
$env:ARM_SUBSCRIPTION_ID = "your-subscription-id"
$env:ARM_TENANT_ID = "your-tenant-id"
$env:ARM_CLIENT_ID = "your-client-id"
$env:ARM_CLIENT_SECRET = "your-client-secret"

# Run ephemeral test (creates, validates, destroys)
..\..\bin\terraship.exe validate . --mode ephemeral-sandbox --policy ..\..\policies\sample-policy.yml --provider azure
```

### Output Formats

#### Human-Readable (Default)
```powershell
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml
```

#### JSON (For Scripts/CI)
```powershell
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --output json > report.json
```

#### SARIF (For Security Tools)
```powershell
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --output sarif > results.sarif
```

## Test Results by Scenario

### ✅ Expected to PASS (Compliant Resources)

| Resource | Policy Rules Satisfied |
|----------|----------------------|
| `compliant_rg` | All required tags, naming convention, cost center |
| `compliant_storage` | Encryption, HTTPS only, private access |
| `compliant_container` | Private access, no public exposure |
| `compliant_disk` | Encryption enabled |
| `compliant_sql_server` | Private access, secure configuration |
| `compliant_db` | TDE enabled, backup configured |
| `compliant_vm` | Private subnet, secure configuration |
| `compliant_kv` | Private access, purge protection |
| `compliant_app` | HTTPS enforced, TLS 1.2 |

### ❌ Expected to FAIL (Non-Compliant Resources)

| Resource | Violations |
|----------|-----------|
| `missing_tags_rg` | Missing required tags: Environment, Owner, Project |
| `missing_costcenter_rg` | Missing CostCenter tag |
| `BAD_NAMING_RG` | Invalid naming convention (uppercase) |
| `insecure_storage` | HTTPS not enforced, public access allowed |
| `public_container` | Public blob access enabled |
| `unencrypted_disk` | No encryption settings |
| `insecure_sql_server` | Public network access enabled |
| `insecure_db` | TDE disabled, no backup configuration |
| `insecure_kv` | Public access, no purge protection |
| `overpermissive_owner` | Owner role too broad (least privilege violation) |
| `insecure_app` | HTTP allowed, old TLS version |

## Understanding the Results

### Severity Levels

- **🔴 Error**: Must fix - blocks deployment in production
- **🟡 Warning**: Should fix - best practice violation
- **🔵 Info**: Nice to have - informational only

### Common Violations Explained

#### Missing Required Tags
```
❌ azurerm_resource_group.missing_tags_rg
   Rule: required-tags
   Severity: error
   Message: Resources must have Environment, Owner, and Project tags
   Fix: Add the following tags:
     tags = {
       Environment = "Production"
       Owner       = "Your Team"
       Project     = "Your Project"
     }
```

#### HTTPS Not Enforced
```
❌ azurerm_storage_account.insecure_storage
   Rule: azure-storage-https-only
   Severity: error
   Message: Storage account should enforce HTTPS only
   Fix: Set enable_https_traffic_only = true
```

#### Encryption Disabled
```
❌ azurerm_managed_disk.unencrypted_disk
   Rule: encryption-at-rest
   Severity: error
   Message: Encryption at rest must be enabled
   Fix: Add encryption_settings { enabled = true }
```

## Integration with CI/CD

### GitHub Actions Example

Create `.github/workflows/azure-validation.yml`:

```yaml
name: Validate Azure Infrastructure

on:
  pull_request:
    paths:
      - 'examples/azure/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v2
      
      - name: Validate with Terraship
        run: |
          cd examples/azure
          ../../bin/terraship validate . \
            --policy ../../policies/sample-policy.yml \
            --provider azure \
            --output sarif > results.sarif
      
      - name: Upload Results
        uses: github/codeql-action/upload-sarif@v2
        if: always()
        with:
          sarif_file: examples/azure/results.sarif
```

## Troubleshooting

### Error: "Terraform not found"
```powershell
# Install Terraform
choco install terraform -y

# Verify installation
terraform version
```

### Error: "Azure provider authentication"
```powershell
# Login to Azure
az login

# List subscriptions
az account list --output table

# Set active subscription
az account set --subscription "your-subscription-id"
```

### Error: "Policy file not found"
```powershell
# Check if running from correct directory
pwd

# Should be in: .../terraship/examples/azure
# Policy should be at: ../../policies/sample-policy.yml

# Verify file exists
Test-Path "..\..\policies\sample-policy.yml"
```

### Validation Takes Too Long
```powershell
# Use dry-run mode (doesn't connect to Azure)
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml

# Or skip integration tests
go test -short ./...
```

## Next Steps

1. **Customize Policy**: Edit `../../policies/sample-policy.yml` for your needs
2. **Add More Scenarios**: Create additional test cases in new `.tf` files
3. **Integrate CI/CD**: Add Terraship to your pipeline
4. **Production Use**: Apply learnings to real infrastructure

## Resources

- [Terraship Documentation](../../README.md)
- [Policy Writing Guide](../../docs/policy-guide.md)
- [Azure Provider Documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs)
- [Testing Guide](../../TESTING.md)

## Support

- 📖 **Docs**: See main README.md
- 🐛 **Issues**: GitHub Issues
- 💬 **Questions**: GitHub Discussions
