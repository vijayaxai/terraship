# 🚀 Terraship End User Guide

## How End Users Interact with Terraship

Terraship provides **4 main ways** for users to validate their Terraform infrastructure:

---

## 1. 💻 Command-Line Interface (CLI)

### Most Common Usage: Developers and DevOps Engineers

#### Installation (One-Time)
```powershell
# Download the binary
# Windows: Download terraship.exe
# Or build from source: go install github.com/vijayaxai/terraship/cmd/terraship@latest
```

#### Daily Workflow

**Scenario 1: Quick Validation Before Commit**
```powershell
# Navigate to your Terraform project
cd C:\projects\my-terraform-infra

# Validate against company policy
terraship validate . --policy C:\policies\company-security.yml

# Output shows immediate feedback:
# ✓ 45 resources passed
# ✗ 3 resources failed
# - aws_s3_bucket.logs: Missing required tag 'Environment'
# - aws_db_instance.main: Encryption not enabled
```

**Scenario 2: Check Existing Infrastructure for Drift**
```powershell
# Validate what's currently in AWS/Azure/GCP
terraship validate . --mode validate-existing --policy ./policy.yml

# This connects to your cloud and checks:
# - Does actual infrastructure match Terraform state?
# - Are there compliance violations?
# - Has someone made manual changes?
```

**Scenario 3: Test in Temporary Sandbox**
```powershell
# Create temporary resources, test them, then destroy
terraship validate ./terraform --mode ephemeral-sandbox

# Terraship will:
# 1. terraform init
# 2. terraform apply (creates resources)
# 3. Validate against policy
# 4. Check cloud for drift
# 5. terraform destroy (cleanup)
```

**Scenario 4: Different Output Formats**
```powershell
# Human-readable (default)
terraship validate . --policy policy.yml

# JSON for parsing by scripts
terraship validate . --policy policy.yml --output json > report.json

# SARIF for security tools (GitHub Advanced Security, etc.)
terraship validate . --policy policy.yml --output sarif > results.sarif
```

---

## 2. 🔌 VS Code Extension

### Most Common Usage: Developers Writing Terraform

#### Installation (One-Time)
1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "Terraship"
4. Click Install

#### Daily Workflow

**Scenario 1: Real-Time Validation While Coding**
```terraform
# You're editing main.tf in VS Code:

resource "aws_s3_bucket" "data" {
  bucket = "my-data-bucket"
  # As you type, Terraship shows warnings in the editor:
  # ⚠️ Missing required tag: Environment
  # ⚠️ Encryption not enabled
}
```

**Scenario 2: Validate Current File**
1. Right-click in Terraform file
2. Select "Terraship: Validate Current File"
3. Results appear in VS Code panel:
   - ✓ Compliant resources
   - ✗ Violations with line numbers (clickable)
   - 💡 Quick fix suggestions

**Scenario 3: Validate Entire Workspace**
1. Open Command Palette (Ctrl+Shift+P)
2. Type: "Terraship: Validate Workspace"
3. View comprehensive report in webview
4. Click violations to jump to code

**Scenario 4: Configure Policy**
```json
// .vscode/settings.json
{
  "terraship.policyPath": "./policies/team-policy.yml",
  "terraship.validateOnSave": true,
  "terraship.cloudProvider": "aws",
  "terraship.outputFormat": "human"
}
```

---

## 3. 🔄 GitHub Actions (CI/CD Pipeline)

### Most Common Usage: Automated Checks on Pull Requests

#### Setup (One-Time)
Create `.github/workflows/terraform-validation.yml`:

```yaml
name: Validate Terraform

on:
  pull_request:
    paths:
      - 'terraform/**'
      - 'policies/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Validate Infrastructure
        uses: terraship/terraship-action@v1
        with:
          terraform-directory: ./terraform
          policy-path: ./policies/security.yml
          mode: validate-existing
          output-format: sarif
          
      - name: Upload Results
        uses: github/codeql-action/upload-sarif@v2
        if: always()
        with:
          sarif_file: terraship-report.txt
```

