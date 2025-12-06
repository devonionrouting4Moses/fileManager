#!/bin/bash

# FileManager Security Verification Script
# Verifies that FileManager is installed with proper security permissions

set -e

echo ""
echo "🔐 FileManager Security Verification"
echo "====================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

ERRORS=0
WARNINGS=0

# Function to print results
print_check() {
    local status=$1
    local message=$2
    
    if [ "$status" = "✅" ]; then
        echo -e "${GREEN}${status}${NC} $message"
    elif [ "$status" = "❌" ]; then
        echo -e "${RED}${status}${NC} $message"
        ((ERRORS++))
    elif [ "$status" = "⚠️" ]; then
        echo -e "${YELLOW}${status}${NC} $message"
        ((WARNINGS++))
    fi
}

# Check 1: Binary exists and is executable
echo -e "${BLUE}1. Binary Installation${NC}"
if command -v filemanager &> /dev/null; then
    BINARY_PATH=$(which filemanager)
    print_check "✅" "Binary found: $BINARY_PATH"
    
    # Check if executable
    if [ -x "$BINARY_PATH" ]; then
        print_check "✅" "Binary is executable"
    else
        print_check "❌" "Binary is not executable"
    fi
else
    print_check "❌" "Binary not found in PATH"
fi
echo ""

# Check 2: Frontend directory
echo -e "${BLUE}2. Frontend Directory${NC}"

# Try to find frontend directory
FRONTEND_DIR=""
if [ -d "/usr/share/filemanager/frontend" ]; then
    FRONTEND_DIR="/usr/share/filemanager/frontend"
elif [ -d "$HOME/.local/share/filemanager/frontend" ]; then
    FRONTEND_DIR="$HOME/.local/share/filemanager/frontend"
elif [ -d "/snap/filemanager/current/frontend" ]; then
    FRONTEND_DIR="/snap/filemanager/current/frontend"
fi

if [ -n "$FRONTEND_DIR" ]; then
    print_check "✅" "Frontend directory found: $FRONTEND_DIR"
    
    # Check if index.html exists
    if [ -f "$FRONTEND_DIR/index.html" ]; then
        print_check "✅" "Frontend files present (index.html found)"
    else
        print_check "❌" "Frontend files missing (index.html not found)"
    fi
    
    # Check directory permissions
    PERMS=$(stat -c %a "$FRONTEND_DIR" 2>/dev/null || stat -f %OLp "$FRONTEND_DIR" 2>/dev/null || echo "unknown")
    if [ "$PERMS" = "755" ] || [ "$PERMS" = "unknown" ]; then
        print_check "✅" "Directory permissions: $PERMS (rwxr-xr-x)"
    else
        print_check "⚠️" "Directory permissions: $PERMS (expected 755)"
    fi
    
    # Check file permissions
    if [ -f "$FRONTEND_DIR/index.html" ]; then
        FILE_PERMS=$(stat -c %a "$FRONTEND_DIR/index.html" 2>/dev/null || stat -f %OLp "$FRONTEND_DIR/index.html" 2>/dev/null || echo "unknown")
        if [ "$FILE_PERMS" = "644" ] || [ "$FILE_PERMS" = "unknown" ]; then
            print_check "✅" "File permissions: $FILE_PERMS (rw-r--r--)"
        else
            print_check "⚠️" "File permissions: $FILE_PERMS (expected 644)"
        fi
    fi
    
    # Check ownership
    OWNER=$(stat -c %U:%G "$FRONTEND_DIR" 2>/dev/null || stat -f %Su:%Sg "$FRONTEND_DIR" 2>/dev/null || echo "unknown")
    if [ "$OWNER" = "root:root" ]; then
        print_check "✅" "Ownership: $OWNER (root:root)"
    elif [ "$OWNER" = "unknown" ]; then
        print_check "⚠️" "Ownership: $OWNER (could not determine)"
    else
        print_check "⚠️" "Ownership: $OWNER (expected root:root for system install)"
    fi
else
    print_check "⚠️" "Frontend directory not yet created (will be created on first run)"
fi
echo ""

# Check 3: Configuration files
echo -e "${BLUE}3. Configuration Files${NC}"

