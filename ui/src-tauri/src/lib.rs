use std::{sync::Arc, thread};

use mpvipc::{Mpv, MpvDataType};
use tauri::{
    async_runtime::{block_on, Mutex},
    AppHandle, State,
};
use tauri_plugin_shell::{process::CommandEvent, ShellExt};

#[derive(Default)]
struct AppStateInner {
    mpv_started: bool,
    playback_started: bool,
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

fn start_mpv(app_handle: AppHandle, state: AppState, auth: String) {
    tauri::async_runtime::spawn(async move {
        let mut mpv = app_handle.shell().command("mpv");
        mpv = mpv.args([
            "--idle",
            "--quiet",
            "--save-position-on-quit=no",
            "--input-ipc-server=/tmp/mpv.sock",
            &format!("--http-header-fields=Authorization: Bearer {auth}"),
        ]);

        let (mut rx, mut _child) = mpv.spawn().unwrap();
        state.lock().await.mpv_started = true;
        state.lock().await.playback_started = false;

        spawn_mpv_ipc_control(state.clone());

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

fn spawn_mpv_ipc_control(state: AppState) {
    thread::spawn(move || {
        let mut mpv = loop {
            if let Ok(mpv) = Mpv::connect("/tmp/mpv.sock") {
                break mpv;
            }
        };
        println!("Connected to mpv");

        mpv.observe_property(1, "idle-active").unwrap();

        loop {
            let Ok(event) = mpv.event_listen() else {
                println!("Error listening to mpv event, exiting");
                break;
            };
            println!("{:?}", event);

            match event {
                mpvipc::Event::PropertyChange {
                    id: 1,
                    property: mpvipc::Property::Unknown { name: _, data },
                } => {
                    let mut state = block_on(state.lock());

                    if let MpvDataType::Bool(false) = data {
                        state.playback_started = true;
                    } else if state.playback_started {
                        println!("Playlist empty, exiting");
                        mpv.kill().unwrap();
                        break;
                    }
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
    let mpv = loop {
        if let Ok(mpv) = Mpv::connect("/tmp/mpv.sock") {
            break mpv;
        }
    };

    let mut loadfile_params = Vec::new();

    if let Some(video) = video.as_deref() {
        loadfile_params.push(video);
        loadfile_params.push("append-play");
        loadfile_params.push("-1");
    }

    let mut options_params: String = "aid=1,".to_string();
    if let Some(inst) = inst.as_deref() {
        options_params = format!("{options_params}audio-file={inst},");
    }

    if let Some(sub) = sub.as_deref() {
        options_params = format!("{options_params}sub-file={sub},");
    }

    loadfile_params.push(options_params.as_str());

    mpv.run_command_raw("loadfile", &loadfile_params).unwrap();
}
