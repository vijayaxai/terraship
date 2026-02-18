# Manual Terraform Provider Download and Installation Script
# Run this if corporate firewall blocks direct Terraform provider downloads

Write-Host "=== Terraform Provider Manual Download Guide ===" -ForegroundColor Cyan
Write-Host ""

# Provider versions (update these as needed)
$awsVersion = "5.31.0"
$azureVersion = "3.85.0"

Write-Host "OPTION 1: Direct Download (If your browser works)" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host ""
Write-Host "AWS Provider v$awsVersion (Windows AMD64):"
Write-Host "https://releases.hashicorp.com/terraform-provider-aws/$awsVersion/terraform-provider-aws_${awsVersion}_windows_amd64.zip" -ForegroundColor Green
Write-Host ""
Write-Host "Azure Provider v$azureVersion (Windows AMD64):"
Write-Host "https://releases.hashicorp.com/terraform-provider-azurerm/$azureVersion/terraform-provider-azurerm_${azureVersion}_windows_amd64.zip" -ForegroundColor Green
Write-Host ""

Write-Host "OPTION 2: Mirror Site Downloads" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host ""
Write-Host "Try Terraform Registry API (alternate method):"
Write-Host "1. AWS: https://registry.terraform.io/v1/providers/hashicorp/aws/$awsVersion/download/windows/amd64" -ForegroundColor Green
Write-Host "2. Azure: https://registry.terraform.io/v1/providers/hashicorp/azurerm/$azureVersion/download/windows/amd64" -ForegroundColor Green
Write-Host ""

Write-Host "OPTION 3: Use Terraform Mirror" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "Configure Terraform to use a mirror that might work better with your network:"
Write-Host ""
$cliConfig = @"
provider_installation {
  network_mirror {
    url = "https://terraform-mirror.yandex.net/"
  }
}
"@
Write-Host "Create file: $env:APPDATA\terraform.rc with content:" -ForegroundColor Cyan
Write-Host $cliConfig -ForegroundColor Gray
Write-Host ""

Write-Host "OPTION 4: Download with PowerShell (Trying now...)" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host ""

# Try downloading with PowerShell with different settings
$downloadDir = "$PSScriptRoot\terraform-providers"
New-Item -ItemType Directory -Force -Path $downloadDir | Out-Null

Write-Host "Attempting to download AWS provider..." -ForegroundColor Cyan

try {
    # Method 1: Try with System.Net.WebClient
    $webClient = New-Object System.Net.WebClient
    $webClient.Headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
    $awsUrl = "https://releases.hashicorp.com/terraform-provider-aws/$awsVersion/terraform-provider-aws_${awsVersion}_windows_amd64.zip"
    $awsZip = "$downloadDir\aws-provider.zip"
    
    Write-Host "Downloading from: $awsUrl" -ForegroundColor Gray
    $webClient.DownloadFile($awsUrl, $awsZip)
    
    if (Test-Path $awsZip) {
        Write-Host "✅ AWS provider downloaded successfully to: $awsZip" -ForegroundColor Green
        
        # Extract and install
        $extractPath = "$downloadDir\aws-extracted"
        Expand-Archive -Path $awsZip -DestinationPath $extractPath -Force
        
        $pluginDir = "$env:APPDATA\terraform.d\plugins\registry.terraform.io\hashicorp\aws\$awsVersion\windows_amd64"
        New-Item -ItemType Directory -Force -Path $pluginDir | Out-Null
        
        Get-ChildItem -Path $extractPath -File | Copy-Item -Destination $pluginDir -Force
        Write-Host "✅ AWS provider installed to: $pluginDir" -ForegroundColor Green
    }
} catch {
    Write-Host "❌ WebClient download failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    Write-Host "MANUAL STEPS REQUIRED:" -ForegroundColor Yellow
    Write-Host "1. Open your browser (Edge, Chrome, Firefox)" -ForegroundColor White
    Write-Host "2. Copy this URL and paste in browser:" -ForegroundColor White
    Write-Host "   https://releases.hashicorp.com/terraform-provider-aws/$awsVersion/terraform-provider-aws_${awsVersion}_windows_amd64.zip" -ForegroundColor Cyan
    Write-Host "3. Download the ZIP file" -ForegroundColor White
    Write-Host "4. Save it to: $downloadDir" -ForegroundColor White
    Write-Host "5. Run this script again, or extract manually to:" -ForegroundColor White
    Write-Host "   $env:APPDATA\terraform.d\plugins\registry.terraform.io\hashicorp\aws\$awsVersion\windows_amd64" -ForegroundColor Cyan
    Write-Host ""
}

Write-Host ""
Write-Host "OPTION 5: Use Terraform Plugin Cache" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "If you have another machine with internet access:" -ForegroundColor White
Write-Host "1. On that machine, run: terraform init in any Terraform project" -ForegroundColor White
Write-Host "2. Copy the plugins from: %APPDATA%\terraform.d\plugins\" -ForegroundColor White
Write-Host "3. Paste to the same location on this machine" -ForegroundColor White
Write-Host ""

Write-Host "OPTION 6: Skip Integration Tests" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "Your Terraship project is already working! You can:" -ForegroundColor White
Write-Host "1. Run unit tests only: go test -short ./..." -ForegroundColor Cyan
Write-Host "2. Use the CLI: .\bin\terraship.exe --help" -ForegroundColor Cyan
Write-Host "3. Test validation: .\bin\terraship.exe validate . --policy policies\sample-policy.yml" -ForegroundColor Cyan
Write-Host ""
Write-Host "Integration tests are optional and require Terraform providers." -ForegroundColor Gray
Write-Host ""

Read-Host "Press Enter to exit"
