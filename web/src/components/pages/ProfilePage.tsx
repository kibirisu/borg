import { useContext, useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { type LoaderFunctionArgs, Outlet, useLoaderData } from "react-router";
import type { AppClient } from "../../lib/client";
import AppContext from "../../lib/state";
import Sidebar from "../common/Sidebar";
import anonAvatar from "../../assets/Anonomous.jpg";

export const loader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    // Backend profile endpoints are unimplemented; pass handle for display only.
    return { handle: params.handle };
  };

export default function ProfilePage() {
  const { handle } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const appState = useContext(AppContext);
  const tokenUsername = appState?.username ?? "";
  const tokenUserId = appState?.userId ?? null;
  const derivedUsername = useMemo(() => {
    if (tokenUsername) {
      return tokenUsername;
    }
    return handle ? String(handle) : "";
  }, [handle, tokenUsername]);
  const derivedBio =
    tokenUserId !== null
      ? `u r logged in `
      : "Profile data is unavailable until the API endpoint is implemented.";
  const [isFollowed, setIsFollowed] = useState(false);

  const { data: profileData } = useQuery({
    queryKey: ["profile", handle ?? derivedUsername],
    queryFn: async () => ({
      username: derivedUsername,
      bio: derivedBio,
    }),
  });

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="grid grid-cols-[1fr_256px] gap-6">
        <main className="px-6 py-6 space-y-6">
          <section className="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
            <div className="flex items-start gap-4">
              <div className="avatar">
                <div className="w-20 rounded-full overflow-hidden border border-gray-200 shadow-sm">
                  <img
                    src={anonAvatar}
                    alt="User avatar"
                    className="w-full h-full object-cover"
                  />
                </div>
              </div>
              <div className="flex-1">
                <p className="text-gray-500">@{profileData?.username}</p>
                <p className="text-2xl font-semibold text-gray-800">
                  {profileData?.bio}
                </p>
                <div className="mt-4 flex items-center gap-8 text-sm text-gray-600">
                  <span>
                    Followers: <strong className="text-gray-900">—</strong>
                  </span>
                  <span>
                    Following: <strong className="text-gray-900">—</strong>
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
