use std::path::PathBuf;
use std::process::{Command, Child};
use std::sync::Mutex;
use tauri::Manager;

fn find_server() -> PathBuf {
    if let Ok(exe) = std::env::current_exe() {
        let p = exe.parent().unwrap().join("binaries").join("mox-server.exe");
        if p.exists() { return p; }
    }
    let dev = PathBuf::from("../binaries/mox-server-x86_64-pc-windows-msvc.exe");
    if dev.exists() { return dev.canonicalize().unwrap_or(dev); }
    PathBuf::from("mox-server.exe")
}

struct AppState {
    server: Mutex<Option<Child>>,
}

pub fn run() {
    tauri::Builder::default()
        .setup(|app| {
            let path = find_server();
            let child = Command::new(&path).spawn()
                .expect(&format!("Failed to start backend at {:?}", path));
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