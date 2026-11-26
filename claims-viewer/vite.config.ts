import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

// base "./" so dist/ works from any path, including a directory listing
// served by `python3 -m http.server`.
// The dev server and the preview listen on every address, so a viewer run on
// a machine in a tailnet or a LAN is reachable at http://<that address>:5173/.
export default defineConfig({
  plugins: [react()],
  base: "./",
  server: { host: true, port: 5173, strictPort: true },
  preview: { host: true, port: 5173, strictPort: true },
  test: { environment: "node", include: ["src/**/*.test.ts"] },
});
