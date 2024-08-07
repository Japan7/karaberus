export const isTauri = () => "__TAURI__" in window;
export const isBrowser = () => !isTauri();
