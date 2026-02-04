import { defineConfig } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";
import { fileURLToPath } from "node:url";
import { dirname, resolve } from "node:path";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

export default defineConfig({
  plugins: [pluginReact()],
  source: {
    entry: {
      index: "./src/main.tsx",
    },
  },
  resolve: {
    alias: {
      "@": resolve(__dirname, "./src"),
      "@wailsjs": resolve(__dirname, "./wailsjs"),
    },
  },
  html: {
    template: "./index.html",
  },
  output: {
    distPath: {
      root: "dist",
    },
    assetPrefix: "/",
  },
  server: {
    port: 34115,
    strictPort: false,
  },
});
