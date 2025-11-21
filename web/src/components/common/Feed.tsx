import { MessageCircle, Repeat, Heart, Share2 } from "lucide-react";
import ReactMarkdown from "react-markdown";
import type { components } from "../../lib/api/v1";
import type { Client } from "../../lib/api/client";
import { useLoaderData, type LoaderFunctionArgs } from "react-router";
import { useQuery } from "@tanstack/react-query";

export const loader =
  (client: Client) =>
    async ({ params }: LoaderFunctionArgs) => {
      const userId = parseInt(String(params.handle));
      const queryParams = { params: { path: { id: userId } } };
      const opts = client.$api.queryOptions(
        "get",
        "/api/users/{id}/posts",
        queryParams,
      );
      console.log(opts.queryKey);
      return { opts: opts };
    };

export default function Feed() {
  const { opts: opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data, isPending } = useQuery(opts);
  if (isPending) {
    return <></>;
  }
  return (
    <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
      {data?.map((post) => (
        <Post key={post.id} data={post} />
      ))}
    </div>
  );
}

interface Data {
  data: components["schemas"]["Post"];
}

const Post = (post: Data) => {
  return (
    <div className="border-b border-gray-200 p-4 hover:bg-gray-50 transition-colors">
      <div className="flex space-x-3">
        <div className="flex-1">
          <div className="flex items-center space-x-1"></div>
          <div className="prose max-w-none text-gray-800">
            <ReactMarkdown>{post.data.content}</ReactMarkdown>
          </div>

          <div className="flex justify-between mt-3 text-gray-500 text-sm max-w-md">
            <button
              type="button"
              className="flex items-center space-x-1 hover:text-blue-500 transition"
            >
              <MessageCircle size={16} /> <span>{post.data.commentCount}</span>
            </button>
            <button
              type="button"
              className="flex items-center space-x-1 hover:text-green-500 transition"
            >
              <Repeat size={16} /> <span>{post.data.shareCount}</span>
            </button>
            <button
              type="button"
              className="flex items-center space-x-1 hover:text-pink-500 transition"
            >
              <Heart size={16} /> <span>{post.data.likeCount}</span>
            </button>
            <button
              type="button"
              className="flex items-center space-x-1 hover:text-gray-700 transition"
            >
              <Share2 size={16} />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
