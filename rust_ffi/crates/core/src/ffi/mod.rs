use crate::common::{c_str_to_string, OperationResult};
use crate::operations;
use std::os::raw::c_char;

/// FFI wrapper for create_folder
#[no_mangle]
pub extern "C" fn create_folder(path: *const c_char) -> OperationResult {
    match c_str_to_string(path) {
        Ok(path_str) => match operations::create::create_folder(&path_str) {
            Ok(msg) => OperationResult::success(&msg),
            Err(e) => OperationResult::error(&e.to_string()),
        },
        Err(e) => OperationResult::error(&e.to_string()),
    }
}

/// FFI wrapper for create_file
#[no_mangle]
pub extern "C" fn create_file(path: *const c_char) -> OperationResult {
    match c_str_to_string(path) {
        Ok(path_str) => match operations::create::create_file(&path_str) {
            Ok(msg) => OperationResult::success(&msg),
            Err(e) => OperationResult::error(&e.to_string()),
        },
        Err(e) => OperationResult::error(&e.to_string()),
    }
}

/// FFI wrapper for rename_path
#[no_mangle]
pub extern "C" fn rename_path(old_path: *const c_char, new_path: *const c_char) -> OperationResult {
    let old_str = match c_str_to_string(old_path) {
        Ok(s) => s,
        Err(e) => return OperationResult::error(&e.to_string()),
    };
    
    let new_str = match c_str_to_string(new_path) {
        Ok(s) => s,
        Err(e) => return OperationResult::error(&e.to_string()),
    };

    match operations::rename::rename_path(&old_str, &new_str) {
        Ok(msg) => OperationResult::success(&msg),
        Err(e) => OperationResult::error(&e.to_string()),
    }
}

/// FFI wrapper for delete_path
#[no_mangle]
pub extern "C" fn delete_path(path: *const c_char) -> OperationResult {
    match c_str_to_string(path) {
        Ok(path_str) => match operations::delete::delete_path(&path_str) {
            Ok(msg) => OperationResult::success(&msg),
            Err(e) => OperationResult::error(&e.to_string()),
        },
        Err(e) => OperationResult::error(&e.to_string()),
    }
}

/// FFI wrapper for change_permissions
#[no_mangle]
pub extern "C" fn change_permissions(path: *const c_char, mode: u32) -> OperationResult {
    match c_str_to_string(path) {
        Ok(path_str) => match operations::file_permissions::change_permissions(&path_str, mode) {
            Ok(msg) => OperationResult::success(&msg),
            Err(e) => OperationResult::error(&e.to_string()),
        },
        Err(e) => OperationResult::error(&e.to_string()),
    }
}

/// FFI wrapper for move_path
#[no_mangle]
pub extern "C" fn move_path(src: *const c_char, dst: *const c_char) -> OperationResult {
    let src_str = match c_str_to_string(src) {
        Ok(s) => s,
        Err(e) => return OperationResult::error(&e.to_string()),
    };
    
    let dst_str = match c_str_to_string(dst) {
        Ok(s) => s,
        Err(e) => return OperationResult::error(&e.to_string()),
    };

    match operations::move_ops::move_path(&src_str, &dst_str) {
        Ok(msg) => OperationResult::success(&msg),
        Err(e) => OperationResult::error(&e.to_string()),
    }
}

/// FFI wrapper for copy_path
#[no_mangle]
pub extern "C" fn copy_path(src: *const c_char, dst: *const c_char) -> OperationResult {
    let src_str = match c_str_to_string(src) {
        Ok(s) => s,
        Err(e) => return OperationResult::error(&e.to_string()),
    };
    
    let dst_str = match c_str_to_string(dst) {
        Ok(s) => s,
        Err(e) => return OperationResult::error(&e.to_string()),
    };

    match operations::copy::copy_path(&src_str, &dst_str) {
        Ok(msg) => OperationResult::success(&msg),
        Err(e) => OperationResult::error(&e.to_string()),
    }
}