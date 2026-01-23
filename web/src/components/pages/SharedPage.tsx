import { useQuery } from "@tanstack/react-query";
import { useContext, useState } from "react";
import { useLoaderData } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import AppContext from "../../lib/state";
import PostComposerOverlay from "../common/PostComposerOverlay";
import { PostItem, type PostPresentable } from "../common/PostItem";
import Sidebar from "../common/Sidebar";

export const loader = (client: AppClient) => async () => {
  const opts = client.$api.queryOptions("get", "/api/timelines/reblogged", {});
  await client.queryClient.ensureQueryData(opts);
  return { opts };
};

export default function SharedPage() {
  const client = useContext(ClientContext);
  const appState = useContext(AppContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;

  const { data, isPending } = useQuery(opts);
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [selectedPost, setSelectedPost] = useState<PostPresentable | null>(
    null,
  );
  const userId = appState?.userId ?? null;

  const handlePostSelect = (post: PostPresentable) => {
    setSelectedPost(post);
    setIsComposerOpen(true);
  };

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
    const replyToId =
      selectedPost?.data?.reblog?.id ?? selectedPost?.data?.id ?? null;
    await client.fetchClient.POST("/api/statuses", {
      body: { status: content, in_reply_to_id: replyToId },
    });
    await client.queryClient.invalidateQueries({
      queryKey: ["account-statuses", userId],
    });
    await client.queryClient.invalidateQueries({
      queryKey: ["get", "/api/timelines/home", {}],
    });
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="grid grid-cols-[1fr_256px] gap-6">
        <main className="px-6 py-6 space-y-6">
          <section className="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
            <h1 className="text-2xl font-semibold text-gray-800">Shared</h1>
            <p className="text-gray-500">
              Your reshared posts, clean and easy to find.
            </p>
          </section>
          <section className="rounded-2xl bg-transparent min-h-[400px]">
            {isPending && <p className="text-center text-gray-500">Loading…</p>}
            {!isPending &&
              client &&
              data?.map((post: components["schemas"]["Status"]) => (
                <div
                  key={post.id}
                  className="mb-3 rounded-2xl border border-gray-200 bg-white shadow-sm overflow-hidden"
                >
                  <PostItem
                    post={{ data: post }}
                    client={client}
                    onSelect={handlePostSelect}
                  />
                </div>
              ))}
            {!isPending && client && (!data || data.length === 0) && (
              <p className="text-center text-gray-500">Nothing shared yet.</p>
            )}
            {!client && (
              <p className="text-center text-gray-500">
                Client is not ready yet. Please try again.
              </p>
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
