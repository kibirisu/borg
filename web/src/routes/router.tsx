import React from "react";
import { createBrowserRouter } from "react-router";
import App from "../App";
import MainFeed, { loader as mainFeedLoader } from "../components/feed/MainFeed";
import User, { loader as userLoader } from "../components/profile/UserProfile";
import type { Client } from "../lib/api/client";
import Feed, { loader as feedLoader } from "../components/common/Feed";

function ErrorPage() {
  return (
    <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
      <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
        Profile
      </header>
      <div className="p-6 text-red-600">
        <h2 className="text-xl font-bold mb-2">Error loading user profile</h2>
        <p>The user you're looking for might not exist or there was an error loading the profile.</p>
        <a href="/" className="text-blue-500 hover:underline mt-4 inline-block">
          Go back to home
        </a>
      </div>
    </div>
  );
}

export const newRouter = (client: Client) =>
  createBrowserRouter([
    {
      path: "/",
      Component: App,
      children: [
        { index: true, Component: MainFeed, loader: mainFeedLoader(client) },
        {
          path: "profile/:handle",
          Component: User,
          loader: userLoader(client),
          errorElement: <ErrorPage />,
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
