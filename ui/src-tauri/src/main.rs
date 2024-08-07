// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![play])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

#[tauri::command]
async fn play(video: String, sub: String, auth: String) {
    std::process::Command::new("mpv")
        .args(&[
            &format!("--http-header-fields=Authorization: Bearer {auth}"),
            &format!("--sub-file={sub}"),
            &video,
        ])
        .spawn()
        .unwrap();
}
