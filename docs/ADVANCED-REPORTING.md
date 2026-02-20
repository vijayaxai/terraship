# Terraship Advanced Reporting Guide

## Overview

Terraship v1.1.0+ includes advanced reporting capabilities with multiple output formats and interactive features.

## Usage

### 1. **Human-Readable Report** (Default)

```bash
terraship validate ./terraform
```

Output:
```
═══════════════════════════════════════════════════════════
                TERRASHIP VALIDATION REPORT
═══════════════════════════════════════════════════════════

SUMMARY:
  Total Resources:    27
  ✓ Passed:           10
  ✗ Failed:           16
  ⚠ Warnings:         1

📊 Compliance Score: 37.0%
⏱  Validation completed: 2026-02-19 11:15 AM
```

---

### 2. **HTML Report** (Interactive)

```bash
# Generate interactive HTML report
terraship validate ./terraform --output html --output-file report.html

# Open in browser
open report.html
```

**Features:**
- ✅ Expandable resource details
- ✅ Color-coded status (green/red/yellow)
- ✅ Search functionality
- ✅ Filter by status (passed/failed/warnings)
- ✅ Professional styling
- ✅ Mobile-responsive design

---

### 3. **Advanced HTML Report**

```bash
# With dark mode, charts, and all features
terraship validate ./terraform \
  --output html \
  --output-file report-advanced.html \
  --html-advanced \
  --include-history

# Open in browser
open report-advanced.html
```

**Advanced Features:**
- 🌙 Dark/Light mode toggle
- 📈 Timeline charts (7-day validation history)
- 🔄 Compare current run vs previous
- 🔗 Team integrations (Slack, Email, Webhooks)
- 📊 Analytics dashboard
- 🔍 Advanced search and filtering

---

### 4. **PDF Report**

```bash
# Generate PDF report
terraship validate ./terraform --output pdf --output-file report.pdf

# Open PDF
open report.pdf
```

**Prerequisites:**
- Requires `wkhtmltopdf` installed

**Installation:**
```bash
# macOS
brew install wkhtmltopdf

# Ubuntu/Debian
sudo apt-get install wkhtmltopdf

# Windows
choco install wkhtmltopdf
```

**Fallback:** If `wkhtmltopdf` not available, generates HTML that you can print as PDF:
```
Open in browser → Ctrl+P (or Cmd+P) → Save as PDF
```

---

### 5. **JSON Report** (For CI/CD)

```bash
# Generate JSON report
terraship validate ./terraform --output json --output-file report.json

# View report
cat report.json | jq
```

**Output Example:**
```json
{
  "timestamp": "2026-02-19T11:15:00Z",
  "total_resources": 27,
  "passed_resources": 10,
  "failed_resources": 16,
  "warning_resources": 1,
  "compliance_percent": 37.0,
  "validation_passed": false,
  "resources": [...]
}
```

---

### 6. **SARIF Report** (GitHub Code Scanning)

```bash
# Generate SARIF report
terraship validate ./terraform --output sarif --output-file report.sarif

# Upload to GitHub
gh codeql database upload-results report.sarif
```

**GitHub Actions Integration:**
```yaml
- name: Validate with Terraship
  run: terraship validate ./terraform --output sarif --output-file terraship.sarif

- name: Upload to GitHub Code Scanning
  uses: github/codeql-action/upload-sarif@v2
  with:
    sarif_file: terraship.sarif
```

---

## Advanced Features

### Compare with Previous Run

```bash
# Compare current validation with previous results
terraship validate ./terraform \
  --compare previous-report.json \
  --output html \
  --output-file comparison.html
```

Shows:
- What resources improved ✓
- What resources regressed ✗
- Compliance trend
- Rule changes

---

### Generate Multiple Formats

```bash
# Generate all formats at once
terraship validate ./terraform \
  --output html,pdf,json,sarif

# Output files created:
# - terraship-report.html
# - terraship-report.pdf
# - terraship-report.json
# - terraship-report.sarif
```

---

### Include Validation History

```bash
# Include 7-day validation history in report
terraship validate ./terraform \
  --output html \
  --include-history \
  --output-file report-with-history.html
```

Shows:
- 📈 Line chart of passed/failed/warnings over time
- 📊 Trend analysis (improving/regressing)
- 📅 Historical data points

---

## Integration Examples

### GitHub Actions Workflow

