use anyhow::{bail, Context, Result};
use clap::{Parser, Subcommand};
use std::{
    env,
    ffi::OsString,
    fs,
    path::{Path, PathBuf},
    process::{Command, Stdio},
};

const CHAFA_VERSION: &str = "1.18.2";
const CHAFA_URL: &str =
    "https://github.com/hpjansson/chafa/releases/download/1.18.2/chafa-1.18.2.tar.xz";

#[derive(Parser)]
#[command(name = "xtask")]
#[command(about = "Project build tasks")]
struct Cli {
    #[command(subcommand)]
    command: CommandKind,
}

#[derive(Subcommand)]
enum CommandKind {
    /// Build Chafa and then the Rust application
    Build {
        /// Build in release mode
        #[arg(long)]
        release: bool,

        /// Override the Cargo target triple
        #[arg(long)]
        target: Option<String>,

        /// Force rebuilding Chafa
        #[arg(long)]
        force: bool,
    },

    /// Build only Chafa
    Chafa {
        /// Override the Cargo target triple
        #[arg(long)]
        target: Option<String>,

        /// Force rebuilding Chafa
        #[arg(long)]
        force: bool,
    },

    /// Check that native build prerequisites are available
    Doctor,

    /// Remove locally built Chafa artifacts
    Clean,
}

#[derive(Debug)]
struct BuildInfo {
    root: PathBuf,
    native: PathBuf,
    chafa_source: PathBuf,
    chafa_prefix: PathBuf,
    target: String,
}

fn main() -> Result<()> {
    let cli = Cli::parse();

    match cli.command {
        CommandKind::Build {
            release,
            target,
            force,
        } => {
            let info = build_info(target)?;
            build_chafa(&info, force)?;
            build_rust(&info, release)?;
        }

        CommandKind::Chafa { target, force } => {
            let info = build_info(target)?;
            build_chafa(&info, force)?;
        }

        CommandKind::Doctor => doctor()?,

        CommandKind::Clean => clean()?,
    }

    Ok(())
}

fn build_info(target: Option<String>) -> Result<BuildInfo> {
    let root = project_root()?;
    let native = root.join(".native");

    let target = match target {
        Some(target) => target,
        None => default_target()?,
    };

    let chafa_source = native
        .join("chafa-source")
        .join(format!("chafa-{CHAFA_VERSION}"));

    let chafa_prefix = native
        .join("chafa")
        .join(&target)
        .join(CHAFA_VERSION);

    Ok(BuildInfo {
        root,
        native,
        chafa_source,
        chafa_prefix,
        target,
    })
}

fn project_root() -> Result<PathBuf> {
    let manifest_dir = PathBuf::from(env!("CARGO_MANIFEST_DIR"));

    manifest_dir
        .parent()
        .map(Path::to_path_buf)
        .context("xtask must be in the project workspace")
}

fn default_target() -> Result<String> {
    let host = rustc_host()?;

    if cfg!(windows) {
        /*
         * Chafa's upstream Windows support is MinGW-based, so don't
         * accidentally try to build the Chafa dependency with MSVC.
         */
        Ok("x86_64-pc-windows-gnu".to_string())
    } else {
        Ok(host)
    }
}

fn rustc_host() -> Result<String> {
    let output = Command::new("rustc")
        .args(["-vV"])
        .output()
        .context("failed to execute rustc")?;

    if !output.status.success() {
        bail!("rustc -vV failed");
    }

    let stdout = String::from_utf8(output.stdout)
        .context("rustc produced invalid UTF-8")?;

    stdout
        .lines()
        .find_map(|line| line.strip_prefix("host: "))
        .map(str::to_owned)
        .context("could not determine Rust host target")
}

