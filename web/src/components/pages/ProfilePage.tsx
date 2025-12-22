import { useSuspenseQuery } from "@tanstack/react-query";
import { useState } from "react";
import { type LoaderFunctionArgs, Outlet, useLoaderData } from "react-router";
import type { AppClient } from "../../lib/client";
import Sidebar from "../common/Sidebar";

export const loader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    const userId = parseInt(String(params.handle), 10);
    if (Number.isNaN(userId)) {
      throw new Response("User handle should be a numeric id", {
        status: 400,
      });
    }
    const queryParams = { params: { path: { id: userId } } };
    const userOpts = client.$api.queryOptions(
      "get",
      "/api/users/{id}",
      queryParams,
    );
    await client.queryClient.ensureQueryData(userOpts);
    return { userOpts };
  };

export default function ProfilePage() {
  const { userOpts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data: userData } = useSuspenseQuery(userOpts);
  const [isFollowed, setIsFollowed] = useState(false);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="grid grid-cols-[1fr_256px] gap-6">
        <main className="px-6 py-6 space-y-6">
          <section className="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
            <div className="flex items-start gap-4">
              <div className="avatar placeholder">
                <div className="bg-neutral text-neutral-content w-20 rounded-full">
                  <span className="text-3xl">
                    {getUserInitials(userData.username)}
                  </span>
                </div>
              </div>
              <div className="flex-1">
                <p className="text-gray-500">@{userData.username}</p>
                <p className="text-2xl font-semibold text-gray-800">
                  {userData.bio || "Без опису"}
                </p>
                <div className="mt-4 flex items-center gap-8 text-sm text-gray-600">
                  <span>
                    Followers:{" "}
                    <strong className="text-gray-900">
                      {userData.followersCount}
                    </strong>
                  </span>
                  <span>
                    Following:{" "}
                    <strong className="text-gray-900">
                      {userData.followingCount}
                    </strong>
                  </span>
                </div>
              </div>
              <button
                type="button"
                onClick={() => setIsFollowed((prev) => !prev)}
                className={`btn rounded-[12px] ${isFollowed ? "btn-outline btn-secondary" : "btn-primary"}`}
              >
                {isFollowed ? "Unfollow" : "Follow"}
              </button>
            </div>
          </section>
            {/* POSTS */}
          <section className="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
              <Outlet />
          </section>
        </main>
        <Sidebar />
      </div>
    </div>
  );
}

function getUserInitials(username: string): string {
  if (!username) return "";
  return username.replace(/^@/, "").slice(0, 2).toUpperCase();
}
