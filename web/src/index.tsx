import React from "react";
import ReactDOM from "react-dom/client";
import {
  createBrowserRouter,
  RouterProvider,
  useLoaderData,
} from "react-router";
import App from "./App";
import MainFeed from "./components/feed/MainFeed";
import UserProfile from "./components/profile/UserProfile";
import { sampleUsers, type User } from "./components/feed/feedData";
import type { LoaderFunctionArgs } from "react-router";
import { newClient, type Client } from "./lib/api/client";
import { QueryClientProvider, useSuspenseQuery } from "@tanstack/react-query";

const client = newClient();

const Component = () => {
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data } = useSuspenseQuery(opts);
  return <>{data.username}</>;
};

const loader =
  (client: Client) =>
    async ({ params }: LoaderFunctionArgs) => {
      const userId = parseInt(String(params.handle));
      const opts = client.$api.queryOptions("get", "/api/users/{id}", {
        params: { path: { id: userId } },
      });
      await client.queryClient.ensureQueryData(opts);
      return { userId: userId, opts: opts };
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
        path: "foo/:handle",
        Component: Component,
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
      <QueryClientProvider client={client.queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </React.StrictMode>,
  );
}
