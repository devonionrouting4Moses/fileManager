# FileManager Windows Security Permissions Script
# Run as Administrator to set proper ACLs on FileManager installation
# 
# Usage: powershell -ExecutionPolicy Bypass -File Set-FileManagerPermissions.ps1
# Or:    powershell -ExecutionPolicy Bypass -File Set-FileManagerPermissions.ps1 -InstallPath "C:\Custom\Path\FileManager"

param(
    [string]$InstallPath = "$env:ProgramFiles\FileManager"
)

# Verify admin rights
if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "❌ ERROR: This script must be run as Administrator" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please run PowerShell as Administrator and try again:"
    Write-Host "  1. Right-click PowerShell"
    Write-Host "  2. Select 'Run as Administrator'"
    Write-Host "  3. Run: powershell -ExecutionPolicy Bypass -File Set-FileManagerPermissions.ps1"
    exit 1
}

Write-Host ""
Write-Host "🔒 FileManager Windows Security Permissions" -ForegroundColor Green
Write-Host "=============================================" -ForegroundColor Green
Write-Host ""

# Verify installation path exists
if (-not (Test-Path $InstallPath)) {
    Write-Host "❌ ERROR: Installation path not found: $InstallPath" -ForegroundColor Red
    exit 1
}

Write-Host "📂 Installation Path: $InstallPath" -ForegroundColor Cyan
Write-Host ""

# Get the frontend directory
$FrontendDir = Join-Path $InstallPath "frontend"

if (-not (Test-Path $FrontendDir)) {
    Write-Host "⚠️  WARNING: Frontend directory not found: $FrontendDir" -ForegroundColor Yellow
    Write-Host "Skipping frontend permissions..."
    Write-Host ""
} else {
    Write-Host "🔐 Setting permissions on frontend directory..." -ForegroundColor Yellow
    Write-Host ""
    
    try {
        # Get current ACL
        $Acl = Get-Acl $FrontendDir
        
        # Remove inheritance and set explicit permissions
        Write-Host "  • Removing inherited permissions..." -ForegroundColor Gray
        $Acl.SetAccessRuleProtection($true, $false)
        
        # Clear existing rules
        $Acl.Access | ForEach-Object {
            $Acl.RemoveAccessRule($_) | Out-Null
        }
        
        # Grant Administrators full control
        Write-Host "  • Granting Administrators full control..." -ForegroundColor Gray
        $AdminRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
            "BUILTIN\Administrators",
            "FullControl",
            "ContainerInherit,ObjectInherit",
            "None",
            "Allow"
        )
        $Acl.AddAccessRule($AdminRule)
        
        # Grant Users read and execute only (no write)
        Write-Host "  • Granting Users read & execute (no write)..." -ForegroundColor Gray
        $UserRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
            "BUILTIN\Users",
            "ReadAndExecute",
            "ContainerInherit,ObjectInherit",
            "None",
            "Allow"
        )
        $Acl.AddAccessRule($UserRule)
        
        # Apply permissions
        Write-Host "  • Applying permissions..." -ForegroundColor Gray
        Set-Acl $FrontendDir $Acl
        
        Write-Host ""
        Write-Host "✅ Frontend permissions set successfully" -ForegroundColor Green
        Write-Host ""
        Write-Host "📋 Permission Summary:" -ForegroundColor Cyan
        Write-Host "  • Administrators: Full Control (Read, Write, Execute)" -ForegroundColor Gray
        Write-Host "  • Users: Read & Execute (no write)" -ForegroundColor Gray
        Write-Host ""
        
    } catch {
        Write-Host "❌ ERROR: Failed to set permissions" -ForegroundColor Red
        Write-Host "Details: $_" -ForegroundColor Red
        exit 1
    }
}

# Set binary permissions
$BinaryPath = Join-Path $InstallPath "filemanager.exe"

