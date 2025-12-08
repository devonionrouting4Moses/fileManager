cask "filemanager" do
  version "2.0.2"
  
  if Hardware::CPU.intel?
    sha256 "c17f21ac63d941a79776cbf98ab352e0a442f71bf2073605733f1ed5fb4bc295"
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v#{version}/FileManager-#{version}-macos-amd64.dmg"
  else
    sha256 "bf2a7ae6b3e5c994db5852a0d3faf3056ac034471cee2f95743b0334788ba02c"
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v#{version}/FileManager-#{version}-macos-arm64.dmg"
  end

  name "FileManager"
  desc "Modern file manager with Rust+Go backend and web interface"
  homepage "https://github.com/devonionrouting4Moses/fileManager"

  livecheck do
    url :url
    strategy :github_latest
  end

  app "FileManager.app"

  uninstall quit: "com.devchigarlicmoses.filemanager"

  zap trash: [
    "~/Library/Application Support/FileManager",
    "~/Library/Preferences/com.devchigarlicmoses.filemanager.plist",
    "~/Library/Caches/FileManager",
    "~/Library/Logs/FileManager",
    "~/Library/Saved Application State/com.devchigarlicmoses.filemanager.savedState",
  ]

  caveats <<~EOS
    FileManager has been installed!
    
    To launch the app:
      - Open from Applications folder, or
      - Run: open -a FileManager
    
    The web interface will be available at:
      http://localhost:8080
    
    Note: On first launch, you may need to:
      1. Right-click the app and select "Open"
      2. Grant Full Disk Access in System Preferences → Privacy & Security
  EOS
end
