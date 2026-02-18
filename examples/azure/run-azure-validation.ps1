# Load Azure Credentials and Run Terraship Validation
# Usage: .\run-azure-validation.ps1

Write-Host "🔐 Loading Azure credentials..." -ForegroundColor Cyan

# Check if credentials file exists
$credsFile = "..\..\azure-credentials.env"
if (-not (Test-Path $credsFile)) {
    Write-Host "❌ Error: azure-credentials.env not found!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please create $credsFile from the template:" -ForegroundColor Yellow
    Write-Host "  1. Copy azure-credentials.env.template to azure-credentials.env" -ForegroundColor White
    Write-Host "  2. Fill in your Azure credentials" -ForegroundColor White
    Write-Host "  3. Run this script again" -ForegroundColor White
    Write-Host ""
    Write-Host "Or use Azure CLI:" -ForegroundColor Yellow
    Write-Host "  az login" -ForegroundColor Green
    Write-Host "  az account set --subscription 'your-subscription-id'" -ForegroundColor Green
    exit 1
}

# Load credentials from file
Write-Host "📄 Reading credentials from: $credsFile" -ForegroundColor Gray
Get-Content $credsFile | ForEach-Object {
    if ($_ -match '^\s*([^#][^=]+)=(.+)$') {
        $name = $matches[1].Trim()
        $value = $matches[2].Trim().Trim('"')
        
        # Skip if value is placeholder
        if ($value -notmatch 'your-.*-here' -and $value -ne 'your-subscription-id' -and $value -ne 'your-tenant-id') {
            Set-Item -Path "env:$name" -Value $value
            Write-Host "  ✓ Set $name" -ForegroundColor Green
        }
    }
}

# Verify credentials are set
$credentialsSet = $true
@('ARM_SUBSCRIPTION_ID', 'ARM_TENANT_ID') | ForEach-Object {
    if (-not (Test-Path "env:$_") -or (Get-Item "env:$_").Value -eq '') {
        Write-Host "  ❌ Missing: $_" -ForegroundColor Red
        $credentialsSet = $false
    }
}

if (-not $credentialsSet) {
    Write-Host ""
    Write-Host "❌ Credentials incomplete. Please update azure-credentials.env" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "✅ Credentials loaded successfully!" -ForegroundColor Green
Write-Host ""

# Display what will be validated
Write-Host "📊 Validation Configuration:" -ForegroundColor Cyan
Write-Host "  Directory: $(Get-Location)" -ForegroundColor Gray
Write-Host "  Policy: ..\..\policies\sample-policy.yml" -ForegroundColor Gray
Write-Host "  Provider: Azure" -ForegroundColor Gray
Write-Host "  Subscription: $env:ARM_SUBSCRIPTION_ID" -ForegroundColor Gray
Write-Host ""

# Run Terraship validation
Write-Host "🚢 Running Terraship validation..." -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor DarkGray
Write-Host ""

# Run validation and display to console
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --provider azure --verbose

$exitCode = $LASTEXITCODE

Write-Host ""
Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor DarkGray

# Generate reports automatically
Write-Host ""
Write-Host "📝 Generating validation reports..." -ForegroundColor Cyan

Write-Host "  → validation-report.txt (human-readable)" -ForegroundColor Gray
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --provider azure --output-file validation-report.txt 2>$null

Write-Host "  → validation-report.json (for CI/CD)" -ForegroundColor Gray
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --provider azure --output json --output-file validation-report.json 2>$null

Write-Host "  → validation-report.sarif (for GitHub Security)" -ForegroundColor Gray
..\..\bin\terraship.exe validate . --policy ..\..\policies\sample-policy.yml --provider azure --output sarif --output-file validation-report.sarif 2>$null

Write-Host ""
if (Test-Path "validation-report.txt") {
    $reports = Get-ChildItem validation-report.* | Select-Object Name, @{Name="Size";Expression={"{0:N0} KB" -f ($_.Length / 1KB)}}, LastWriteTime
    Write-Host "✅ Reports saved:" -ForegroundColor Green
    $reports | Format-Table -AutoSize | Out-String | Write-Host
} else {
    Write-Host "⚠️  Warning: Reports may not have been generated" -ForegroundColor Yellow
}

Write-Host "═══════════════════════════════════════════════════════════════" -ForegroundColor DarkGray

if ($exitCode -eq 0) {
    Write-Host "✅ Validation completed successfully!" -ForegroundColor Green
} else {
    Write-Host "❌ Validation failed - Check reports for details" -ForegroundColor Red
}

exit $exitCode
