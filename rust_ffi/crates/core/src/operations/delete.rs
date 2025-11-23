use crate::common::FsResult;
use std::fs;
use std::path::Path;

/// Delete a file or directory at the specified path
/// Recursively deletes directories and their contents
pub fn delete_path(path: &str) -> FsResult<String> {
    let path_obj = Path::new(path);
    
    if path_obj.is_dir() {
        fs::remove_dir_all(path)?;
    } else {
        fs::remove_file(path)?;
    }
    
    Ok(format!("Deleted: {}", path))
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;

    #[test]
    fn test_delete_file() {
        let test_path = "/tmp/test_delete_file.txt";
        fs::File::create(test_path).unwrap();
        
        let result = delete_path(test_path);
        assert!(result.is_ok());
        assert!(!std::path::Path::new(test_path).exists());
    }

    #[test]
    fn test_delete_directory() {
        let test_path = "/tmp/test_delete_dir";
        fs::create_dir_all(test_path).unwrap();
        
        let result = delete_path(test_path);
        assert!(result.is_ok());
        assert!(!std::path::Path::new(test_path).exists());
    }
}