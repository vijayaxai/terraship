# 🚢 Terraship - Build & Deploy Guide

## Prerequisites

### Required Tools
- **Go 1.22+**: [Install Go](https://go.dev/doc/install)
- **Terraform 1.6+**: [Install Terraform](https://www.terraform.io/downloads)
- **Make**: Usually pre-installed on Unix systems
- **Git**: For version control

### Optional Tools
- **golangci-lint**: For code linting
- **Node.js 20+**: For VS Code extension
- **Docker**: For containerized builds

## Quick Build

```bash
# Clone the repository
git clone https://github.com/vijayaxai/terraship.git
cd terraship

# Download dependencies
make deps

# Build the binary
make build

# The binary will be in ./bin/terraship
./bin/terraship --version
```

## Development Build

### 1. Setup Development Environment

```bash
# Install dependencies
go mod download
go mod tidy

# Install golangci-lint (optional but recommended)
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Windows (using Chocolatey)
choco install golangci-lint
```

### 2. Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run only short tests (skip integration)
make test-short

# View coverage report
make coverage
# Open coverage.html in browser
```

### 3. Run Linters

```bash
# Run all linters
make lint

# Format code
make fmt

# Run go vet
make vet
```

### 4. Build Binary

```bash
# Build for current platform
make build

# Build with version information
VERSION=1.0.0 make build

# Install to GOPATH/bin
make install
```

## Cross-Platform Builds

### Build for All Platforms

```bash
# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o bin/terraship-linux-amd64 ./cmd/terraship
GOOS=darwin GOARCH=amd64 go build -o bin/terraship-darwin-amd64 ./cmd/terraship
GOOS=darwin GOARCH=arm64 go build -o bin/terraship-darwin-arm64 ./cmd/terraship
GOOS=windows GOARCH=amd64 go build -o bin/terraship-windows-amd64.exe ./cmd/terraship
```

### Using Build Script

Create a `build.sh` script:

```bash
#!/bin/bash
set -e

VERSION=${VERSION:-dev}
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME"

echo "Building Terraship $VERSION..."

platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS="/" read -r -a parts <<< "$platform"
    GOOS="${parts[0]}"
    GOARCH="${parts[1]}"
    
    output="bin/terraship-$GOOS-$GOARCH"
    if [ "$GOOS" = "windows" ]; then
        output="$output.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "$output" ./cmd/terraship
done

echo "Build complete!"
```

Make it executable and run:

```bash
chmod +x build.sh
VERSION=1.0.0 ./build.sh
```

## VS Code Extension Build

### Prerequisites

```bash
cd vscode-extension
npm install
```

### Development

```bash
# Compile TypeScript
npm run compile

# Watch for changes
npm run watch

# Run linter
npm run lint
```

### Package Extension

```bash
# Install vsce
npm install -g @vscode/vsce

# Package extension
npm run package

# This creates terraship-vscode-0.1.0.vsix
```

### Install Locally

```bash
code --install-extension terraship-vscode-0.1.0.vsix
```

## Docker Build

### Create Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o terraship ./cmd/terraship

FROM alpine:latest
RUN apk --no-cache add ca-certificates terraform

COPY --from=builder /app/terraship /usr/local/bin/
ENTRYPOINT ["terraship"]
```

### Build Docker Image

```bash
docker build -t terraship:latest .

# Run with Docker
docker run --rm \
  -v $(pwd):/workspace \
  -e AWS_PROFILE=default \
  terraship:latest validate /workspace
```

## CI/CD Integration

### GitHub Actions

The project includes `.github/workflows/ci.yml` that automatically:
- Runs linters
- Executes tests
- Builds binaries
- Runs integration tests
- Builds VS Code extension

### Manual Release

1. **Tag the release**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

2. **Build release binaries**:
   ```bash
   VERSION=1.0.0 ./build.sh
   ```

3. **Create GitHub Release**:
   - Go to GitHub Releases
   - Create new release
   - Upload binaries from `bin/`
   - Add release notes

## Testing the Build

### Smoke Test

```bash
# Version check
./bin/terraship --version

# Help command
./bin/terraship --help

# Validate command help
./bin/terraship validate --help

# Test with example
cd examples/aws
../../bin/terraship validate . --policy ../../policies/sample-policy.yml --mode validate-existing
```

### Integration Test

```bash
# Run with actual Terraform
cd examples/aws
terraform init
../../bin/terraship validate . --policy ../../policies/sample-policy.yml
```

## Distribution

### Homebrew (macOS)

Create a Homebrew formula:

```ruby
class Terraship < Formula
  desc "Multi-cloud Terraform validation tool"
  homepage "https://github.com/vijayaxai/terraship"
  url "https://github.com/vijayaxai/terraship/releases/download/v1.0.0/terraship-darwin-amd64"
  sha256 "..."
  version "1.0.0"

  def install
    bin.install "terraship-darwin-amd64" => "terraship"
  end

  test do
    system "#{bin}/terraship", "--version"
  end
end
```

### Chocolatey (Windows)

Create a Chocolatey package with `terraship.nuspec`.

### Snap (Linux)

Create a `snapcraft.yaml` for Snap packaging.

## Troubleshooting

### Build Fails with "Package Not Found"

```bash
# Clean and retry
go clean -cache
go mod tidy
go mod download
make build
```

### Cross-Compilation Errors

```bash
# Ensure CGO is disabled for cross-compilation
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/terraship
```

### VS Code Extension Issues

```bash
# Clean node_modules
cd vscode-extension
rm -rf node_modules
npm install
npm run compile
```

### Test Failures

```bash
# Run verbose tests
go test -v ./...

# Run specific test
go test -v -run TestRulesEngine ./internal/rules/...

# Skip integration tests
go test -short ./...
```

## Performance Optimization

### Build with Optimizations

```bash
go build -ldflags="-s -w" -o bin/terraship ./cmd/terraship
```

Flags:
- `-s`: Omit symbol table
- `-w`: Omit DWARF debug info

This reduces binary size by ~30%.

### Profile the Application

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## Next Steps

After building:

1. ✅ Run the test suite
2. ✅ Test with example Terraform code
3. ✅ Install locally for development
4. ✅ Create a release when ready
5. ✅ Publish VS Code extension to marketplace

## Resources

- [Go Build Documentation](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)
- [Cross Compilation Guide](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04)
- [VS Code Extension Publishing](https://code.visualstudio.com/api/working-with-extensions/publishing-extension)

---

Happy building! 🏗️
