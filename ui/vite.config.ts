import path from "path";
import react from "@vitejs/plugin-react";
import { defineConfig, ConfigEnv, UserConfig, loadEnv } from "vite";
import viteCompression from "vite-plugin-compression";
import { wrapperEnv } from "./src/lib/utils";

export default defineConfig(({ mode }: ConfigEnv): UserConfig => {
  const env = loadEnv(mode, process.cwd());
  const viteEnv = wrapperEnv(env);

  return {
    plugins: [
      react(),
      viteEnv.VITE_BUILD_GZIP &&
        viteCompression({
          verbose: true,
          disable: false,
          threshold: 10240,
          algorithm: "gzip",
          ext: ".gz",
        }),
    ],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: {
      port: viteEnv.VITE_PORT,
      host: "0.0.0.0",
      proxy: {
        // "/api": "http://127.0.0.1:8090",
        "/api": {
          target: viteEnv.VITE_API_URL,
          changeOrigin: true,
          ws: true,
        },
      },
    },
    esbuild: {
      pure: viteEnv.VITE_DROP_LOG ? ["console.log", "debugger"] : [],
    },
    build: {
      outDir: "../internal/web/dist",
      minify: "esbuild",
      rollupOptions: {
        output: {
          chunkFileNames: "assets/js/[name]-[hash].js",
          entryFileNames: "assets/js/[name]-[hash].js",
          assetFileNames: "assets/[ext]/[name]-[hash].[ext]",
        },
      },
    },
  };
});
