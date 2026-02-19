#!/bin/bash
# Terraship VS Code Extension Publishing Script
# This script publishes the extension to the VS Code Marketplace

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 Terraship VS Code Extension Publishing${NC}"
echo "=================================================="
echo ""

# Check if vsce is installed
if ! command -v vsce &> /dev/null; then
    echo -e "${YELLOW}⚠️  vsce not found, installing...${NC}"
    npm install -g @vscode/vsce
fi

# Navigate to extension directory
cd vscode-extension

# Check if .vsix already exists
if [ -f "terraship-vscode-0.2.0.vsix" ]; then
    echo -e "${GREEN}✓ Package found: terraship-vscode-0.2.0.vsix${NC}"
else
    echo -e "${YELLOW}📦 Building extension...${NC}"
    npm run compile
    vsce package
fi

echo ""
echo -e "${YELLOW}📋 Before publishing, ensure:${NC}"
echo "   ✓ You have a Microsoft account"
echo "   ✓ You created a publisher on https://marketplace.visualstudio.com/manage"
echo "   ✓ Publisher ID: 'terraship'"
echo "   ✓ You have a Personal Access Token (PAT)"
echo ""

echo -e "${YELLOW}🔐 To get a Personal Access Token:${NC}"
echo "   1. Go to https://dev.azure.com"
echo "   2. Click your profile → User settings → Personal access tokens"
echo "   3. Click 'New Token'"
echo "   4. Settings:"
echo "      - Name: Terraship Extension Publishing"
echo "      - Scopes: Marketplace → Manage"
echo "      - Expiration: 90 days"
echo "   5. Copy the token (you'll only see it once!)"
echo ""

read -p "Continue with publishing? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}Logging in to marketplace...${NC}"
    vsce login terraship
    
    echo -e "${GREEN}Publishing extension...${NC}"
    vsce publish
    
    echo ""
    echo -e "${GREEN}✅ Extension published successfully!${NC}"
    echo -e "${GREEN}📱 Available at: https://marketplace.visualstudio.com/items?itemName=terraship.terraship-vscode${NC}"
else
    echo -e "${YELLOW}Publishing cancelled.${NC}"
fi
