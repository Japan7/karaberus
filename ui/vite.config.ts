import { defineConfig, type PluginOption } from "vite";
import pages from "vite-plugin-pages";
import solid from "vite-plugin-solid";
import { viteStaticCopy } from "vite-plugin-static-copy";

export default defineConfig(({ mode }) => {
  const plugins: PluginOption[] = [solid(), pages()];

  if (mode.indexOf("tauri") !== 0) {
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

  return {
    plugins,
    server: {
      proxy: {
        "/api": "http://localhost:8888",
      },
    },
  };
});