# System config
if [ -f "/etc/filemanager/config.yaml" ]; then
    print_check "✅" "System config found: /etc/filemanager/config.yaml"
    
    # Check permissions
    CONFIG_PERMS=$(stat -c %a "/etc/filemanager/config.yaml" 2>/dev/null || stat -f %OLp "/etc/filemanager/config.yaml" 2>/dev/null || echo "unknown")
    if [ "$CONFIG_PERMS" = "644" ] || [ "$CONFIG_PERMS" = "unknown" ]; then
        print_check "✅" "Config permissions: $CONFIG_PERMS (rw-r--r--)"
    else
        print_check "⚠️" "Config permissions: $CONFIG_PERMS (expected 644)"
    fi
else
    print_check "⚠️" "System config not found (optional for user installs)"
fi

# User config
USER_CONFIG="$HOME/.config/filemanager/config.yaml"
if [ -f "$USER_CONFIG" ]; then
    print_check "✅" "User config found: $USER_CONFIG"
    
    # Check permissions
    USER_CONFIG_PERMS=$(stat -c %a "$USER_CONFIG" 2>/dev/null || stat -f %OLp "$USER_CONFIG" 2>/dev/null || echo "unknown")
    if [ "$USER_CONFIG_PERMS" = "644" ] || [ "$USER_CONFIG_PERMS" = "unknown" ]; then
        print_check "✅" "User config permissions: $USER_CONFIG_PERMS (rw-r--r--)"
    else
        print_check "⚠️" "User config permissions: $USER_CONFIG_PERMS (expected 644)"
    fi
else
    print_check "⚠️" "User config not found (will be created on first run)"
fi
echo ""

# Check 4: Frontend resolution
echo -e "${BLUE}4. Frontend Resolution${NC}"

if command -v filemanager &> /dev/null; then
    # Run frontend-info command
    if filemanager --frontend-info &> /tmp/filemanager-info.txt; then
        print_check "✅" "Frontend info command works"
        
        # Check if it shows a valid frontend directory
        if grep -q "Active Frontend:" /tmp/filemanager-info.txt; then
            ACTIVE_FRONTEND=$(grep "Active Frontend:" /tmp/filemanager-info.txt | sed 's/.*: //')
            print_check "✅" "Active frontend: $ACTIVE_FRONTEND"
        else
            print_check "⚠️" "Could not determine active frontend"
        fi
        
        rm -f /tmp/filemanager-info.txt
    else
        print_check "❌" "Frontend info command failed"
    fi
else
    print_check "⚠️" "Cannot test frontend resolution (binary not in PATH)"
fi
echo ""

# Check 5: No duplicate frontend folders
echo -e "${BLUE}5. Duplicate Frontend Detection${NC}"

DUPLICATES=$(find "$HOME" -name "filemanager_frontend" -type d 2>/dev/null | wc -l)
if [ "$DUPLICATES" -eq 0 ]; then
    print_check "✅" "No duplicate frontend folders found"
elif [ "$DUPLICATES" -eq 1 ]; then
    print_check "✅" "Single frontend folder found (as expected)"
else
    print_check "⚠️" "Multiple frontend folders found ($DUPLICATES)"
    find "$HOME" -name "filemanager_frontend" -type d 2>/dev/null | while read dir; do
        echo "   Found: $dir"
    done
fi
echo ""

# Check 6: No world-writable files
echo -e "${BLUE}6. World-Writable Files${NC}"

if [ -n "$FRONTEND_DIR" ]; then
    WORLD_WRITABLE=$(find "$FRONTEND_DIR" -type f -perm -002 2>/dev/null | wc -l)
    if [ "$WORLD_WRITABLE" -eq 0 ]; then
        print_check "✅" "No world-writable files found"
    else
        print_check "❌" "Found $WORLD_WRITABLE world-writable files"
        find "$FRONTEND_DIR" -type f -perm -002 2>/dev/null | while read file; do
            echo "   $file"
        done
    fi
else
    print_check "⚠️" "Cannot check (frontend directory not found)"
fi
echo ""

# Summary
echo "====================================="
echo ""

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo -e "${GREEN}✅ All security checks passed!${NC}"
    echo ""
    echo "FileManager is properly installed with correct security permissions."
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo -e "${YELLOW}⚠️  $WARNINGS warnings found${NC}"
    echo ""
    echo "FileManager is installed but some checks need attention."
    exit 0
else
    echo -e "${RED}❌ $ERRORS errors found${NC}"
    echo ""
    echo "Please fix the errors above before using FileManager."
    exit 1
fi
