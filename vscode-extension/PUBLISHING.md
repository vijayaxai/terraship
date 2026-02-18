# 📦 Publishing Terraship VS Code Extension (Beta)

## Prerequisites

### 1. Create a Visual Studio Marketplace Account

1. Go to https://marketplace.visualstudio.com/manage
2. Sign in with your Microsoft account
3. Click **"Create Publisher"**
4. Fill in details:
   - **Publisher ID:** `terraship` (must match package.json)
   - **Display Name:** Terraship
   - **Email:** your-email@domain.com

### 2. Get a Personal Access Token (PAT)

1. Go to https://dev.azure.com
2. Click your profile → **User settings** → **Personal access tokens**
3. Click **New Token**
4. Settings:
   - **Name:** `Terraship Extension Publishing`
   - **Organization:** All accessible organizations
   - **Expiration:** 90 days (or custom)
   - **Scopes:** 
     - ✅ **Marketplace** → **Manage** (full access)
5. Copy the token (you'll only see it once!)

---

## 🚀 Publishing Steps

### Step 1: Install Required Tools

```powershell
cd vscode-extension

# Install Node.js dependencies
npm install

# Install VS Code Extension Manager (vsce)
npm install -g @vscode/vsce
```

### Step 2: Update Extension for Beta

```powershell
# Edit package.json - already shows version 0.1.0 (perfect for beta!)
# Version format: 0.x.y = beta/preview
```

Your current version `0.1.0` indicates:
- **0** = Major (pre-release/beta)
- **1** = Minor
- **0** = Patch

### Step 3: Build the Extension

```powershell
# Compile TypeScript
npm run compile

# Package extension into .vsix file
vsce package
```

This creates: `terraship-vscode-0.1.0.vsix`

### Step 4: Login to Marketplace

```powershell
# Login with your Personal Access Token
vsce login terraship

# Enter your PAT when prompted
```

### Step 5: Publish to Marketplace

```powershell
# Publish as beta/preview
vsce publish
```

**Or publish specific version:**
```powershell
vsce publish 0.1.0
```

---

## 📝 Before Publishing Checklist

### Required Files

- [x] `package.json` - Extension manifest ✅
- [ ] `README.md` - User documentation
- [ ] `CHANGELOG.md` - Version history
- [ ] `LICENSE` - License file
- [ ] `icon.png` - Extension icon (128x128)

### Update README.md

Create `vscode-extension/README.md`:

```markdown
# 🚢 Terraship VS Code Extension (Beta)

Multi-cloud Terraform validation and policy checking for AWS, Azure, and GCP.

## Features

- ✅ Real-time policy validation
- ✅ Multi-cloud support (AWS, Azure, GCP)
- ✅ Inline error reporting
- ✅ Quick fix suggestions

## Installation

1. Install from VS Code Marketplace
2. Configure policy path in settings
3. Open any `.tf` file
4. Run: `Terraship: Validate Workspace`

## Beta Notice

⚠️ This is a BETA release. Features may change.

Report issues: https://github.com/vijayaxai/terraship/issues
```

### Add Icon

Place `icon.png` (128x128) in `vscode-extension/`:

```powershell
# You can generate one or use a placeholder
# For now, skip it (optional for beta)
```

---

## 🎯 Publishing Options

### Option 1: Public Beta (Marketplace)

**Pros:**
- Anyone can install
- Searchable in VS Code
- Automatic updates

**Cons:**
- Public visibility
- Need Microsoft account

**Command:**
```powershell
vsce publish
```

### Option 2: Private Beta (.vsix Distribution)

**Pros:**
- Full control over who gets it
- No marketplace account needed
- Test with select users

**Cons:**
- Manual distribution
- No automatic updates

**Command:**
```powershell
# Package only (don't publish)
vsce package

# Share terraship-vscode-0.1.0.vsix with testers
# They install: code --install-extension terraship-vscode-0.1.0.vsix
```

### Option 3: Internal Marketplace (Enterprise)

**Pros:**
- Company-only distribution
- Private and secure

**Cons:**
- Requires enterprise setup

---

## 📋 Step-by-Step Publishing Commands

```powershell
# 1. Navigate to extension folder
cd "c:\Users\vijayamalraj.arulx\OneDrive - HCL Technologies Ltd\Documents\AI Project\terraship\vscode-extension"

# 2. Install dependencies
npm install

# 3. Compile TypeScript
npm run compile

# 4. Test locally first
code --install-extension terraship-vscode-0.1.0.vsix

# 5. Package for distribution
vsce package

# 6. (Optional) Publish to marketplace
# First time: vsce login terraship
vsce publish
```

---

## 🧪 Testing Beta Before Publishing

### Local Testing

```powershell
# Package extension
vsce package

# Install in VS Code
code --install-extension terraship-vscode-0.1.0.vsix

# Test with your Azure examples
# Open: examples/azure/main.tf
# Run: Ctrl+Shift+P → "Terraship: Validate Workspace"
```

### Share with Team (Private Beta)

```powershell
# 1. Package extension
vsce package

# 2. Share file with team
# Send: terraship-vscode-0.1.0.vsix

# 3. Team installs:
code --install-extension terraship-vscode-0.1.0.vsix

# 4. Configure in VS Code settings:
{
  "terraship.policyPath": "./policies/sample-policy.yml",
  "terraship.cloudProvider": "azure"
}
```

---

## 🔄 Updating Beta Version

```powershell
# Bump version
npm version patch  # 0.1.0 → 0.1.1
# or
npm version minor  # 0.1.0 → 0.2.0

# Rebuild and publish
npm run compile
vsce publish
```

---

## 📊 Version Strategy

| Version | Stage | Description |
|---------|-------|-------------|
| **0.1.0** | Alpha | Initial internal testing |
| **0.2.0** | Beta | Limited user testing |
| **0.9.0** | RC | Release candidate |
| **1.0.0** | GA | General availability |

Current: **0.1.0 Beta** ✅

---

## 🐛 Beta Testing Checklist

- [ ] Extension loads in VS Code
- [ ] Commands appear in command palette
- [ ] Settings are configurable
- [ ] Validation runs successfully
- [ ] Results display correctly
- [ ] Works with Azure/AWS/GCP
- [ ] No error messages in console

---

## 🚀 Quick Start (Private Beta)

**Recommended for first release:**

```powershell
# 1. Build extension
cd vscode-extension
npm install
npm run compile
vsce package

# 2. Test yourself
code --install-extension terraship-vscode-0.1.0.vsix

# 3. Share with 2-3 team members
# Send them: terraship-vscode-0.1.0.vsix

# 4. Gather feedback
# Fix issues

# 5. Then publish to marketplace
vsce login terraship
vsce publish
```

---

## 📞 Support

- **Issues:** https://github.com/vijayaxai/terraship/issues
- **Docs:** See README.md
- **Beta Feedback:** Create GitHub issue with "beta" label

---

## Next Steps

1. **Now:** Package extension (`vsce package`)
2. **Test:** Install locally and validate
3. **Beta:** Share with 2-3 users
4. **Iterate:** Fix bugs based on feedback
5. **Publish:** Release to marketplace when stable
