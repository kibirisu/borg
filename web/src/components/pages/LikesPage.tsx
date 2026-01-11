import { useContext, useState } from "react";
import { useLoaderData } from "react-router";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import { PostItem, type PostPresentable } from "../common/PostItem";
import PostComposerOverlay from "../common/PostComposerOverlay";
import Sidebar from "../common/Sidebar";

export const loader = (client: AppClient) => async () => {
  return {};
};

export default function LikesPage() {
  const client = useContext(ClientContext);
  useLoaderData();
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [selectedPost, setSelectedPost] = useState<PostPresentable | null>(
    null,
  );

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
              Posts you have liked will appear here. For now we&apos;re showing
              a general feed.
            </p>
          </section>
          {/* Likes by post ID are handled via the form above; feed removed. */}
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
