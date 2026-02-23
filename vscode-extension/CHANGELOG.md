# Changelog

All notable changes to the Terraship VS Code extension will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.1] - 2026-02-23

### Fixed
- **Ephemeral-Sandbox Cleanup** - Fixed resource leak when terraform apply fails
  - Ensure `terraform destroy` automatically runs even if apply encounters errors
  - Prevents orphaned Azure/AWS/GCP resources
  - CLI bug fix automatically propagated to extension through version sync

### Changed
- Version aligned with CLI v1.3.1 (patch release)
- Updated underlying Terraship engine with garbage collection fix

## [0.5.0] - 2026-02-23

### Changed
- **Version Alignment** - Extension version now synced with CLI v1.3.0
- **Enhanced Policy Support** - Now supports 41 granular production-ready policies (25 new rules)
  - Encryption specificity rules
  - Authentication & access control
  - Audit & compliance
  - Network security
  - Database hardening
  - Comprehensive tagging
  - Cost optimization

### Improved
- Documentation updated to reflect new policy capabilities
- Installation instructions now show v0.5.0 as latest version
- Better support for granular compliance checking

## [0.4.1] - 2026-02-20

### Fixed
- **Version Alignment** - Extension version now synced with CLI v1.2.1
- **Documentation** - Updated README with specific version installation instructions

### Changed
- Updated README to remove Beta notice - now stable release
- Enhanced installation documentation with version selection examples
- Improved prerequisites section with clearer setup guidance

### Performance
- Stable release with validated advanced HTML reporting features

## [0.4.0] - 2026-02-20

### Added
- **Comprehensive HTML Report Features** - Enhanced HTML reporting with interactive dashboards:
  - 🔍 Real-time search functionality for resources
  - 📊 Status and type filtering with interactive dropdowns
  - 📈 Chart.js visualizations (compliance breakdown doughnut chart + timeline trends)
  - 🌙 Dark mode toggle with localStorage persistence
  - 🆚 Comparison view for previous validation runs
  - 📱 Fully responsive design (desktop, tablet, mobile)
  - 💾 Print-friendly CSS for PDF export from browser

### Changed
- HTML report template completely redesigned for production-ready output
- Enhanced CSS with CSS variables for consistent theming
- Improved interactive resource filtering with real-time result count updates
- Better expandable/collapsible resource sections with status indicators

### Performance
- Embedded comprehensive ~12KB template as minified single-line string
- Optimized template loading for faster report generation
- Chart.js loaded from CDN for lightweight extension package

## [0.1.8] - 2026-02-19

### Added
- **terraship init Command** - CLI now has `terraship init` command to automatically generate policy files
  - Creates `policies/terraship-policy.yml` with 8 comprehensive security rules
  - Supports custom directory and filename options
  - Shows helpful next steps after policy creation

### Fixed
- **Policy File Error Messages** - Enhanced error guidance when policy file is missing:
  - Error message now suggests running `terraship init`
  - Provides examples of creating custom policies
  - Links to validation help documentation

### Changed
- Updated README with comprehensive Prerequisites section
- Added clear environment variable documentation for all cloud providers (Azure, AWS, GCP)
- Included PowerShell and Bash examples for credential setup
- Reorganized Quick Start section with policy initialization first
- Improved Configuration section with marked required vs optional variables
- Added VS Code Extension credential configuration in main README

### Documentation
- Added Prerequisites section covering Terraform, Cloud CLI, and SSH key setup
- Comprehensive credential setup methods in Quick Start (Option 1, 2, 3)
- Clarified environment variable requirements per cloud provider
- Added PATH configuration troubleshooting tips

## [0.1.7] - 2026-02-18

### Added
- **Cloud Credential Configuration Settings** - Added new settings for Azure, AWS, and GCP credentials:
  - `terraship.azureSubscriptionId` - Azure Subscription ID
  - `terraship.azureTenantId` - Azure Tenant ID
  - `terraship.awsProfile` - AWS Profile name
  - `terraship.gcpProject` - GCP Project ID

### Changed
- Extension now passes configured credentials to CLI automatically
- Updated documentation with credential setting examples

## [0.1.6] - 2026-02-18

### Added
- **Go Module Installation** - Users can now install CLI with `go install github.com/vijayaxai/terraship/cmd/terraship@latest`
- **Output Format Documentation** - Added comprehensive guides for JSON, human-readable, and SARIF formats
- **Credential Configuration Settings** - Added new settings for Azure, AWS, and GCP credentials:
  - `terraship.azureSubscriptionId` - Azure Subscription ID
  - `terraship.azureTenantId` - Azure Tenant ID
  - `terraship.awsProfile` - AWS Profile name
  - `terraship.gcpProject` - GCP Project ID
- **v1.0.0 Release Tag** - Go module properly versioned and published

### Fixed
- **CLI Import Paths** - Corrected imports from `terraship/terraship` to `vijayaxai/terraship`
- **Source Repository** - All cmd/terraship files now properly committed to git
- **Extension Independence** - Extension no longer depends on terraship project folder structure

### Changed
- Updated README with output formats section and CLI installation examples
- Enhanced documentation with format comparison
- Extension now passes configured credentials to CLI automatically

## [0.1.5] - 2026-02-18

### Changed
- Updated GitHub repository reference from `terraship/terraship` to `vijayaxai/terraship`
- Updated all documentation links to point to new repository
- Updated Go module path to `github.com/vijayaxai/terraship`

## [0.1.4] - 2026-02-18

### Fixed
- **ENOENT Error Handling**: Improved error messages when terraship CLI not found
- **Windows Path Support**: Auto-adds `.exe` extension on Windows for executable path
- **User Guidance**: Added helpful error dialog with step-by-step configuration instructions
- **One-Click Settings**: Added "Open Settings" button in error dialog for quick path configuration

### Changed
- Enhanced error detection to distinguish between missing executable and other errors
- Updated README with troubleshooting guide for ENOENT errors
- Improved documentation with clear Windows/macOS/Linux instructions

### Documentation
- Added troubleshooting section for "spawn terraship ENOENT" errors
- Clarified executablePath configuration for all platforms
- Updated quick start guide with proper configuration examples

## [0.1.0] - 2026-02-17 (BETA)

### Added
- Initial beta release
- Multi-cloud Terraform validation (AWS, Azure, GCP)
- Policy-based compliance checking
- Inline error reporting in editor
- Command palette integration
- Configurable policy paths
- Multiple validation modes (validate-existing, ephemeral-sandbox)
- Cloud provider auto-detection
- Settings for customization
- Validation on demand (workspace and file level)

### Known Issues
- Drift detection requires deployed resources
- Some encryption rule checks need refinement
- Performance not optimized for very large workspaces (>100 files)

### Coming Soon
- Auto-fix suggestions for violations
- Real-time validation on typing
- Quick fix integration
- Status bar indicators
- Report export functionality

## [Unreleased]

### Planned for 0.2.0
- Auto-fix for common violations
- Improved error messages
- Performance optimizations
- Better Azure authentication handling
- Custom rule templates

### Planned for 1.0.0 (GA)
- Production-ready stability
- Complete test coverage
- Performance benchmarks
- Comprehensive documentation
- Enterprise features (SSO, audit logs)

---

## Version Strategy

- **0.x.y** - Beta/Preview releases
- **1.x.y** - General Availability (GA)
- **x.0.0** - Major releases with breaking changes
- **x.y.0** - Minor releases with new features
- **x.y.z** - Patch releases with bug fixes
