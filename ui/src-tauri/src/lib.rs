mod cmds;
mod mpv;

use std::{env, sync::Arc};

use mpv::Mpv;
use tauri::{async_runtime::Mutex, Wry};
use tauri_plugin_store::StoreCollection;

#[derive(Default)]
struct AppStateInner {
    mpv: Option<Mpv>,
}

type AppState = Arc<Mutex<AppStateInner>>;

type AppStore = StoreCollection<Wry>;

const STORE_BIN: &str = "store.bin";

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_store::Builder::new().build())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_os::init())
        .plugin(tauri_plugin_updater::Builder::new().build())
        .manage(Arc::new(Mutex::new(AppStateInner::default())))
        .invoke_handler(tauri::generate_handler![cmds::play_mpv])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
