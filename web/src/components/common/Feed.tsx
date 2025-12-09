import { useContext } from "react";
import { ClientContext } from "../../lib/client";
import { useQuery } from "@tanstack/react-query";
import { Heart, MessageCircle, Repeat, Share2 } from "lucide-react";
import ReactMarkdown from "react-markdown";
import { Form, Link, type LoaderFunctionArgs, useFetcher, useLoaderData } from "react-router";
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
  const client = useContext(ClientContext);
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
        <Post key={post.id} post={{ data: post }} client={client!} />
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

interface PostProps {
  post: PostPresentable;
  client: AppClient;
}

export const Post = ({ post, client }: PostProps) => {
  const { mutate: like } = client.$api.useMutation(
      "post",
      "/api/posts/{id}/likes",
      {
        onSuccess: () => {
          client.queryClient.invalidateQueries({ queryKey: ["get", "/api/posts", {}] });
          client.queryClient.invalidateQueries({ queryKey: ["get", "/api/posts/{id}/comments", { id: post.data.id }], });
        },
      }
  );

  const likeAction = async () => {
    const newCommentOps: components["schemas"]["NewLike"] = {
      postID: post.data.id,
      userID: 1, //TODO
    };

    like({ params: { 
            path: { id: post.data.id }
        }, body: newCommentOps });
  };

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
              <form
                action={likeAction}
              >
                <button
                  type="submit"
                  className="flex items-center space-x-1 hover:text-pink-500 transition"
                >
                  <Heart size={16} /> 
                  <span>{post.data.likeCount}</span>
                </button>
              </form>
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
