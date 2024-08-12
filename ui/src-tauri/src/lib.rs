pub mod mpv;

use mpv::*;
use std::{sync::Arc};

use tauri::{
    async_runtime::{Mutex},
    AppHandle, State,
};
use tauri_plugin_shell::{process::CommandEvent, ShellExt};

#[derive(Default)]
struct AppStateInner {
    mpv_started: bool,
}

type AppState = Arc<Mutex<AppStateInner>>;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_os::init())
        .plugin(tauri_plugin_updater::Builder::new().build())
        .manage(Arc::new(Mutex::new(AppStateInner::default())))
        .invoke_handler(tauri::generate_handler![play_mpv])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

fn get_socket() -> String {
    if cfg!(windows) {
        return "karaberus-mpv".to_string();
    } else {
        return "/tmp/karaberus-mpv.sock".to_string();
    }
}

fn start_mpv(app_handle: AppHandle, state: AppState, auth: String) {
    tauri::async_runtime::spawn(async move {
        let mut mpv = app_handle.shell().command("mpv");
        mpv = mpv.args([
            "--idle=once",
            "--quiet",
            "--save-position-on-quit=no",
            &format!("--input-ipc-server={}", get_socket()),
            &format!("--http-header-fields=Authorization: Bearer {auth}"),
        ]);

        let (mut rx, mut _child) = mpv.spawn().unwrap();
        state.lock().await.mpv_started = true;

        while let Some(event) = rx.recv().await {
            match event {
                CommandEvent::Terminated(_) => {
                    state.lock().await.mpv_started = false;
                }
                CommandEvent::Stdout(line) => {
                    print!("{}", String::from_utf8(line).unwrap());
                }
                CommandEvent::Stderr(line) => {
                    eprint!("{}", String::from_utf8(line).unwrap());
                }
                _ => {}
            }
        }
    });
}

#[tauri::command]
async fn play_mpv(
    app_handle: AppHandle,
    state: State<'_, AppState>,
    auth: String,
    video: Option<String>,
    inst: Option<String>,
    sub: Option<String>,
) -> Result<(), ()> {
    let app_state = state.lock().await;

    if !app_state.mpv_started {
        start_mpv(app_handle, state.inner().clone(), auth);
    }

    tauri::async_runtime::spawn(add_to_mpv_playlist(video, inst, sub));

    Ok(())
}

async fn add_to_mpv_playlist(
    video: Option<String>,
    inst: Option<String>,
    sub: Option<String>,
) {
    let mpv = Mpv {
        socket: get_socket(),
    };

    let mut loadfile = LoadFile::default();

    if let Some(video) = video.as_deref() {
        loadfile.url = video.to_string();
        loadfile.flags = "append-play".to_string();
    }

    loadfile.options.insert("aid".to_string(), "1".to_string());

    if let Some(inst) = inst.as_deref() {
        loadfile.options.insert("audio-file".to_string(), inst.to_string());
    }

    if let Some(sub) = sub.as_deref() {
        loadfile.options.insert("sub-file".to_string(), sub.to_string());
    }

    mpv.loadfile(loadfile).await;
}
