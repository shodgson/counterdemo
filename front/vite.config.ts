import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { resolve } from "path";

const root = resolve(__dirname, "src");
const outDir = resolve(__dirname, "dist");

export default defineConfig({
  root,
  plugins: [vue()],
  resolve: {
    alias: {
      "@/": `${root}/app/`,
    },
  },
  build: {
    outDir,
    emptyOutDir: true,
    rollupOptions: {
      input: {
        main: resolve(root, "", "index.html"),
        app: resolve(root, "app", "index.html"),
      },
    },
  },
  envDir: "../",
});
