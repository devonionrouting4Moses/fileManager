# FileManager Installation Script for Windows
# Run with: powershell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "  FileManager v2.0.0 Installation" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Check if running as Administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "❌ Error: This script requires Administrator privileges." -ForegroundColor Red
    Write-Host ""
    Write-Host "Right-click on PowerShell and select 'Run as Administrator', then run:" -ForegroundColor Yellow
    Write-Host "  powershell -ExecutionPolicy Bypass -File install.ps1" -ForegroundColor Yellow
    Write-Host ""
    pause
    exit 1
}

# Installation directory
$installDir = "$env:ProgramFiles\FileManager"
Write-Host "📂 Installation directory: $installDir" -ForegroundColor Green

# Create installation directory
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
    Write-Host "✅ Created installation directory" -ForegroundColor Green
}

# Copy files
Write-Host ""
Write-Host "📋 Copying files..." -ForegroundColor Yellow

Copy-Item "filemanager.exe" -Destination $installDir -Force
Write-Host "  ✓ filemanager.exe" -ForegroundColor Gray

if (Test-Path "fs_operations_core.dll") {
    Copy-Item "fs_operations_core.dll" -Destination $installDir -Force
    Write-Host "  ✓ fs_operations_core.dll" -ForegroundColor Gray
}

if (Test-Path "fsops.exe") {
    Copy-Item "fsops.exe" -Destination $installDir -Force
    Write-Host "  ✓ fsops.exe" -ForegroundColor Gray
}

if (Test-Path "filemanager_frontend") {
    Copy-Item "filemanager_frontend" -Destination $installDir -Recurse -Force
    Write-Host "  ✓ filemanager_frontend/" -ForegroundColor Gray
}

# Add to PATH
Write-Host ""
Write-Host "🔧 Adding to system PATH..." -ForegroundColor Yellow

$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($currentPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$installDir", "Machine")
    Write-Host "✅ Added to PATH" -ForegroundColor Green
} else {
    Write-Host "✓ Already in PATH" -ForegroundColor Gray
}

# Create Start Menu shortcuts
Write-Host ""
Write-Host "🔗 Creating shortcuts..." -ForegroundColor Yellow

$startMenuFolder = "$env:ProgramData\Microsoft\Windows\Start Menu\Programs\FileManager"
if (-not (Test-Path $startMenuFolder)) {
    New-Item -ItemType Directory -Path $startMenuFolder | Out-Null
}

$WshShell = New-Object -ComObject WScript.Shell

# Main shortcut
$shortcut = $WshShell.CreateShortcut("$startMenuFolder\FileManager.lnk")
$shortcut.TargetPath = "$installDir\filemanager.exe"
$shortcut.WorkingDirectory = $installDir
$shortcut.Description = "FileManager - Modern file management tool"
$shortcut.Save()
Write-Host "  ✓ Start Menu: FileManager" -ForegroundColor Gray

# Web mode shortcut
$shortcutWeb = $WshShell.CreateShortcut("$startMenuFolder\FileManager (Web Mode).lnk")
$shortcutWeb.TargetPath = "$installDir\filemanager.exe"
$shortcutWeb.Arguments = "--web"
$shortcutWeb.WorkingDirectory = $installDir
$shortcutWeb.Description = "FileManager - Web Interface Mode"
$shortcutWeb.Save()
Write-Host "  ✓ Start Menu: FileManager (Web Mode)" -ForegroundColor Gray

# Desktop shortcut (optional)
$createDesktop = Read-Host "Create desktop shortcut? (Y/n)"
if ($createDesktop -ne "n" -and $createDesktop -ne "N") {
    $desktopShortcut = $WshShell.CreateShortcut("$env:Public\Desktop\FileManager.lnk")
    $desktopShortcut.TargetPath = "$installDir\filemanager.exe"
    $desktopShortcut.WorkingDirectory = $installDir
    $desktopShortcut.Save()
    Write-Host "  ✓ Desktop shortcut created" -ForegroundColor Gray
}

# Create uninstaller
Write-Host ""
Write-Host "📝 Creating uninstaller..." -ForegroundColor Yellow

$uninstallScript = @"
`$installDir = "$installDir"
Remove-Item -Path "`$installDir" -Recurse -Force
Remove-Item -Path "$startMenuFolder" -Recurse -Force
Remove-Item -Path "`$env:Public\Desktop\FileManager.lnk" -ErrorAction SilentlyContinue
`$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
`$newPath = (`$currentPath.Split(';') | Where-Object { `$_ -ne "`$installDir" }) -join ';'
[Environment]::SetEnvironmentVariable("Path", `$newPath, "Machine")
Write-Host "✅ FileManager uninstalled successfully!" -ForegroundColor Green
"@

$uninstallScript | Out-File -FilePath "$installDir\uninstall.ps1" -Encoding UTF8

# Success message
Write-Host ""
Write-Host "======================================" -ForegroundColor Green
Write-Host "  ✅ Installation Complete!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green
Write-Host ""
Write-Host "FileManager has been installed to:" -ForegroundColor White
Write-Host "  $installDir" -ForegroundColor Cyan
Write-Host ""
Write-Host "To start FileManager:" -ForegroundColor White
Write-Host "  • Open a new terminal and type: filemanager" -ForegroundColor Cyan
Write-Host "  • Or find it in Start Menu" -ForegroundColor Cyan
Write-Host "  • Or double-click the desktop shortcut" -ForegroundColor Cyan
Write-Host ""
Write-Host "For web mode:" -ForegroundColor White
Write-Host "  filemanager --web" -ForegroundColor Cyan
Write-Host ""
Write-Host "To uninstall:" -ForegroundColor White
Write-Host "  powershell -ExecutionPolicy Bypass -File `"$installDir\uninstall.ps1`"" -ForegroundColor Cyan
Write-Host ""
pause