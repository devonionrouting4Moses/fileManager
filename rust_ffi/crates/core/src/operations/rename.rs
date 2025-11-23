use crate::common::FsResult;
use std::fs;

/// Rename a file or directory from old_path to new_path
pub fn rename_path(old_path: &str, new_path: &str) -> FsResult<String> {
    fs::rename(old_path, new_path)?;
    Ok(format!("Renamed: {} -> {}", old_path, new_path))
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use std::path::Path;

    #[test]
    fn test_rename_file() {
        let old_path = "/tmp/test_rename_old.txt";
        let new_path = "/tmp/test_rename_new.txt";
        
        fs::File::create(old_path).unwrap();
        
        let result = rename_path(old_path, new_path);
        assert!(result.is_ok());
        assert!(!Path::new(old_path).exists());
        assert!(Path::new(new_path).exists());
        
        let _ = fs::remove_file(new_path);
    }
}