fn build_chafa(info: &BuildInfo, force: bool) -> Result<()> {
    let marker = info.chafa_prefix.join(".built");

    if marker.exists() && !force {
        println!(
            "Chafa {} already built for {}",
            CHAFA_VERSION, info.target
        );
        return Ok(());
    }

    prepare_platform_dependencies(&info.target)?;

    download_chafa(&info.native, &info.chafa_source)?;

    fs::create_dir_all(&info.chafa_prefix)
        .with_context(|| {
            format!(
                "failed to create Chafa prefix {}",
                info.chafa_prefix.display()
            )
        })?;

    let source = info.chafa_source.clone();
    let prefix = info.chafa_prefix.clone();

    println!(
        "Building Chafa {} for {}",
        CHAFA_VERSION, info.target
    );

    if cfg!(windows) {
        build_chafa_windows(&source, &prefix)?;
    } else {
        build_chafa_unix(&source, &prefix)?;
    }

    fs::write(&marker, format!("{}\n", CHAFA_VERSION))
        .context("failed to write Chafa build marker")?;

    println!("Chafa installed into {}", prefix.display());

    Ok(())
}

fn download_chafa(native: &Path, source: &Path) -> Result<()> {
    if source.exists() {
        return Ok(());
    }

    fs::create_dir_all(native)
        .with_context(|| format!("failed to create {}", native.display()))?;

    let archive = native.join(format!("chafa-{CHAFA_VERSION}.tar.xz"));

    if !archive.exists() {
        println!("Downloading Chafa {CHAFA_VERSION}");

        if cfg!(windows) {
            let status = Command::new("curl.exe")
                .args([
                    "-L",
                    "--fail",
                    "--retry",
                    "3",
                    "-o",
                    archive.to_str().context("invalid archive path")?,
                    CHAFA_URL,
                ])
                .status()
                .context("failed to execute curl.exe")?;

            if !status.success() {
                bail!("curl failed while downloading Chafa");
            }
        } else if command_exists("curl") {
            let status = Command::new("curl")
                .args([
                    "-L",
                    "--fail",
                    "--retry",
                    "3",
                    "-o",
                    archive.to_str().context("invalid archive path")?,
                    CHAFA_URL,
                ])
                .status()
                .context("failed to execute curl")?;

            if !status.success() {
                bail!("curl failed while downloading Chafa");
            }
        } else if command_exists("fetch") {
            let status = Command::new("fetch")
                .args([
                    "-o",
                    archive.to_str().context("invalid archive path")?,
                    CHAFA_URL,
                ])
                .status()
                .context("failed to execute fetch")?;

            if !status.success() {
                bail!("fetch failed while downloading Chafa");
            }
        } else {
            bail!("neither curl nor fetch is available");
        }
    }

    println!("Extracting Chafa");

    let source_parent = source
        .parent()
        .context("invalid Chafa source path")?;

    fs::create_dir_all(source_parent)?;

    let status = Command::new("tar")
        .args([
            "-xf",
            archive.to_str().context("invalid archive path")?,
            "-C",
            source_parent
                .to_str()
                .context("invalid source directory")?,
        ])
        .status()
        .context("failed to execute tar")?;

    if !status.success() {
        bail!("failed to extract Chafa archive");
    }

    if !source.exists() {
        bail!(
            "expected Chafa source directory {} after extraction",
            source.display()
        );
    }

    Ok(())
}

fn build_chafa_unix(source: &Path, prefix: &Path) -> Result<()> {
    run_shell(
        "sh",
        &[
            "-c",
            &format!(
                "set -eu
                 cd {source}
                 ./configure \
                    --prefix={prefix} \
                    --without-tools \
                    --disable-shared \
                    --enable-static
                 make -j$(getconf _NPROCESSORS_ONLN 2>/dev/null || sysctl -n hw.ncpu || echo 1)
                 make install",
                source = shell_quote(source),
                prefix = shell_quote(prefix),
            ),
        ],
        None,
    )
}

