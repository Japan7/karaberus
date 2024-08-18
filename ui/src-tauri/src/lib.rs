mod mpv;

use std::{env, sync::Arc};

use mpv::{LoadFile, Mpv};
use tauri::{async_runtime::Mutex, AppHandle, Emitter, Manager, State};
use tauri_plugin_shell::{process::CommandEvent, ShellExt};

#[derive(Default)]
struct AppStateInner {
    mpv: Option<Mpv>,
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

#[tauri::command]
async fn play_mpv(
    app_handle: AppHandle,
    state: State<'_, AppState>,
    video: Option<String>,
    inst: Option<String>,
    sub: Option<String>,
    token: String,
) -> Result<(), ()> {
    let mut app_state = state.lock().await;
    let mpv = app_state.mpv.get_or_insert_with(|| {
        start_mpv(
            app_handle.clone(),
            state.inner().clone(),
            get_mpv_socket(&app_handle),
            token,
        )
    });
    tauri::async_runtime::spawn(add_to_mpv_playlist(mpv.clone(), video, inst, sub));
    Ok(())
}

fn start_mpv(app_handle: AppHandle, state: AppState, socket: String, token: String) -> Mpv {
    let mut mpv = app_handle.shell().command("mpv");
    mpv = mpv.args([
        "--idle=once",
        "--quiet",
        "--save-position-on-quit=no",
        &format!("--input-ipc-server={socket}"),
        &format!("--http-header-fields=Authorization: Bearer {token}"),
    ]);
    let (mut rx, _) = mpv.spawn().unwrap();

    tauri::async_runtime::spawn(async move {
        while let Some(event) = rx.recv().await {
            match event {
                CommandEvent::Terminated(_) => {
                    state.lock().await.mpv = None;
                }
                CommandEvent::Stdout(line) => {
                    let line = String::from_utf8(line).unwrap();
                    print!("{}", &line);
                    let _ = app_handle.emit("mpv-stdout", &line);
                }
                CommandEvent::Stderr(line) => {
                    let line = String::from_utf8(line).unwrap();
                    eprint!("{}", &line);
                    let _ = app_handle.emit("mpv-stderr", &line);
                }
                _ => {}
            }
        }
    });

    mpv::Mpv { socket }
}

fn get_mpv_socket(app_handle: &AppHandle) -> String {
    let env_var = env::var("KARABERUS_MPV_SOCKET");
    match (env_var, cfg!(windows)) {
        (Ok(env_var), _) => env_var,
        (_, true) => "karaberus-mpv".to_string(),
        (_, false) => {
            let base_directory = if cfg!(target_os = "linux") {
                tauri::path::BaseDirectory::Runtime
            } else if cfg!(target_os = "macos") {
                tauri::path::BaseDirectory::Temp
            } else {
                tauri::path::BaseDirectory::LocalData
            };
            app_handle
                .path()
                .resolve("karaberus-mpv.sock", base_directory)
                .unwrap()
                .to_str()
                .unwrap()
                .to_string()
        }
    }
}

async fn add_to_mpv_playlist(
    mpv: Mpv,
    video: Option<String>,
    inst: Option<String>,
    sub: Option<String>,
) {
    let mut loadfile = LoadFile::default();

    if let Some(video) = video.as_deref() {
        loadfile.url = video.to_string();
        loadfile.flags = "append-play".to_string();
    }

    loadfile.options.insert("aid".to_string(), "1".to_string());

    if let Some(inst) = inst.as_deref() {
        loadfile
            .options
            .insert("audio-file".to_string(), inst.to_string());
    }

    if let Some(sub) = sub.as_deref() {
        loadfile
            .options
            .insert("sub-file".to_string(), sub.to_string());
    }

    mpv.loadfile(loadfile).await;
}
