import { createBrowserRouter } from "react-router";
import App from "../App";
import MainFeed from "../components/feed/MainFeed";
import User, { loader as userLoader } from "../components/profile/UserProfile";
import type { Client } from "../lib/api/client";
import Feed, { loader as feedLoader } from "../components/common/Feed";

export const newRouter = (client: Client) =>
  createBrowserRouter([
    {
      path: "/",
      Component: App,
      children: [
        { index: true, Component: MainFeed },
        {
          path: "profile/:handle",
          Component: User,
          loader: userLoader(client),
          children: [
            {
              index: true,
              Component: Feed,
              loader: feedLoader(client),
            },
          ],
        },
      ],
    },
  ]);
