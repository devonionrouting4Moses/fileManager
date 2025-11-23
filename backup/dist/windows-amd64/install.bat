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
