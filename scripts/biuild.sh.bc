#!/bin/bash

# FileManager Cross-Platform Build Script
# Builds for Linux, macOS, and Windows

VERSION="0.1.0"
APP_NAME="filemanager"
BUILD_DIR="dist"

echo "🏗️  Building FileManager v${VERSION} for multiple platforms..."
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Clean previous builds
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# Build Rust library first
echo "${BLUE}📦 Building Rust library...${NC}"
cd ../rust_ffi
cargo build --release -p fs-operations-core
if [ $? -ne 0 ]; then
    echo "❌ Rust build failed"
    exit 1
fi
cd ../scripts
echo "${GREEN}✅ Rust library built${NC}"
echo ""

# Function to build for a platform
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT_NAME="${APP_NAME}"
    local PLATFORM_DIR="${BUILD_DIR}/${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${APP_NAME}.exe"
    fi
    
    echo "${BLUE}🔨 Building for ${GOOS}/${GOARCH}...${NC}"
    
    mkdir -p ${PLATFORM_DIR}
    
    # Build Go binary with Windows-specific settings
    cd ../file_manager
	if [ "$GOOS" = "windows" ]; then
		echo "Building for Windows - using Windows-specific implementation"
		# First build the Windows executable
		env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ \
		CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
		go build -ldflags="-s -w" -o ../scripts/${PLATFORM_DIR}/${OUTPUT_NAME} ./cmd/app
		
		# If that fails, try the more complex build with excluded files
		if [ $? -ne 0 ]; then
			echo "Standard build failed, trying with file exclusions..."
			find . -name "*.go" -not -name "*_windows.go" -not -name "*_test.go" -not -name "*_unix.go" -not -name "helper_windows.go" -not -name "main.go" -not -name "version.go" -not -name "webserver.go" -not -name "templates.go" > build_files.txt
			echo "cmd/app/main.go" >> build_files.txt
			echo "pkg/version/version.go" >> build_files.txt
			echo "internal/handler/webserver.go" >> build_files.txt
			
			env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ \
			CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
			go build -ldflags="-s -w" -o ../scripts/${PLATFORM_DIR}/${OUTPUT_NAME} $(cat build_files.txt)
			rm -f build_files.txt
		fi
	else
		env GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=1 \
		CGO_LDFLAGS="-L../rust_ffi/target/release -lfs_operations_core -ldl -lpthread -lm" \
		go build -ldflags="-s -w" -o ../scripts/${PLATFORM_DIR}/${OUTPUT_NAME} ./cmd/app
	fi
    cd ../scripts
    
    if [ $? -eq 0 ]; then
        # Copy Rust library
        if [ "$GOOS" = "linux" ]; then
            cp ../rust_ffi/target/release/libfs_operations_core.so ${PLATFORM_DIR}/
        elif [ "$GOOS" = "darwin" ]; then
            cp ../rust_ffi/target/release/libfs_operations_core.dylib ${PLATFORM_DIR}/ 2>/dev/null || cp ../rust_ffi/target/release/libfs_operations_core.so ${PLATFORM_DIR}/libfs_operations_core.dylib
        elif [ "$GOOS" = "windows" ]; then
            cp ../rust_ffi/target/release/fs_operations_core.dll ${PLATFORM_DIR}/ 2>/dev/null || echo "Note: Windows build requires manual DLL"
        fi
        
        # Copy documentation
        cp README.md ${PLATFORM_DIR}/ 2>/dev/null
        cp USAGE_GUIDE.md ${PLATFORM_DIR}/ 2>/dev/null
        cp ADVANCED_FEATURES.md ${PLATFORM_DIR}/ 2>/dev/null
        
        # Create install script
        if [ "$GOOS" = "windows" ]; then
            create_windows_installer ${PLATFORM_DIR}
        else
            create_unix_installer ${PLATFORM_DIR} ${GOOS}
        fi
        
        # Create archive
        cd ${BUILD_DIR}
        if [ "$GOOS" = "windows" ]; then
            zip -r ${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}.zip ${GOOS}-${GOARCH}/ > /dev/null
        else
            tar -czf ${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}.tar.gz ${GOOS}-${GOARCH}/
        fi
        cd - > /dev/null
        
        echo "${GREEN}✅ Built for ${GOOS}/${GOARCH}${NC}"
    else
        echo "❌ Failed to build for ${GOOS}/${GOARCH}"
    fi
    echo ""
}

