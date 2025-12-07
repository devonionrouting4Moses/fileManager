@echo off
REM FileManager Installation Script

echo Installing FileManager...
set INSTALL_DIR=%ProgramFiles%\FileManager

if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copy executable
copy /Y filemanager.exe "%INSTALL_DIR%\"

REM Copy Rust library
copy /Y fs_operations_core.dll "%INSTALL_DIR%\" 2>nul

REM Add to PATH
setx /M PATH "%PATH%;%INSTALL_DIR%"

echo.
echo ✅ Installation complete!
echo FileManager installed to: %INSTALL_DIR%
echo.
echo Run 'filemanager' from command prompt to start
pause
