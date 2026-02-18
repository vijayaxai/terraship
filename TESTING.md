# 🧪 Terraship Testing Guide

## Prerequisites Installation

### 1. Install Go 1.22+

**Windows (using Chocolatey)**:
```powershell
# Install Chocolatey if not installed
# Run as Administrator
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install Go
choco install golang -y

# Refresh environment
refreshenv
```

**Windows (Manual)**:
1. Download from: https://go.dev/dl/
2. Install the MSI installer
3. Restart PowerShell
4. Verify: `go version`

**Alternative - Use winget**:
```powershell
winget install GoLang.Go
```

### 2. Install Terraform

```powershell
# Using Chocolatey
choco install terraform -y

# Or download from: https://www.terraform.io/downloads
```

### 3. Verify Installation

```powershell
go version          # Should show: go version go1.22.x
terraform version   # Should show: Terraform v1.x.x
```

## Quick Test (Without Building)

### ⚠️ Corporate Network? Start Here!

If `go mod download` fails with **403 Forbidden** errors:

```powershell
# STEP 1: Try alternative Go proxy
go env -w GOPROXY=https://goproxy.cn,direct

# STEP 2: Retry download
go mod download

# STEP 3: If still fails, try Athens proxy
go env -w GOPROXY=https://athens.azurefd.net,direct
go mod download

# STEP 4: If all proxies fail, check syntax first
go build -o NUL .\cmd\terraship
# This will show if code is valid (ignore missing package errors)
```

