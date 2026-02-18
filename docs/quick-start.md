# Terraship Quick Start Guide

This guide will help you get started with Terraship in 5 minutes.

## Prerequisites

- Terraform 1.6+ installed
- Go 1.22+ (for building from source)
- Cloud provider credentials configured
- Basic understanding of Terraform

## Step 1: Install Terraship

### Option A: Download Pre-built Binary

```bash
# Linux/macOS
curl -L https://github.com/vijayaxai/terraship/releases/latest/download/terraship-$(uname -s)-$(uname -m) -o /usr/local/bin/terraship
chmod +x /usr/local/bin/terraship

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/vijayaxai/terraship/releases/latest/download/terraship-windows-amd64.exe" -OutFile "terraship.exe"
```

### Option B: Install with Go

```bash
go install github.com/vijayaxai/terraship/cmd/terraship@latest
```

### Option C: Build from Source

```bash
git clone https://github.com/vijayaxai/terraship.git
cd terraship
make build
sudo cp bin/terraship /usr/local/bin/
```

Verify installation:

```bash
terraship --version
```

## Step 2: Set Up Cloud Credentials

### AWS

```bash
export AWS_REGION=us-east-1
export AWS_PROFILE=default
# Or use AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
```

### Azure

```bash
export AZURE_SUBSCRIPTION_ID=your-subscription-id
export AZURE_TENANT_ID=your-tenant-id
# Or use Azure CLI: az login
```

### GCP

```bash
export GCP_PROJECT=your-project-id
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
# Or use: gcloud auth application-default login
```

## Step 3: Prepare Your Terraform Code

Navigate to your Terraform directory:

```bash
cd /path/to/terraform
```

Your directory should contain `.tf` files:

```
terraform/
├── main.tf
├── variables.tf
└── outputs.tf
```

## Step 4: Create a Policy File

Create a basic policy file `policy.yml`:

```yaml
version: "1.0"
name: "Basic Security Policy"
description: "Essential security checks"

rules:
  - name: "required-tags"
    description: "Ensure resources have required tags"
    severity: "warning"
    enabled: true
    resource_types:
      - "*"
    conditions:
      tags.required:
        - "Environment"
        - "Owner"
    message: "Resources should have Environment and Owner tags"
    remediation: "Add tags to your resources"

  - name: "encryption-enabled"
    description: "Encryption must be enabled"
    severity: "error"
    enabled: true
    resource_types:
      - "aws_s3_bucket"
      - "azurerm_storage_account"
      - "google_storage_bucket"
    conditions:
      encryption.enabled: true
    message: "Storage resources must have encryption enabled"
    remediation: "Enable encryption in your resource configuration"
```

Or use the sample policy:

```bash
mkdir -p policies
curl -o policies/sample-policy.yml https://raw.githubusercontent.com/terraship/terraship/main/policies/sample-policy.yml
```

## Step 5: Run Validation

### Validate Existing Infrastructure

```bash
terraship validate . --policy policy.yml --mode validate-existing
```

This will:
1. Initialize Terraform
2. Generate a plan
3. Validate resources against policies
4. Check for drift (if resources exist)
5. Display results

### Create Ephemeral Sandbox

```bash
terraship validate . --policy policy.yml --mode ephemeral-sandbox
```

This will:
1. Initialize Terraform
2. Apply the configuration
3. Validate the created resources
4. Destroy the resources (unless `--no-destroy` is used)

### Customize Output

```bash
# JSON output
terraship validate . --policy policy.yml --output json

# SARIF output (for CI/CD)
terraship validate . --policy policy.yml --output sarif --output-file report.sarif

# Verbose mode
terraship validate . --policy policy.yml --verbose
```

## Step 6: Interpret Results

### Successful Validation

```
═══════════════════════════════════════════════════════════════
                    TERRASHIP VALIDATION REPORT                  
═══════════════════════════════════════════════════════════════

SUMMARY:
  Total Resources:    5
  ✓ Passed:           5
  ✗ Failed:           0
  ⚠ Warnings:         0

✓ VALIDATION PASSED
```

### Failed Validation

```
SUMMARY:
  Total Resources:    5
  ✓ Passed:           3
  ✗ Failed:           2
  ⚠ Warnings:         1

✗ VALIDATION FAILED

DETAILED RESULTS:
───────────────────────────────────────────────────────────────

✗ aws_s3_bucket.example (aws_s3_bucket)
  Provider: aws
  Policy Checks:
    ✗ encryption-at-rest [error]
      Message: Encryption at rest must be enabled
      - Encryption is not enabled
      💡 Remediation: Enable server-side encryption for your resource
```

## Step 7: Fix Issues

Based on the validation results, update your Terraform code:

```hcl
# Before
resource "aws_s3_bucket" "example" {
  bucket = "my-bucket"
}

# After
resource "aws_s3_bucket" "example" {
  bucket = "my-bucket"
  
  tags = {
    Environment = "production"
    Owner       = "devops@example.com"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "example" {
  bucket = aws_s3_bucket.example.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}
```

## Step 8: Re-validate

Run validation again:

```bash
terraship validate . --policy policy.yml
```

## Using with Terratest

Create a test file `infrastructure_test.go`:

```go
package test

import (
    "testing"
    "github.com/vijayaxai/terraship/pkg/terraship"
)

func TestInfraCompliance(t *testing.T) {
    result := terraship.ValidateExisting(t, ".", "./policy.yml")
    result.AssertCompliant(t)
    result.AssertNoDrift(t)
}
```

Run the test:

```bash
go test -v
```

## Using in CI/CD

### GitHub Actions

Create `.github/workflows/terraform.yml`:

```yaml
name: Terraform Validation

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Validate with Terraship
        uses: terraship/terraship-action@v1
        with:
          terraform-directory: ./terraform
          policy-path: ./policy.yml
          mode: validate-existing
```

## Next Steps

- **Customize policies**: Edit `policy.yml` to match your requirements
- **Add more rules**: See the [Policy Guide](policy-guide.md)
- **Integrate with CI/CD**: Use the GitHub Action or CLI in your pipeline
- **Install VS Code extension**: Get real-time validation in your editor
- **Explore examples**: Check the `examples/` directory

## Common Issues

### "Terraform binary not found"

Ensure Terraform is installed and in your PATH:

```bash
terraform --version
```

### "No cloud provider detected"

Manually specify the provider:

```bash
terraship validate . --provider aws --policy policy.yml
```

### "Policy file not found"

Use the full path to your policy file:

```bash
terraship validate . --policy /full/path/to/policy.yml
```

### "Credentials validation failed"

Check your cloud provider credentials:

```bash
# AWS
aws sts get-caller-identity

# Azure
az account show

# GCP
gcloud auth list
```

## Getting Help

- 📖 [Full Documentation](../README.md)
- 💬 [GitHub Discussions](https://github.com/vijayaxai/terraship/discussions)
- 🐛 [Report Issues](https://github.com/vijayaxai/terraship/issues)
- 📧 [Email Support](mailto:support@terraship.io)

---

Ready to ship secure infrastructure! 🚢
