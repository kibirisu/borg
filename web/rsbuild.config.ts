import { defineConfig } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";

export default defineConfig({
  output: {
    cleanDistPath: {
      keep: [/dist\/docs.html/],
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
      },
      "/auth": {
        target: "http://localhost:8080",
      },
    },
  },
  plugins: [pluginReact()],
});
