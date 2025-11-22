import { defineConfig } from "@rsbuild/core";
import { pluginReact } from "@rsbuild/plugin-react";

export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
        secure: false,
        timeout: 10000,
        // Don't fail on connection errors - let the app handle it
        onError: (err, req, res) => {
          console.error("Proxy error for", req.url, ":", err.message);
          if (res && !res.headersSent) {
            res.writeHead(503, {
              "Content-Type": "text/plain",
            });
            res.end("Backend server is not available. Please make sure it's running on port 8080.");
          }
        },
      },
    },
  },
  plugins: [pluginReact()],
});
