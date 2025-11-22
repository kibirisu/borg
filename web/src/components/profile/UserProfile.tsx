// If you're working on this please be familiar how react query works!!!
// https://tanstack.com/query/latest

import { Outlet, useLoaderData, type LoaderFunctionArgs } from "react-router";
import type { Client } from "../../lib/api/client";
import { useSuspenseQuery } from "@tanstack/react-query";
import { useState } from "react";
import type { components } from "../../lib/api/v1";

export const loader =
  (client: Client) =>
    async ({ params }: LoaderFunctionArgs) => {
      if (!params.handle) {
        throw new Response("Handle parameter is required", { status: 400 });
      }
      // Remove @ if present
      const username = String(params.handle).replace(/^@/, '');
      const userQueryParams = { params: { path: { username } } };
      const userOpts = client.$api.queryOptions(
        "get",
        "/api/users/by-username/{username}",
        userQueryParams,
      );
      
      try {
        await client.queryClient.ensureQueryData(userOpts);
        const userData = client.queryClient.getQueryData(userOpts.queryKey) as components["schemas"]["User"] | undefined;
        
        // Prefetch posts using user ID
        if (userData?.id) {
          const postQueryParams = { params: { path: { id: userData.id } } };
          const postOpts = client.$api.queryOptions(
            "get",
            "/api/users/{id}/posts",
            postQueryParams,
          );
          client.queryClient.prefetchQuery(postOpts);
        }
      } catch (error) {
        console.error("Error loading user:", error);
        // Re-throw to let error boundary handle it
        throw error;
      }
      
      return { opts: userOpts };
    };

export default function User() {
  const { opts: opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data: userData, error } = useSuspenseQuery(opts);
  const [isFollowed, setIsFollowed] = useState(false);

  if (error) {
    return (
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Profile
        </header>
        <div className="p-6 text-red-600">
          Error loading user profile: {error instanceof Error ? error.message : String(error)}
        </div>
      </div>
    );
  }

  if (!userData) {
    return (
      <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
        <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
          Profile
        </header>
        <div className="p-6">User not found</div>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
      <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
        Profile
      </header>
      <div className="p-6">
        <div className="flex flex-col items-center space-y-4">
          <div className="avatar avatar-placeholder">
            <div className="bg-neutral text-neutral-content w-24 rounded-full">
              <span className="text-3xl">
                {getUserInitials(userData.username || "")}
              </span>
            </div>
          </div>
          <div className="text-center">
            <div className="text-gray-600">{userData.username || "Unknown"}</div>
            <div className="text-2xl font-semibold text-gray-700">
              {userData.bio || "No bio"}
            </div>
            <div className="mt-2 text-sm text-gray-500">
              {userData.followersCount || 0} followers · {userData.followingCount || 0} following
            </div>
            <button
              type="button"
              onClick={() => setIsFollowed(!isFollowed)}
              className={`btn mt-4 ${isFollowed ? "btn-outline btn-secondary" : "btn-primary"}`}
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
