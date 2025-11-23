use std::ffi::{CStr, CString};
use std::os::raw::c_char;
use thiserror::Error;

/// Custom error type for file system operations
#[derive(Error, Debug)]
pub enum FsError {
    #[error("Invalid UTF-8 in path: {0}")]
    InvalidUtf8(String),
    
    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),
    
    #[error("Path error: {0}")]
    PathError(String),
    
    #[error("Permission error: {0}")]
    PermissionError(String),
}

/// Result type for file system operations
pub type FsResult<T> = Result<T, FsError>;

/// C-compatible operation result
#[repr(C)]
pub struct OperationResult {
    pub success: i32,
    pub message: *mut c_char,
}

impl OperationResult {
    /// Create a success result with a message
    pub fn success(msg: &str) -> Self {
        OperationResult {
            success: 1,
            message: CString::new(msg).unwrap().into_raw(),
        }
    }

    /// Create an error result with a message
    pub fn error(msg: &str) -> Self {
        OperationResult {
            success: 0,
            message: CString::new(msg).unwrap().into_raw(),
        }
    }
}

/// Helper function to safely convert C string to Rust string
pub fn c_str_to_string(ptr: *const c_char) -> FsResult<String> {
    let c_str = unsafe { CStr::from_ptr(ptr) };
    c_str
        .to_str()
        .map(|s| s.to_string())
        .map_err(|e| FsError::InvalidUtf8(e.to_string()))
}

/// Free the memory allocated for OperationResult
#[no_mangle]
pub extern "C" fn free_result(result: OperationResult) {
    if !result.message.is_null() {
        unsafe {
            let _ = CString::from_raw(result.message);
        }
    }
}