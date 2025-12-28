import { useQuery } from "@tanstack/react-query";
import { useContext, useState } from "react";
import { useLoaderData } from "react-router";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import type { components } from "../../lib/api/v1";
import { PostItem, type PostPresentable } from "../common/PostItem";
import PostComposerOverlay from "../common/PostComposerOverlay";
import Sidebar from "../common/Sidebar";

export const loader = (client: AppClient) => async () => {
  // const opts = client.$api.queryOptions("get", "/api/posts", {});
  // await client.queryClient.ensureQueryData(opts);
  // return { opts };
  return { opts: undefined };
};

export default function LikesPage() {
  const client = useContext(ClientContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const queryArgs = opts
    ? { ...opts }
    : { queryKey: [], queryFn: async () => [], enabled: false };
  const { data, isPending } = useQuery(queryArgs);
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [selectedPost, setSelectedPost] = useState<PostPresentable | null>(null);

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

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="grid grid-cols-[1fr_256px] gap-6">
        <main className="px-6 py-6 space-y-6">
          <section className="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
            <h1 className="text-2xl font-semibold text-gray-800">Likes</h1>
            <p className="text-gray-500">
              Posts you have liked will appear here. For now we&apos;re showing a
              general feed.
            </p>
          </section>
          <section className="bg-white rounded-2xl border border-gray-200 p-4 space-y-4 min-h-[400px]">
            {isPending && opts && (
              <p className="text-center text-gray-500">Loading…</p>
            )}
            {!isPending && opts &&
              data?.map((post: components["schemas"]["Post"]) => (
                <PostItem
                  key={post.id}
                  post={{ data: post }}
                  client={client!}
                  onSelect={handlePostSelect}
                />
              ))}
            {!isPending && opts && !data?.length && (
              <p className="text-center text-gray-500">
                Nothing liked yet.
              </p>
            )}
            {!opts && (
              <p className="text-center text-gray-500">
                Likes feed is not available yet.
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
      />
    </div>
  );
}
