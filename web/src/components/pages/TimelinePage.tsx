import { useQuery } from "@tanstack/react-query";
import { useContext, useState } from "react";
import { useLoaderData, useNavigate } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import PostComposerOverlay from "../common/PostComposerOverlay";
import { PostItem, type PostPresentable } from "../common/PostItem";
import Sidebar from "../common/Sidebar";

export const loader = (client: AppClient) => async () => {
  const opts = client.$api.queryOptions("get", "/api/timelines/home", {});
  await client.queryClient.ensureQueryData(opts);
  return { opts };
};

export default function TimelinePage() {
  const client = useContext(ClientContext);
  const navigate = useNavigate();
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data, isPending } = useQuery(opts);
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [selectedPost, setSelectedPost] = useState<PostPresentable | null>(
    null,
  );

  const handlePostSelect = (post: PostPresentable) => {
    setSelectedPost(post);
    setIsComposerOpen(true);
  };

  const handleCommentClick = (post: PostPresentable) => {
    if ("id" in post.data) {
      navigate(`/post/${post.data.id}`);
    }
  };

  const openComposerForNewPost = () => {
    setSelectedPost(null);
    setIsComposerOpen(true);
  };

  const closeComposer = () => {
    setIsComposerOpen(false);
    setSelectedPost(null);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="grid grid-cols-[1fr_256px] gap-6">
        <main className="px-6 py-6 space-y-6">
          <section className="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
            <h1 className="text-2xl font-semibold text-gray-800">Timeline</h1>
            <p className="text-gray-500">
              Updates from accounts you follow.
            </p>
          </section>
          <section className="bg-white rounded-2xl border border-gray-200 p-4 space-y-4 min-h-[400px]">
            {isPending && <p className="text-gray-500 text-center">Loading…</p>}
            {!isPending &&
              data?.map((post: components["schemas"]["Status"]) => (
                <PostItem
                  key={post.id}
                  post={{ data: post }}
                  client={client!}
                  onSelect={handlePostSelect}
                  onCommentClick={handleCommentClick}
                />
              ))}
            {!isPending && !data?.length && (
              <p className="text-center text-gray-500">Nothing here yet.</p>
            )}
          </section>
        </main>
        <Sidebar onPostClick={openComposerForNewPost} />
      </div>
      <PostComposerOverlay
        isOpen={isComposerOpen}
        onClose={closeComposer}
        replyTo={selectedPost}
      />
    </div>
  );
}
