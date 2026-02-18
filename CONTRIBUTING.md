# Contributing to Terraship

Thank you for your interest in contributing to Terraship! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please read it before contributing.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce**
- **Expected vs actual behavior**
- **Environment details** (OS, Go version, Terraform version)
- **Sample code or configuration** if applicable

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear title**
- **Provide detailed description**
- **Explain why this enhancement would be useful**
- **Provide examples** if applicable

### Pull Requests

1. **Fork the repository**
2. **Create a branch** from `develop`: `git checkout -b feature/my-feature`
3. **Make your changes** following our coding standards
4. **Add tests** for your changes
5. **Run tests and linters**: `make test lint`
6. **Commit with clear messages** following [Conventional Commits](https://www.conventionalcommits.org/)
7. **Push to your fork** and submit a pull request

## Development Setup

### Prerequisites

- Go 1.22 or later
- Terraform 1.6+ installed
- Make
- golangci-lint
- Node.js (for VS Code extension development)

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/terraship.git
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

## Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` and `goimports`
- Write clear, concise comments
- Add tests for new functionality
- Keep functions small and focused
- Use meaningful variable names

### Testing

- Write unit tests for all new code
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies
- Add integration tests for cloud provider interactions

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Example:
```
feat(aws): add S3 bucket encryption validation

Add validation rule to check if S3 buckets have encryption enabled.
Supports both SSE-S3 and SSE-KMS encryption methods.

Closes #123
```

## Project Structure

```
terraship/
├── cmd/terraship/          # CLI application
│   └── commands/           # Cobra commands
├── pkg/terraship/          # Public API
├── internal/
│   ├── cloud/             # Cloud provider adapters
│   │   ├── aws/
│   │   ├── azure/
│   │   └── gcp/
│   ├── terraform/         # Terraform client
│   ├── rules/             # Policy engine
│   ├── core/              # Core validation logic
│   └── output/            # Output formatters
├── policies/              # Sample policies
├── vscode-extension/      # VS Code extension
├── action/                # GitHub Action
└── tests/                 # Integration tests
```

## Adding a New Cloud Provider

1. Create a new package under `internal/cloud/`
2. Implement the `cloud.Adapter` interface
3. Add provider detection logic
4. Add resource validation methods
5. Write unit tests
6. Add integration tests
7. Update documentation

## Adding a New Validation Rule

1. Define the rule in a policy YAML file
2. Add condition evaluation logic in `internal/rules/engine.go`
3. Write unit tests
4. Add examples to `policies/sample-policy.yml`
5. Document in the policy guide

## Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run only unit tests
go test -short ./...

# Run specific package tests
go test ./internal/rules/...

# Run integration tests (requires cloud credentials)
go test -tags=integration ./tests/integration/...
```

## Documentation

- Update README.md for user-facing changes
- Add/update code comments
- Update API documentation
- Add examples if applicable

## Release Process

Releases are handled by maintainers:

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create a release branch
4. Tag the release
5. Build and publish binaries
6. Update documentation

## Questions?

- Open an issue for discussion
- Join our community chat
- Check existing documentation

Thank you for contributing to Terraship! 🚢
