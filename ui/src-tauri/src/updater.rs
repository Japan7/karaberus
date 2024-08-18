//! https://github.com/tauri-apps/tauri/blob/1.6/core/tauri/src/updater/mod.rs

use tauri::AppHandle;
use tauri_plugin_dialog::DialogExt;
use tauri_plugin_updater::{Update, UpdaterExt};

pub async fn check_update_with_dialog(app_handle: AppHandle) {
    let updater = app_handle.updater().unwrap();
    if let Some(update) = updater.check().await.unwrap() {
        prompt_for_install(app_handle, update);
    }
}

fn prompt_for_install(app_handle: AppHandle, update: Update) {
    let package_info = app_handle.package_info().clone();
    let app_name = package_info.name;

    let title = format!(r#"A new version of {app_name} is available!"#);
    let message = format!(
        r#"{app_name} {} is now available -- you have {}.

Would you like to install it now?

Release Notes:
{}"#,
        update.version,
        update.current_version,
        update.body.clone().unwrap_or_default()
    );

    app_handle
        .dialog()
        .message(message)
        .title(title)
        .ok_button_label("Yes")
        .cancel_button_label("No")
        .show(|should_update| {
            if should_update {
                tauri::async_runtime::spawn(install_update(app_handle, update));
            }
        });
}

async fn install_update(app_handle: AppHandle, update: Update) {
    let title = "Ready to Restart";
    let message = "The installation was successful, do you want to restart the application now?";
    update
        .download_and_install(
            |_, _| {},
            || {
                app_handle
                    .dialog()
                    .message(message)
                    .title(title)
                    .ok_button_label("Yes")
                    .cancel_button_label("No")
                    .show(move |should_exit| {
                        if should_exit {
                            app_handle.restart();
                        }
                    });
            },
        )
        .await
        .unwrap();
}
