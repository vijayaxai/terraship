# Terraship VS Code Extension

🚢 Multi-cloud Terraform validation and policy checking for AWS, Azure, and GCP.

## ✨ Latest Version

**Current Version:** 0.4.1 - Stable Release with Advanced HTML Reporting

Features now include interactive HTML reports with search, filters, charts, dark mode, and compliance tracking.

## 🌟 Features

- ✅ **Real-time Policy Validation** - Check Terraform against your policies as you code
- ✅ **Multi-Cloud Support** - Works with AWS, Azure, and GCP
- ✅ **Inline Error Reporting** - See violations directly in your editor
- ✅ **Quick Fix Suggestions** - Get remediation guidance for each issue
- ✅ **Customizable Policies** - Define your own security and compliance rules
- ✅ **Multiple Output Formats** - Human-readable, JSON, and SARIF reports

## 📦 Installation

### From Marketplace (Coming Soon)

1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "Terraship"
4. Click Install

### Manual Installation (Beta)

```powershell
code --install-extension terraship-vscode-0.1.0.vsix
```

## ⚙️ Prerequisites

**⚠️ REQUIRED: Install Terraship CLI first!**

The extension needs the Terraship CLI to work. Without it, you'll get `spawn terraship ENOENT` error.

### Option 1: Go Install (Recommended)

If you have Go 1.16+ installed:

```bash
go install github.com/vijayaxai/terraship/cmd/terraship@latest
```

Binary installs to `$GOPATH/bin` (usually `~/go/bin` or `C:\Users\<username>\go\bin`). Make sure it's in PATH.

**Or install a specific version:**
```bash
go install github.com/vijayaxai/terraship/cmd/terraship@v1.2.1
```

Verify:
```bash
terraship --version
```

### Option 2: Download Pre-Built Binary

#### Windows
```powershell
Invoke-WebRequest -Uri "https://github.com/vijayaxai/terraship/releases/latest/download/terraship-windows-amd64.exe" -OutFile "$env:USERPROFILE\bin\terraship.exe"

# Add to PATH or set full path in extension settings
# In VS Code: "terraship.executablePath": "C:\\Users\\YourName\\bin\\terraship.exe"

terraship --version
```

#### macOS / Linux
```bash
curl -L https://github.com/vijayaxai/terraship/releases/latest/download/terraship-$(uname -s)-$(uname -m) -o /usr/local/bin/terraship
chmod +x /usr/local/bin/terraship

terraship --version
```

### Option 3: Build from Source

```bash
git clone https://github.com/vijayaxai/terraship
cd terraship
go build -o bin/terraship ./cmd/terraship

# Copy to PATH location or use full path in settings
cp bin/terraship /usr/local/bin/  # macOS/Linux
# or
copy bin\terraship.exe C:\Users\<username>\bin\  # Windows
```

## 🚀 Quick Start

### 1. Create Your Policy File

Create a `policy.yml` with your security rules:

```yaml
version: "1.0"
name: "Security Policy"

rules:
  - name: "required-tags"
    severity: "error"
    resource_types: ["azurerm_*", "aws_*"]
    conditions:
      tags.required: ["Environment", "Owner"]
    message: "All resources must have required tags"
    
  - name: "enforce-encryption"
    severity: "error"
    resource_types: ["azurerm_storage_account"]
    conditions:
      enable_https_traffic_only: true
    message: "Storage accounts must enforce HTTPS"
```

### 2. Configure Extension

Open VS Code Settings (Ctrl+,) and search for "terraship":

```json
{
  "terraship.policyPath": "./policies/policy.yml",
  "terraship.cloudProvider": "azure",
  "terraship.executablePath": "terraship",
  "terraship.validateOnSave": false
}
```

**On Windows:** If Terraship is not in PATH, set the full path:
```json
{
  "terraship.executablePath": "C:\\Users\\YourName\\bin\\terraship.exe"
}
```

### 3. Validate Your Terraform

Open any `.tf` file, then:

**Option A: Command Palette**
1. Press `Ctrl+Shift+P`
2. Type "Terraship"
3. Select "Terraship: Validate Workspace"

**Option B: Right-Click**
1. Right-click in a `.tf` file
2. Select "Terraship: Validate Current File"

### 4. View Results

Results appear in:
- **Problems Panel** - See violations with line numbers
- **Output Panel** - View detailed validation report
- **Inline** - Hover over code to see issues

## 📊 Output Formats

Terraship generates validation reports in three formats:

### Human-Readable (Default)
Console output with formatted summary and violation details.

### JSON Format
Structured format for programmatic access and CI/CD integration:
```json
{
  "total_resources": 15,
  "passed_resources": 12,
  "failed_resources": 3,
  "resources": [...]
}
```

