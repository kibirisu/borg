// If you're working on this please be familiar how react query works!!!
// https://tanstack.com/query/latest

import { Outlet, useLoaderData, type LoaderFunctionArgs } from "react-router";
import type { Client } from "../../lib/api/client";
import { useSuspenseQuery } from "@tanstack/react-query";
import { useState } from "react";

export const loader =
  (client: Client) =>
    async ({ params }: LoaderFunctionArgs) => {
      const userId = parseInt(String(params.handle));
      const queryParams = { params: { path: { id: userId } } };
      const userOpts = client.$api.queryOptions(
        "get",
        "/api/users/{id}",
        queryParams,
      );
      const postOpts = client.$api.queryOptions(
        "get",
        "/api/users/{id}/posts",
        queryParams,
      );
      client.queryClient.prefetchQuery(postOpts);
      await client.queryClient.ensureQueryData(userOpts);
      console.log(postOpts.queryKey);
      return { opts: userOpts };
    };

export default function User() {
  const { opts: opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data: userData } = useSuspenseQuery(opts);
  const [isFollowed, setIsFollowed] = useState(false);

  return (
    <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
      <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black"></header>
      <div className="p-6">
        <div className="flex flex-col items-center space-y-4">
          <div className="avatar avatar-placeholder">
            <div className="bg-neutral text-neutral-content w-24 rounded-full">
              <span className="text-3xl">
                {getUserInitials(userData.username)}
              </span>
            </div>
          </div>
          <div className="text-center">
            <div className="text-gray-600">{userData.username}</div>
            <div className="text-2xl font-semibold text-gray-700">
              {userData.bio}
            </div>
            <button
              type="button"
              onClick={() => setIsFollowed(!isFollowed)}
              className={`btn ${isFollowed ? "btn-outline btn-secondary" : "btn-primary"}`}
            >
              {isFollowed ? "Unfollow" : "Follow"}
            </button>
          </div>
        </div>
      </div>
      <Outlet />
    </div>
  );
}

function getUserInitials(username: string): string {
  if (!username) return "";
  return username.replace(/^@/, "").slice(0, 2).toUpperCase();
}
