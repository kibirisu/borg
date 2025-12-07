import { useQuery } from "@tanstack/react-query";
import { Heart, MessageCircle, Repeat, Share2 } from "lucide-react";
import ReactMarkdown from "react-markdown";
import { Link, type LoaderFunctionArgs, useLoaderData } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";

export const loader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    const userId = parseInt(String(params.handle));
    const queryParams = { params: { path: { id: userId } } };
    const opts = client.$api.queryOptions(
      "get",
      "/api/users/{id}/posts",
      queryParams,
    );
    await client.queryClient.ensureQueryData(opts);
    return { opts: opts };
  };

export default function Feed() {
  const { opts } = useLoaderData() as Awaited<
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

interface PostData {
  data: components["schemas"]["Post"];
}
interface CommentData {
  data: components["schemas"]["Comment"];
}

export type PostPresentable = PostData | CommentData;

export const Post = (post: PostPresentable) => {
  return (
    <div className="border-b border-gray-200 p-4 hover:bg-gray-50 transition-colors">
      <div className="flex space-x-3">
        <div className="flex-1">
          {"username" in post.data && (
            <div className="flex items-center space-x-1 mb-2">
              <Link
                to={`/profile/${post.data.userID}`}
                className="hover:underline font-semibold text-gray-900"
              >
                {post.data.username}
              </Link>
            </div>
          )}
          <div className="flex items-center space-x-1"></div>
          <div className="prose max-w-none text-gray-800">
            <ReactMarkdown>{post.data.content}</ReactMarkdown>
          </div>

          <div className="flex justify-between mt-3 text-gray-500 text-sm max-w-md">
            {"commentCount" in post.data && (
              <Link to={`/post/${post.data.id}`}>
                <button
                  type="button"
                  className="flex items-center space-x-1 hover:text-blue-500 transition"
                >
                  <MessageCircle size={16} />{" "}
                  <span>{post.data.commentCount}</span>
                </button>
              </Link>
            )}
            {"shareCount" in post.data && (
              <button
                type="button"
                className="flex items-center space-x-1 hover:text-green-500 transition"
              >
                <Repeat size={16} /> <span>{post.data.shareCount}</span>
              </button>
            )}
            {"likeCount" in post.data && (
              <button
                type="button"
                className="flex items-center space-x-1 hover:text-pink-500 transition"
              >
                <Heart size={16} /> <span>{post.data.likeCount}</span>
              </button>
            )}
            {"shareCount" in post.data && (
              <button
                type="button"
                className="flex items-center space-x-1 hover:text-gray-700 transition"
              >
                <Share2 size={16} />
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
