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
