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

export default function SharedPage() {
  const client = useContext(ClientContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  if (!opts) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="grid grid-cols-[1fr_256px] gap-6">
          <main className="px-6 py-6 space-y-6">
            <section className="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
              <h1 className="text-2xl font-semibold text-gray-800">Shared</h1>
              <p className="text-gray-500">
                Posts feed is not available yet. Check back soon.
              </p>
            </section>
          </main>
          <Sidebar onPostClick={() => {}} />
        </div>
      </div>
    );
  }
  const { data, isPending } = useQuery<components["schemas"]["Post"][]>(
    opts as any,
  );
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
            <h1 className="text-2xl font-semibold text-gray-800">Shared</h1>
            <p className="text-gray-500">
              Posts you share will live here. We&apos;ll plug in the sharing
              logic soon, so we are reusing a general feed for now.
            </p>
          </section>
          <section className="bg-white rounded-2xl border border-gray-200 p-4 space-y-4 min-h-[400px]">
            {isPending && (
              <p className="text-center text-gray-500">Loading…</p>
            )}
            {!isPending &&
              data?.map((post: components["schemas"]["Post"]) => (
                <PostItem
                  key={post.id}
                  post={{ data: post }}
                  client={client!}
                  onSelect={handlePostSelect}
                />
              ))}
            {!isPending && !data?.length && (
              <p className="text-center text-gray-500">
                Nothing shared yet.
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
