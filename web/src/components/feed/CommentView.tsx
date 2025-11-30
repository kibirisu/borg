import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { type LoaderFunctionArgs, Outlet, useLoaderData } from "react-router";
import type { Client } from "../../lib/client";
import { Post } from "../common/Feed";

export const loader =
  (client: Client) =>
  async ({ params }: LoaderFunctionArgs) => {
    const postId = parseInt(String(params.postId));
    console.log(postId);
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
    return { opts: postOpts };
  };
export const commentsLoader =
  (client: Client) =>
  async ({ params }: LoaderFunctionArgs) => {
    const postId = parseInt(String(params.postId));
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
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const postData = useSuspenseQuery(opts);
  return (
    <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
      <header className="p-4 border-b border-gray-300 text-xl font-bold sticky top-0 bg-white/80 backdrop-blur z-10 text-black">
        Post
      </header>

      <div className="border-b border-gray-200 p-4 bg-gray-50">
        <Post {...postData} />
      </div>
      <Outlet />
    </div>
  );
}

export function CommentsFeed() {
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof commentsLoader>>
  >;
  const { data, isPending } = useQuery(opts);
  if (isPending) {
    return <></>;
  }
  return (
    <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
      {data?.map((comment) => (
        <Post key={comment.id} data={comment} />
      ))}
    </div>
  );
}
