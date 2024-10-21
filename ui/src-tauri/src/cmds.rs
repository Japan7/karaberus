use std::env;

use serde::Serialize;
use tauri::{ipc::Channel, AppHandle, Manager, State};
use tauri_plugin_shell::{process::CommandEvent, ShellExt};
use tauri_plugin_store::StoreExt;

use crate::{
    mpv::{LoadFile, Mpv},
    AppState, STORE_BIN,
};

#[derive(Clone, Serialize)]
#[serde(rename_all = "camelCase", tag = "event", content = "data")]
pub enum LogEvent {
    #[serde(rename_all = "camelCase")]
    Stdout(String),
    #[serde(rename_all = "camelCase")]
    Stderr(String),
}

#[tauri::command]
pub async fn play_mpv(
    app_handle: AppHandle,
    state: State<'_, AppState>,
    video: Option<String>,
    inst: Option<String>,
    sub: Option<String>,
    title: String,
    on_log: Channel<LogEvent>,
) -> tauri::Result<()> {
    let mpv = state
        .lock()
        .await
        .mpv
        .get_or_insert_with(|| {
            let socket = get_mpv_socket(&app_handle);
            let token = get_player_token(&app_handle);
            start_mpv(&app_handle, socket, token, on_log)
        })
        .clone();
    add_to_mpv_playlist(&mpv, video, inst, sub, title).await;
    Ok(())
}

fn start_mpv(
    app_handle: &AppHandle,
    socket: String,
    token: String,
    on_log: Channel<LogEvent>,
) -> Mpv {
    let mut mpv = app_handle.shell().command("mpv");
    mpv = mpv.args([
        "--idle=once",
        "--quiet",
        "--save-position-on-quit=no",
        &format!("--input-ipc-server={socket}"),
        &format!("--http-header-fields=Authorization: Bearer {token}"),
    ]);
    let (mut rx, _) = mpv.spawn().unwrap();

    let app_handle = app_handle.clone();
    tauri::async_runtime::spawn(async move {
        let state = app_handle.state::<AppState>();
        while let Some(event) = rx.recv().await {
            match event {
                CommandEvent::Terminated(_) => {
                    state.lock().await.mpv.take();
                }
                CommandEvent::Stdout(line) => {
                    let line = String::from_utf8(line).unwrap();
                    print!("{}", &line);
                    let _ = on_log.send(LogEvent::Stdout(line));
                }
                CommandEvent::Stderr(line) => {
                    let line = String::from_utf8(line).unwrap();
                    eprint!("{}", &line);
                    let _ = on_log.send(LogEvent::Stderr(line));
                }
                _ => {}
            }
        }
    });

    Mpv { socket }
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

fn get_player_token(app_handle: &AppHandle) -> String {
    let store = app_handle.store(STORE_BIN).unwrap();
    store
        .get("player_token")
        .unwrap()
        .as_str()
        .unwrap()
        .to_string()
}

async fn add_to_mpv_playlist(
    mpv: &Mpv,
    video: Option<String>,
    inst: Option<String>,
    sub: Option<String>,
    title: String,
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

    loadfile
        .options
        .insert("force-media-title".to_string(), title);

    mpv.loadfile(loadfile).await;
}
