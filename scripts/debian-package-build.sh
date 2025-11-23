#!/bin/bash
set -e

# Configuration
PACKAGE_NAME="filemanager"
VERSION="0.1.2"
MAINTAINER="Moses Muranja <moses.muranja@strathmore.edu>"
DESCRIPTION="A modern file manager with Rust+Go backend"
HOMEPAGE="https://github.com/devonionrouting4Moses/fileManager"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building Debian package for ${PACKAGE_NAME} v${VERSION}${NC}"

# Detect architecture
ARCH=$(dpkg --print-architecture)
echo -e "${YELLOW}Architecture: ${ARCH}${NC}"

# Create package structure
PKG_DIR="${PACKAGE_NAME}_${VERSION}_${ARCH}"
rm -rf "$PKG_DIR"
mkdir -p "$PKG_DIR/DEBIAN"
mkdir -p "$PKG_DIR/usr/local/bin"
mkdir -p "$PKG_DIR/usr/local/lib"
mkdir -p "$PKG_DIR/usr/share/applications"
mkdir -p "$PKG_DIR/usr/share/doc/${PACKAGE_NAME}"
mkdir -p "$PKG_DIR/usr/share/man/man1"

# Copy binary
echo -e "${YELLOW}Copying binary...${NC}"
if [ -f "filemanager" ]; then
    cp filemanager "$PKG_DIR/usr/local/bin/"
    chmod 755 "$PKG_DIR/usr/local/bin/filemanager"
else
    echo -e "${RED}Error: filemanager binary not found${NC}"
    exit 1
fi

# Copy Rust library
echo -e "${YELLOW}Copying Rust library...${NC}"
if [ -f "target/release/libfilemanager.so" ]; then
    cp target/release/libfilemanager.so "$PKG_DIR/usr/local/lib/"
    chmod 644 "$PKG_DIR/usr/local/lib/libfilemanager.so"
fi

# Create control file
cat > "$PKG_DIR/DEBIAN/control" << EOF
Package: ${PACKAGE_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: ${MAINTAINER}
Depends: libc6 (>= 2.31)
Homepage: ${HOMEPAGE}
Description: ${DESCRIPTION}
 FileManager is a cross-platform file manager built with Rust and Go,
 offering high performance and modern features for file management.
 .
 Features:
  - Fast file operations
  - Intuitive user interface
  - Cross-platform compatibility
EOF

# Create postinst script
cat > "$PKG_DIR/DEBIAN/postinst" << 'EOF'
#!/bin/bash
set -e

# Update library cache
ldconfig

echo "FileManager installed successfully!"
echo "Run 'filemanager' to start the application."

exit 0
EOF
chmod 755 "$PKG_DIR/DEBIAN/postinst"

# Create postrm script
cat > "$PKG_DIR/DEBIAN/postrm" << 'EOF'
#!/bin/bash
set -e

# Update library cache
ldconfig

exit 0
EOF
chmod 755 "$PKG_DIR/DEBIAN/postrm"

# Create .desktop file
cat > "$PKG_DIR/usr/share/applications/filemanager.desktop" << EOF
[Desktop Entry]
Version=1.0
Type=Application
Name=FileManager
Comment=Modern file manager
Exec=/usr/local/bin/filemanager
Icon=filemanager
Terminal=false
Categories=System;FileTools;FileManager;
Keywords=file;manager;files;
EOF

# Copy documentation
if [ -f "README.md" ]; then
    cp README.md "$PKG_DIR/usr/share/doc/${PACKAGE_NAME}/"
fi

if [ -f "LICENSE" ]; then
    cp LICENSE "$PKG_DIR/usr/share/doc/${PACKAGE_NAME}/copyright"
fi

# Create changelog
cat > "$PKG_DIR/usr/share/doc/${PACKAGE_NAME}/changelog.Debian" << EOF
${PACKAGE_NAME} (${VERSION}) unstable; urgency=medium

  * New release ${VERSION}

 -- ${MAINTAINER}  $(date -R)
EOF
gzip -9 "$PKG_DIR/usr/share/doc/${PACKAGE_NAME}/changelog.Debian"

# Create man page
cat > "$PKG_DIR/usr/share/man/man1/filemanager.1" << EOF
.TH FILEMANAGER 1 "$(date '+%B %Y')" "FileManager ${VERSION}" "User Commands"
.SH NAME
filemanager \- modern file manager with Rust+Go backend
.SH SYNOPSIS
.B filemanager
[\fIOPTIONS\fR]
.SH DESCRIPTION
FileManager is a cross-platform file manager built with Rust and Go,
offering high performance and modern features for file management.
.SH OPTIONS
.TP
\fB\-\-version\fR
Display version information
.TP
\fB\-\-help\fR
Display help message
.SH AUTHOR
Written by FileManager contributors.
.SH REPORTING BUGS
Report bugs at: ${HOMEPAGE}/issues
EOF
gzip -9 "$PKG_DIR/usr/share/man/man1/filemanager.1"

# Set proper permissions
find "$PKG_DIR" -type d -exec chmod 755 {} \;
find "$PKG_DIR" -type f -exec chmod 644 {} \;
chmod 755 "$PKG_DIR/usr/local/bin/filemanager"
chmod 755 "$PKG_DIR/DEBIAN/postinst"
chmod 755 "$PKG_DIR/DEBIAN/postrm"

# Build the package
echo -e "${YELLOW}Building .deb package...${NC}"
dpkg-deb --build --root-owner-group "$PKG_DIR"

# Verify the package
echo -e "${YELLOW}Verifying package...${NC}"
dpkg-deb --info "${PKG_DIR}.deb"
dpkg-deb --contents "${PKG_DIR}.deb"

# Run lintian if available
if command -v lintian &> /dev/null; then
    echo -e "${YELLOW}Running lintian checks...${NC}"
    lintian "${PKG_DIR}.deb" || true
fi

echo -e "${GREEN}✅ Package built successfully: ${PKG_DIR}.deb${NC}"
echo -e "${YELLOW}To install: sudo dpkg -i ${PKG_DIR}.deb${NC}"
echo -e "${YELLOW}To upload to PPA: dput ppa:your-username/ppa ${PACKAGE_NAME}_${VERSION}_source.changes${NC}"