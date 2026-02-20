## v1.2.0 - Advanced HTML Reporting & Resource Data Population Fix

**Release Date:** February 20, 2026

### 🎉 Major Features

#### Advanced HTML Reporting (NEW)
- **Interactive Compliance Dashboard**: Real-time compliance score with visual compliance status
- **Dark Mode Toggle**: Professional dark theme for comfortable night viewing with persistent storage
- **Smart Search**: Real-time search by resource name, type, or provider
- **Status & Type Filters**: Quick filtering dropdown for resource status and resource types
- **Chart.js Visualizations**:
  - Compliance doughnut chart showing Passed/Failed/Warning distribution
  - 7-day validation timeline to track trends
- **Expandable Resource Details**: Collapsible resource sections with validation checks
- **Remediation Guidance**: Quick fix suggestions for each failed check
- **Responsive Design**: Optimized for desktop, tablet, and mobile devices
- **Print-to-PDF**: Export reports directly from browser

#### CLI Enhancements
```bash
# Generate advanced HTML with all features
terraship validate ./terraform --output html --html-advanced

# Include validation history (timeline of past runs)
terraship validate ./terraform --output html --include-history

# Compare with previous validation runs
terraship validate ./terraform --output html --compare previous-report.json
```

### 🐛 Bug Fixes

#### Critical: Resource Data Population (HIGH PRIORITY)
- **Fixed**: Resources now display with proper names, types, and providers in HTML reports
- **Root Cause**: `convertResourcesToOutputFormat()` was mapping incorrect fields from cloud validation results
- **Details**:
  - Corrected field mapping: `result.Passed` (boolean) instead of nonexistent `result.Status`
  - Fixed `result.Details` array handling instead of string concatenation
  - All validation details and remediation guidance now fully rendered
  - Resources previously displayed as numbers, now show meaningful data

### 📦 Releases

#### CLI Binary v1.2.0
Available for download:
- **macOS (Intel)**: `terraship-darwin-amd64`
- **macOS (Apple Silicon)**: `terraship-darwin-arm64`
- **Linux (x64)**: `terraship-linux-amd64`
- **Windows**: `terraship-windows-amd64.exe`

#### VS Code Extension v0.4.0
- **Status**: Published to VS Code Marketplace
- **Features**: Same advanced HTML reporting capabilities
- **Search**: Look for "Terraship" in VS Code Extensions marketplace
- **Link**: [VS Code Marketplace - Terraship Extension](https://marketplace.visualstudio.com/items?itemName=terraship.terraship-vscode)

### 📥 Installation

#### Go Install (CLI)
```bash
go install github.com/vijayaxai/terraship/cmd/terraship@v1.2.0
terraship --version  # Should show: Terraship v1.2.0
```

#### VS Code Extension
1. Open VS Code → Extensions (Ctrl+Shift+X)
2. Search for "Terraship"
3. Click Install

### 📊 Sample Usage

```bash
# Basic validation with advanced HTML
terraship validate ./terraform --output html --html-advanced

# Generate all formats
terraship validate ./terraform --output html,json,sarif --html-advanced

# Save to custom file
terraship validate ./terraform --output html --output-file compliance-report.html --html-advanced

# Compare with previous results
terraship validate ./terraform --output html --html-advanced --compare previous-report.json
```

### 🧪 Testing

All validation modes tested and working:
- ✅ Validate existing infrastructure
- ✅ Ephemeral sandbox mode
- ✅ Multi-cloud providers (AWS, Azure, GCP)
- ✅ HTML report generation with resource data population
- ✅ JSON, SARIF, and PDF export formats
- ✅ Search and filter functionality
- ✅ Dark mode toggle
- ✅ Comparison views

### 🔄 Migration from v1.1.0

If upgrading from v1.1.0, no breaking changes:
- All existing policies continue to work
- All existing commands remain compatible
- Existing reports remain valid

To use new features:
```bash
# Simply add the flag to existing commands
terraship validate ./terraform --output html --html-advanced
```

### 📝 Documentation

- **Quick Start**: Run `terraship init` to create a sample policy
- **Policy Examples**: See `policies/sample-policy.yml`
- **HTML Reports**: New `--html-advanced` flag enables all interactive features
- **GitHub Integration**: Use SARIF output with GitHub Code Scanning

### 🙏 Acknowledgments

- Built with Go, Chart.js, and Terraform
- Cloud SDKs: AWS SDK, Azure SDK, Google Cloud SDK
- CLI framework: Cobra

### 📞 Feedback & Support

- Report issues: [GitHub Issues](https://github.com/vijayaxai/terraship/issues)
- Questions: Open a discussion on GitHub
- Security concerns: security@terraship.io

---

**Made with ❤️ by the Terraship team**
