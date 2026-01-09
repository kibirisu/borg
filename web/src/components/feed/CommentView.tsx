import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { useContext } from "react";
import { type LoaderFunctionArgs, useLoaderData } from "react-router";
import ClientContext, { type AppClient } from "../../lib/client";
import type { components } from "../../lib/api/v1";
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
  const postData = postOpts ? useSuspenseQuery(postOpts) : null;
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="w-full bg-white border border-gray-200 overflow-hidden divide-y divide-gray-200 shadow-sm">
        <div className="bg-white">
          {postData && postData.data ? (
            <PostItem post={{ data: postData.data as any }} client={client!} />
          ) : (
            <div className="p-6 text-center text-gray-600">Post not found.</div>
          )}
        </div>
        <div className="p-4 bg-gray-50">
          <CommentForm />
        </div>
        <div className="bg-white">
          <CommentsFeed opts={commentOpts} postId={postId} />
        </div>
      </div>
    </div>
  );
}

export function CommentsFeed({
  opts,
  postId,
}: {
  opts?: any;
  postId?: number;
}) {
  const client = useContext(ClientContext);
  if (!opts) {
    return (
      <div className="p-6 text-center text-gray-600">
        Comments are not available yet.
      </div>
    );
  }
  const { data, isPending } = useQuery<components["schemas"]["Comment"][]>(
    opts as any,
  );
  if (isPending) {
    return <></>;
  }
  return (
    <div className="divide-y divide-gray-200">
      {data && data.length > 0 ? (
        data.map((comment: components["schemas"]["Comment"]) => (
          <PostItem
            key={comment.id}
            post={{ data: comment }}
            client={client!}
          />
        ))
      ) : (
        <div className="p-6 text-center text-gray-600">No comments yet.</div>
      )}
    </div>
  );
}
