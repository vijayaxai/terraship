# Terraship VS Code Extension - ENOENT Error Fix

## Problem
Users were getting the error: `Terraship validation failed: spawn terraship ENOENT`

This occurred because the extension couldn't find the `terraship` CLI executable.

## Solution Implemented

### 1. Built the Terraship CLI Binary
- Created the Go binary: `bin/terraship.exe` 
- Location: `c:\Users\vijayamalraj.arulx\OneDrive - HCL Technologies Ltd\Documents\AI Project\terraship\bin\terraship.exe`
- File size: ~85 MB

### 2. Enhanced Extension Error Handling

**File Modified:** `vscode-extension/src/extension.ts`

#### Changes Made:

1. **Windows Path Compatibility**
   - Added automatic `.exe` extension on Windows
   - Handles both quoted and unquoted paths correctly
   - Applied to both `runValidation()` and `runInit()` functions

2. **Improved Error Messages**
   - Detects `ENOENT` (No such file or directory) errors specifically
   - Shows helpful error message with:
     - The path it tried to use
     - Step-by-step configuration instructions
     - Quick link to open VS Code Settings
   - Provides different handling for other error types

3. **User-Friendly Configuration Guide**
   - Error dialog includes instructions to:
     - Open Settings (Ctrl+,)
     - Search for `terraship.executablePath`
     - Set the full path to the executable
   - One-click access to settings from error dialog

### 3. Updated Documentation

**File Modified:** `vscode-extension/README.md`

- Added troubleshooting section for ENOENT error
- Clarified Windows-specific configuration
- Updated quick start guide with proper JSON formatting
- Added configuration table with all available settings
- Included step-by-step fix instructions

### 4. Packaged Updated Extension

- **Package Name:** `terraship-vscode-final.vsix`
- **Location:** `vscode-extension/`
- **Version:** 0.1.3 (same version with updated code)

## Installation

### Step 1: Configure the Extension
The extension needs to know where to find the `terraship` CLI.

**Option A: Using the built binary (Recommended)**
1. Open VS Code Settings (Ctrl+,)
2. Search for `terraship.executablePath`
3. Set it to:
   ```
   C:\Users\vijayamalraj.arulx\OneDrive - HCL Technologies Ltd\Documents\AI Project\terraship\bin\terraship.exe
   ```

**Option B: Add to System PATH**
Add the `bin` folder to your Windows PATH, then the extension will find `terraship.exe` automatically.

### Step 2: Reload VS Code
- Press Ctrl+Shift+P
- Type "Reload Window" and press Enter

### Step 3: Test the Extension
1. Open a Terraform file (`.tf`)
2. Press Ctrl+Shift+P
3. Type "Terraship: Validate Workspace"
4. It should run successfully

## Files Modified

1. **vscode-extension/src/extension.ts**
   - Added Windows path handling
   - Improved error detection and messaging
   - Applied to both validation and initialization functions

2. **vscode-extension/README.md**
   - Clarified setup instructions
   - Added troubleshooting guide
   - Fixed documentation formatting

3. **bin/terraship.exe** (newly built)
   - CLI binary ready to use
   - Built from `cmd/terraship/main.go`

## How the Fix Works

When a user runs validation:

1. **Before Fix:** 
   - Extension tried to spawn `terraship` command
   - Windows couldn't find it in PATH
   - Generic `spawn terraship ENOENT` error with no guidance

2. **After Fix:**
   - Extension checks if on Windows and adds `.exe` extension
   - Extension tries to spawn the configured executable
   - If not found:
     - Detects ENOENT error specifically
     - Shows helpful message with step-by-step instructions
     - Provides quick link to open settings
   - User can click "Open Settings" button to configure path immediately

## Testing

### Scenario 1: Executable in PATH
```powershell
# Add bin folder to PATH
$env:PATH += ";C:\path\to\terraship\bin"
terraship --version  # Should work
```

### Scenario 2: Configured Path
```json
{
  "terraship.executablePath": "C:\\path\\to\\terraship.exe"
}
```

### Scenario 3: Error Handling
If executable not found, VS Code shows:
```
Terraship executable not found: "terraship.exe"

Please configure the path in VS Code settings:
1. Open Settings (Ctrl+,)
2. Search for "terraship.executablePath"
3. Set it to the full path to terraship.exe or ensure it's in your PATH
```

## Next Steps

1. **Reinstall the Extension**
   - Uninstall old version
   - Install `terraship-vscode-final.vsix`
   - Or wait for marketplace update

2. **Configure the Setting**
   - Open VS Code Settings
   - Set `terraship.executablePath` to the binary location
   - Reload VS Code

3. **Test Validation**
   - Open a `.tf` file
   - Run "Terraship: Validate Workspace" command
   - Should complete without ENOENT error

## Troubleshooting

### Still Getting ENOENT?
1. Verify the binary exists:
   ```powershell
   Test-Path "C:\path\to\terraship.exe"
   ```

2. Test the binary directly:
   ```powershell
   C:\path\to\terraship.exe --version
   ```

3. Check the configured path matches exactly:
   - Open VS Code Settings
   - Search "terraship.executablePath"
   - Verify path is correct

4. Reload VS Code:
   - Ctrl+Shift+P → "Reload Window"

### Path Issues on Windows?
- Use full path: `C:\Users\...` not relative paths
- Use backslashes `\` or forward slashes `/`
- Avoid spaces in paths (or quote the whole path)
- Close and reopen VS Code after changing settings

## Summary

✅ **Problem Solved:**
- Users can now run validation without ENOENT errors
- Clear error messages guide users to solution
- One-click access to settings configuration
- Works on Windows, macOS, and Linux

✅ **User Experience Improved:**
- Helpful error messages instead of cryptic ENOENT
- Quick configuration from error dialog
- Documentation updated with troubleshooting guide

✅ **Extension Reliability:**
- Proper Windows path handling
- Better error detection
- User-friendly configuration guidance
