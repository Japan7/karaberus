use std::{sync::Arc, thread};

use mpvipc::{Mpv, PlaylistAddOptions};
use tauri::{
    async_runtime::{block_on, Mutex},
    State,
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

fn start_mpv(
    app_handle: tauri::AppHandle,
    state: AppState,
    auth: String,
    video: Option<String>,
    inst: Option<String>,
    sub: Option<String>,
) {
    tauri::async_runtime::spawn(async move {
        let mut mpv = app_handle.shell().command("mpv");
        mpv = mpv.args([
            "--keep-open",
            "--save-position-on-quit=no",
            "--input-ipc-server=/tmp/mpv.sock",
            &format!("--http-header-fields=Authorization: Bearer {auth}"),
        ]);
        if let Some(sub) = sub {
            mpv = mpv.arg(&format!("--sub-file={sub}"));
        }
        mpv = match (video, inst) {
            (Some(video), Some(inst)) => mpv.args([&format!("--external-file={inst}"), &video]),
            (Some(video), None) => mpv.arg(&video),
            (None, Some(inst)) => mpv.arg(&inst),
            _ => mpv,
        };

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
                CommandEvent::Terminated(_) => {
                    state.lock().await.mpv_started = false;
                }
                _ => {}
            }
        }
    });
}

fn spawn_mpv_ipc_control(state: AppState) {
    thread::spawn(move || {
        let mut mpv = loop {
            match Mpv::connect("/tmp/mpv.sock") {
                Ok(mpv) => break mpv,
                _ => {}
            }
        };
        println!("Connected to mpv");

        mpv.observe_property(1, "eof-reached").unwrap();
        let mut started_once = false;

        loop {
            let Ok(event) = mpv.event_listen() else {
                println!("Error listening to mpv event, exiting");
                break;
            };

            println!("{:?}", event);
            match event {
                mpvipc::Event::StartFile => {
                    started_once = true;
                }

                mpvipc::Event::PropertyChange { id: 1, property: _ } => {
                    if !started_once {
                        continue;
                    }

                    let mut state = block_on(state.lock());

                    if state.playlist.is_empty() {
                        println!("Playlist empty, exiting");
                        continue;
                    }

                    let entry = state.playlist.remove(0);
                    println!("Playing next entry: {:?}", entry);

                    mpv.run_command(mpvipc::MpvCommand::LoadFile {
                        file: entry.video.unwrap(),
                        option: PlaylistAddOptions::Append,
                    })
                    .unwrap();
                    if let Some(inst) = entry.inst {
                        mpv.run_command(mpvipc::MpvCommand::LoadFile {
                            file: inst,
                            option: PlaylistAddOptions::Append,
                        })
                        .unwrap();
                    }
                    if let Some(sub) = entry.sub {
                        mpv.run_command(mpvipc::MpvCommand::LoadFile {
                            file: sub,
                            option: PlaylistAddOptions::Append,
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
    app_handle: tauri::AppHandle,
    state: State<'_, AppState>,
    auth: String,
    video: Option<String>,
    inst: Option<String>,
    sub: Option<String>,
) -> Result<(), ()> {
    if !state.lock().await.mpv_started {
        start_mpv(app_handle, state.inner().clone(), auth, video, inst, sub);
    } else {
        let mut state = state.lock().await;
        state.playlist.push(PlaylistEntry { video, inst, sub });
    }
    Ok(())
}
