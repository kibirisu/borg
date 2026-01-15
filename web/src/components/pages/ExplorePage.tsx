import { useMutation, useQuery } from "@tanstack/react-query";
import { useContext, useState } from "react";
import { useLoaderData, useNavigate } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import AppContext from "../../lib/state";
import PostComposerOverlay from "../common/PostComposerOverlay";
import { PostItem, type PostPresentable } from "../common/PostItem";
import Sidebar from "../common/Sidebar";

function FoundUserItem({
  account,
}: {
  account: components["schemas"]["Account"];
}) {
  const display = account.displayName || account.username;
  const handle = account.acct || `@${account.username}`;
  const initial = display?.slice(0, 1).toUpperCase() || "?";

  return (
    <div className="bg-white rounded-2xl border border-gray-200 p-4 shadow-sm">
      <div className="flex items-center gap-3">
        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-indigo-100 text-indigo-600 font-semibold">
          {initial}
        </div>
        <div className="flex flex-col">
          <a
            href={`/profile/${account.id}`}
            className="text-base font-semibold text-gray-900 hover:text-indigo-600"
          >
            {display}
          </a>
          <span className="text-sm text-gray-500">{handle}</span>
        </div>
      </div>
    </div>
  );
}

export const loader = (client: AppClient) => async () => {
  const opts = client.$api.queryOptions("get", "/api/posts", {});
  await client.queryClient.ensureQueryData(opts);
  return { opts };
};

export default function ExplorePage() {
  const client = useContext(ClientContext);
  const appState = useContext(AppContext);
  const navigate = useNavigate();
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data, isPending } = useQuery(opts);
  const [searchTerm, setSearchTerm] = useState("");
  const [searchError, setSearchError] = useState("");
  const [searchResult, setSearchResult] = useState<
    components["schemas"]["Account"] | null
  >(null);
  const [isComposerOpen, setIsComposerOpen] = useState(false);
  const [selectedPost, setSelectedPost] = useState<PostPresentable | null>(
    null,
  );

  const lookupMutation = useMutation({
    mutationFn: async (acct: string) => {
      if (!client) {
        throw new Error("Search client is not ready yet.");
      }

      const response = await client.fetchClient.GET("/api/accounts/lookup", {
        params: { query: { acct } },
      });

      if (response.error) {
        throw new Error("Error during fetching client.");
      }

      if (!response.data) {
        throw new Error("No user found for this handle.");
      }

      return response.data;
    },
  });

  const handleSearch = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const trimmed = searchTerm.trim();
    const handlePattern = /^@[A-Za-z0-9._-]+@[A-Za-z0-9.-]+(?::\d{2,5})?$/;
    if (!handlePattern.test(trimmed)) {
      setSearchError("Format must be @user@host or @user@host:port");
      return;
    }
    if (!client) {
      setSearchError("Search client is not ready yet.");
      return;
    }
    setSearchError("");
    setSearchResult(null);
    try {
      const result = await lookupMutation.mutateAsync(trimmed);
      setSearchResult(result);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Unable to perform search.";
      setSearchError(message);
    }
  };

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

  const handleCreatePost = async (content: string) => {
    const userId = appState?.userId ?? null;
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
                placeholder="@user@host:port"
                className="block w-full rounded-xl border border-gray-200 bg-gray-50 p-3 pl-9 text-sm text-gray-900 placeholder:text-gray-500 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none"
                value={searchTerm}
                onChange={onSearchChange}
              />
            </div>
            <button
              type="submit"
              disabled={lookupMutation.isPending}
              className={`inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-xl text-white shadow-sm focus:outline-none focus:ring-4 focus:ring-indigo-300 ${
                lookupMutation.isPending
                  ? "bg-indigo-300 cursor-not-allowed"
                  : "bg-indigo-600 hover:bg-indigo-700"
              }`}
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
            <div
              className="flex items-start sm:items-center p-4 mb-4 text-sm text-fg-warning rounded-base bg-warning-soft text-yellow-800 bg-yellow-50 border border-yellow-200 rounded-lg"
              role="alert"
            >
              <svg
                className="w-4 h-4 me-2 shrink-0 mt-0.5 sm:mt-0"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  stroke="currentColor"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M10 11h2v5m-2 0h4m-2.592-8.5h.01M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
                />
              </svg>
              <p>
                <span className="font-medium me-1">Warning!</span>
                {searchError === "Error during fetching client."
                  ? "Sorry, we coudn't find this profile. Are you sure you entered the username and host correctly?"
                  : searchError}
              </p>
            </div>
          )}
          {searchResult && (
            <div className="max-w-4xl mx-auto">
              <p className="text-sm font-medium text-gray-500 mb-2">
                Search result
              </p>
              <FoundUserItem account={searchResult} />
            </div>
          )}
          <section className="space-y-2 min-h-[400px]">
            {isPending && <p className="text-gray-500 text-center">Loading…</p>}
            {!isPending &&
              client &&
              data?.map((post: components["schemas"]["Post"]) => (
                <PostItem
                  key={post.id}
                  post={{ data: post }}
                  client={client}
                  onSelect={handlePostSelect}
                  onCommentClick={handleCommentClick}
                />
              ))}
            {!isPending && !client && (
              <p className="text-center text-gray-500">
                Client is not ready yet. Please try again.
              </p>
            )}
            {!isPending && client && !data?.length && (
              <p className="text-center text-gray-500">Nothing posted yet.</p>
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
