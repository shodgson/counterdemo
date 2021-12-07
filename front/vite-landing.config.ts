import { defineConfig } from "vite";
import { resolve } from "path";
import { minifyHtml, injectHtml } from "vite-plugin-html";

const root = resolve(__dirname, "src/landing");
const outDir = resolve(__dirname, "dist");

export default defineConfig({
  root,
  envDir: "../",
  plugins: [injectHtml()],
  build: {
    outDir,
    rollupOptions: {
      input: resolve(root, "index.html"),
    },
  },
});