if (Test-Path $BinaryPath) {
    Write-Host "🔐 Setting permissions on binary..." -ForegroundColor Yellow
    Write-Host ""
    
    try {
        $Acl = Get-Acl $BinaryPath
        $Acl.SetAccessRuleProtection($true, $false)
        
        # Clear existing rules
        $Acl.Access | ForEach-Object {
            $Acl.RemoveAccessRule($_) | Out-Null
        }
        
        # Grant Administrators full control
        $AdminRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
            "BUILTIN\Administrators",
            "FullControl",
            "None",
            "None",
            "Allow"
        )
        $Acl.AddAccessRule($AdminRule)
        
        # Grant Users read and execute
        $UserRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
            "BUILTIN\Users",
            "ReadAndExecute",
            "None",
            "None",
            "Allow"
        )
        $Acl.AddAccessRule($UserRule)
        
        Set-Acl $BinaryPath $Acl
        
        Write-Host "✅ Binary permissions set successfully" -ForegroundColor Green
        Write-Host ""
        
    } catch {
        Write-Host "⚠️  WARNING: Could not set binary permissions" -ForegroundColor Yellow
        Write-Host "Details: $_" -ForegroundColor Gray
        Write-Host ""
    }
}

# Set config directory permissions
$ConfigDir = Join-Path $InstallPath "config.yaml"

if (Test-Path $ConfigDir) {
    Write-Host "🔐 Setting permissions on config file..." -ForegroundColor Yellow
    Write-Host ""
    
    try {
        $Acl = Get-Acl $ConfigDir
        $Acl.SetAccessRuleProtection($true, $false)
        
        # Clear existing rules
        $Acl.Access | ForEach-Object {
            $Acl.RemoveAccessRule($_) | Out-Null
        }
        
        # Grant Administrators full control
        $AdminRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
            "BUILTIN\Administrators",
            "FullControl",
            "None",
            "None",
            "Allow"
        )
        $Acl.AddAccessRule($AdminRule)
        
        # Grant Users read only
        $UserRule = New-Object System.Security.AccessControl.FileSystemAccessRule(
            "BUILTIN\Users",
            "ReadAndExecute",
            "None",
            "None",
            "Allow"
        )
        $Acl.AddAccessRule($UserRule)
        
        Set-Acl $ConfigDir $Acl
        
        Write-Host "✅ Config file permissions set successfully" -ForegroundColor Green
        Write-Host ""
        
    } catch {
        Write-Host "⚠️  WARNING: Could not set config permissions" -ForegroundColor Yellow
        Write-Host "Details: $_" -ForegroundColor Gray
        Write-Host ""
    }
}

# Verify permissions
Write-Host "🔍 Verifying permissions..." -ForegroundColor Yellow
Write-Host ""

if (Test-Path $FrontendDir) {
    $Acl = Get-Acl $FrontendDir
    Write-Host "📂 Frontend Directory: $FrontendDir" -ForegroundColor Cyan
    Write-Host "   Owner: $($Acl.Owner)" -ForegroundColor Gray
    Write-Host "   Access Rules:" -ForegroundColor Gray
    $Acl.Access | ForEach-Object {
        Write-Host "     • $($_.IdentityReference): $($_.FileSystemRights)" -ForegroundColor Gray
    }
    Write-Host ""
}

Write-Host "=============================================" -ForegroundColor Green
Write-Host "✅ Security configuration complete!" -ForegroundColor Green
Write-Host ""
Write-Host "🚀 FileManager is now properly secured:" -ForegroundColor Cyan
Write-Host "   • Administrators can manage files" -ForegroundColor Gray
Write-Host "   • Regular users can run the application" -ForegroundColor Gray
Write-Host "   • Frontend files are read-only for users" -ForegroundColor Gray
Write-Host ""
Write-Host "💡 Next steps:" -ForegroundColor Cyan
Write-Host "   1. Run: filemanager --frontend-info" -ForegroundColor Gray
Write-Host "   2. Test: filemanager --web" -ForegroundColor Gray
Write-Host ""
