use std::{sync::Arc, thread};

use mpvipc::{Mpv, MpvCommand, MpvDataType, PlaylistAddOptions};
use tauri::{
    async_runtime::{block_on, Mutex},
    AppHandle, State,
};
use tauri_plugin_shell::{process::CommandEvent, ShellExt};

#[derive(Debug)]
struct PlaylistEntry {
    video: Option<String>,
    inst: Option<String>,
    sub: Option<String>,
}

#[derive(Default)]
struct AppStateInner {
    mpv_started: bool,
    playlist: Vec<PlaylistEntry>,
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

        spawn_mpv_ipc_control(state.clone());

        while let Some(event) = rx.recv().await {
            match event {
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

        mpv.observe_property(1, "eof-reached").unwrap();

        loop {
            let Ok(event) = mpv.event_listen() else {
                println!("Error listening to mpv event, exiting");
                break;
            };
            println!("{:?}", event);

            match event {
                mpvipc::Event::StartFile => {
                    let mut state = block_on(state.lock());
                    let entry = state.playlist.remove(0);
                    println!("Adding additional tracks for: {:?}", entry);
                    if let (Some(_), Some(inst)) = (entry.video, entry.inst) {
                        mpv.run_command_raw("audio-add", &[&inst, "auto"]).unwrap();
                    }
                    if let Some(sub) = entry.sub {
                        mpv.run_command_raw("sub-add", &[&sub, "select"]).unwrap();
                    }
                }

                mpvipc::Event::PropertyChange {
                    id: 1,
                    property: mpvipc::Property::Unknown { name: _, data },
                } => {
                    if let MpvDataType::Bool(false) = data {
                        continue;
                    }

                    let mut state = block_on(state.lock());

                    if state.playlist.is_empty() {
                        println!("Playlist empty, exiting");
                        mpv.kill().unwrap();
                        state.mpv_started = false;
                        break;
                    }

                    let entry = &state.playlist[0];
                    println!("Loading next entry media: {:?}", entry);

                    if let (Some(file), None) | (None, Some(file)) = (&entry.video, &entry.inst) {
                        mpv.run_command(MpvCommand::LoadFile {
                            file: file.to_string(),
                            option: PlaylistAddOptions::Replace,
                        })
                        .unwrap();
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
    let mut app_state = state.lock().await;
    app_state.playlist.push(PlaylistEntry { video, inst, sub });
    if !app_state.mpv_started {
        start_mpv(app_handle, state.inner().clone(), auth);
    }
    Ok(())
}
