import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { resolve } from "path";

const root = resolve(__dirname, "src", "app");
const outDir = resolve(__dirname, "dist", "app");

export default defineConfig({
  root,
  base: "/app/",
  plugins: [vue()],
  resolve: {
    alias: [{ find: "@", replacement: resolve(root) }],
  },
  build: {
    outDir,
    rollupOptions: {
      input: resolve(root, "index.html"),
    },
  },
  server: {
    host: true,
    port: 3030,
  },
  envDir: "../",
});
