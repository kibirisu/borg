import {
  type QueryKey,
  type UseQueryOptions,
  type UseSuspenseQueryOptions,
  useQuery,
  useSuspenseQuery,
} from "@tanstack/react-query";
import { useContext } from "react";
import { type LoaderFunctionArgs, useLoaderData } from "react-router";
import type { components } from "../../lib/api/v1";
import ClientContext, { type AppClient } from "../../lib/client";
import { PostItem } from "../common/PostItem";
import CommentForm from "./CommentForm";

export const loader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    if (!params.postId) {
      return { postOpts: undefined, commentOpts: undefined, postId: undefined };
    }
    const postId = Number(params.postId);
    const queryParams = { params: { path: { id: postId } } };
    const postOpts = client.$api.queryOptions(
      "get",
      "/api/posts/{id}",
      queryParams,
    );
    const commentOpts = client.$api.queryOptions(
      "get",
      "/api/posts/{id}/comments",
      queryParams,
    );
    client.queryClient.prefetchQuery(commentOpts);
    await client.queryClient.ensureQueryData(postOpts);
    return { postOpts, commentOpts, postId };
  };
export const commentsLoader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    if (!params.postId) {
      return { opts: undefined };
    }
    const postId = Number(params.postId);
    const queryParams = { params: { path: { id: postId } } };
    const commentOpts = client.$api.queryOptions(
      "get",
      "/api/posts/{id}/comments",
      queryParams,
    );
    await client.queryClient.ensureQueryData(commentOpts);
    return { opts: commentOpts };
  };

/**
 * View a single post (enlarged) and display its comments below.
 */
export default function CommentView() {
  const client = useContext(ClientContext);
  const { postOpts, commentOpts, postId } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;

  const postQueryOptions = (postOpts ??
    ({
      queryKey: ["post-view-disabled", postId],
      queryFn: async () => null,
      // Suspense queries are always “enabled”; fallback keeps the shape but
      // returns null immediately so UI can handle missing data.
    } satisfies UseSuspenseQueryOptions<
      components["schemas"]["Post"] | null,
      Error,
      components["schemas"]["Post"] | null,
      QueryKey
    >)) as UseSuspenseQueryOptions<
    components["schemas"]["Post"] | null,
    Error,
    components["schemas"]["Post"] | null,
    QueryKey
  >;

  const postData = useSuspenseQuery(postQueryOptions);

  if (!client) {
    return null;
  }
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <div className="w-full bg-white border border-gray-200 shadow-sm">
        {postData?.data ? (
          <PostItem
            post={{ data: postData.data as components["schemas"]["Post"] }}
            client={client}
          />
        ) : (
          <div className="p-6 text-center text-gray-600">Post not found.</div>
        )}
      </div>
      <div className="flex-1 overflow-y-auto bg-white border-x border-b border-gray-200">
        <CommentsFeed opts={commentOpts} postId={postId} />
      </div>
      <div className="sticky bottom-0 bg-white border-x border-b border-gray-200">
        <CommentForm />
      </div>
    </div>
  );
}

export function CommentsFeed({
  opts,
  postId: _postId,
}: {
  opts?:
    | UseQueryOptions<
        components["schemas"]["Comment"][],
        Error,
        components["schemas"]["Comment"][],
        QueryKey
      >
    | undefined;
  postId?: number;
}) {
  const client = useContext(ClientContext);

  const queryOptions = (opts ??
    ({
      queryKey: ["comments-feed-disabled", _postId],
      queryFn: async () => [] as components["schemas"]["Comment"][],
      enabled: false,
    } satisfies UseQueryOptions<
      components["schemas"]["Comment"][],
      Error,
      components["schemas"]["Comment"][],
      QueryKey
    >)) as UseQueryOptions<
    components["schemas"]["Comment"][],
    Error,
    components["schemas"]["Comment"][],
    QueryKey
  >;

  const { data, isPending } = useQuery(queryOptions);

  if (!opts) {
    return (
      <div className="p-6 text-center text-gray-600">
        Comments are not available yet.
      </div>
    );
  }

  if (isPending) {
    return (
      <div className="p-6 text-center text-gray-600">Loading comments…</div>
    );
  }
  if (!client) {
    return (
      <div className="p-6 text-center text-gray-600">
        Client is not ready yet. Please try again.
      </div>
    );
  }

  return (
    <div className="divide-y divide-gray-200">
      {data && data.length > 0 ? (
        data.map((comment: components["schemas"]["Comment"]) => (
          <PostItem key={comment.id} post={{ data: comment }} client={client} />
        ))
      ) : (
        <div className="p-6 text-center text-gray-600">No comments yet.</div>
      )}
    </div>
  );
}
