import { execFileSync } from "node:child_process";
import path from "node:path";

import legacyPlugin from "@vitejs/plugin-legacy";
import reactPlugin from "@vitejs/plugin-react";
import fs from "fs-extra";
import { type Plugin, defineConfig } from "vite";

const appVersion = execFileSync("git", ["describe", "--match", "v[0-9]*", "--always", "--tags", "--abbrev=8"])?.toString();

const preserveFilesPlugin = (filesToPreserve: string[]): Plugin => {
  return {
    name: "preserve-files",
    apply: "build",
    buildStart() {
      // 在构建开始时将要保留的文件或目录移动到临时位置
      filesToPreserve.forEach((file) => {
        const srcPath = path.resolve(__dirname, file);
        const tempPath = path.resolve(__dirname, `node_modules`, `.tmp`, `build_${file}`);
        if (fs.existsSync(srcPath)) {
          fs.moveSync(srcPath, tempPath, { overwrite: true });
        }
      });
    },
    closeBundle() {
      // 在构建完成后将临时位置的文件或目录移回原来的位置
      filesToPreserve.forEach((file) => {
        const srcPath = path.resolve(__dirname, file);
        const tempPath = path.resolve(__dirname, `node_modules`, `.tmp`, `build_${file}`);
        if (fs.existsSync(tempPath)) {
          fs.moveSync(tempPath, srcPath, { overwrite: true });
        }
      });
    },
  };
};

export default defineConfig({
  define: {
    __APP_VERSION__: JSON.stringify(appVersion),
  },
  plugins: [
    reactPlugin({}),
    legacyPlugin({
      targets: ["defaults", "not IE 11"],
    }),
    preserveFilesPlugin(["dist/.gitkeep"]),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/api": "http://127.0.0.1:8090",
    },
  },
});
