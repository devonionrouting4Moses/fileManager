use crate::common::FsResult;
use std::fs;
use std::path::Path;

/// Copy a file or directory from source to destination
/// Recursively copies directories and their contents
pub fn copy_path(src: &str, dst: &str) -> FsResult<String> {
    let src_path = Path::new(src);
    
    if src_path.is_dir() {
        copy_dir_all(src, dst)?;
    } else {
        fs::copy(src, dst)?;
    }
    
    Ok(format!("Copied: {} -> {}", src, dst))
}

/// Recursively copy a directory and all its contents
fn copy_dir_all(src: &str, dst: &str) -> std::io::Result<()> {
    fs::create_dir_all(dst)?;
    
    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let file_type = entry.file_type()?;
        let src_path = entry.path();
        let dst_path = Path::new(dst).join(entry.file_name());
        
        if file_type.is_dir() {
            copy_dir_all(
                src_path.to_str().unwrap(),
                dst_path.to_str().unwrap()
            )?;
        } else {
            fs::copy(&src_path, &dst_path)?;
        }
    }
    
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;

    #[test]
    fn test_copy_file() {
        let src = "/tmp/test_copy_src.txt";
        let dst = "/tmp/test_copy_dst.txt";
        
        fs::write(src, "test content").unwrap();
        
        let result = copy_path(src, dst);
        assert!(result.is_ok());
        assert!(Path::new(src).exists());
        assert!(Path::new(dst).exists());
        
        let _ = fs::remove_file(src);
        let _ = fs::remove_file(dst);
    }
}