# 🚢 Terraship

[![CI](https://github.com/vijayaxai/terraship/actions/workflows/ci.yml/badge.svg)](https://github.com/vijayaxai/terraship/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/vijayaxai/terraship)](https://goreportcard.com/report/github.com/vijayaxai/terraship)
[![License](https://img.shields.io/github/license/terraship/terraship)](LICENSE)

**Terraship** is a production-grade, multi-cloud Terraform validation tool that validates infrastructure against policies and detects drift. It works with **existing infrastructure** without requiring deployments and supports AWS, Azure, and GCP.

## ✨ Features

- 🔍 **Policy-Based Validation** - YAML-driven rules for security, compliance, and best practices
- ☁️ **Multi-Cloud Support** - AWS, Azure, and GCP with auto-detection
- 🔄 **Drift Detection** - Compare planned state with actual cloud resources
- 🎯 **Two Validation Modes**:
  - `validate-existing`: Validate existing infrastructure without applying changes
  - `ephemeral-sandbox`: Create temporary infrastructure, validate, and destroy
- 📊 **Multiple Output Formats** - Human-readable, JSON, and SARIF
- 🧪 **Terratest Integration** - First-class Go API for testing
- 🔌 **VS Code Extension** - Validate from your editor
- ⚡ **GitHub Action** - Integrate into CI/CD pipelines

## 🚀 Quick Start

### Installation

```bash
# macOS / Linux
curl -L https://github.com/vijayaxai/terraship/releases/latest/download/terraship-$(uname -s)-$(uname -m) -o /usr/local/bin/terraship
chmod +x /usr/local/bin/terraship

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/vijayaxai/terraship/releases/latest/download/terraship-windows-amd64.exe" -OutFile "terraship.exe"

# Go install
go install github.com/vijayaxai/terraship/cmd/terraship@latest
```

### Initialize Policy

Before validating, create a policy file:

```bash
# Create sample policy in current directory
terraship init

# Or in a specific directory
terraship init ./my-project
```

This creates `policies/terraship-policy.yml` with comprehensive security and compliance rules.

### Setup Cloud Credentials

Terraship needs credentials to connect to your cloud provider. Choose the method that works best for you:

#### Option 1: Azure CLI (Recommended for Azure)

```bash
# Login with Azure CLI
az login

# Export subscription ID (optional, but recommended)
# PowerShell
$env:AZURE_SUBSCRIPTION_ID = "your-subscription-id"

# Bash / macOS / Linux
export AZURE_SUBSCRIPTION_ID="your-subscription-id"

# Validate
terraship validate ./terraform
```

#### Option 2: Environment Variables

**Azure:**
```bash
# PowerShell
$env:AZURE_SUBSCRIPTION_ID = "d30ec219-d601-414b-98b6-230b6e520d37"
$env:AZURE_TENANT_ID = "2111de49-6a33-4187-af6d-96575525e6ef"

# Bash / macOS / Linux
export AZURE_SUBSCRIPTION_ID="d30ec219-d601-414b-98b6-230b6e520d37"
export AZURE_TENANT_ID="2111de49-6a33-4187-af6d-96575525e6ef"
```

**AWS:**
```bash
# PowerShell
$env:AWS_REGION = "us-east-1"
$env:AWS_PROFILE = "my-profile"  # Or use AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY

# Bash / macOS / Linux
export AWS_REGION="us-east-1"
export AWS_PROFILE="my-profile"
```

**GCP:**
```bash
# PowerShell
$env:GCP_PROJECT = "my-gcp-project"
$env:GOOGLE_APPLICATION_CREDENTIALS = "$env:USERPROFILE\.config\gcloud\application_default_credentials.json"

# Bash / macOS / Linux
export GCP_PROJECT="my-gcp-project"
export GOOGLE_APPLICATION_CREDENTIALS="$HOME/.config/gcloud/application_default_credentials.json"
```

#### Option 3: VS Code Extension Settings

If using the VS Code extension, configure credentials in your settings:

```json
{
  "terraship.azureSubscriptionId": "d30ec219-d601-414b-98b6-230b6e520d37",
  "terraship.azureTenantId": "2111de49-6a33-4187-af6d-96575525e6ef",
  "terraship.awsProfile": "my-profile",
  "terraship.gcpProject": "my-gcp-project"
}
```

### Basic Usage

```bash
# Validate with default policy
terraship validate ./terraform

# Validate existing infrastructure (no apply)
terraship validate ./terraform --mode validate-existing --policy ./policies/terraship-policy.yml

# Create ephemeral sandbox for testing
terraship validate ./terraform --mode ephemeral-sandbox

# Specify cloud provider and output format
terraship validate ./terraform --provider aws --output json

# Save results to file
terraship validate ./terraform --output json --output-file report.json
```

## 🎯 Getting Started with a New Terraform Project

### Scenario: You have a new Terraform project and want to validate it with Terraship

**Step 1: Navigate to your Terraform project**
```bash
cd your-terraform-project
ls -la
# You should see .tf files here
```

**Step 2: Create policy file**
```bash
terraship init
```
✅ This creates: `policies/terraship-policy.yml` with 8 security/compliance rules

**Step 3: Set up cloud credentials** (Choose one method)

**Option A - Azure CLI (Easiest):**
```bash
az login
```

**Option B - Environment Variables:**
```powershell
# PowerShell
$env:AZURE_SUBSCRIPTION_ID = "your-subscription-id"
```

**Step 4: Run validation from terminal**
```bash
terraship validate ./terraform
```
📊 View the validation report in terminal

**Step 5 (Optional): Use VS Code Extension for continuous validation**

1. Install extension: Open VS Code → Extensions (Ctrl+Shift+X) → Search "Terraship" → Install
2. Configure extension: Press `Ctrl + ,` → Search "terraship" → Set your cloud provider
3. Validate in editor: 
   - Press `Ctrl+Shift+P` → Type "Terraship: Validate"
   - Or set `"terraship.validateOnSave": true` for auto-validation

---

## 📋 Policy File Explained

When you run `terraship init`, it creates a policy file with the same 8 built-in rules from [policies/sample-policy.yml](policies/sample-policy.yml):

```yaml
# policies/terraship-policy.yml (generated by terraship init)

version: "1.0"
name: "Multi-Cloud Security and Compliance Policy"
description: "Comprehensive policy for AWS, Azure, and GCP resources"

rules:
  - name: "required-tags"           # ← All resources must have Environment, Owner, Project tags
  - name: "encryption-at-rest"      # ← All storage must be encrypted
  - name: "block-public-access"     # ← No public access to sensitive resources
  - name: "enable-versioning"       # ← Enable versioning for storage resources
  - name: "enable-logging"          # ← Enable audit logs
  - name: "backup-enabled"          # ← Enable backups for databases
  - name: "iam-least-privilege"     # ← No wildcard IAM permissions
  - name: "use-private-subnet"      # ← Deploy resources in private networks
```

**Same rules** in both:
- `terraship init` → Creates `policies/terraship-policy.yml` with these 8 rules
- `policies/sample-policy.yml` → Reference policy file with the same 8 rules

**You can customize this file** to add/remove rules for your organization's standards. See [policies/sample-policy.yml](policies/sample-policy.yml) for the full detailed policy with descriptions and remediation steps.



## 📝 Policy Configuration

Create a `policy.yml` file to define validation rules:

```yaml
version: "1.0"
name: "Security Policy"
description: "Security and compliance rules"

rules:
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

  - name: "encryption-at-rest"
    description: "Ensure encryption at rest is enabled"
    severity: "error"
    category: "security"
    enabled: true
    resource_types:
      - "aws_s3_bucket"
      - "aws_ebs_volume"
      - "azurerm_storage_account"
    conditions:
      encryption.enabled: true
    message: "Encryption at rest must be enabled"
    remediation: "Enable server-side encryption for your resource"
```

See [policies/sample-policy.yml](policies/sample-policy.yml) for a comprehensive example.

## 🧪 Terratest Integration

Use Terraship in your Terratest test suites:

```go
package test

import (
    "testing"
    "github.com/vijayaxai/terraship/pkg/terraship"
)

func TestInfrastructureCompliance(t *testing.T) {
    t.Parallel()

    opts := terraship.Options{
        TerraformDir:  "./terraform",
        PolicyPath:    "./policy.yml",
        Mode:          "validate-existing",
        CloudProvider: "aws",
    }

    result := terraship.Validate(t, opts)
    
    // Assert validation passed
    result.AssertCompliant(t)
    
    // Assert no drift detected
    result.AssertNoDrift(t)
    
    // Print detailed report
    t.Log(result.PrintReport())
}
```

## 🔌 GitHub Action

Add to your `.github/workflows/terraform-validation.yml`:

```yaml
name: Terraform Validation

on:
  pull_request:
    paths:
      - 'terraform/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Validate with Terraship
        uses: terraship/terraship-action@v1
        with:
          terraform-directory: ./terraform
          policy-path: ./policies/security-policy.yml
          mode: validate-existing
          output-format: sarif
          
      - name: Upload SARIF results
        uses: github/codeql-action/upload-sarif@v2
        if: always()
        with:
          sarif_file: terraship-report.txt
```

## 🎨 VS Code Extension

### Prerequisites

**⚠️ Install Terraship CLI first:**

```bash
# macOS / Linux
curl -L https://github.com/vijayaxai/terraship/releases/latest/download/terraship-$(uname -s)-$(uname -m) -o /usr/local/bin/terraship
chmod +x /usr/local/bin/terraship

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/vijayaxai/terraship/releases/latest/download/terraship-windows-amd64.exe" -OutFile "$env:USERPROFILE\bin\terraship.exe"
# Add to PATH or configure executablePath in VS Code settings

# Verify installation
terraship --version
```

### Install Extension

Install from the [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=terraship.terraship-vscode) or:

```bash
code --install-extension terraship.terraship-vscode
```

### Configure Extension

Open VS Code Settings (Ctrl+,) and configure:

```json
{
  "terraship.policyPath": "./policy.yml",
  "terraship.cloudProvider": "azure",
  "terraship.executablePath": "terraship"  // Or full path: "C:\\path\\to\\terraship.exe"
}
```

### Features
- Validate current workspace or file
- Real-time validation on save (optional)
- View detailed results in webview
- Quick fixes and remediation suggestions

## 🏗️ Architecture

```
terraship/
├── cmd/terraship/          # CLI application
├── pkg/terraship/          # Public Go API for Terratest
├── internal/
│   ├── cloud/             # Cloud provider adapters
│   │   ├── aws/           # AWS SDK integration
│   │   ├── azure/         # Azure SDK integration
│   │   └── gcp/           # GCP SDK integration
│   ├── terraform/         # Terraform operations
│   ├── rules/             # Policy rule engine
│   ├── core/              # Core validation logic
│   └── output/            # Output formatters
├── policies/              # Sample policies
├── vscode-extension/      # VS Code extension
└── action/                # GitHub Action
```

## 🔧 Prerequisites

Before running Terraship validation, ensure you have:

1. **Terraform** installed and in PATH
   - Verify: `terraform -version`
   - If terraform is not in PATH, add it or use `terraform` command directly

2. **Cloud CLI and Authentication**
   - **Azure**: `az login` (install [Azure CLI](https://learn.microsoft.com/cli/azure/install-azure-cli))
   - **AWS**: `aws configure` (install [AWS CLI](https://aws.amazon.com/cli/))
   - **GCP**: `gcloud auth application-default login` (install [Google Cloud SDK](https://cloud.google.com/sdk/docs/install))

3. **SSH Key** (for certain resource types)
   - Most cloud deployments need SSH keys
   - Generate with: `ssh-keygen -t rsa -b 4096`

## 🔧 Configuration

### Environment Variables

For more details, see the **Setup Cloud Credentials** section in Quick Start above.

#### Azure
- `AZURE_SUBSCRIPTION_ID` - Azure subscription ID *(required)*
- `AZURE_TENANT_ID` - Azure tenant ID *(optional, recommended)*
- `AZURE_CLIENT_ID` - Service principal client ID *(for non-interactive auth)*
- `AZURE_CLIENT_SECRET` - Service principal secret *(for non-interactive auth)*
- `AZURE_CLOUD_ENVIRONMENT` - Azure cloud environment (e.g., `AzurePublicCloud`, `AzureUSGovernmentCloud`)

#### AWS
- `AWS_REGION` - AWS region (default: `us-east-1`)
- `AWS_PROFILE` - AWS profile name
- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret key

#### GCP
- `GCP_PROJECT` or `GOOGLE_CLOUD_PROJECT` - GCP project ID *(required)*
- `GOOGLE_APPLICATION_CREDENTIALS` - Path to service account key
- `GOOGLE_CLOUD_REGION` - GCP region

**Example: Setting Environment Variables**

PowerShell:
```powershell
$env:AZURE_SUBSCRIPTION_ID = "d30ec219-d601-414b-98b6-230b6e520d37"
$env:PATH = "C:\terraform;" + $env:PATH  # Add Terraform to PATH
terraship validate .
```

Bash / macOS / Linux:
```bash
export AZURE_SUBSCRIPTION_ID="d30ec219-d601-414b-98b6-230b6e520d37"
export PATH="/usr/local/terraform:$PATH"  # Add Terraform to PATH
terraship validate .
```

## 📊 Output Formats

### Human (Default)
```
═══════════════════════════════════════════════════════════════
                    TERRASHIP VALIDATION REPORT                  
═══════════════════════════════════════════════════════════════

SUMMARY:
  Total Resources:    15
  ✓ Passed:           12
  ✗ Failed:           3
  ⚠ Warnings:         2
  ⨯ Errors:           0
  ↔ Drift Detected:   1

✗ VALIDATION FAILED
```

### JSON
```json
{
  "total_resources": 15,
  "passed_resources": 12,
  "failed_resources": 3,
  "reports": [...]
}
```

### SARIF
SARIF 2.1.0 format for integration with GitHub Code Scanning and other tools.

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/vijayaxai/terraship.git
cd terraship

# Install dependencies
make deps

# Run tests
make test

# Build
make build

# Run linters
make lint
```

## 🔒 Security

Please report security vulnerabilities to security@terraship.io. See [SECURITY.md](SECURITY.md) for details.

## 📄 License

Terraship is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Built with [Terratest](https://terratest.gruntwork.io/)
- Cloud SDKs: [AWS SDK for Go](https://aws.amazon.com/sdk-for-go/), [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go), [Google Cloud Go](https://cloud.google.com/go)
- CLI framework: [Cobra](https://github.com/spf13/cobra)

## 📚 Documentation

- [User Guide](docs/user-guide.md)
- [Policy Writing Guide](docs/policy-guide.md)
- [API Reference](docs/api-reference.md)
- [Examples](examples/)

## 🗺️ Roadmap

- [ ] Support for Terraform Cloud/Enterprise
- [ ] Custom rule functions
- [ ] Web UI for policy management
- [ ] Integration with OPA (Open Policy Agent)
- [ ] Support for additional cloud providers (Alibaba Cloud, Oracle Cloud)
- [ ] Remediation automation

---

Made with ❤️ by the Terraship team
