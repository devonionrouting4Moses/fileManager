; FileManager Windows Installer
; NSIS Installer Script
; Version: 2.0.2

!include "MUI2.nsh"
!include "x64.nsh"

; Basic Settings
Name "FileManager v2.0.2"
OutFile "FileManager-2.0.2-Setup.exe"
VIProductVersion "2.0.2.0"
VIAddVersionKey "ProductName" "FileManager"
VIAddVersionKey "ProductVersion" "2.0.2"
VIAddVersionKey "FileVersion" "2.0.2"
VIAddVersionKey "FileDescription" "FileManager - File manager with web interface"
VIAddVersionKey "LegalCopyright" "2025"
InstallDir "$PROGRAMFILES\FileManager"
InstallDirRegKey HKCU "Software\FileManager" ""

; Request admin privileges
RequestExecutionLevel admin

; MUI Settings
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

; Installer Sections
Section "Install"
    SetOutPath "$INSTDIR"
    
    ; Install binary
    File "filemanager.exe"
    
    ; Install frontend (bundled)
    SetOutPath "$INSTDIR\frontend"
    File /r "frontend\*.*"
    
    ; Set frontend directory permissions
    ExecWait 'powershell -ExecutionPolicy Bypass -File "$INSTDIR\Set-Permissions.ps1"'
    
    ; Create config directory
    CreateDirectory "$INSTDIR\config"
    
    ; Create config file
    FileOpen $0 "$INSTDIR\config.yaml" w
    FileWrite $0 "# FileManager Configuration$\r$\n"
    FileWrite $0 "# Frontend directory path$\r$\n"
    FileWrite $0 "frontend_dir: $INSTDIR\frontend$\r$\n"
    FileClose $0
    
    ; Create uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"
    
    ; Create Start Menu shortcuts
    CreateDirectory "$SMPROGRAMS\FileManager"
    CreateShortcut "$SMPROGRAMS\FileManager\FileManager.lnk" "$INSTDIR\filemanager.exe" "--web"
    CreateShortcut "$SMPROGRAMS\FileManager\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
    
    ; Create Desktop shortcut
    CreateShortcut "$DESKTOP\FileManager.lnk" "$INSTDIR\filemanager.exe" "--web"
    
    ; Store installation folder
    WriteRegStr HKCU "Software\FileManager" "" $INSTDIR
    
    ; Write uninstall information
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\FileManager" "DisplayName" "FileManager"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\FileManager" "UninstallString" "$INSTDIR\Uninstall.exe"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\FileManager" "InstallLocation" "$INSTDIR"
    WriteRegStr HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\FileManager" "DisplayVersion" "2.0.2"
    
    MessageBox MB_OK "FileManager installed successfully!$\n$\nFrontend location: $INSTDIR\frontend"
SectionEnd

; Uninstaller Section
Section "Uninstall"
    ; Remove files
    RMDir /r "$INSTDIR\frontend"
    Delete "$INSTDIR\filemanager.exe"
    Delete "$INSTDIR\config.yaml"
    Delete "$INSTDIR\Set-Permissions.ps1"
    Delete "$INSTDIR\Uninstall.exe"
    
    ; Remove directories
    RMDir "$INSTDIR\config"
    RMDir "$INSTDIR"
    
    ; Remove Start Menu shortcuts
    RMDir /r "$SMPROGRAMS\FileManager"
    
    ; Remove Desktop shortcut
    Delete "$DESKTOP\FileManager.lnk"
    
    ; Remove registry entries
    DeleteRegKey HKCU "Software\FileManager"
    DeleteRegKey HKCU "Software\Microsoft\Windows\CurrentVersion\Uninstall\FileManager"
    
    MessageBox MB_OK "FileManager uninstalled successfully!"
SectionEnd
