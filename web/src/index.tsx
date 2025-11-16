import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router";
import App from "./App";
import UserProfile from "./components/profile/UserProfile";
import { sampleUsers, type User } from "./components/feed/feedData";
import MainFeed from "./components/feed/MainFeed";

const router = createBrowserRouter([
  {
    path: "/",
    Component: App,
    children: [
      { index: true, Component: MainFeed },
      {
        path: "profile/:handle",
        Component: UserProfile,
        loader: ({ params }) => {
          const handleParam = params.handle ? `@${params.handle}` : undefined;
          return sampleUsers.find((u: User) => u.username === handleParam);
        },
      },
    ],
  },
]);

const rootEl = document.getElementById("root");
if (rootEl) {
  const root = ReactDOM.createRoot(rootEl);
  root.render(
    <React.StrictMode>
      <RouterProvider router={router} />
    </React.StrictMode>,
  );
}
