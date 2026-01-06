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
  onSelect?: (post: PostPresentable) => void;
  showActions?: boolean;
  onCommentClick?: (post: PostPresentable) => void;
}

export const PostItem = ({
  post,
  client,
  onSelect,
  showActions = false,
  onCommentClick,
}: PostProps) => {
  const likeAction = async () => {
    console.warn("Likes API is not available yet");
  };

  const shareAction = async () => {
    console.warn("Shares API is not available yet");
  };

  const handleSelect = () => {
    if (onSelect) {
      onSelect(post);
    }
  };

  return (
    <div
      className="border-b border-gray-200 p-4 hover:bg-gray-50 transition-colors cursor-pointer"
      onClick={handleSelect}
    >
      <div className="flex space-x-3">
        <div className="flex-1">
          <div className="flex items-start justify-between mb-2">
            {"username" in post.data && (
              <div className="flex items-center space-x-1">
                <Link
                  to={`/profile/${post.data.userID}`}
                  className="hover:underline font-semibold text-gray-900"
                  onClick={(event) => event.stopPropagation()}
                >
                  {post.data.username}
                </Link>
              </div>
            )}
            {showActions && (
              <div className="flex items-center gap-2">
                <button
                  type="button"
                  className="text-black bg-white box-border border border-black hover:bg-gray-100 hover:cursor-pointer shadow-xs font-medium leading-5 rounded-full text-sm px-4 py-2.5 focus:outline-none"
                  onClick={(event) => {
                    event.stopPropagation();
                  }}
                >
                  <i className="bi bi-pencil mr-1" aria-hidden="true"></i>
                  Edit
                </button>
                <button
                  type="button"
                  className="text-black bg-white box-border border border-black hover:bg-gray-100 hover:cursor-pointer shadow-xs font-medium leading-5 rounded-full text-sm px-4 py-2.5 focus:outline-none"
                  onClick={(event) => {
                    event.stopPropagation();
                  }}
                >
                  <i className="bi bi-trash3 mr-1" aria-hidden="true"></i>
                  Delete
                </button>
              </div>
            )}
          </div>
          <div className="flex items-center space-x-1"></div>
          <div className="prose max-w-none text-gray-800">
            <ReactMarkdown>{post.data.content}</ReactMarkdown>
          </div>

          <div className="flex justify-between mt-3 text-gray-500 text-sm max-w-md">
            {"commentCount" in post.data && (
              <button
                type="button"
                className="flex items-center space-x-1 hover:text-blue-500 transition"
                onClick={(event) => {
                  event.stopPropagation();
                  if (onCommentClick) {
                    onCommentClick(post);
                  }
                }}
              >
                <MessageCircle size={16} /> <span>{post.data.commentCount}</span>
              </button>
            )}
            {"shareCount" in post.data && (
              <form
                action={shareAction}
                onClick={(event) => event.stopPropagation()}
              >
                <button
                  type="submit"
                  className="flex items-center space-x-1 hover:text-green-500 transition"
                >
                  <Repeat size={16} /> <span>{post.data.shareCount}</span>
                </button>
              </form>
            )}
            {"likeCount" in post.data && (
              <form
                action={likeAction}
                onClick={(event) => event.stopPropagation()}
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
                onClick={(event) => event.stopPropagation()}
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
