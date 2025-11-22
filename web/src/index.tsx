import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router";
import newClient from "./lib/api/client";
import { QueryClientProvider } from "@tanstack/react-query";
import { newRouter } from "./routes/router";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

const client = newClient();
const router = newRouter(client);

const rootEl = document.getElementById("root");
if (rootEl) {
  const root = ReactDOM.createRoot(rootEl);
  root.render(
    <React.StrictMode>
      <QueryClientProvider client={client.queryClient}>
        <RouterProvider router={router} />
        <ReactQueryDevtools />
      </QueryClientProvider>
    </React.StrictMode>,
  );
}
