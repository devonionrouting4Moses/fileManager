use clap::{Parser, Subcommand};
use fs_operations_core::*;

#[derive(Parser)]
#[command(name = "fsops")]
#[command(about = "File system operations CLI", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Create a new folder
    CreateFolder {
        /// Path to the folder
        path: String,
    },
    /// Create a new file
    CreateFile {
        /// Path to the file
        path: String,
    },
    /// Rename a file or folder
    Rename {
        /// Source path
        old_path: String,
        /// Destination path
        new_path: String,
    },
    /// Delete a file or folder
    Delete {
        /// Path to delete
        path: String,
    },
    /// Change file permissions (Unix only)
    Chmod {
        /// Path to the file/folder
        path: String,
        /// Permission mode (octal, e.g., 755)
        #[arg(value_parser = parse_octal)]
        mode: u32,
    },
    /// Move a file or folder
    Move {
        /// Source path
        src: String,
        /// Destination path
        dst: String,
    },
    /// Copy a file or folder
    Copy {
        /// Source path
        src: String,
        /// Destination path
        dst: String,
    },
}

fn parse_octal(s: &str) -> Result<u32, std::num::ParseIntError> {
    u32::from_str_radix(s, 8)
}

fn main() -> anyhow::Result<()> {
    let cli = Cli::parse();

    let result = match cli.command {
        Commands::CreateFolder { path } => {
            create_folder(&path)
        }
        Commands::CreateFile { path } => {
            create_file(&path)
        }
        Commands::Rename { old_path, new_path } => {
            rename_operation(&old_path, &new_path)
        }
        Commands::Delete { path } => {
            delete_operation(&path)
        }
        Commands::Chmod { path, mode } => {
            change_perms(&path, mode)
        }
        Commands::Move { src, dst } => {
            move_operation(&src, &dst)
        }
        Commands::Copy { src, dst } => {
            copy_operation(&src, &dst)
        }
    };

    match result {
        Ok(msg) => {
            println!("✓ {}", msg);
            Ok(())
        }
        Err(e) => {
            eprintln!("✗ Error: {}", e);
            std::process::exit(1);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_octal() {
        assert_eq!(parse_octal("755").unwrap(), 0o755);
        assert_eq!(parse_octal("644").unwrap(), 0o644);
    }
}