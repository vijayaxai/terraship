# Quick Validation Script - Automatically Generates Reports
# Usage: .\validate.ps1

$ErrorActionPreference = "Continue"

Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "              TERRASHIP QUICK VALIDATION                        " -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

# Set Azure credentials
$env:AZURE_SUBSCRIPTION_ID = "d30ec219-d601-414b-98b6-230b6e520d37"
$env:AZURE_TENANT_ID = "2111de49-6a33-4187-af6d-96575525e6ef"

# Check if Azure CLI is logged in
Write-Host "🔐 Checking Azure authentication..." -ForegroundColor Yellow
$azAccount = az account show 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "⚠️  Not logged into Azure - Run: az login" -ForegroundColor Red
    Write-Host ""
    exit 1
}
Write-Host "✅ Azure authentication OK" -ForegroundColor Green
Write-Host ""

# Run validation
Write-Host "🚢 Running Terraship validation..." -ForegroundColor Yellow
Write-Host "─────────────────────────────────────────────────────────────── " -ForegroundColor DarkGray
Write-Host ""

..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --provider azure

$exitCode = $LASTEXITCODE

Write-Host ""
Write-Host "─────────────────────────────────────────────────────────────── " -ForegroundColor DarkGray

# Generate reports
Write-Host ""
Write-Host "📊 Generating reports..." -ForegroundColor Yellow

$timestamp = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
$reportDir = ".\reports"

# Create reports directory if it doesn't exist
if (-not (Test-Path $reportDir)) {
    New-Item -ItemType Directory -Path $reportDir -Force | Out-Null
}

Write-Host "  → $reportDir\validation-report-$timestamp.txt" -ForegroundColor Gray
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --provider azure --output-file "$reportDir\validation-report-$timestamp.txt" 2>$null

Write-Host "  → $reportDir\validation-report-$timestamp.json" -ForegroundColor Gray
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --provider azure --output json --output-file "$reportDir\validation-report-$timestamp.json" 2>$null

Write-Host "  → $reportDir\validation-report-$timestamp.sarif" -ForegroundColor Gray
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --provider azure --output sarif --output-file "$reportDir\validation-report-$timestamp.sarif" 2>$null

# Also create "latest" versions
Copy-Item "$reportDir\validation-report-$timestamp.txt" "$reportDir\validation-report-latest.txt" -Force
Copy-Item "$reportDir\validation-report-$timestamp.json" "$reportDir\validation-report-latest.json" -Force
Copy-Item "$reportDir\validation-report-$timestamp.sarif" "$reportDir\validation-report-latest.sarif" -Force

Write-Host ""
Write-Host "✅ Reports saved to: $reportDir" -ForegroundColor Green

# Show report summary
$reports = Get-ChildItem "$reportDir\validation-report-$timestamp.*" | 
    Select-Object Name, @{Name="Size";Expression={"{0:N1} KB" -f ($_.Length / 1KB)}}, LastWriteTime

$reports | Format-Table -AutoSize | Out-String | ForEach-Object { Write-Host $_ -ForegroundColor Gray }

Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor Cyan

if ($exitCode -eq 0) {
    Write-Host "✅ VALIDATION PASSED" -ForegroundColor Green
} else {
    Write-Host "❌ VALIDATION FAILED - Check reports for details" -ForegroundColor Red
    Write-Host "   View latest report: code $reportDir\validation-report-latest.txt" -ForegroundColor Yellow
}

Write-Host ""
exit $exitCode
