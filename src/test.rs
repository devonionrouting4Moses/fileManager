use std::ffi::{CStr, CString, OsStr};
use std::fs;
use std::io;
use std::os::raw::c_char;
use std::path::{Path, PathBuf};

#[cfg(unix)]
use std::os::unix::fs::PermissionsExt;

#[cfg(windows)]
use std::os::windows::fs::MetadataExt;

/// Result structure for FFI operations
#[repr(C)]
pub struct OperationResult {
    success: i32,
    message: *mut c_char,
}

impl OperationResult {
    /// Creates a success result with a message
    fn success(msg: impl AsRef<str>) -> Self {
        Self {
            success: 1,
            message: Self::create_message(msg.as_ref()),
        }
    }

    /// Creates an error result with a message
    fn error(msg: impl AsRef<str>) -> Self {
        Self {
            success: 0,
            message: Self::create_message(msg.as_ref()),
        }
    }

    /// Helper to create a C string message
    fn create_message(msg: &str) -> *mut c_char {
        CString::new(msg)
            .unwrap_or_else(|_| CString::new("Invalid message content").unwrap())
            .into_raw()
    }
}

/// Helper trait to safely convert C strings to Rust strings
trait CStrExt {
    fn to_string_result(&self) -> Result<&str, &'static str>;
}

impl CStrExt for CStr {
    fn to_string_result(&self) -> Result<&str, &'static str> {
        self.to_str().map_err(|_| "Invalid UTF-8 encoding")
    }
}

/// Macro to reduce boilerplate for dual-path FFI functions
macro_rules! dual_path_operation {
    ($func_name:ident, $op:expr, $success_msg:expr) => {
        #[no_mangle]
        pub extern "C" fn $func_name(
            src: *const c_char,
            dst: *const c_char,
        ) -> OperationResult {
            match dual_path_helper(src, dst, $op) {
                Ok((src_str, dst_str)) => {
                    OperationResult::success(format!($success_msg, src_str, dst_str))
                }
                Err(e) => OperationResult::error(e),
            }
        }
    };
}

/// Helper function for operations requiring two paths
fn dual_path_helper<F>(
    src: *const c_char,
    dst: *const c_char,
    operation: F,
) -> Result<(String, String), String>
where
    F: FnOnce(&Path, &Path) -> io::Result<()>,
{
    let src_str = parse_c_path(src)?;
    let dst_str = parse_c_path(dst)?;
    
    let src_path = Path::new(src_str);
    let dst_path = Path::new(dst_str);
    
    operation(src_path, dst_path)
        .map_err(|e| format!("Operation failed: {}", e))?;
    
    Ok((src_str.to_string(), dst_str.to_string()))
}

/// Safely parse a C string pointer to a Rust string
fn parse_c_path(path: *const c_char) -> Result<&'static str, String> {
    if path.is_null() {
        return Err("Null pointer provided".to_string());
    }
    
    unsafe { CStr::from_ptr(path) }
        .to_string_result()
        .map_err(|e| e.to_string())
}

/// Executes a single-path operation with proper error handling
fn single_path_operation<F>(
    path: *const c_char,
    operation: F,
    success_msg: impl Fn(&str) -> String,
) -> OperationResult
where
    F: FnOnce(&Path) -> io::Result<()>,
{
    match parse_c_path(path) {
        Ok(path_str) => {
            let path_obj = Path::new(path_str);
            match operation(path_obj) {
                Ok(_) => OperationResult::success(success_msg(path_str)),
                Err(e) => OperationResult::error(format!("Operation failed: {}", e)),
            }
        }
        Err(e) => OperationResult::error(e),
    }
}

#[no_mangle]
pub extern "C" fn create_folder(path: *const c_char) -> OperationResult {
    single_path_operation(
        path,
        |p| fs::create_dir_all(p),
        |s| format!("Folder created: {}", s),
    )
}

#[no_mangle]
pub extern "C" fn create_file(path: *const c_char) -> OperationResult {
    single_path_operation(
        path,
        |p| fs::File::create(p).map(|_| ()),
        |s| format!("File created: {}", s),
    )
}

#[no_mangle]
pub extern "C" fn delete_path(path: *const c_char) -> OperationResult {
    single_path_operation(
        path,
        |p| {
            if p.is_dir() {
                fs::remove_dir_all(p)
            } else {
                fs::remove_file(p)
            }
        },
        |s| format!("Deleted: {}", s),
    )
}

