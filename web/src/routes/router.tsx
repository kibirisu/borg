import { useContext } from "react";
import {
  createBrowserRouter,
  createContext,
  RouterContextProvider,
  RouterProvider as Provider,
} from "react-router";
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
import ProfilePage, {
  loader as profileLoader,
} from "../components/pages/ProfilePage";
import ExplorePage, {
  loader as exploreLoader,
} from "../components/pages/ExplorePage";
import Root from "../components/Root";
import ClientContext, { type AppClient } from "../lib/client";
import { Home } from "lucide-react";
import { SignIn } from "../components/auth/SignIn";
import { SignUp } from "../components/auth/SignUp";
import { signInAction } from "../components/auth/signInAction";
import { signUpAction } from "../components/auth/signUpAction";

const RouterProvider = () => {
  const client = useContext(ClientContext);
  return <Provider router={router(client!)} />;
};

export default RouterProvider;

// We may use react router context strategy instead of passing client object to each route loader
export const routerContext = createContext<AppClient | undefined>(undefined);

function router(client: AppClient) {
  return createBrowserRouter(
    [
      {
        path: "/",
        Component: Root,
        action: addPostAction(client),
        loader: ({ context }) => {
          context.get(routerContext); // we collect object from context like this
          return null;
        },
        children: [
          {
            path: "home",
            Component: Home,
          },
          {
            path: "signin",
            Component: SignIn,
            action: signInAction(client),
          },
          {
            path: "signup",
            Component: SignUp,
            action: signUpAction(client),
          },
          {
            path: "explore",
            Component: ExplorePage,
            loader: exploreLoader(client),
            errorElement: "error",
          },
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
            Component: ProfilePage,
            loader: profileLoader(client),
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
    ],
    {
      // We add base context to all routes
      getContext() {
        const context = new RouterContextProvider();
        context.set(routerContext, client);
        return context;
      },
    },
  );
}
