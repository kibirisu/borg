import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router";
import App from "./App";
import MainFeed from "./components/feed/MainFeed";
import UserProfile from "./components/profile/UserProfile";
import { sampleUsers, type User } from "./components/feed/feedData";
import type { LoaderFunctionArgs } from "react-router";
import { newClient, type Client } from "./lib/api/client";

const client = newClient();

const loader =
  (client: Client) =>
    async ({ params }: LoaderFunctionArgs) => {
      const foo = client.$api.queryOptions("get", "/api/users/{id}", {
        params: { path: { id: 1 } },
      });
      console.log(foo);
      console.log(params);
      return foo;
    };

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
      {
        path: "foo",
        element: <span>bar</span>,
        loader: loader(client),
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
