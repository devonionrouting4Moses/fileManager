use std::ffi::{CStr, CString};
use std::fs;
use std::os::raw::c_char;
use std::os::unix::fs::PermissionsExt;
use std::path::Path;

#[repr(C)]
pub struct OperationResult {
    success: i32,
    message: *mut c_char,
}

impl OperationResult {
    fn success(msg: &str) -> Self {
        OperationResult {
            success: 1,
            message: CString::new(msg).unwrap().into_raw(),
        }
    }

    fn error(msg: &str) -> Self {
        OperationResult {
            success: 0,
            message: CString::new(msg).unwrap().into_raw(),
        }
    }
}

#[no_mangle]
pub extern "C" fn create_folder(path: *const c_char) -> OperationResult {
    let c_str = unsafe { CStr::from_ptr(path) };
    let path_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 path"),
    };

    match fs::create_dir_all(path_str) {
        Ok(_) => OperationResult::success(&format!("Folder created: {}", path_str)),
        Err(e) => OperationResult::error(&format!("Failed to create folder: {}", e)),
    }
}

#[no_mangle]
pub extern "C" fn create_file(path: *const c_char) -> OperationResult {
    let c_str = unsafe { CStr::from_ptr(path) };
    let path_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 path"),
    };

    match fs::File::create(path_str) {
        Ok(_) => OperationResult::success(&format!("File created: {}", path_str)),
        Err(e) => OperationResult::error(&format!("Failed to create file: {}", e)),
    }
}

#[no_mangle]
pub extern "C" fn rename_path(old_path: *const c_char, new_path: *const c_char) -> OperationResult {
    let old_c_str = unsafe { CStr::from_ptr(old_path) };
    let new_c_str = unsafe { CStr::from_ptr(new_path) };
    
    let old_str = match old_c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 in old path"),
    };
    
    let new_str = match new_c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 in new path"),
    };

    match fs::rename(old_str, new_str) {
        Ok(_) => OperationResult::success(&format!("Renamed: {} -> {}", old_str, new_str)),
        Err(e) => OperationResult::error(&format!("Failed to rename: {}", e)),
    }
}

#[no_mangle]
pub extern "C" fn delete_path(path: *const c_char) -> OperationResult {
    let c_str = unsafe { CStr::from_ptr(path) };
    let path_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 path"),
    };

    let path_obj = Path::new(path_str);
    
    let result = if path_obj.is_dir() {
        fs::remove_dir_all(path_str)
    } else {
        fs::remove_file(path_str)
    };

    match result {
        Ok(_) => OperationResult::success(&format!("Deleted: {}", path_str)),
        Err(e) => OperationResult::error(&format!("Failed to delete: {}", e)),
    }
}

#[no_mangle]
pub extern "C" fn change_permissions(path: *const c_char, mode: u32) -> OperationResult {
    let c_str = unsafe { CStr::from_ptr(path) };
    let path_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 path"),
    };

    match fs::metadata(path_str) {
        Ok(metadata) => {
            let mut permissions = metadata.permissions();
            permissions.set_mode(mode);
            
            match fs::set_permissions(path_str, permissions) {
                Ok(_) => OperationResult::success(&format!("Permissions changed: {} (0{:o})", path_str, mode)),
                Err(e) => OperationResult::error(&format!("Failed to change permissions: {}", e)),
            }
        }
        Err(e) => OperationResult::error(&format!("Failed to get metadata: {}", e)),
    }
}

#[no_mangle]
pub extern "C" fn move_path(src: *const c_char, dst: *const c_char) -> OperationResult {
    let src_c_str = unsafe { CStr::from_ptr(src) };
    let dst_c_str = unsafe { CStr::from_ptr(dst) };
    
    let src_str = match src_c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 in source path"),
    };
    
    let dst_str = match dst_c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 in destination path"),
    };

    match fs::rename(src_str, dst_str) {
        Ok(_) => OperationResult::success(&format!("Moved: {} -> {}", src_str, dst_str)),
        Err(e) => OperationResult::error(&format!("Failed to move: {}", e)),
    }
}

#[no_mangle]
pub extern "C" fn copy_path(src: *const c_char, dst: *const c_char) -> OperationResult {
    let src_c_str = unsafe { CStr::from_ptr(src) };
    let dst_c_str = unsafe { CStr::from_ptr(dst) };
    
    let src_str = match src_c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 in source path"),
    };
    
    let dst_str = match dst_c_str.to_str() {
        Ok(s) => s,
        Err(_) => return OperationResult::error("Invalid UTF-8 in destination path"),
    };

    let src_path = Path::new(src_str);
    
    let result = if src_path.is_dir() {
        copy_dir_all(src_str, dst_str)
    } else {
        fs::copy(src_str, dst_str).map(|_| ())
    };

    match result {
        Ok(_) => OperationResult::success(&format!("Copied: {} -> {}", src_str, dst_str)),
        Err(e) => OperationResult::error(&format!("Failed to copy: {}", e)),
    }
}

fn copy_dir_all(src: &str, dst: &str) -> std::io::Result<()> {
    fs::create_dir_all(dst)?;
    
    for entry in fs::read_dir(src)? {
        let entry = entry?;
        let file_type = entry.file_type()?;
        let src_path = entry.path();
        let dst_path = Path::new(dst).join(entry.file_name());
        
        if file_type.is_dir() {
            copy_dir_all(src_path.to_str().unwrap(), dst_path.to_str().unwrap())?;
        } else {
            fs::copy(&src_path, &dst_path)?;
        }
    }
    
    Ok(())
}

#[no_mangle]
pub extern "C" fn free_result(result: OperationResult) {
    if !result.message.is_null() {
        unsafe {
            let _ = CString::from_raw(result.message);
        }
    }
}