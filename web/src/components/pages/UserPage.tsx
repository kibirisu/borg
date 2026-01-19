import { useQuery } from "@tanstack/react-query";
import { useContext, useMemo, useState } from "react";
import {
  type LoaderFunctionArgs,
  Outlet,
  useLoaderData,
  useNavigate,
} from "react-router";
import anonAvatar from "../../assets/Anonomous.jpg";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import AppContext from "../../lib/state";
import PostComposerOverlay from "../common/PostComposerOverlay";
import { PostItem } from "../common/PostItem";
import Sidebar from "../common/Sidebar";

export const loader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    // Pass handle for routing; data is loaded via queries.
    return { handle: params.handle };
  };

export default function UserPage() {
  const { handle } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const navigate = useNavigate();
  const appState = useContext(AppContext);
  const client = useContext(ClientContext);
  const tokenUsername = appState?.username ?? "";
  const tokenUserId = appState?.userId ?? null;
  const derivedUsername = useMemo(() => {
    if (tokenUsername) {
      return tokenUsername;
    }
    return handle ? String(handle) : "";
  }, [handle, tokenUsername]);
  const fallbackDisplay =
    tokenUserId !== null
      ? "Signed-in profile"
      : "Profile data is unavailable until the API endpoint is implemented.";

  const userId = useMemo(() => {
    if (tokenUserId !== null) return String(tokenUserId);
    if (handle) return String(handle);
    return null;
  }, [handle, tokenUserId]);

  const { data: profileData } = useQuery<
    components["schemas"]["Account"] | null
  >({
    queryKey: ["profile", userId ?? derivedUsername],
    enabled: Boolean(client) && userId !== null,
    queryFn: async () => {
      const id = userId;
      console.log("[UserPage] fetching profile for id", id);
      if (!id) {
        return null;
      }
      const res = await client!.fetchClient.GET("/api/accounts/{id}", {
        params: { path: { id: String(id) } },
      });
      if (res.error || !res.data) {
        console.warn("[UserPage] profile fetch failed");
        return null;
      }
      return res.data;
    },
  });
  const profileHandleRaw = profileData?.acct ?? derivedUsername;
  const profileHandle = profileHandleRaw
    ? profileHandleRaw.startsWith("@")
      ? profileHandleRaw.slice(1)
      : profileHandleRaw
    : "";
  const profileDisplay =
    profileData?.display_name ||
    profileData?.username ||
    derivedUsername ||
    fallbackDisplay;
  const followersCount = profileData?.followers_count;
  const followingCount = profileData?.following_count;

  const {
    data: posts,
    isPending: postsPending,
    isError: postsError,
  } = useQuery<components["schemas"]["Status"][]>({
    queryKey: ["account-statuses", userId],
    enabled: Boolean(client) && userId !== null,
    queryFn: async () => {
      const res = await client!.fetchClient.GET("/api/accounts/{id}/statuses", {
        params: { path: { id: userId! } },
      });
      if (res.error) {
        throw new Error("Failed to fetch posts");
      }
      return res.data ?? [];
    },
  });
  const displayPosts = useMemo(() => {
    if (!posts) return [];
    const filtered =
      userId !== null
        ? posts.filter((post) => post.account?.id === userId)
        : posts;
    const seen = new Set<string>();
    const unique: components["schemas"]["Status"][] = [];
    for (const post of filtered) {
      const key = String(post.id);
      if (seen.has(key)) {
        continue;
      }
      seen.add(key);
      unique.push(post);
    }
    return unique;
  }, [posts, userId]);

  const [isComposerOpen, setComposerOpen] = useState(false);
  const { data: followers } = useQuery<components["schemas"]["Account"][]>({
    queryKey: ["followers", userId],
    enabled: Boolean(client) && userId !== null,
    queryFn: async () => {
      if (!client || userId === null) {
        throw new Error("Client or user not ready");
      }
      const res = await client.fetchClient.GET("/api/accounts/{id}/followers", {
        params: { path: { id: userId } },
      });
      if (res.error) {
        throw new Error("Failed to fetch followers");
      }
      return (res.data as components["schemas"]["Account"][]) ?? [];
    },
  });

  const { data: following } = useQuery<components["schemas"]["Account"][]>({
    queryKey: ["following", userId],
    enabled: Boolean(client) && userId !== null,
    queryFn: async () => {
      if (!client || userId === null) {
        throw new Error("Client or user not ready");
      }
      const res = await client.fetchClient.GET("/api/accounts/{id}/following", {
        params: { path: { id: userId } },
      });
      if (res.error) {
        throw new Error("Failed to fetch following");
      }
      return (res.data as components["schemas"]["Account"][]) ?? [];
    },
  });
  const openComposer = () => {
    console.log("[UserPage] open composer", { userId });
    setComposerOpen(true);
  };
  const closeComposer = () => setComposerOpen(false);

  const handleCreatePost = async (content: string) => {
    if (!client || userId === null) {
      throw new Error("User not authenticated");
    }
    await client.fetchClient.POST("/api/statuses", {
      body: { status: content, in_reply_to_id: null },
    });
    await client.queryClient.invalidateQueries({
      queryKey: ["account-statuses", userId],
    });
  };

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
                <p className="text-gray-500">
                  {profileHandle ? `@${profileHandle}` : "@"}
                </p>
                <p className="text-2xl font-semibold text-gray-800">
                  {profileDisplay}
                </p>
                <div className="mt-4 flex items-center gap-8 text-sm text-gray-600">
                  <span>
                    Followers:{" "}
                    <strong className="text-gray-900">
                      {followersCount ??
                        (followers ? followers.length : "—")}
                    </strong>
                  </span>
                  <span>
                    Following:{" "}
                    <strong className="text-gray-900">
                      {followingCount ??
                        (following ? following.length : "—")}
                    </strong>
                  </span>
                </div>
              </div>
            </div>
          </section>
          {/* POSTS */}
          <section className="rounded-2xl bg-transparent">
            {postsPending && (
              <div className="p-4 text-sm text-gray-500">Loading posts…</div>
            )}
            {postsError && (
              <div className="p-4 text-sm text-red-600">
                Failed to load posts.
              </div>
            )}
            {!postsPending && !postsError && (
              <div className="space-y-3">
                {displayPosts.length > 0 ? (
                  displayPosts.map((post) => (
                    <div
                      key={post.id}
                      className="rounded-2xl border border-gray-200 bg-white shadow-sm overflow-hidden"
                    >
                      <PostItem
                        post={{ data: post }}
                        client={client!}
                        showActions
                        onCommentClick={(p) => {
                          if ("id" in p.data) {
                            navigate(`/post/${p.data.id}`);
                          }
                        }}
                      />
                    </div>
                  ))
                ) : (
                  <div className="p-4 text-sm text-gray-500">No posts yet.</div>
                )}
              </div>
            )}
            <Outlet />
          </section>
        </main>
        <Sidebar onPostClick={openComposer} />
      </div>
      <PostComposerOverlay
        isOpen={isComposerOpen}
        onClose={closeComposer}
        replyTo={null}
        onSubmit={handleCreatePost}
      />
    </div>
  );
}