```yaml
name: Infrastructure Validation

on:
  pull_request:
    paths:
      - 'terraform/**'
  schedule:
    - cron: '0 9 * * *'  # Daily at 9 AM

jobs:
  validate:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Validate Terraform
        run: |
          terraship validate ./terraform \
            --output html,json,sarif \
            --include-history
      
      - name: Upload SARIF to GitHub
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: terraship-report.sarif
      
      - name: Archive HTML Report
        uses: actions/upload-artifact@v4
        with:
          name: validation-report
          path: terraship-report.html
          
      - name: Post Results to PR
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const report = fs.readFileSync('terraship-report.json', 'utf8');
            const data = JSON.parse(report);
            
            const comment = `## 🚢 Terraship Validation Results
            
            - **Compliance:** ${data.compliance_percent.toFixed(1)}%
            - **Passed:** ${data.passed_resources} ✓
            - **Failed:** ${data.failed_resources} ✗
            - **Warnings:** ${data.warning_resources} ⚠
            
            [📊 View Full Report](artifacts/validation-report.html)`;
            
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });
```

---

### Local Development Setup

```bash
# Install Terraship
go install github.com/vijayaxai/terraship/cmd/terraship@latest

# Generate report
cd my-terraform-project
terraship validate . --output html --html-advanced

# Open in browser
open terraship-report.html

# Watch for changes (optional)
while inotifywait -e modify -r ./terraform; do
  terraship validate . --output html
done
```

---

## Report Contents

### HTML Report Includes

1. **Summary Section**
   - Total resources
   - Passed resources
   - Failed resources
   - Warning resources
   - Compliance percentage

2. **Detailed Results**
   - Resource name
   - Resource type
   - Provider
   - Policy check results
   - Remediation steps

3. **Compare View**
   - Current vs previous run
   - Changes highlighted
   - Trend analysis

4. **Integrations**
   - Slack setup
   - Email configuration
   - Webhook URL
   - GitHub Actions docs
   - Cloud storage options

---

## Troubleshooting

### PDF Generation Fails

**Issue:** `wkhtmltopdf not found`

**Solution:**
```bash
# Install the tool
brew install wkhtmltopdf  # macOS

# Or use HTML export
terraship validate . --output html
# Open in browser and print as PDF
```

---

### Report Too Large

**Issue:** HTML report is very large (many resources)

**Solution:**
```bash
# Use JSON format instead
terraship validate . --output json

# Filter report with jq
cat report.json | jq '.resources[] | select(.status=="failed")'
```

---

### Can't Open HTML Report

**Issue:** Report opens but styling doesn't appear

**Solution:**
- Ensure you're opening the `.html` file directly (not through a server)
- Try a different browser (Chrome, Firefox, Safari)
- Check that JavaScript is enabled

---

## Command Reference

```bash
terraship validate [terraform-dir] [flags]

Flags:
  --output string           Output format: human, json, html, pdf, sarif
                           (comma-separated for multiple formats)
                           Default: human

  --output-file string     Output file path
                           Auto-named if not specified:
                           - terraship-report.html
                           - terraship-report.pdf
                           - terraship-report.json
                           - terraship-report.sarif

  --html-advanced         Enable advanced HTML features
                          (dark mode, charts, search, compare)
                          Default: false

  --include-history       Include 7-day validation history
                          Default: false

  --compare string        Compare with previous results
                          Accepts: JSON report file path

  --policy string         Custom policy file path
                          Default: policies/terraship-policy.yml

  --provider string       Cloud provider: aws, azure, gcp
                          Default: auto-detect

  --mode string          Validation mode: validate-existing, ephemeral-sandbox
                          Default: validate-existing
```

---

## Examples

```bash
# 1. Quick validation
terraship validate ./terraform

# 2. Generate interactive report
terraship validate ./terraform --output html

# 3. Advanced HTML with charts
terraship validate ./terraform \
  --output html \
  --html-advanced \
  --include-history

# 4. Compare runs
terraship validate ./terraform \
  --output html \
  --compare previous-report.json

# 5. All formats for archive
terraship validate ./terraform \
  --output html,pdf,json,sarif

# 6. CI/CD integration
terraship validate ./terraform \
  --output sarif \
  --output-file report.sarif

# 7. Custom policy
terraship validate ./terraform \
  --output html \
  --policy ./policies/strict-policy.yml

# 8. Specific provider
terraship validate ./terraform \
  --output html \
  --provider azure
```

---

## Next Steps

- ✅ Try generating an HTML report
- ✅ Explore the interactive features
- ✅ Set up GitHub Actions integration
- ✅ Configure team integrations (Slack, Email)
- ✅ Archive reports for compliance tracking
