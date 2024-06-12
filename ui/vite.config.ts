import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';
// import devtools from 'solid-devtools/vite';

export default defineConfig(({ mode }) => {
  let api_endpoint = (mode == "production") ? "/api" : process.env.KARABERUS_API_ENDPOINT ?? "http://localhost:8888/api"

  return {
    plugins: [
      /* 
      Uncomment the following line to enable solid-devtools.
      For more info see https://github.com/thetarnav/solid-devtools/tree/main/packages/extension#readme
      */
      // devtools(),
      solidPlugin(),
    ],
    define: {
      "import.meta.env.API_ENDPOINT": JSON.stringify(api_endpoint),
    },
    server: {
      port: 3000,
    },
    build: {
      target: 'esnext',
    },
  }
});
