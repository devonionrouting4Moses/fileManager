// Core library implementation
// Re-export all modules
pub mod common;
pub mod operations;
pub mod ffi;

// Re-export main types for convenience
pub use common::{FsError, FsResult};
pub use operations::{
    copy::copy_path as copy_operation,
    create::{create_file, create_folder},
    delete::delete_path as delete_operation,
    file_permissions::change_permissions as change_perms,
    move_ops::move_path as move_operation,
    rename::rename_path as rename_operation,
};