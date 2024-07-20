import suidPlugin from "@suid/vite-plugin";
import { defineConfig } from "vite";
import mkcert from "vite-plugin-mkcert";
import solid from "vite-plugin-solid";

export default defineConfig({
  plugins: [solid(), suidPlugin(), mkcert()],
  server: {
    proxy: {
      "/api": "http://localhost:8888",
    },
  },
});