fn build_chafa_windows(source: &Path, prefix: &Path) -> Result<()> {
    let msys_root = find_msys2()?;

    let bash = msys_root.join("usr").join("bin").join("bash.exe");

    let source_msys = windows_to_msys_path(&msys_root, source)?;
    let prefix_msys = windows_to_msys_path(&msys_root, prefix)?;

    let script = format!(
        r#"
set -eu

export PATH="/mingw64/bin:/usr/bin:$PATH"

cd {source}

./configure \
    --host=x86_64-w64-mingw32 \
    --prefix={prefix} \
    --without-tools \
    --disable-shared \
    --enable-static \
    --disable-Bsymbolic

make -j$(nproc)
make install
"#,
        source = shell_quote(&source_msys),
        prefix = shell_quote(&prefix_msys),
    );

    let status = Command::new(bash)
        .args(["-lc", &script])
        .status()
        .context("failed to execute MSYS2 bash")?;

    if !status.success() {
        bail!("Chafa build failed under MSYS2");
    }

    Ok(())
}

fn build_rust(info: &BuildInfo, release: bool) -> Result<()> {
    let mut command = Command::new("cargo");

    command.args(["build", "--target", &info.target]);

    if release {
        command.arg("--release");
    }

    let pkgconfig = info
        .chafa_prefix
        .join("lib")
        .join("pkgconfig");

    let mut pkgconfig_path = OsString::new();

    pkgconfig_path.push(pkgconfig);

    if let Some(existing) = env::var_os("PKG_CONFIG_PATH") {
        pkgconfig_path.push(if cfg!(windows) { ";" } else { ":" });
        pkgconfig_path.push(existing);
    }

    command.env("PKG_CONFIG_PATH", pkgconfig_path);
    command.env("PKG_CONFIG_ALLOW_SYSTEM_CFLAGS", "1");
    command.env("PKG_CONFIG_ALLOW_SYSTEM_LIBS", "1");

    if cfg!(windows) {
        let msys_root = find_msys2()?;

        let mingw_bin = msys_root.join("mingw64").join("bin");
        let usr_bin = msys_root.join("usr").join("bin");

        prepend_path(&mut command, mingw_bin);
        prepend_path(&mut command, usr_bin);

        command.env(
            "PKG_CONFIG",
            msys_root
                .join("mingw64")
                .join("bin")
                .join("pkg-config.exe"),
        );
    }

    println!("Building Rust target {}", info.target);

    run_command(command)
}

fn prepare_platform_dependencies(target: &str) -> Result<()> {
    if cfg!(windows) {
        if !target.starts_with("x86_64-pc-windows-gnu") {
            bail!(
                "Windows Chafa builds currently require x86_64-pc-windows-gnu; got {target}"
            );
        }

        let msys_root = find_msys2()?;

        println!(
            "Using MSYS2 at {}",
            msys_root.display()
        );

        let required = [
            "gcc",
            "make",
            "autoconf",
            "automake",
            "libtool",
            "pkgconf",
            "mingw-w64-x86_64-gcc",
            "mingw-w64-x86_64-glib2",
            "mingw-w64-x86_64-pkgconf",
        ];

        let status = Command::new(msys_root.join("usr").join("bin").join("bash.exe"))
            .args([
                "-lc",
                &format!(
                    "pacman -S --needed --noconfirm {}",
                    required.join(" ")
                ),
            ])
            .status()
            .context("failed to run MSYS2 pacman")?;

        if !status.success() {
            bail!("MSYS2 dependency installation failed");
        }

        ensure_rust_target(target)?;
    } else {
        ensure_rust_target(target)?;

        if !command_exists("pkg-config") {
            bail!("pkg-config is required but was not found");
        }

        if !command_exists("gcc") && !command_exists("cc") {
            bail!("a C compiler is required but was not found");
        }

        if !command_exists("make") {
            bail!("make is required but was not found");
        }
    }

    Ok(())
}

fn ensure_rust_target(target: &str) -> Result<()> {
    let output = Command::new("rustup")
        .args(["target", "list", "--installed"])
        .output()
        .context("failed to execute rustup")?;

    if !output.status.success() {
        bail!("rustup target list failed");
    }

    let installed = String::from_utf8(output.stdout)
        .context("rustup produced invalid UTF-8")?;

    if installed.lines().any(|line| line.trim() == target) {
        return Ok(());
    }

    println!("Installing Rust target {target}");

    let status = Command::new("rustup")
        .args(["target", "add", target])
        .status()
        .context("failed to execute rustup target add")?;

    if !status.success() {
        bail!("failed to install Rust target {target}");
    }

    Ok(())
}

