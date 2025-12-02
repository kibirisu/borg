import { createBrowserRouter } from "react-router";
import App from "../App";
import Feed, { loader as feedLoader } from "../components/common/Feed";
import { action as addCommentAction } from "../components/feed/CommentForm";
import CommentView, {
  CommentsFeed,
  commentsLoader,
  loader as masterPostLoader,
} from "../components/feed/CommentView";
import MainFeed, {
  loader as mainFeedLoader,
} from "../components/feed/MainFeed";
import User, { loader as userLoader } from "../components/profile/UserProfile";
import type { Client } from "../lib/client";

export default function newRouter(client: Client) {
  return createBrowserRouter([
    {
      path: "/",
      Component: App,
      children: [
        {
          path: "",
          Component: MainFeed,
          children: [
            {
              index: true,
              Component: Feed,
              loader: mainFeedLoader(client),
            },
          ],
        },
        {
          path: "profile/:handle",
          Component: User,
          loader: userLoader(client),
          errorElement: "error",
          children: [
            {
              index: true,
              Component: Feed,
              loader: feedLoader(client),
            },
          ],
        },
        {
          path: "post/:postId",
          Component: CommentView,
          loader: masterPostLoader(client),
          action: addCommentAction(client),
          errorElement: "error",
          children: [
            {
              index: true,
              Component: CommentsFeed,
              loader: commentsLoader(client),
            },
          ],
        },
      ],
    },
  ]);
}
