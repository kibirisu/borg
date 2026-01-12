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
    // Backend profile endpoints are unimplemented; pass handle for display only.
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
  const derivedBio =
    tokenUserId !== null
      ? `u r logged in `
      : "Profile data is unavailable until the API endpoint is implemented.";

  const { data: profileData } = useQuery({
    queryKey: ["profile", tokenUserId ?? handle ?? derivedUsername],
    enabled: Boolean(client) && (tokenUserId !== null || Boolean(handle)),
    queryFn: async () => {
      const id =
        tokenUserId !== null ? tokenUserId : handle ? Number(handle) : null;
      console.log("[UserPage] fetching profile for id", id);
      if (!id) {
        return { username: derivedUsername, bio: derivedBio };
      }
      const res = await client!.fetchClient.GET("/api/users/{id}", {
        params: { path: { id: Number(id) } },
      });
      if (res.error || !res.data) {
        console.warn("[UserPage] profile fetch failed");
        return {
          username: derivedUsername,
          bio: derivedBio,
        };
      }
      return {
        username: res.data.username,
        bio: res.data.bio ?? derivedBio,
      };
    },
  });

  const userId = useMemo(() => {
    if (tokenUserId !== null) return tokenUserId;
    if (handle && !Number.isNaN(Number(handle))) return Number(handle);
    return null;
  }, [handle, tokenUserId]);

  const {
    data: posts,
    isPending: postsPending,
    isError: postsError,
  } = useQuery<components["schemas"]["Post"][]>({
    queryKey: ["user-posts", userId],
    enabled: Boolean(client) && userId !== null,
    queryFn: async () => {
      const res = await client!.fetchClient.GET("/api/users/{id}/posts", {
        params: { path: { id: userId! } },
      });
      if (res.error) {
        throw new Error("Failed to fetch posts");
      }
      return res.data ?? [];
    },
  });

  const [isComposerOpen, setComposerOpen] = useState(false);
  const { data: followers } = useQuery<components["schemas"]["User"][]>({
    queryKey: ["followers", userId],
    enabled: Boolean(client) && userId !== null,
    queryFn: async () => {
      if (!client || userId === null) {
        throw new Error("Client or user not ready");
      }
      const res = await client.fetchClient.GET("/api/users/{id}/followers", {
        params: { path: { id: userId } },
      });
      if (res.error) {
        throw new Error("Failed to fetch followers");
      }
      return (res.data as components["schemas"]["User"][]) ?? [];
    },
  });

  const { data: following } = useQuery<components["schemas"]["User"][]>({
    queryKey: ["following", userId],
    enabled: Boolean(client) && userId !== null,
    queryFn: async () => {
      if (!client || userId === null) {
        throw new Error("Client or user not ready");
      }
      const res = await client.fetchClient.GET("/api/users/{id}/following", {
        params: { path: { id: userId } },
      });
      if (res.error) {
        throw new Error("Failed to fetch following");
      }
      return (res.data as components["schemas"]["User"][]) ?? [];
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
    await client.fetchClient.POST("/api/posts", {
      body: { userID: userId, content },
    });
    await client.queryClient.invalidateQueries({
      queryKey: ["user-posts", userId],
    });
    await client.queryClient.invalidateQueries({
      queryKey: ["get", "/api/posts", {}],
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
                <p className="text-gray-500">@{profileData?.username}</p>
                <p className="text-2xl font-semibold text-gray-800">
                  {profileData?.bio}
                </p>
                <div className="mt-4 flex items-center gap-8 text-sm text-gray-600">
                  <span>
                    Followers:{" "}
                    <strong className="text-gray-900">
                      {followers ? followers.length : "—"}
                    </strong>
                  </span>
                  <span>
                    Following:{" "}
                    <strong className="text-gray-900">
                      {following ? following.length : "—"}
                    </strong>
                  </span>
                </div>
              </div>
            </div>
          </section>
          {/* POSTS */}
          <section className="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
            {postsPending && (
              <div className="p-4 text-sm text-gray-500">Loading posts…</div>
            )}
            {postsError && (
              <div className="p-4 text-sm text-red-600">
                Failed to load posts.
              </div>
            )}
            {!postsPending && !postsError && (
              <>
                {posts && posts.length > 0 ? (
                  posts.map((post) => (
                    <PostItem
                      key={post.id}
                      post={{ data: post }}
                      client={client!}
                      showActions
                      onCommentClick={(p) => {
                        if ("id" in p.data) {
                          navigate(`/post/${p.data.id}`);
                        }
                      }}
                    />
                  ))
                ) : (
                  <div className="p-4 text-sm text-gray-500">No posts yet.</div>
                )}
              </>
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
