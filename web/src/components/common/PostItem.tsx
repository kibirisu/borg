import { Heart, MessageCircle, Repeat, Share2 } from "lucide-react";
import ReactMarkdown from "react-markdown";
import { Link } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";

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

export const PostItem = ({ post, client }: PostProps) => {
  const { mutateAsync: like } = client.$api.useMutation(
    "post",
    "/api/posts/{id}/likes",
    {
      onSuccess: () => {
        client.queryClient.invalidateQueries({
          queryKey: ["get", "/api/posts", {}],
        });
        client.queryClient.invalidateQueries({
          queryKey: [
            "get",
            "/api/posts/{id}",
            { params: { path: { id: post.data.id } } },
          ],
        });
      },
    },
  );

  const likeAction = async () => {
    const newCommentOps: components["schemas"]["NewLike"] = {
      postID: post.data.id,
      userID: 1, //TODO: API needs refactor, the user data will be passed with JWT
    };

    await like({
      params: {
        path: { id: post.data.id },
      },
      body: newCommentOps,
    });
  };

  const { mutateAsync: share } = client.$api.useMutation(
    "post",
    "/api/posts/{id}/shares",
    {
      onSuccess: () => {
        client.queryClient.invalidateQueries({
          queryKey: ["get", "/api/posts", {}],
        });
        client.queryClient.invalidateQueries({
          queryKey: [
            "get",
            "/api/posts/{id}",
            { params: { path: { id: post.data.id } } },
          ],
        });
      },
    },
  );

  const shareAction = async () => {
    const newCommentOps: components["schemas"]["NewShare"] = {
      postID: post.data.id,
      userID: 1, //TODO
    };

    await share({
      params: {
        path: { id: post.data.id },
      },
      body: newCommentOps,
    });
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
              <form action={shareAction}>
                <button
                  type="submit"
                  className="flex items-center space-x-1 hover:text-green-500 transition"
                >
                  <Repeat size={16} /> <span>{post.data.shareCount}</span>
                </button>
              </form>
            )}
            {"likeCount" in post.data && (
              <form action={likeAction}>
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
