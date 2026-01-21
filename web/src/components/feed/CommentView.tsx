import {
  type UseQueryOptions,
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
    const routePostId = String(params.postId);
    const routeParams = { params: { path: { id: routePostId } } };
    const routePostOpts = client.$api.queryOptions(
      "get",
      "/api/statuses/{id}",
      routeParams,
    );
    const postData = await client.queryClient.ensureQueryData(routePostOpts);
    const canonicalPostId = postData?.reblog?.id ?? routePostId;
    const canonicalParams = { params: { path: { id: canonicalPostId } } };
    const postOpts = client.$api.queryOptions(
      "get",
      "/api/statuses/{id}",
      canonicalParams,
    );
    const commentOpts = client.$api.queryOptions(
      "get",
      "/api/statuses/{id}/replies",
      canonicalParams,
    );
    client.queryClient.prefetchQuery(commentOpts);
    if (canonicalPostId !== routePostId) {
      await client.queryClient.ensureQueryData(postOpts);
    }
    return { postOpts, commentOpts, postId: canonicalPostId };
  };
export const commentsLoader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    if (!params.postId) {
      return { opts: undefined };
    }
    const postId = String(params.postId);
    const queryParams = { params: { path: { id: postId } } };
    const commentOpts = client.$api.queryOptions(
      "get",
      "/api/statuses/{id}/replies",
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

  const postQueryOptions =
    postOpts ??
    ({
      queryKey: ["post-view-disabled", postId],
      queryFn: async () => null,
      // Suspense queries are always “enabled”; fallback keeps the shape but
      // returns null immediately so UI can handle missing data.
    } satisfies Parameters<typeof useSuspenseQuery>[0]);

  const postData = useSuspenseQuery(postQueryOptions as any);

  if (!client) {
    return null;
  }
  return (
    <div className="min-h-screen bg-gray-50 pb-28">
      <div className="px-6 pt-6">
        <button
          type="button"
          onClick={() => window.history.back()}
          aria-label="Go back"
          className="inline-flex items-center justify-center border border-black text-black rounded-[7px] text-sm p-2.5"
        >
          <i className="bi bi-arrow-left" />
        </button>
      </div>
      <div className="mt-4 w-full bg-white border border-gray-200 overflow-hidden divide-y divide-gray-200 shadow-sm">
        <div className="bg-white">
          {postData && postData.data ? (
            <PostItem
              post={{ data: postData.data as components["schemas"]["Status"] }}
              client={client}
            />
          ) : (
            <div className="p-6 text-center text-gray-600">Post not found.</div>
          )}
        </div>
        <div className="bg-gray-100">
          <CommentsFeed opts={commentOpts} postId={postId} />
        </div>
      </div>
      <CommentForm postId={postId} />
    </div>
  );
}

export function CommentsFeed({
  opts,
  postId: _postId,
}: {
  opts?:
    | UseQueryOptions<
        components["schemas"]["Status"][],
        any,
        components["schemas"]["Status"][],
        any
      >
    | any;
  postId?: string;
}) {
  const client = useContext(ClientContext);

  const queryOptions =
    opts ??
    ({
      queryKey: ["comments-feed-disabled", _postId],
      queryFn: async () => [] as components["schemas"]["Status"][],
      enabled: false,
    } satisfies Parameters<typeof useQuery>[0]);

  const { data, isPending } = useQuery<components["schemas"]["Status"][]>(
    queryOptions as any,
  );

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
  return (
    <div className="p-4">
      {data && data.length > 0 ? (
        <div className="space-y-3">
          {data.map((comment: components["schemas"]["Status"]) => (
            comment && (
              <div
                key={comment.id}
                className="rounded-2xl border border-gray-200 bg-white shadow-sm overflow-hidden"
              >
                <PostItem post={{ data: comment }} client={client!} />
              </div>
            )
          ))}
        </div>
      ) : (
        <div className="p-6 text-center text-gray-600">No comments yet.</div>
      )}
    </div>
  );
}
