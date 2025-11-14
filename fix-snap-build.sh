#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Fixing Snap Build Configuration      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if we're in the right directory
if [ ! -f "snap/snapcraft.yaml" ]; then
    echo -e "${RED}Error: snap/snapcraft.yaml not found${NC}"
    echo -e "${YELLOW}Please run this script from the project root directory${NC}"
    exit 1
fi

# Create snap/local directory
echo -e "${YELLOW}Creating snap/local directory...${NC}"
mkdir -p snap/local

# Create desktop file
echo -e "${YELLOW}Creating desktop file...${NC}"
cat > snap/local/filemanager.desktop << 'EOF'
[Desktop Entry]
Version=1.0
Type=Application
Name=FileManager
GenericName=File Manager
Comment=Modern file manager with Rust+Go backend
Icon=${SNAP}/usr/share/icons/hicolor/256x256/apps/filemanager.png
Exec=filemanager
Terminal=false
Categories=System;FileTools;FileManager;Utility;
Keywords=file;manager;files;browser;explorer;
StartupNotify=true
StartupWMClass=FileManager
EOF
echo -e "${GREEN}✓ Desktop file created${NC}"

# Check if icon exists
if [ ! -f "snap/local/filemanager.png" ]; then
    echo -e "${YELLOW}Creating placeholder icon...${NC}"
    
    # Check if ImageMagick is available
    if command -v convert &> /dev/null; then
        convert -size 256x256 xc:"#4A90E2" \
                -gravity center \
                -pointsize 48 \
                -fill white \
                -annotate +0+0 "FM" \
                snap/local/filemanager.png
        echo -e "${GREEN}✓ Placeholder icon created${NC}"
    else
        echo -e "${YELLOW}⚠ ImageMagick not found. Creating text placeholder...${NC}"
        echo -e "${YELLOW}Please replace snap/local/filemanager.png with a real 256x256 PNG icon${NC}"
        # Create empty file as placeholder
        touch snap/local/filemanager.png
    fi
fi

# Fix snapcraft.yaml
echo -e "${YELLOW}Updating snapcraft.yaml...${NC}"

# Backup original
cp snap/snapcraft.yaml snap/snapcraft.yaml.backup

# Update the desktop-file part with source-type
if grep -q "source-type: local" snap/snapcraft.yaml; then
    echo -e "${GREEN}✓ snapcraft.yaml already has source-type: local${NC}"
else
    # Use sed to add source-type after the source line in desktop-file part
    sed -i '/desktop-file:/,/source: snap\/local\// {
        /source: snap\/local\//a\    source-type: local
    }' snap/snapcraft.yaml
    echo -e "${GREEN}✓ Added source-type: local to snapcraft.yaml${NC}"
fi

# Verify the fix
echo -e "${YELLOW}Verifying configuration...${NC}"

if grep -A 3 "desktop-file:" snap/snapcraft.yaml | grep -q "source-type: local"; then
    echo -e "${GREEN}✓ Configuration verified${NC}"
else
    echo -e "${RED}✗ Configuration verification failed${NC}"
    echo -e "${YELLOW}Manually edit snap/snapcraft.yaml and add 'source-type: local' under the desktop-file part${NC}"
fi

echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Fix Complete!                         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}Next steps:${NC}"
echo "1. Review snap/snapcraft.yaml to ensure it's correct"
echo "2. Replace snap/local/filemanager.png with your actual icon (256x256 PNG)"
echo "3. Build the snap:"
echo -e "   ${YELLOW}snapcraft${NC}"
echo ""
echo "4. Test locally:"
echo -e "   ${YELLOW}sudo snap install --dangerous filemanager_0.1.2_amd64.snap${NC}"
echo ""
echo -e "${BLUE}Backup created at: snap/snapcraft.yaml.backup${NC}"