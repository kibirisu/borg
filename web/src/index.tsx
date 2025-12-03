import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router";
import newClient, { checkToken } from "./lib/client";
import { type AppState, AppStateProvider } from "./lib/state";
import newRouter from "./routes/router";

const client = newClient();
const router = newRouter(client);
const state: AppState = { $api: client.$api, username: checkToken() };

const rootEl = document.getElementById("root");
if (rootEl) {
  const root = ReactDOM.createRoot(rootEl);
  root.render(
    <React.StrictMode>
      <AppStateProvider state={state}>
        <QueryClientProvider client={client.queryClient}>
          <RouterProvider router={router} />
          <ReactQueryDevtools />
        </QueryClientProvider>
      </AppStateProvider>
    </React.StrictMode>,
  );
}
