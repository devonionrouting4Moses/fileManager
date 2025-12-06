class Filemanager < Formula
  desc "File manager with web interface"
  homepage "https://github.com/devonionrouting4Moses/fileManager"
  url "https://github.com/devonionrouting4Moses/fileManager/archive/v2.0.2.tar.gz"
  version "2.0.2"
  sha256 "CHECKSUM_HERE"  # Update with actual checksum
  license "MIT"

  depends_on "go" => :build

  def install
    # Build the binary
    cd "file_manager" do
      system "go", "build", "-o", "filemanager", "./cmd/app"
      
      # Install binary
      bin.install "filemanager"
    end
    
    # Install frontend (bundled)
    (share/"filemanager/frontend").install Dir["filemanager_frontend/*"]
    
    # Create config directory
    (etc/"filemanager").mkpath
    
    # Create system config
    (etc/"filemanager/config.yaml").write <<~EOS
      # FileManager System Configuration
      # Frontend directory path
      frontend_dir: #{share}/filemanager/frontend
    EOS
  end
  
  def post_install
    # Set frontend directory permissions
    frontend_dir = "#{share}/filemanager/frontend"
    
    # Make directories readable
    system "chmod", "-R", "755", frontend_dir
    
    # Make files readable
    system "find", frontend_dir, "-type", "f", "-exec", "chmod", "644", "{}", "+"
  end
  
  test do
    system "#{bin}/filemanager", "--version"
  end
end
