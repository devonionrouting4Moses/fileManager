use crate::common::FsResult;
use std::fs;

/// Create a new folder at the specified path
/// Creates all parent directories if they don't exist
pub fn create_folder(path: &str) -> FsResult<String> {
    fs::create_dir_all(path)?;
    Ok(format!("Folder created: {}", path))
}

/// Create a new file at the specified path
pub fn create_file(path: &str) -> FsResult<String> {
    fs::File::create(path)?;
    Ok(format!("File created: {}", path))
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::path::Path;

    #[test]
    fn test_create_folder() {
        let test_path = "/tmp/test_rust_ffi_folder";
        let result = create_folder(test_path);
        assert!(result.is_ok());
        assert!(Path::new(test_path).exists());
        let _ = std::fs::remove_dir(test_path);
    }

    #[test]
    fn test_create_file() {
        let test_path = "/tmp/test_rust_ffi_file.txt";
        let result = create_file(test_path);
        assert!(result.is_ok());
        assert!(Path::new(test_path).exists());
        let _ = std::fs::remove_file(test_path);
    }
}