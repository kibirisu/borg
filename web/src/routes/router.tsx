import { useContext } from "react";
import { createBrowserRouter, RouterProvider as Provider } from "react-router";
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
import { action as addPostAction } from "../components/feed/NewPostForm";
import User, { loader as userLoader } from "../components/profile/UserProfile";
import Root from "../components/Root";
import { type AppClient, ClientContext } from "../lib/client";

export const RouterProvider = () => {
  const client = useContext(ClientContext);
  if (!client) {
    throw Error();
  }
  const router = newRouter(client);
  return <Provider router={router} />;
};

function newRouter(client: AppClient) {
  return createBrowserRouter([
    {
      path: "/",
      Component: Root,
      action: addPostAction(client),
      children: [
        {
          path: "",
          Component: MainFeed,
          errorElement: "error",
          children: [
            {
              index: true,
              action: addPostAction(client),
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
