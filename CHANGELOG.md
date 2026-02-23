# Changelog

All notable changes to Terraship are documented in this file.

## [1.3.0] - 2026-02-23

### 🎯 Major Release: Production-Ready Security & Compliance Rules

#### ✨ New Features
- **25 Granular Production-Ready Rules** - Expanded policy framework from 16 to 41 comprehensive rules
  - Encryption specificity (KMS, TLS, database encryption in transit)
  - Authentication & access control (MFA, root hardening, credential rotation)
  - Audit & compliance (CloudTrail multi-region, immutable logs, 90-day retention)
  - Network security (restrict SSH/RDP, NAT gateway, WAF requirements)
  - Database hardening (delete protection, backup retention, enhanced monitoring)
  - Comprehensive tagging & governance (5 required tags)
  - Cost optimization (latest instances, auto-scaling, cross-region replication)

#### 📊 Policy Improvements
- Total rules increased: 16 → 41 comprehensive policies
- ERROR severity rules: 6 → 14 critical checks
- WARNING severity rules: 7 → 16 recommended checks
- Production-ready compliance mappings
- Multi-cloud coverage for AWS, Azure, GCP

#### 📥 Installation
```bash
go install github.com/vijayaxai/terraship/cmd/terraship@v1.3.0
terraship --version  # Shows: Terraship 1.3.0
```

---

## [1.2.1] - 2026-02-23

### 🔧 Patch Release: Stability & Performance

#### ✨ Improvements
- Stabilized advanced HTML reporting features
- Optimized resource data population performance
- Improved error handling and edge cases
- Enhanced CLI messaging and user feedback

#### 🐛 Bug Fixes
- Fixed potential issues with large terraform configurations
- Improved resource filtering accuracy
- Fixed chart rendering on edge cases
- Enhanced dark mode CSS for better readability

#### 📦 Releases
- **CLI v1.2.1** - Stable patch release
- **VS Code Extension v0.4.1** - Matching patch release

#### 📥 Installation
```bash
go install github.com/vijayaxai/terraship/cmd/terraship@v1.2.1
terraship --version  # Shows: v1.2.1
```

---

## [1.2.0] - 2026-02-20

### 🎉 Major Release: Advanced HTML Reporting & Critical Bug Fix

#### 🐛 Critical Bug Fixes
- **Fixed Resource Data Population** - Resources now display with proper names, types, and providers
  - Root cause: Incorrect field mapping in `convertResourcesToOutputFormat()`
  - Previously: Resources appeared as numbers instead of meaningful data
  - Now: Full resource details with validation checks and remediation guidance
  - Impact: All HTML reports now fully functional with complete resource information

#### ✨ New Advanced HTML Features
- **Real-Time Resource Search** - Search by name, type, or provider
- **Status & Type Filters** - Quick filtering with dropdown menus
- **Dark Mode Toggle** - Professional dark theme with persistent storage
- **Chart.js Visualizations**:
  - Compliance doughnut chart (Passed/Failed/Warnings)
  - 7-day validation timeline chart
- **Remediation Guidance** - Quick fix suggestions for each failed check
- **Comparison Views** - Track improvements across validation runs
- **Responsive Design** - Works on desktop, tablet, and mobile

#### 🔧 Technical Improvements
- Fixed `result.Passed` (boolean) mapping instead of nonexistent `result.Status`
- Fixed `result.Details` (array) handling instead of string concatenation
- Improved data flow from core validation to output formatters
- Enhanced template rendering for resource details

#### 📦 Releases
- **CLI v1.2.0** - With advanced HTML reporting
- **VS Code Extension v0.4.0** - Same advanced features available in marketplace

#### 📥 Installation
```bash
go install github.com/vijayaxai/terraship/cmd/terraship@v1.2.0
terraship --version  # Shows: v1.2.0
```

#### 📊 Usage Examples
```bash
# Generate advanced HTML with all features
terraship validate ./terraform --output html --html-advanced

# Track validation history
terraship validate ./terraform --output html --include-history

# Compare with previous run
terraship validate ./terraform --output html --compare previous-report.json
```

---

## VS Code Extension [0.3.1] - 2026-02-20

### 🐛 Fixes & Updates
- **Fixed Binary Distribution**: Recompiled CLI with full HTML/PDF reporting support
- **HTML Report Generation**: CLI now fully supports `--output html` flag
- **PDF Export**: CLI now fully supports `--output pdf` flag  
- **Multiple Formats**: Can now generate multiple report formats in one command: `--output html,pdf,json,sarif`

### ✨ Features Coming to Extension
- Integration with new HTML report viewer
- Support for advanced HTML features (dark mode, charts, search)
- Comparison view for validation history
- Direct report file handling

---

## [1.1.0] - 2026-02-19

### 🚀 Major Features

#### Comprehensive Reporting System
- **Interactive HTML Reports** - Beautiful, responsive web-based validation reports with:
  - Expandable resource details and policy checks
  - Real-time filtering by status (Passed/Failed/Warnings)
  - Compliance dashboard with percentage scoring
  - Side-by-side comparison with previous validation runs
  - Dark mode toggle (with `--html-advanced` flag)
  - Timeline charts showing validation history (with `--html-advanced` flag)
  - Print-friendly styling for PDF export from browser

- **PDF Export** - Professional PDF reports with:
  - Native PDF generation via `wkhtmltopdf` (auto-detected)
  - Intelligent fallback to HTML when tool unavailable
  - Platform-specific installation instructions (macOS/Ubuntu/Windows)
  - Browser print-to-PDF support as fallback

