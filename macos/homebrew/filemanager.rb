class Filemanager < Formula
  desc "Modern file manager with Rust+Go backend and web interface"
  homepage "https://github.com/devonionrouting4Moses/fileManager"
  version "2.0.2"
  license "MIT"

  if Hardware::CPU.intel?
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v2.0.2/filemanager-2.0.2-macos-amd64.tar.gz"
    sha256 "763af349eaee1a9dba16f75047dbc1d3bb178fa76d8e67f90785a59088b228cd"
  elsif Hardware::CPU.arm?
    url "https://github.com/devonionrouting4Moses/fileManager/releases/download/v2.0.2/filemanager-2.0.2-macos-arm64.tar.gz"
    sha256 "e6d36f40af273c47d29fd274a1688f6f1b8acd9fc03bbd3b3624f4e13d762f86"
  end

  def install
    # Install binary
    bin.install "filemanager"
    
    # Install Rust dynamic library
    lib.install "libfs_operations_core.dylib"
    
    # Install frontend if present
    if File.directory?("frontend")
      (share/"filemanager/frontend").install Dir["frontend/*"]
    end
    
    # Create config directory
    (etc/"filemanager").mkpath
    
    # Create default system config
    (etc/"filemanager/config.yaml").write <<~EOS
      # FileManager System Configuration
      # Frontend directory path
      frontend_dir: #{share}/filemanager/frontend
      
      # Server configuration
      host: 127.0.0.1
      port: 8080
      
      # Default settings
      theme: auto
      show_hidden: false
    EOS
  end

  def post_install
    # Set frontend directory permissions
    frontend_dir = "#{share}/filemanager/frontend"
    
    if File.directory?(frontend_dir)
      # Make directories readable
      system "chmod", "-R", "755", frontend_dir
      
      # Make files readable
      system "find", frontend_dir, "-type", "f", "-exec", "chmod", "644", "{}", "+"
    end
    
    # Verify library is accessible
    dylib_path = "#{lib}/libfs_operations_core.dylib"
    if File.exist?(dylib_path)
      system "chmod", "755", dylib_path
    end
  end

  def caveats
    <<~EOS
      FileManager has been installed!
      
      To start the file manager:
        filemanager
      
      The web interface will be available at:
        http://localhost:8080
      
      Configuration file:
        #{etc}/filemanager/config.yaml
      
      Frontend files:
        #{share}/filemanager/frontend
      
      To run as a background service:
        brew services start filemanager
    EOS
  end

  service do
    run [opt_bin/"filemanager"]
    keep_alive true
    log_path var/"log/filemanager/stdout.log"
    error_log_path var/"log/filemanager/stderr.log"
  end

  test do
    # Test that binary exists and is executable
    assert_predicate bin/"filemanager", :exist?
    assert_predicate bin/"filemanager", :executable?
    
    # Test that library exists
    assert_predicate lib/"libfs_operations_core.dylib", :exist?
    
    # Test version output (if implemented)
    system "#{bin}/filemanager", "--version" rescue nil
  end
end
