#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   FileManager Packaging Setup         ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if running on Linux
if [[ "$OSTYPE" != "linux-gnu"* ]]; then
    echo -e "${RED}This script must be run on Linux${NC}"
    exit 1
fi

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install DEB packaging tools
install_deb_tools() {
    echo -e "${YELLOW}Installing DEB packaging tools...${NC}"
    sudo apt update
    sudo apt install -y \
        devscripts \
        debhelper \
        dh-make \
        build-essential \
        lintian \
        ubuntu-dev-tools \
        dput-ng
    echo -e "${GREEN}✓ DEB tools installed${NC}"
}

# Install Snap tools
install_snap_tools() {
    echo -e "${YELLOW}Installing Snap tools...${NC}"
    if ! command_exists snapcraft; then
        sudo snap install snapcraft --classic
    fi
    if ! command_exists multipass; then
        sudo snap install multipass --classic
    fi
    echo -e "${GREEN}✓ Snap tools installed${NC}"
}

# Setup GPG key for signing
setup_gpg() {
    echo -e "${YELLOW}Checking GPG setup...${NC}"
    if ! gpg --list-keys | grep -q "@"; then
        echo -e "${BLUE}No GPG key found. Creating one...${NC}"
        echo -e "${YELLOW}Please provide the following information:${NC}"
        read -p "Your name: " name
        read -p "Your email: " email
        
        gpg --batch --generate-key <<EOF
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: $name
Name-Email: $email
Expire-Date: 0
%no-protection
%commit
EOF
        
        KEY_ID=$(gpg --list-keys --keyid-format LONG "$email" | grep pub | awk '{print $2}' | cut -d'/' -f2)
        echo -e "${GREEN}✓ GPG key created: $KEY_ID${NC}"
        echo -e "${YELLOW}Upload to keyserver: gpg --send-keys $KEY_ID --keyserver keyserver.ubuntu.com${NC}"
    else
        echo -e "${GREEN}✓ GPG key already exists${NC}"
    fi
}

# Create snap directory structure
setup_snap_structure() {
    echo -e "${YELLOW}Setting up Snap directory structure...${NC}"
    mkdir -p snap/local
    
    # Create desktop file if it doesn't exist
    if [ ! -f "snap/local/filemanager.desktop" ]; then
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
        echo -e "${GREEN}✓ Created desktop file${NC}"
    fi
    
    # Remind about icon
    if [ ! -f "snap/local/filemanager.png" ]; then
        echo -e "${YELLOW}⚠ Please add a 256x256 PNG icon at snap/local/filemanager.png${NC}"
    fi
    
    # Create snapcraft.yaml if it doesn't exist
    if [ ! -f "snap/snapcraft.yaml" ]; then
        echo -e "${BLUE}Creating snapcraft.yaml...${NC}"
        cp snapcraft.yaml snap/snapcraft.yaml 2>/dev/null || true
        echo -e "${GREEN}✓ Snap structure created${NC}"
    fi
}

# Create GitHub workflow directory
setup_github_workflows() {
    echo -e "${YELLOW}Setting up GitHub workflows...${NC}"
    mkdir -p .github/workflows
    
    if [ ! -f ".github/workflows/publish-packages.yml" ]; then
        echo -e "${YELLOW}⚠ Add publish-packages.yml to .github/workflows/${NC}"
    else
        echo -e "${GREEN}✓ GitHub workflows configured${NC}"
    fi
}

# Main menu
show_menu() {
    echo ""
    echo -e "${BLUE}What would you like to setup?${NC}"
    echo "1) DEB packaging tools"
    echo "2) Snap packaging tools"
    echo "3) Both (DEB + Snap)"
    echo "4) GPG signing key"
    echo "5) Complete setup (All of the above)"
    echo "6) Exit"
    echo ""
    read -p "Enter choice [1-6]: " choice
    
    case $choice in
        1)
            install_deb_tools
            ;;
        2)
            install_snap_tools
            setup_snap_structure
            ;;
        3)
            install_deb_tools
            install_snap_tools
            setup_snap_structure
            ;;
        4)
            setup_gpg
            ;;
        5)
            install_deb_tools
            install_snap_tools
            setup_gpg
            setup_snap_structure
            setup_github_workflows
            ;;
        6)
            echo -e "${GREEN}Goodbye!${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}Invalid choice${NC}"
            show_menu
            ;;
    esac
}

# Run setup
show_menu

echo ""
echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║   Setup Complete!                      ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Build DEB package: ./debian-package-build.sh"
echo "2. Build Snap package: snapcraft"
echo "3. Test locally before publishing"
echo "4. Configure GitHub secrets for automatic publishing"
echo ""
echo -e "${BLUE}For detailed instructions, see PACKAGING.md${NC}"