### SARIF Format
Standardized machine-readable format compatible with:
- GitHub Code Scanning
- GitLab Security Scanning
- Azure DevOps
- Other SARIF-compatible tools

**Example CLI usage:**
```bash
terraship validate ./terraform --output json --output-file report.json
terraship validate ./terraform --output sarif --output-file report.sarif
terraship validate ./terraform --output human
```

## ⚙️ Configuration

| Setting | Description | Default |
|---------|-------------|---------|
| `terraship.policyPath` | Path to your policy YAML file | `./policies/sample-policy.yml` |
| `terraship.cloudProvider` | Cloud provider (aws, azure, gcp, or empty for auto-detect) | `""` |
| `terraship.mode` | Validation mode | `validate-existing` |
| `terraship.validateOnSave` | Auto-validate on file save | `false` |
| `terraship.executablePath` | Path to Terraship CLI executable | `terraship` |
| `terraship.azureSubscriptionId` | Azure Subscription ID | `""` |
| `terraship.azureTenantId` | Azure Tenant ID | `""` |
| `terraship.awsProfile` | AWS Profile name | `""` |
| `terraship.gcpProject` | GCP Project ID | `""` |

## 🔐 Credential Configuration

### Azure Credentials

Set in VS Code Settings (Ctrl+,):
```json
{
  "terraship.azureSubscriptionId": "d30ec219-d601-414b-98b6-230b6e520d37",
  "terraship.azureTenantId": "2111de49-6a33-4187-af6d-96575525e6ef"
}
```

Or via environment variables (takes precedence):
```bash
$env:AZURE_SUBSCRIPTION_ID="your-id"
$env:AZURE_TENANT_ID="your-id"
```

### AWS Credentials

Set in VS Code Settings:
```json
{
  "terraship.awsProfile": "my-profile"
}
```

Or via environment variables:
```bash
$env:AWS_PROFILE="my-profile"
$env:AWS_ACCESS_KEY_ID="your-key"
$env:AWS_SECRET_ACCESS_KEY="your-secret"
```

### GCP Credentials

Set in VS Code Settings:
```json
{
  "terraship.gcpProject": "my-project-id"
}
```

Or via environment variables:
```bash
$env:GCP_PROJECT="my-project-id"
$env:GOOGLE_APPLICATION_CREDENTIALS="/path/to/key.json"
```

### Troubleshooting: "spawn terraship ENOENT"

This error means the extension cannot find the Terraship CLI. Fix it:

1. **Verify CLI is installed:**
   ```bash
   terraship --version
   ```

2. **If not in PATH, set full path in VS Code Settings:**
   ```json
   {
     "terraship.executablePath": "/usr/local/bin/terraship"  // macOS/Linux
     // or
     "terraship.executablePath": "C:\\Users\\YourName\\bin\\terraship.exe"  // Windows
   }
   ```

3. **Reload VS Code** after changing settings (Ctrl+Shift+P → "Reload Window")


## 📝 Example Policy

Create a `policy.yml` file:

```yaml
version: "1.0"
name: "Security Policy"

rules:
  - name: "required-tags"
    severity: "error"
    resource_types: ["azurerm_*", "aws_*"]
    conditions:
      tags.required: ["Environment", "Owner"]
    message: "All resources must have Environment and Owner tags"
```

## 🎯 Use Cases

### For Developers
- Catch policy violations before commit
- Learn cloud best practices while coding
- No context switching - validate in VS Code

### For DevOps Teams
- Enforce infrastructure standards
- Automated compliance checking
- Consistent validation across team

### For Security Teams
- Prevent misconfigurations early
- Track compliance violations
- Enforce encryption and access controls

## 🐛 Known Issues (Beta)

- [ ] Drift detection requires deployed resources
- [ ] Some encryption checks need refinement
- [ ] Performance optimization for large workspaces

## 📊 Roadmap

- [ ] Auto-fix for common violations
- [ ] Custom rule functions
- [ ] Integration with GitHub/GitLab CI
- [ ] Real-time validation as you type
- [ ] Terraform state file analysis

## 🤝 Contributing

Found a bug or have a feature request?

1. Check existing issues: https://github.com/vijayaxai/terraship/issues
2. Create new issue with "beta" label
3. Provide:
   - VS Code version
   - Extension version
   - Terraform version
   - Steps to reproduce

## 📄 License

MIT License - See LICENSE file

## 🔗 Links

- **GitHub:** https://github.com/vijayaxai/terraship
- **Documentation:** See project README
- **Issues:** https://github.com/vijayaxai/terraship/issues
- **Changelog:** See CHANGELOG.md

## 💬 Support

- **Questions:** Create a GitHub Discussion
- **Bugs:** Create a GitHub Issue
- **Email:** support@terraship.io (coming soon)

---

**Made with ❤️ by the Terraship Team**

*Beta testers rock! 🎸 Thanks for helping us improve.*
