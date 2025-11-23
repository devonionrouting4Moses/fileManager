use crate::common::FsResult;
use std::fs;

#[cfg(unix)]
use std::os::unix::fs::PermissionsExt;

/// Change file or directory permissions (Unix only)
#[cfg(unix)]
pub fn change_permissions(path: &str, mode: u32) -> FsResult<String> {
    let metadata = fs::metadata(path)?;
    let mut permissions = metadata.permissions();
    permissions.set_mode(mode);
    
    fs::set_permissions(path, permissions)?;
    Ok(format!("Permissions changed: {} (0{:o})", path, mode))
}

/// Change file or directory permissions (Windows stub)
#[cfg(not(unix))]
pub fn change_permissions(path: &str, _mode: u32) -> FsResult<String> {
    use crate::common::FsError;
    Err(FsError::PermissionError(
        "Permission changes are only supported on Unix systems".to_string()
    ))
}

#[cfg(test)]
#[cfg(unix)]
mod tests {
    use super::*;
    use std::fs;

    #[test]
    fn test_change_permissions() {
        let test_path = "/tmp/test_permissions.txt";
        fs::File::create(test_path).unwrap();
        
        let result = change_permissions(test_path, 0o644);
        assert!(result.is_ok());
        
        let metadata = fs::metadata(test_path).unwrap();
        assert_eq!(metadata.permissions().mode() & 0o777, 0o644);
        
        let _ = fs::remove_file(test_path);
    }
}