- **Machine Formats**
  - **JSON Export** - Structured data for CI/CD pipelines with compliance calculations
  - **SARIF 2.1.0** - GitHub Code Scanning integration for automatic security alerts
  - **Human Format** - Default terminal output with colored compliance summary

### 🎯 New CLI Capabilities

#### Enhanced `validate` Command
```bash
# Generate interactive HTML report
terraship validate ./terraform --output html

# Generate PDF report
terraship validate ./terraform --output pdf

# Generate multiple formats at once
terraship validate ./terraform --output html,pdf,json,sarif

# Compare with previous validation
terraship validate ./terraform --compare previous-report.json

# Advanced HTML features
terraship validate ./terraform --output html --html-advanced

# Include validation history
terraship validate ./terraform --output html --include-history
```

#### New Flags
- `--output` - Output format: human, json, html, pdf, sarif (comma-separated for multiple)
- `--output-file` - Custom output filename (auto-named if not specified)
- `--html-advanced` - Enable advanced HTML features (dark mode, charts)
- `--include-history` - Include 7-day validation history in reports
- `--compare` - Compare with previous validation results (JSON file path)

### 📊 Report Features

| Feature | Human | HTML | PDF | JSON | SARIF |
|---------|-------|------|-----|------|-------|
| Terminal Output | ✓ | - | - | - | - |
| Interactive UI | - | ✓ | ✓ | - | - |
| Compliance Score | ✓ | ✓ | ✓ | ✓ | - |
| Resource Details | ✓ | ✓ | ✓ | ✓ | ✓ |
| Filtering | - | ✓ | ✓ | - | - |
| Comparison | - | ✓ | ✓ | - | - |
| History Timeline | - | ✓ | ✓ | - | - |
| Dark Mode | - | ✓ | ✓ | - | - |
| GitHub Integration | - | - | - | - | ✓ |

### 🔧 Internal Changes

- Created `internal/output/html_reporter.go` - Interactive HTML report generation
- Created `internal/output/pdf_reporter.go` - PDF export with fallback support
- Created `internal/output/validation_result.go` - Result types and export formats
- Enhanced `cmd/terraship/commands/validate.go` - Multi-format report generation
- Removed duplicate `validate_advanced.go` - Consolidated into `validate.go`

### 📚 Documentation

- Added comprehensive `docs/ADVANCED-REPORTING.md` (~500 lines)
  - 6 output format examples with CLI commands
  - GitHub Actions CI/CD workflow examples
  - Integration examples (Slack, Email, Webhooks)
  - Troubleshooting guide
  - Local development setup
  - 8+ real-world usage scenarios

- Updated `README.md`:
  - New "📊 Reporting" section with all examples
  - Installation instructions for PDF tools
  - GitHub Actions workflow examples
  - Updated Features section highlighting reporting capabilities

### ✅ Testing

- Added 9 comprehensive tests for reporting features:
  - JSON export format and compliance calculation
  - SARIF 2.1.0 format validation
  - Compliance percentage calculations (edge cases)
  - HTML report generation
  - PDF reporter initialization
  - Installation instructions availability
- All tests pass with 100% passing rate (14/14 total tests)

### 🏗️ Architecture

- Consolidated reporting pipeline: `ValidationResult` → format-specific generators
- Graceful fallback strategies (PDF → HTML with print instructions)
- No external Go dependencies added (uses only stdlib + Cobra)
- Modular design allows easy addition of new export formats

### 🐛 Bug Fixes

- Fixed error handling in validate command for missing policy files
- Improved CLI error messages with actionable suggestions
- Added color-coded output for report generation status

## [0.1.8] - 2026-02-18

### 🎉 Features

- Added `terraship init` command to generate sample policy file
- VS Code Extension v0.1.8 published to marketplace
- Improved error messages with policy file guidance
- Added support for 8 core security rules

### 📚 Documentation

- Created Getting Started guide
- Added policy file explanation
- Documented VS Code extension configuration

## [0.1.0] - 2026-01-15

### Initial Release

- Multi-cloud Terraform validation (AWS, Azure, GCP)
- Policy-based rule engine (YAML)
- Two validation modes: validate-existing, ephemeral-sandbox
- Drift detection
- Terraform CLI integration
- Basic policy rules (8 core rules)

---

## Upgrade Guide

### From 0.1.8 to 1.1.0

**Breaking Changes:** None

**New Capabilities:**
1. Try the new HTML reports:
   ```bash
   terraship validate ./terraform --output html
   ```

2. Export to JSON for scripting:
   ```bash
   terraship validate ./terraform --output json
   ```

3. Integrate with GitHub Code Scanning:
   ```bash
   terraship validate ./terraform --output sarif
   # Upload the terraship-report.sarif to GitHub
   ```

4. See the new documentation:
   ```bash
   cat docs/ADVANCED-REPORTING.md
   ```

## Planned Features

### Upcoming in v1.2.0
- [ ] Slack webhook integration for alert notifications
- [ ] Email report delivery
- [ ] Team licensing and premium rules
- [ ] Custom rule support (user-defined policies)
- [ ] Analytics dashboard with historical trends

### Upcoming in v2.0.0
- [ ] Cloud-native policy server
- [ ] API for programmatic access
- [ ] Advanced compliance frameworks (CIS, SOC2, PCI-DSS)
- [ ] Multi-account/multi-environment support
- [ ] Custom dashboards and reporting

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE)