#### Daily Workflow

**When Developer Creates Pull Request:**

1. Developer pushes Terraform changes
2. GitHub Action triggers automatically
3. Terraship validates changes
4. Results appear in PR:
   - ✅ All checks passed → PR can be merged
   - ❌ Policy violations → PR blocked, shows errors
   - Security findings appear in GitHub Security tab

**Example PR Comment:**
```
🚢 Terraship Validation Results

Summary:
✓ 23 resources passed
✗ 2 resources failed

Failures:
1. terraform/s3.tf:15 - aws_s3_bucket.logs
   ❌ Missing required tag: CostCenter
   💡 Fix: Add tag { CostCenter = "engineering" }

2. terraform/rds.tf:42 - aws_db_instance.prod
   ❌ Encryption at rest not enabled
   💡 Fix: Add storage_encrypted = true
```

---

## 4. 🧪 Terratest Integration (Automated Testing)

### Most Common Usage: QA Engineers and Test Automation

#### Setup (One-Time)
Install Terraship Go package:
```bash
go get github.com/vijayaxai/terraship/pkg/terraship
```

#### Daily Workflow

**Scenario 1: Unit Test for Infrastructure**
```go
// test/infrastructure_test.go
package test

import (
    "testing"
    "github.com/vijayaxai/terraship/pkg/terraship"
)

func TestProductionInfrastructure(t *testing.T) {
    t.Parallel()

    opts := terraship.Options{
        TerraformDir:  "../terraform/prod",
        PolicyPath:    "../policies/prod-policy.yml",
        Mode:          "validate-existing",
        CloudProvider: "aws",
    }

    // Validate infrastructure
    result := terraship.Validate(t, opts)
    
    // Assert all checks passed
    result.AssertCompliant(t)
    
    // Assert no drift
    result.AssertNoDrift(t)
    
    // Check specific resource
    result.AssertResourceCompliant(t, "aws_s3_bucket.data")
}
```

**Scenario 2: Integration Test with Ephemeral Resources**
```go
func TestNewFeature_EphemeralEnvironment(t *testing.T) {
    opts := terraship.Options{
        TerraformDir:  "../terraform/feature-x",
        PolicyPath:    "../policies/security.yml",
        Mode:          "ephemeral-sandbox",  // Creates and destroys
        CloudProvider: "aws",
    }

    result := terraship.Validate(t, opts)
    
    // Terraship will:
    // 1. Create resources in AWS
    // 2. Validate them
    // 3. Destroy them automatically
    // 4. Return results
    
    if result.HasViolations() {
        t.Errorf("Found %d violations", result.ViolationCount())
    }
}
```

**Run Tests:**
```bash
# Run all tests
go test ./test/... -v

# Run specific test
go test ./test/... -run TestProductionInfrastructure -v

# Generate coverage
go test ./test/... -cover
```

---

## 📊 Real-World User Scenarios

### Scenario A: Solo Developer

**Tools Used:** CLI + VS Code Extension

**Workflow:**
1. Write Terraform in VS Code → See violations in real-time
2. Fix issues as you code
3. Before commit: `terraship validate . --policy ./policy.yml`
4. Commit only when clean

### Scenario B: Small Team (3-10 developers)

**Tools Used:** CLI + VS Code Extension + GitHub Actions

**Workflow:**
1. Each developer has VS Code extension
2. GitHub Action validates all PRs
3. Team lead reviews policy violations in PR
4. Merge only when validation passes

### Scenario C: Enterprise Organization

**Tools Used:** All 4 tools

**Workflow:**
1. **Developers:** VS Code extension for immediate feedback
2. **CI/CD:** GitHub Actions validate every PR
3. **QA:** Terratest suite runs nightly
4. **Ops:** CLI for ad-hoc production checks
5. **Security Team:** Reviews SARIF reports in GitHub Security

---

## 🎯 End User Benefits

### For Developers
- ✅ Catch issues before code review
- ✅ Learn best practices through validation messages
- ✅ No context switching (validate in VS Code)

