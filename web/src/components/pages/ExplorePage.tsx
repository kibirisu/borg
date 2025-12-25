import { useQuery } from "@tanstack/react-query";
import { useContext, useState } from "react";
import { useLoaderData } from "react-router";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import { PostItem, type PostPresentable } from "../common/PostItem";
import PostComposerOverlay from "../common/PostComposerOverlay";
import Sidebar from "../common/Sidebar";

export const loader = (client: AppClient) => async () => {
  const opts = client.$api.queryOptions("get", "/api/posts", {});
  await client.queryClient.ensureQueryData(opts);
  return { opts };
};

export default function ExplorePage() {
  const client = useContext(ClientContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data, isPending } = useQuery(opts);
  const [searchTerm, setSearchTerm] = useState("");
  const [searchError, setSearchError] = useState("");
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [selectedPost, setSelectedPost] = useState<PostPresentable | null>(null);

  const handleSearch = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const trimmed = searchTerm.trim();
    const handlePattern = /^@[A-Za-z0-9._-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$/;
    if (!handlePattern.test(trimmed)) {
      setSearchError("Format must be @user@instance.com");
      return;
    }
    setSearchError("");
    // TODO: replace with actual search logic once API endpoint is ready
    console.info("Searching for handle:", trimmed);
    console.log("Search input value:", trimmed);
  };

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

  const onSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(event.target.value);
    if (searchError) {
      setSearchError("");
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="grid grid-cols-[1fr_256px] gap-6">
        <main className="px-6 py-6 space-y-6">
          <section className="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
            <h1 className="text-2xl font-semibold text-gray-800">
              Explore trending posts
            </h1>
            <p className="text-gray-500">
              See what others are talking about right now.
            </p>
          </section>
          <form
            onSubmit={handleSearch}
            className="max-w-md mx-auto flex items-center space-x-2"
          >
            <label htmlFor="explore-search" className="sr-only">
              Search
            </label>
            <div className="relative w-full">
              <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                <svg
                  className="h-4 w-4 text-gray-500"
                  aria-hidden="true"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke="currentColor"
                    strokeLinecap="round"
                    strokeWidth="2"
                    d="m21 21-3.5-3.5M17 10a7 7 0 1 1-14 0 7 7 0 0 1 14 0Z"
                  />
                </svg>
              </div>
              <input
                id="explore-search"
                type="search"
                placeholder="@user@instance.com"
                className="block w-full rounded-xl border border-gray-200 bg-gray-50 p-3 pl-9 text-sm text-gray-900 placeholder:text-gray-500 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none"
                value={searchTerm}
                onChange={onSearchChange}
              />
            </div>
            <button
              type="submit"
              className="inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-indigo-600 text-white shadow-sm focus:outline-none focus:ring-4 focus:ring-indigo-300 hover:bg-indigo-700"
            >
              <svg
                className="h-5 w-5"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeWidth="2"
                  d="m21 21-3.5-3.5M17 10a7 7 0 1 1-14 0 7 7 0 0 1 14 0Z"
                />
              </svg>
              <span className="sr-only">Search</span>
            </button>
          </form>
          {searchError && (
            <p className="text-center text-sm text-red-600">{searchError}</p>
          )}
          <section className="bg-white rounded-2xl border border-gray-200 p-4 space-y-4 min-h-[400px]">
            {isPending && (
              <p className="text-gray-500 text-center">Loading…</p>
            )}
            {!isPending &&
              data?.map((post) => (
                <PostItem
                  key={post.id}
                  post={{ data: post }}
                  client={client!}
                  onSelect={handlePostSelect}
                />
              ))}
            {!isPending && !data?.length && (
              <p className="text-center text-gray-500">
                Nothing posted yet.
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
