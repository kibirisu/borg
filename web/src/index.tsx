import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router";
import App from "./App";
import Login from "./components/auth/Login";
import Register from "./components/auth/Register";
import UserProfile from "./components/profile/UserProfile";
import { sampleUsers, type User } from "./components/feed/feedData";

const router = createBrowserRouter([
  {
    path: "/",
    Component: App,
  },
  {
    path: "/login",
    Component: Login,
  },
  {
    path: "/register",
    Component: Register,
  },
  {
    path: "/profile/:handle",
    Component: UserProfile,
    loader: ({ params }) => {
      const handleParam = params.handle ? `@${params.handle}` : undefined;
      return sampleUsers.find((u: User) => u.username === handleParam);
    },
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