#[no_mangle]
pub extern "C" fn change_permissions(path: *const c_char, mode: u32) -> OperationResult {
    match parse_c_path(path) {
        Ok(path_str) => {
            let path_obj = Path::new(path_str);
            match set_permissions(path_obj, mode) {
                Ok(_) => OperationResult::success(format!(
                    "Permissions changed: {} (0{:o})",
                    path_str, mode
                )),
                Err(e) => OperationResult::error(format!("Failed to change permissions: {}", e)),
            }
        }
        Err(e) => OperationResult::error(e),
    }
}

/// Permission flags for cross-platform use
#[repr(C)]
pub struct PermissionFlags {
    pub read: bool,
    pub write: bool,
    pub execute: bool,
}

/// Sets permissions on a path (cross-platform)
/// 
/// # Platform-specific behavior
/// - **Unix/Linux/macOS/HarmonyOS**: Uses standard Unix permissions (chmod)
/// - **Windows**: Maps to read-only attribute (limited functionality)
/// 
/// # Arguments
/// * `path` - The file or directory path
/// * `mode` - Unix-style permission bits (e.g., 0o755, 0o644)
fn set_permissions(path: &Path, mode: u32) -> io::Result<()> {
    #[cfg(unix)]
    {
        set_permissions_unix(path, mode)
    }
    
    #[cfg(windows)]
    {
        set_permissions_windows(path, mode)
    }
    
    #[cfg(not(any(unix, windows)))]
    {
        // Fallback for other platforms (like WASI, embedded systems)
        set_permissions_generic(path, mode)
    }
}

#[cfg(unix)]
fn set_permissions_unix(path: &Path, mode: u32) -> io::Result<()> {
    let metadata = fs::metadata(path)?;
    let mut permissions = metadata.permissions();
    permissions.set_mode(mode);
    fs::set_permissions(path, permissions)
}

#[cfg(windows)]
fn set_permissions_windows(path: &Path, mode: u32) -> io::Result<()> {
    use std::fs::OpenOptions;
    
    // Windows doesn't have Unix-style permissions
    // Map Unix permissions to Windows read-only attribute
    // Owner write permission (bit 7): 0o200
    let readonly = (mode & 0o200) == 0;
    
    let metadata = fs::metadata(path)?;
    let mut permissions = metadata.permissions();
    permissions.set_readonly(readonly);
    
    fs::set_permissions(path, permissions)?;
    
    // For directories and executables, we can try to set additional ACLs
    // This requires Windows-specific crates like `windows` or `winapi`
    // For now, we'll just handle the basic read-only flag
    
    Ok(())
}

#[cfg(not(any(unix, windows)))]
fn set_permissions_generic(path: &Path, mode: u32) -> io::Result<()> {
    // Generic fallback: only handle read-only flag
    let readonly = (mode & 0o200) == 0;
    let metadata = fs::metadata(path)?;
    let mut permissions = metadata.permissions();
    permissions.set_readonly(readonly);
    fs::set_permissions(path, permissions)
}

/// Advanced permission setter with detailed control (Unix-like systems only)
#[cfg(unix)]
pub fn set_permissions_detailed(
    path: &Path,
    owner: PermissionFlags,
    group: PermissionFlags,
    others: PermissionFlags,
) -> io::Result<()> {
    let mode = build_unix_mode(&owner, &group, &others);
    set_permissions_unix(path, mode)
}

#[cfg(unix)]
fn build_unix_mode(
    owner: &PermissionFlags,
    group: &PermissionFlags,
    others: &PermissionFlags,
) -> u32 {
    let mut mode: u32 = 0;
    
    // Owner permissions
    if owner.read { mode |= 0o400; }
    if owner.write { mode |= 0o200; }
    if owner.execute { mode |= 0o100; }
    
    // Group permissions
    if group.read { mode |= 0o040; }
    if group.write { mode |= 0o020; }
    if group.execute { mode |= 0o010; }
    
    // Others permissions
    if others.read { mode |= 0o004; }
    if others.write { mode |= 0o002; }
    if others.execute { mode |= 0o001; }
    
    mode
}

