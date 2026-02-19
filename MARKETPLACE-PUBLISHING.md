# 🚀 Publishing Terraship v0.2.0 to VS Code Marketplace

## Package Ready ✅

**File:** `terraship-vscode-0.2.0.vsix`  
**Size:** 1.37 MB  
**Version:** 0.2.0  
**Build Date:** 2026-02-19

---

## Step 1: Create Publisher Account (First Time Only)

If you haven't created a publisher yet:

1. Go to **https://marketplace.visualstudio.com/manage**
2. Sign in with your **Microsoft account** (create one if needed)
3. Click **"Create Publisher"**
4. Fill in:
   - **Publisher ID:** `terraship` (must match `package.json`)
   - **Display Name:** `Terraship`
   - **Email:** your-email@domain.com
5. Click **Create**

✅ **Done!** You now have a publisher account.

---

## Step 2: Generate Personal Access Token (PAT)

You need a PAT to publish programmatically:

1. Go to **https://dev.azure.com**
2. Click your **profile icon** (top-right)
3. Select **User settings**
4. Click **Personal access tokens**
5. Click **New Token**
6. Fill in:
   - **Name:** `Terraship Extension Publishing`
   - **Organization:** Select "All accessible organizations"
   - **Expiration:** 90 days (or custom)
   - **Scopes:** 
     - Check: ✅ **Marketplace** → **Manage**
7. Click **Create**
8. **Copy the token immediately** (you'll only see it once!)
9. Save it somewhere safe (password manager recommended)

✅ **Done!** You have your PAT.

---

## Step 3: Login to VS Code Marketplace via CLI

Use the PAT to authenticate:

```powershell
# Navigate to extension folder
cd "c:\Users\vijayamalraj.arulx\OneDrive - HCL Technologies Ltd\Documents\AI Project\terraship\vscode-extension"

# Login with publisher ID
vsce login terraship

# When prompted, paste your Personal Access Token
```

Expected output:
```
The Personal Access Token (stored in C:\Users\...\vsce) will expire on 2026-06-20.
```

✅ **Done!** You're authenticated.

---

## Step 4: Publish Extension

Now publish the extension to the marketplace:

```powershell
# From vscode-extension directory
vsce publish

# Or specify version explicitly
vsce publish 0.2.0
```

Expected output:
```
 INFO  Publishing terraship/terraship-vscode@0.2.0...
 INFO  Extension URL: https://marketplace.visualstudio.com/items?itemName=terraship.terraship-vscode
 INFO  Hub URL: https://marketplace.visualstudio.com/publishers/terraship/extensions/terraship-vscode
 DONE  Published!
```

✅ **Done!** Extension is now published!

---

## Step 5: Verify Publication

### Check Marketplace

Visit the VS Code Marketplace:
```
https://marketplace.visualstudio.com/items?itemName=terraship.terraship-vscode
```

You should see:
- ✅ Extension name: "Terraship"
- ✅ Version: 0.2.0
- ✅ Publisher: terraship
- ✅ Download button
- ✅ Full README and changelog

### Test Installation

Anyone can now install via:

**In VS Code:**
1. Open Extensions (Ctrl+Shift+X)
2. Search "Terraship"
3. Click Install

**Via CLI:**
```bash
code --install-extension terraship.terraship-vscode
```

---

## 🔄 Updating the Extension (Future)

When you make updates and want to publish a new version:

```powershell
# 1. Update code in src/
# 2. Update version in package.json: "version": "0.2.1"
# 3. Recompile
npm run compile

# 4. Package
vsce package

# 5. Publish (you're already authenticated)
vsce publish
```

Marketplace will automatically update for all users within a few minutes.

---

## 📋 Troubleshooting

### "Authentication failed"
- Verify your PAT hasn't expired
- Make sure you used correct publisher ID
- Re-login: `vsce logout` then `vsce login terraship`

### "File too large"
- The `icon.png` (1.38 MB) is causing warnings but is acceptable
- To reduce: compress the icon or exclude it
- Edit `.vscodeignore` if needed

### "Version already published"
- Marketplace won't allow re-publishing the same version
- Bump version in `package.json`: `0.2.0` → `0.2.1`
- Run `vsce publish` again

### "Publisher not found"
- Verify publisher ID matches `package.json`
- Create publisher at: https://marketplace.visualstudio.com/manage

---

## 📊 After Publishing

### Monitor Usage
- View downloads at: https://marketplace.visualstudio.com/publishers/terraship
- Check ratings and reviews
- Monitor GitHub issues for user feedback

### Next Steps

1. **Announce Release**
   - Tweet about it
   - Post on GitHub Discussions
   - Update README with "Install from Marketplace" link

2. **Gather Feedback**
   - Monitor GitHub issues
   - Respond to marketplace reviews
   - Iterate on features

3. **Plan Next Version**
   - Slack webhook integration
   - Email notifications
   - Advanced reporting UI

---

## 🎯 Release Checklist

Before publishing, verify:

- [x] Extension compiles without errors
- [x] All TypeScript compiled to JavaScript
- [x] Package includes all files
- [x] Version updated to 0.2.0
- [x] CHANGELOG.md has release notes
- [x] README.md contains documentation
- [x] Icon (icon.png) is present
- [x] LICENSE file is included
- [x] Personal Access Token is valid
- [x] Publisher ID is correct

---

## 📚 Helpful Links

- **VS Code Marketplace:** https://marketplace.visualstudio.com
- **Publisher Dashboard:** https://marketplace.visualstudio.com/manage
- **PAT Management:** https://dev.azure.com (User settings → PAT)
- **vsce Documentation:** https://github.com/microsoft/vscode-vsce
- **Extension Manifest:** https://code.visualstudio.com/api/references/extension-manifest
- **Repository:** https://github.com/vijayaxai/terraship

---

## 🎉 Success!

Once published, your extension will be:
- ✅ Searchable on VS Code Marketplace
- ✅ Installable from the Extensions panel
- ✅ Available for download
- ✅ Eligible for ratings and reviews
- ✅ Automatically updated for users

**Users can install with:**
```
VS Code → Extensions → Search "Terraship" → Install
```

---

**Questions?** Check the [PUBLISHING.md](PUBLISHING.md) file or create a GitHub issue.