### For DevOps Engineers
- ✅ Enforce standards across teams
- ✅ Detect drift in production
- ✅ Automate compliance checks

### For Security Teams
- ✅ Prevent security misconfigurations
- ✅ Track violations over time
- ✅ Enforce encryption, tagging, etc.

### For Management
- ✅ Ensure compliance (SOC2, HIPAA, PCI-DSS)
- ✅ Reduce manual reviews
- ✅ Audit trail of all validations

---

## 📝 Policy Management

### Creating Your First Policy

**Step 1: Initialize**
```bash
terraship init
# Creates: terraship-policy.yml with examples
```

**Step 2: Customize Rules**
```yaml
version: "1.0"
name: "My Team Policy"

rules:
  - name: "require-encryption"
    severity: "error"
    resource_types: ["aws_s3_bucket"]
    conditions:
      encryption.enabled: true
    message: "All S3 buckets must have encryption"
```

**Step 3: Test Policy**
```bash
terraship validate ./terraform --policy ./my-policy.yml
```

**Step 4: Share with Team**
- Commit policy to Git
- Reference in CI/CD
- Configure in VS Code settings

---

## 🆘 Common End User Questions

### Q: Do I need cloud credentials?
**A:** 
- For `validate` mode (dry run): NO
- For `validate-existing` mode: YES (read-only)
- For `ephemeral-sandbox` mode: YES (create/destroy)

### Q: Can I run this locally before CI/CD?
**A:** Yes! That's the recommended workflow:
```bash
# Local validation before commit
terraship validate . --policy ./policy.yml
git commit -m "Added S3 bucket with encryption"
```

### Q: What if validation fails?
**A:** Terraship shows:
1. What failed (resource + line number)
2. Why it failed (policy rule violated)
3. How to fix (remediation guidance)

Example:
```
✗ aws_s3_bucket.logs (terraform/s3.tf:15)
  Rule: encryption-at-rest
  Severity: error
  Message: Encryption at rest must be enabled
  Fix: Add server_side_encryption_configuration block
```

### Q: Can I disable specific rules?
**A:** Yes, three ways:

1. **In policy file:**
```yaml
rules:
  - name: "my-rule"
    enabled: false
```

2. **Via CLI flag:**
```bash
terraship validate . --policy policy.yml --disable-rule my-rule
```

3. **In code comment:**
```terraform
# terraship:ignore encryption-at-rest
resource "aws_s3_bucket" "temp" {
  # This bucket is exempt from encryption
}
```

### Q: How long does validation take?
**A:**
- **Dry run (validate):** Seconds (just parses Terraform)
- **Existing infrastructure:** 1-2 minutes (queries cloud APIs)
- **Ephemeral sandbox:** 5-15 minutes (creates + validates + destroys)

### Q: Can I validate without Terraform installed?
**A:** No, Terraship requires:
- Terraform CLI installed
- Valid Terraform configuration
- (Optionally) Cloud credentials for existing resource checks

---

## 🎓 Learning Path for New Users

### Week 1: Basic Usage
1. Install CLI
2. Run `terraship validate` on existing project
3. Fix one violation
4. Install VS Code extension

### Week 2: Policy Creation
1. Run `terraship init`
2. Customize policy for your needs
3. Test against your infrastructure
4. Share with team

### Week 3: CI/CD Integration
1. Add GitHub Action
2. Configure policy enforcement
3. Train team on PR workflow

### Week 4: Advanced Features
1. Try `ephemeral-sandbox` mode
2. Integrate with Terratest
3. Customize output formats
4. Create custom rules

---

## 📞 Getting Help

- 📖 **Documentation:** [docs/](./docs/)
- 💬 **Community:** GitHub Discussions
- 🐛 **Issues:** GitHub Issues
- 📧 **Support:** support@terraship.io

---

**Ready to start?**

```bash
# Download Terraship
# Try it on your project
terraship validate . --policy ./policy.yml

# That's it! 🎉
```