fn doctor() -> Result<()> {
    println!("Host target: {}", rustc_host()?);
    println!("Default target: {}", default_target()?);

    println!(
        "rustc: {}",
        command_exists("rustc")
    );

    println!(
        "cargo: {}",
        command_exists("cargo")
    );

    println!(
        "rustup: {}",
        command_exists("rustup")
    );

    if cfg!(windows) {
        match find_msys2() {
            Ok(path) => {
                println!("MSYS2: {}", path.display());
            }
            Err(_) => {
                println!("MSYS2: NOT FOUND");
            }
        }
    } else {
        println!(
            "pkg-config: {}",
            command_exists("pkg-config")
        );

        println!(
            "gcc/cc: {}",
            command_exists("gcc") || command_exists("cc")
        );

        println!(
            "make: {}",
            command_exists("make")
        );
    }

    Ok(())
}

fn clean() -> Result<()> {
    let root = project_root()?;
    let native = root.join(".native");

    if native.exists() {
        fs::remove_dir_all(&native)
            .with_context(|| format!("failed to remove {}", native.display()))?;
    }

    println!("Removed {}", native.display());

    Ok(())
}

fn find_msys2() -> Result<PathBuf> {
    if let Some(path) = env::var_os("MSYS2_ROOT") {
        let path = PathBuf::from(path);

        if path.join("usr").join("bin").join("bash.exe").exists() {
            return Ok(path);
        }

        bail!(
            "MSYS2_ROOT is set but does not contain usr/bin/bash.exe: {}",
            path.display()
        );
    }

    let candidates = [
        PathBuf::from(r"C:\msys64"),
        PathBuf::from(r"C:\msys2"),
    ];

    for candidate in candidates {
        if candidate.join("usr").join("bin").join("bash.exe").exists() {
            return Ok(candidate);
        }
    }

    bail!(
        "MSYS2 not found; install MSYS2 or set MSYS2_ROOT"
    )
}

fn windows_to_msys_path(msys_root: &Path, path: &Path) -> Result<PathBuf> {
    /*
     * Prefer native MSYS2 path conversion via cygpath when available.
     * This also handles paths outside the MSYS2 installation correctly.
     */
    let cygpath = msys_root.join("usr").join("bin").join("cygpath.exe");

    let output = Command::new(cygpath)
        .args(["-u", &path.to_string_lossy()])
        .output()
        .context("failed to execute MSYS2 cygpath")?;

    if !output.status.success() {
        bail!(
            "cygpath failed for {}",
            path.display()
        );
    }

    let converted = String::from_utf8(output.stdout)
        .context("cygpath produced invalid UTF-8")?;

    Ok(PathBuf::from(converted.trim()))
}

fn prepend_path(command: &mut Command, path: PathBuf) {
    let existing = env::var_os("PATH").unwrap_or_default();

    let mut new_path = OsString::new();
    new_path.push(path);
    new_path.push(";");
    new_path.push(existing);

    command.env("PATH", new_path);
}

fn command_exists(command: &str) -> bool {
    let probe = if cfg!(windows) {
        "where"
    } else {
        "which"
    };

    Command::new(probe)
        .arg(command)
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .status()
        .map(|status| status.success())
        .unwrap_or(false)
}

fn run_command(mut command: Command) -> Result<()> {
    let status = command
        .status()
        .context("failed to execute command")?;

    if !status.success() {
        bail!("command failed with status {status}");
    }

    Ok(())
}

fn run_shell(program: &str, args: &[&str], cwd: Option<&Path>) -> Result<()> {
    let mut command = Command::new(program);
    command.args(args);

    if let Some(cwd) = cwd {
        command.current_dir(cwd);
    }

    run_command(command)
}

fn shell_quote(path: &Path) -> String {
    shell_quote_str(&path.to_string_lossy())
}

fn shell_quote_str(value: &str) -> String {
    format!("'{}'", value.replace('\'', "'\\''"))
}