**Need immediate testing without dependencies?** Jump to [Quick Start Testing (No Setup)](#quick-start-testing-no-setup)

### Step 1: Download Dependencies

```powershell
cd "c:\Users\vijayamalraj.arulx\OneDrive - HCL Technologies Ltd\Documents\AI Project\terraship"

# Download Go modules
go mod download
go mod tidy
```

### Step 2: Run Unit Tests

```powershell
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test -v ./internal/rules/
go test -v ./internal/terraform/
go test -v ./pkg/terraship/

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.txt ./...
go tool cover -html=coverage.txt -o coverage.html
# Open coverage.html in browser
```

### Step 3: Run Short Tests (Skip Integration)

```powershell
# Skip long-running integration tests
go test -short ./...
```

## Build and Test the CLI

### Step 1: Build the Binary

```powershell
# Build for Windows
go build -o bin\terraship.exe .\cmd\terraship

# Verify build
.\bin\terraship.exe --version
.\bin\terraship.exe --help
```

### Step 2: Test with Example

```powershell
# Navigate to example
cd examples\aws

# Initialize Terraform (if you have it)
terraform init

# Test validation (dry run - no cloud access needed)
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml
```

## Testing Individual Components

### 1. Test Rules Engine

```powershell
# Run rules engine tests
go test -v ./internal/rules/ -run TestRulesEngine

# Test specific rule
go test -v ./internal/rules/ -run TestRulesEngine_RequiredTags
```

### 2. Test Terraform Client

```powershell
# Requires Terraform installed
go test -v ./internal/terraform/
```

### 3. Test Cloud Adapters

```powershell
# Note: These may require cloud credentials
# AWS adapter
go test -v ./internal/cloud/aws/

# Azure adapter  
go test -v ./internal/cloud/azure/

# GCP adapter
go test -v ./internal/cloud/gcp/
```

### 4. Test Public API

```powershell
# Test Terratest integration API
go test -v ./pkg/terraship/
```

## Testing Without Cloud Credentials

Most tests can run without cloud credentials. Here's what you can test:

```powershell
# ✅ Rules engine - No credentials needed
go test -v ./internal/rules/

# ✅ Terraform client - Basic tests work without Terraform
go test -v ./internal/terraform/ -short

# ✅ Output formatters - No credentials needed
go test -v ./internal/output/

# ✅ Core validator - Mock tests work
go test -v ./internal/core/ -short

# ⚠️ Cloud adapters - Need credentials for full tests
# But can run with -short flag for basic tests
go test -short -v ./internal/cloud/...
```

## Manual Testing Scenarios

### Scenario 1: Validate Policy File Parsing

```powershell
# Create a test Terraform file
New-Item -Path "test" -ItemType Directory -Force
Set-Content -Path "test\main.tf" -Value @"
resource "aws_s3_bucket" "test" {
  bucket = "my-test-bucket"
}
"@

# Build and run validation
go run .\cmd\terraship validate test --policy .\policies\sample-policy.yml
```

### Scenario 2: Test Output Formats

```powershell
# Human output (default)
.\bin\terraship.exe validate .\examples\aws --policy .\policies\sample-policy.yml

# JSON output
.\bin\terraship.exe validate .\examples\aws --policy .\policies\sample-policy.yml --output json

# SARIF output
.\bin\terraship.exe validate .\examples\aws --policy .\policies\sample-policy.yml --output sarif
```

### Scenario 3: Test CLI Commands

```powershell
# Test help
.\bin\terraship.exe --help
.\bin\terraship.exe validate --help

# Test init command
.\bin\terraship.exe init

# Test with verbose mode
.\bin\terraship.exe validate .\examples\aws --policy .\policies\sample-policy.yml --verbose
```

## VS Code Extension Testing

### Prerequisites

```powershell
# Install Node.js (if not installed)
winget install OpenJS.NodeJS.LTS

# Or using Chocolatey
choco install nodejs-lts -y
```

### Test Extension

```powershell
cd vscode-extension

# Install dependencies
npm install

# Compile TypeScript
npm run compile

# Run linter
npm run lint

# Package extension
npm run package

# Install locally for testing
code --install-extension terraship-vscode-0.1.0.vsix
```

## Automated Test Suite

Run the complete test suite:

```powershell
# Using Make (if you have it)
make test

# Or manually:

# 1. Format check
go fmt ./...

# 2. Vet
go vet ./...

# 3. Run tests
go test -v ./...

# 4. Coverage
go test -coverprofile=coverage.txt -covermode=atomic ./...

# 5. View coverage
go tool cover -html=coverage.txt
```

## Performance Testing

```powershell
# Benchmark tests
go test -bench=. ./...

# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## Integration Testing with Real Infrastructure

⚠️ **Warning**: These tests will interact with your cloud provider and may incur costs!

### AWS Setup

```powershell
# Set AWS credentials
$env:AWS_REGION = "us-east-1"
$env:AWS_PROFILE = "default"

# Or set access keys
$env:AWS_ACCESS_KEY_ID = "your-access-key"
$env:AWS_SECRET_ACCESS_KEY = "your-secret-key"

# Run AWS integration tests
go test -v ./internal/cloud/aws/ -tags=integration
```

### Azure Setup

```powershell
# Login to Azure
az login

# Set subscription
$env:AZURE_SUBSCRIPTION_ID = "your-subscription-id"

# Run Azure integration tests
go test -v ./internal/cloud/azure/ -tags=integration
```

### GCP Setup

```powershell
# Login to GCP
gcloud auth application-default login

# Set project
$env:GCP_PROJECT = "your-project-id"

# Run GCP integration tests
go test -v ./internal/cloud/gcp/ -tags=integration
```

## Troubleshooting Test Failures

### Corporate Network / Proxy Issues (403 Forbidden)

If you get "403 Forbidden" errors when downloading modules (common in corporate networks):

**Option 1: Configure Go to use alternative proxy**
```powershell
# Try Chinese proxy (goproxy.cn)
go env -w GOPROXY=https://goproxy.cn,direct

# Or try goproxy.io
go env -w GOPROXY=https://goproxy.io,direct

# Or try Athens proxy
go env -w GOPROXY=https://athens.azurefd.net,direct

# Then retry download
go mod download
```

**Option 2: Configure corporate proxy settings**
```powershell
# Set HTTP/HTTPS proxy (ask your IT for proxy URL)
$env:HTTP_PROXY = "http://your-proxy:port"
$env:HTTPS_PROXY = "http://your-proxy:port"

# Then retry
go mod download
```

**Option 3: Disable proxy and use direct**
```powershell
go env -w GOPROXY=direct
go mod download
```

**Option 4: Test without downloading dependencies first**
```powershell
# Just verify the code compiles (will show missing import errors)
go build -o NUL .\cmd\terraship

# This checks syntax without downloading modules
```

**Option 5: Download dependencies outside corporate network**
- Use personal internet/hotspot
- Run: `go mod download`
- Commit the `go.sum` file and `$GOPATH/pkg/mod` cache
- Or work from home network

**Option 6: Use Go module vendor**
```powershell
# Download on non-restricted network, then:
go mod vendor

# This copies all dependencies to vendor/ folder
# Then build using vendor:
go build -mod=vendor -o bin\terraship.exe .\cmd\terraship
```

**Check current proxy settings:**
```powershell
go env GOPROXY
go env GOPRIVATE
go env GONOSUMDB
```

### "Package not found"

```powershell
go mod download
go mod tidy
```

### "Terraform not found" in tests

```powershell
# Install Terraform or skip those tests
go test -short ./...
```

### "Cloud credentials" errors

```powershell
# Run without cloud integration
go test -short ./...

# Or set mock credentials for testing
$env:AWS_REGION = "us-east-1"
```

### Test timeout

```powershell
# Increase timeout
go test -timeout 10m ./...
```

## Expected Test Results

When all tests pass, you should see:

```
PASS: internal/rules
PASS: internal/terraform  
PASS: internal/cloud/aws
PASS: internal/cloud/azure
PASS: internal/cloud/gcp
PASS: internal/core
PASS: internal/output
PASS: pkg/terraship
PASS: cmd/terraship/commands

Coverage: ~80%+ expected
```

## Next Steps After Testing

1. ✅ All tests pass → Ready to use!
2. ⚠️ Some tests fail → Check the troubleshooting section
3. 🚀 Ready to deploy → See BUILD.md for release instructions

## Quick Start Testing (No Setup)

If you just want to verify the code structure without running tests:

```powershell
# Check if code compiles
go build -o NUL .\cmd\terraship

# If this succeeds, the code is valid!
```

## Getting Help

- Check test output for specific errors
- Review logs in `go test -v` mode
- See BUILD.md for build issues
- Check CONTRIBUTING.md for development guidelines

---

**Ready to test?** Start with: `go mod download && go test -short ./...`
