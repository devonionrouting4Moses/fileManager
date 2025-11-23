use crate::common::FsResult;
use std::fs;

/// Move a file or directory from source to destination
/// This is essentially a rename operation that can work across filesystems
pub fn move_path(src: &str, dst: &str) -> FsResult<String> {
    fs::rename(src, dst)?;
    Ok(format!("Moved: {} -> {}", src, dst))
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;
    use std::path::Path;

    #[test]
    fn test_move_file() {
        let src = "/tmp/test_move_src.txt";
        let dst = "/tmp/test_move_dst.txt";
        
        fs::File::create(src).unwrap();
        
        let result = move_path(src, dst);
        assert!(result.is_ok());
        assert!(!Path::new(src).exists());
        assert!(Path::new(dst).exists());
        
        let _ = fs::remove_file(dst);
    }
}