import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router";
import App from "./App";
import MainFeed from "./components/feed/MainFeed";
import { newClient } from "./lib/api/client";
import { QueryClientProvider } from "@tanstack/react-query";
import User, { loader } from "./components/profile/UserProfile";

const client = newClient();

const router = createBrowserRouter([
  {
    path: "/",
    Component: App,
    children: [
      { index: true, Component: MainFeed },
      {
        path: "profile/:handle",
        Component: User,
        loader: loader(client),
      },
    ],
  },
]);

const rootEl = document.getElementById("root");
if (rootEl) {
  const root = ReactDOM.createRoot(rootEl);
  root.render(
    <React.StrictMode>
      <QueryClientProvider client={client.queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </React.StrictMode>,
  );
}