/// Cross-platform permission getter
pub fn get_permissions(path: &Path) -> io::Result<u32> {
    #[cfg(unix)]
    {
        let metadata = fs::metadata(path)?;
        let permissions = metadata.permissions();
        Ok(permissions.mode())
    }
    
    #[cfg(windows)]
    {
        let metadata = fs::metadata(path)?;
        let permissions = metadata.permissions();
        // Map Windows readonly to Unix-style
        let mode = if permissions.readonly() {
            0o444 // r--r--r--
        } else {
            0o666 // rw-rw-rw-
        };
        Ok(mode)
    }
    
    #[cfg(not(any(unix, windows)))]
    {
        let metadata = fs::metadata(path)?;
        let permissions = metadata.permissions();
        let mode = if permissions.readonly() { 0o444 } else { 0o666 };
        Ok(mode)
    }
}

/// Check if current platform supports full Unix-style permissions
pub const fn supports_unix_permissions() -> bool {
    cfg!(unix)
}

/// Platform-specific permission constants
pub mod permission_constants {
    // Unix standard permissions
    pub const OWNER_READ: u32 = 0o400;
    pub const OWNER_WRITE: u32 = 0o200;
    pub const OWNER_EXEC: u32 = 0o100;
    pub const OWNER_ALL: u32 = 0o700;
    
    pub const GROUP_READ: u32 = 0o040;
    pub const GROUP_WRITE: u32 = 0o020;
    pub const GROUP_EXEC: u32 = 0o010;
    pub const GROUP_ALL: u32 = 0o070;
    
    pub const OTHERS_READ: u32 = 0o004;
    pub const OTHERS_WRITE: u32 = 0o002;
    pub const OTHERS_EXEC: u32 = 0o001;
    pub const OTHERS_ALL: u32 = 0o007;
    
    // Common combinations
    pub const MODE_755: u32 = 0o755; // rwxr-xr-x
    pub const MODE_644: u32 = 0o644; // rw-r--r--
    pub const MODE_600: u32 = 0o600; // rw-------
    pub const MODE_777: u32 = 0o777; // rwxrwxrwx
}

// Generate dual-path operations using macro
dual_path_operation!(
    rename_path,
    |src: &Path, dst: &Path| fs::rename(src, dst),
    "Renamed: {} -> {}"
);

dual_path_operation!(
    move_path,
    |src: &Path, dst: &Path| fs::rename(src, dst),
    "Moved: {} -> {}"
);

#[no_mangle]
pub extern "C" fn copy_path(src: *const c_char, dst: *const c_char) -> OperationResult {
    match dual_path_helper(src, dst, |src_path, dst_path| {
        if src_path.is_dir() {
            copy_dir_recursive(src_path, dst_path)
        } else {
            fs::copy(src_path, dst_path).map(|_| ())
        }
    }) {
        Ok((src_str, dst_str)) => {
            OperationResult::success(format!("Copied: {} -> {}", src_str, dst_str))
        }
        Err(e) => OperationResult::error(e),
    }
}

/// Recursively copies a directory and its contents
fn copy_dir_recursive(src: &Path, dst: &Path) -> io::Result<()> {
    fs::create_dir_all(dst)?;
    
    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let file_type = entry.file_type()?;
        let dst_path = dst.join(entry.file_name());
        
        if file_type.is_dir() {
            copy_dir_recursive(&entry.path(), &dst_path)?;
        } else if file_type.is_file() {
            fs::copy(entry.path(), dst_path)?;
        }
        // Silently skip symlinks and other special files
    }
    
    Ok(())
}

#[no_mangle]
pub extern "C" fn free_result(result: OperationResult) {
    if !result.message.is_null() {
        unsafe {
            drop(CString::from_raw(result.message));
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::ffi::CString;

    #[test]
    fn test_create_and_delete_folder() {
        let path = CString::new("/tmp/test_folder").unwrap();
        let result = create_folder(path.as_ptr());
        assert_eq!(result.success, 1);
        free_result(result);
        
        let result = delete_path(path.as_ptr());
        assert_eq!(result.success, 1);
        free_result(result);
    }
    
    #[test]
    #[cfg(unix)]
    fn test_unix_permissions() {
        use std::fs::File;
        let test_file = "/tmp/test_perms.txt";
        File::create(test_file).unwrap();
        
        let path = Path::new(test_file);
        assert!(set_permissions(path, 0o644).is_ok());
        
        let mode = get_permissions(path).unwrap();
        assert_eq!(mode & 0o777, 0o644);
        
        fs::remove_file(test_file).unwrap();
    }
    
    #[test]
    fn test_permission_platform_support() {
        // This should compile on all platforms
        let is_unix = supports_unix_permissions();
        #[cfg(unix)]
        assert!(is_unix);
        #[cfg(not(unix))]
        assert!(!is_unix);
    }
}