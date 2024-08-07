// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![play])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

#[tauri::command]
async fn play(auth: String, video: Option<String>, inst: Option<String>, sub: Option<String>) {
    let mut mpv = std::process::Command::new("mpv");
    mpv.arg(&format!(
        "--http-header-fields=Authorization: Bearer {auth}"
    ));
    if let Some(sub) = sub {
        mpv.arg(&format!("--sub-file={sub}"));
    }
    match (video, inst) {
        (Some(video), Some(inst)) => {
            mpv.arg(&format!("--external-file={inst}"));
            mpv.arg(&video);
        }
        (Some(video), None) => {
            mpv.arg(&video);
        }
        (None, Some(inst)) => {
            mpv.arg(&inst);
        }
        _ => {
            return;
        }
    }
    mpv.spawn().expect("failed to spawn mpv");
}
