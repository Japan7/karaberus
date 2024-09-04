import { defineConfig, type PluginOption } from "vite";
import pages from "vite-plugin-pages";
import solid from "vite-plugin-solid";
import { viteStaticCopy } from "vite-plugin-static-copy";

export default defineConfig(({ mode }) => {
  const isTauri = mode.indexOf("tauri") === 0;

  const plugins: PluginOption[] = [solid(), pages()];

  if (!isTauri) {
    plugins.push(
      viteStaticCopy({
        targets: [
          {
            src: "node_modules/@fontsource/amaranth/files/*",
            dest: "amaranth",
          },
          {
            src: "node_modules/libass-wasm/dist/js/*",
            dest: "libass-wasm",
          },
        ],
      }),
    );
  }

  // @ts-expect-error process is a nodejs global
  const host: string = process.env.TAURI_DEV_HOST;

  return {
    plugins,
    // Vite options tailored for Tauri development and only applied in `tauri dev` or `tauri build`
    //
    // 1. prevent vite from obscuring rust errors
    clearScreen: !isTauri,
    // 2. tauri expects a fixed port, fail if that port is not available
    server: {
      strictPort: true,
      host,
      hmr: { host },
      watch: {
        // 3. tell vite to ignore watching `src-tauri`
        ignored: ["**/src-tauri/**"],
      },
      proxy: {
        "/api": "http://localhost:8888",
      },
    },
  };
});
