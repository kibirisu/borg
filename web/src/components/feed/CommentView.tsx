import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { useContext } from "react";
import { type LoaderFunctionArgs, Outlet, useLoaderData } from "react-router";
import ClientContext, { type AppClient } from "../../lib/client";
import type { components } from "../../lib/api/v1";
import { PostItem } from "../common/PostItem";
import CommentForm from "./CommentForm";

export const loader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    // const postId = parseInt(String(params.postId));
    // const queryParams = { params: { path: { id: postId } } };
    // const postOpts = client.$api.queryOptions(
    //   "get",
    //   "/api/posts/{id}",
    //   queryParams,
    // );
    // const commentOpts = client.$api.queryOptions(
    //   "get",
    //   "/api/posts/{id}/comments",
    //   queryParams,
    // );
    // client.queryClient.prefetchQuery(commentOpts);
    // await client.queryClient.ensureQueryData(postOpts);
    // return { opts: postOpts };
    return { opts: undefined };
  };
export const commentsLoader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    // const postId = parseInt(String(params.postId));
    // const queryParams = { params: { path: { id: postId } } };
    // const commentOpts = client.$api.queryOptions(
    //   "get",
    //   "/api/posts/{id}/comments",
    //   queryParams,
    // );
    // await client.queryClient.ensureQueryData(commentOpts);
    // return { opts: commentOpts };
    return { opts: undefined };
  };

/**
 * View a single post (enlarged) and display its comments below.
 */
export default function CommentView() {
  const client = useContext(ClientContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  if (!opts) {
    return (
      <div className="max-w-2xl mx-auto border border-gray-200 bg-white rounded-2xl shadow-sm overflow-hidden">
        <header className="p-4 border-b border-gray-200 text-lg font-semibold bg-gray-50">
          Post
        </header>
        <div className="p-6 text-center text-gray-600">
          Post view is not available yet.
        </div>
      </div>
    );
  }
  const postData = useSuspenseQuery(opts);
  return (
    <div className="max-w-2xl mx-auto bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
      <header className="p-4 border-b border-gray-200 text-xl font-semibold bg-gray-50">
        Post
      </header>
      <div className="border-b border-gray-200 p-4 bg-gray-50">
        {/* <PostItem post={{ ...postData }} client={client!} /> */}
      </div>
      <CommentForm />
      <Outlet />
    </div>
  );
}

export function CommentsFeed() {
  const client = useContext(ClientContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof commentsLoader>>
  >;
  if (!opts) {
    return (
      <div className="max-w-2xl mx-auto bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div className="p-6 text-center text-gray-600">
          Comments are not available yet.
        </div>
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
    <div className="max-w-2xl mx-auto bg-white rounded-2xl shadow-sm border border-gray-200 divide-y">
      {data?.map((comment: components["schemas"]["Comment"]) => (
        <PostItem
          key={comment.id}
          post={{ data: comment }}
          client={client!}
        />
      )) || (
        <div className="p-6 text-center text-gray-600">No comments yet.</div>
      )}
    </div>
  );
}
