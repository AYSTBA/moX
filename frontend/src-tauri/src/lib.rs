use std::path::PathBuf;
use std::process::{Command, Child};
use std::sync::Mutex;
use std::time::Duration;
use tauri::Manager;
#[cfg(windows)]
use std::os::windows::process::CommandExt;

const CREATE_NO_WINDOW: u32 = 0x08000000;

fn find_server() -> PathBuf {
    if let Ok(exe) = std::env::current_exe() {
        let dir = exe.parent().unwrap();
        let along = dir.join("mox-server.exe");
        if along.exists() { return along; }
        let bundled = dir.join("binaries").join("mox-server.exe");
        if bundled.exists() { return bundled; }
    }
    for p in &[
        "../binaries/mox-server-x86_64-pc-windows-msvc.exe",
        "../binaries/mox-server.exe",
        "binaries/mox-server-x86_64-pc-windows-msvc.exe",
        "binaries/mox-server.exe",
    ] {
        let pb = PathBuf::from(p);
        if pb.exists() { return pb.canonicalize().unwrap_or(pb); }
    }
    PathBuf::from("mox-server.exe")
}

fn kill_stale_backend() {
    // Kill any leftover backend process that may hold port 3099
    for name in &["mox-server.exe", "mox-backend.exe"] {
        let _ = Command::new("taskkill")
            .args(["/F", "/IM", name])
            .output();
    }
}

fn wait_for_backend_ready() {
    // Poll port 3099 up to 3 seconds
    for _ in 0..30 {
        if std::net::TcpStream::connect("127.0.0.1:3099").is_ok() {
            return;
        }
        std::thread::sleep(Duration::from_millis(100));
    }
}

struct AppState {
    server: Mutex<Option<Child>>,
}

impl Drop for AppState {
    fn drop(&mut self) {
        if let Ok(mut guard) = self.server.lock() {
            if let Some(mut child) = guard.take() {
                let _ = child.kill();
                let _ = child.wait();
            }
        }
    }
}

fn spawn_backend(path: &PathBuf) -> Child {
    kill_stale_backend();
    let mut cmd = Command::new(path);
    #[cfg(windows)]
    cmd.creation_flags(CREATE_NO_WINDOW);
    let mut child = cmd.spawn()
        .expect(&format!("Failed to start backend at {:?}", path));
    // Small delay then verify process still alive
    std::thread::sleep(Duration::from_millis(300));
    match child.try_wait() {
        Ok(Some(_)) => panic!("Backend exited immediately - port 3099 may be in use"),
        _ => {}
    }
    child
}

pub fn run() {
    tauri::Builder::default()
        .setup(|app| {
            let path = find_server();
            eprintln!("MOX backend path: {:?}", path);
            let child = spawn_backend(&path);
            wait_for_backend_ready();
            eprintln!("MOX backend ready on :3099");
            app.manage(AppState { server: Mutex::new(Some(child)) });
            Ok(())
        })
        .on_window_event(|window, event| {
            if let tauri::WindowEvent::Destroyed = event {
                if let Some(state) = window.try_state::<AppState>() {
                    if let Ok(mut guard) = state.server.lock() {
                        if let Some(mut child) = guard.take() {
                            let _ = child.kill();
                            let _ = child.wait();
                        }
                    }
                }
            }
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
