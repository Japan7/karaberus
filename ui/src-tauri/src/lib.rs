use tauri::{AppHandle, Emitter, Manager};

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_single_instance::init(single_instance_handler))
        .plugin(tauri_plugin_deep_link::init())
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_updater::Builder::new().build())
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

fn single_instance_handler(app: &AppHandle, argv: Vec<String>, _cwd: String) {
    if let Some(window) = app.get_webview_window("main") {
        if !window.is_focused().unwrap_or_default() {
            window.show().unwrap();
            window.set_focus().unwrap();
        }

        // Handle deep linking on Linux/Windows
        if argv.len() == 2 {
            let arg = argv[1].to_string();
            window.emit("deep-link", arg).unwrap();
        }
    }
}