create_unix_installer() {
    local DIR=$1
    local OS=$2
    
    cat > ${DIR}/install.sh << 'EOF'
#!/bin/bash

# FileManager Installation Script

set -e

APP_NAME="filemanager"
INSTALL_DIR="/usr/local/bin"
LIB_DIR="/usr/local/lib"

echo "📦 Installing FileManager..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  This script requires sudo privileges"
    echo "Please run: sudo ./install.sh"
    exit 1
fi

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Copy binary
echo "📋 Copying binary to ${INSTALL_DIR}..."
cp ${APP_NAME} ${INSTALL_DIR}/
chmod +x ${INSTALL_DIR}/${APP_NAME}

# Copy library
echo "📚 Copying library to ${LIB_DIR}..."
if [ "$OS" = "darwin" ]; then
    cp libfilemanager.dylib ${LIB_DIR}/ 2>/dev/null || cp libfilemanager.so ${LIB_DIR}/libfilemanager.dylib
else
    cp libfilemanager.so ${LIB_DIR}/
fi

# Update library cache (Linux only)
if [ "$OS" = "linux" ]; then
    ldconfig
fi

echo ""
echo "✅ FileManager installed successfully!"
echo ""
echo "Usage:"
echo "  filemanager              Start interactive mode"
echo "  filemanager --version    Show version"
echo "  filemanager --update     Check for updates"
echo "  filemanager --help       Show help"
echo ""
EOF
    chmod +x ${DIR}/install.sh
    
    # Create uninstall script
    cat > ${DIR}/uninstall.sh << 'EOF'
#!/bin/bash

# FileManager Uninstallation Script

set -e

APP_NAME="filemanager"
INSTALL_DIR="/usr/local/bin"
LIB_DIR="/usr/local/lib"

echo "🗑️  Uninstalling FileManager..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  This script requires sudo privileges"
    echo "Please run: sudo ./uninstall.sh"
    exit 1
fi

# Remove binary
if [ -f "${INSTALL_DIR}/${APP_NAME}" ]; then
    rm -f ${INSTALL_DIR}/${APP_NAME}
    echo "✅ Removed binary"
fi

# Remove library
rm -f ${LIB_DIR}/libfilemanager.so ${LIB_DIR}/libfilemanager.dylib
echo "✅ Removed library"

# Update library cache (Linux only)
if [ "$(uname -s | tr '[:upper:]' '[:lower:]')" = "linux" ]; then
    ldconfig
fi

echo ""
echo "✅ FileManager uninstalled successfully!"
echo ""
EOF
    chmod +x ${DIR}/uninstall.sh
}

create_windows_installer() {
    local DIR=$1
    
    cat > ${DIR}/install.bat << 'EOF'
@echo off
REM FileManager Installation Script for Windows

echo Installing FileManager...

REM Check for admin privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo This script requires administrator privileges
    echo Please run as administrator
    pause
    exit /b 1
)

REM Create installation directory
set INSTALL_DIR=%ProgramFiles%\FileManager
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copy files
copy /Y filemanager.exe "%INSTALL_DIR%\"
copy /Y filemanager.dll "%INSTALL_DIR%\" 2>nul

REM Add to PATH
setx /M PATH "%PATH%;%INSTALL_DIR%"

echo.
echo FileManager installed successfully!
echo.
echo Usage:
echo   filemanager              Start interactive mode
echo   filemanager --version    Show version
echo   filemanager --update     Check for updates
echo   filemanager --help       Show help
echo.
pause
EOF

    cat > ${DIR}/uninstall.bat << 'EOF'
@echo off
REM FileManager Uninstallation Script for Windows

echo Uninstalling FileManager...

REM Check for admin privileges
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo This script requires administrator privileges
    echo Please run as administrator
    pause
    exit /b 1
)

set INSTALL_DIR=%ProgramFiles%\FileManager

REM Remove files
if exist "%INSTALL_DIR%" (
    rmdir /S /Q "%INSTALL_DIR%"
    echo FileManager removed
)

echo.
echo FileManager uninstalled successfully!
echo.
pause
EOF
}

# Build for all platforms
echo "🚀 Building binaries..."
echo ""

# Linux builds
build_platform "linux" "amd64"
build_platform "linux" "arm64"

# macOS builds
build_platform "darwin" "amd64"
build_platform "darwin" "arm64"

# Windows builds
build_platform "windows" "amd64"

echo ""
echo "✨ Build complete!"
echo ""
echo "📦 Generated packages:"
ls -lh ${BUILD_DIR}/*.{tar.gz,zip} 2>/dev/null | awk '{print "   " $9 " (" $5 ")"}'
echo ""
echo "📂 Distribution packages are in: ${BUILD_DIR}/"
echo ""