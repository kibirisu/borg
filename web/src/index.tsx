import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router";
import App from "./App";
import MainFeed from "./components/feed/MainFeed";
import UserProfile from "./components/profile/UserProfile";
import { sampleUsers, type User } from "./components/feed/feedData";

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
      // {
      //   path: "user/:handle",
      //   Component: UserView,
      //   loader: ({ params }) => {
      //     return parseInt(String(params.handle));
      //   },
      // },
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
