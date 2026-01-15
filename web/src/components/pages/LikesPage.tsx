import { useQuery } from "@tanstack/react-query";
import { useContext, useState } from "react";
import { useLoaderData } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import AppContext from "../../lib/state";
import PostComposerOverlay from "../common/PostComposerOverlay";
import { PostItem } from "../common/PostItem";
import Sidebar from "../common/Sidebar";

export const loader = (_client: AppClient) => async () => {
  return {};
};

export default function LikesPage() {
  const client = useContext(ClientContext);
  const appState = useContext(AppContext);
  useLoaderData();
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [selectedPost, setSelectedPost] = useState<null>(null);
  const userId = appState?.userId ?? null;

  const {
    data: likedPosts,
    isPending: likesPending,
    isError: likesError,
  } = useQuery<components["schemas"]["Post"][]>({
    queryKey: ["user-likes", userId],
    enabled: Boolean(client) && userId !== null,
    queryFn: async () => {
      if (!client || userId === null) {
        throw new Error("Client or user not ready");
      }
      const res = await client.fetchClient.GET("/api/users/{id}/favourites", {
        params: { path: { id: userId } },
      });
      if (res.error) {
        throw new Error("Failed to fetch liked posts");
      }
      return res.data ?? [];
    },
  });

  const openComposerForNewPost = () => {
    setSelectedPost(null);
    setIsComposerOpen(true);
  };

  const closeComposer = () => {
    setIsComposerOpen(false);
    setSelectedPost(null);
  };

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
            <h1 className="text-2xl font-semibold text-gray-800">Likes</h1>
            <p className="text-gray-500">
              Posts you have liked will appear here.
            </p>
          </section>
          <section className="space-y-2">
            {likesPending && (
              <div className="p-4 text-sm text-gray-500">
                Loading liked posts…
              </div>
            )}
            {likesError && (
              <div className="p-4 text-sm text-red-600">
                Failed to load liked posts.
              </div>
            )}
            {!likesPending &&
              !likesError &&
              client &&
              likedPosts &&
              likedPosts.length > 0 &&
              likedPosts.map((post) => (
                <PostItem
                  key={post.id}
                  post={{ data: post }}
                  client={client}
                  likeActive
                />
              ))}
            {!likesPending && !likesError && !client && (
              <div className="p-4 text-sm text-gray-500">
                Client is not ready yet. Please try again.
              </div>
            )}
            {!likesPending &&
              !likesError &&
              client &&
              (!likedPosts || likedPosts.length === 0) && (
                <div className="p-4 text-sm text-gray-500">
                  No liked posts yet.
                </div>
              )}
          </section>
        </main>
        <Sidebar onPostClick={openComposerForNewPost} />
      </div>
      <PostComposerOverlay
        isOpen={isComposerOpen}
        onClose={closeComposer}
        replyTo={selectedPost}
        onSubmit={handleCreatePost}
      />
    </div>
  );
}
