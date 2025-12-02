import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router";
import newClient from "./lib/client";
import type AppState from "./lib/state";
import { createAppContext } from "./lib/state";
import newRouter from "./routes/router";

const client = newClient();
const router = newRouter(client);
const AppContext = createAppContext(client.$api);
const state: AppState = { $api: client.$api };

const rootEl = document.getElementById("root");
if (rootEl) {
  const root = ReactDOM.createRoot(rootEl);
  root.render(
    <React.StrictMode>
      <AppContext value={state}>
        <QueryClientProvider client={client.queryClient}>
          <RouterProvider router={router} />
          <ReactQueryDevtools />
        </QueryClientProvider>
      </AppContext>
    </React.StrictMode>,
  );
}
