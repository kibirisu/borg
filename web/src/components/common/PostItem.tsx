import { Heart, MessageCircle, Repeat, Share2 } from "lucide-react";
import { useContext } from "react";
import ReactMarkdown from "react-markdown";
import { Link } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";
import AppContext from "../../lib/state";

type Status = components["schemas"]["Status"];

export type PostPresentable = {
  data: Status;
};

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
  const appState = useContext(AppContext);
  const currentUserId = appState?.userId ?? null;
  const data = post.data;
  if (!data) {
    return null;
  }
  const renderData = data.reblog ?? data;
  const resharedBy = data.reblog ? data.account : null;
  const authorId = renderData.account?.id ?? data.account?.id;
  const authorName =
    renderData.account?.display_name ||
    renderData.account?.username ||
    data.account?.display_name ||
    data.account?.username;
  const authorHandleRaw =
    renderData.account?.acct ||
    renderData.account?.username ||
    data.account?.acct ||
    data.account?.username;
  const profileId = renderData.account?.id ?? data.account?.id ?? null;
  const authorHandle = authorHandleRaw
    ? authorHandleRaw.startsWith("@")
      ? authorHandleRaw.slice(1)
      : authorHandleRaw
    : "";
  const content = renderData.content;
  const commentCount = renderData.replies_count;
  const shareCount = renderData.reblogs_count;
  const likeCount = renderData.favourites_count;
  const isFavourited = Boolean(renderData.favourited);
  const isReblogged = Boolean(renderData.reblogged);

  const likeAction = async () => {
    if (!currentUserId) {
      console.warn("User not authenticated, cannot like");
      return;
    }
    try {
      const endpoint = renderData.favourited
        ? "/api/statuses/{id}/unfavourite"
        : "/api/statuses/{id}/favourite";
      await client.fetchClient.POST(endpoint, {
        params: { path: { id: String(renderData.id) } },
      });
      client.queryClient.invalidateQueries({
        queryKey: ["account-statuses", authorId],
      });
      client.queryClient.invalidateQueries({
        queryKey: ["get", "/api/accounts/{id}/statuses"],
      });
    } catch (err) {
      console.error("Failed to like post", err);
    }
  };

  const shareAction = async () => {
    if (!currentUserId) {
      console.warn("User not authenticated, cannot share");
      return;
    }
    try {
      const endpoint = renderData.reblogged
        ? "/api/statuses/{id}/unreblog"
        : "/api/statuses/{id}/reblog";
      const targetId = renderData.id;
      await client.fetchClient.POST(endpoint, {
        params: { path: { id: String(targetId) } },
      });
      client.queryClient.invalidateQueries({
        queryKey: ["account-statuses", authorId],
      });
      client.queryClient.invalidateQueries({
        queryKey: ["get", "/api/accounts/{id}/statuses"],
      });
    } catch (err) {
      console.error("Failed to share post", err);
    }
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
      role="button"
      tabIndex={0}
      onKeyDown={(event) => {
        if (event.key === "Enter" || event.key === " ") {
          event.preventDefault();
          handleSelect();
        }
      }}
    >
      <div className="flex space-x-3">
        <div className="flex-1">
          {resharedBy && (
            <div className="mb-2 text-xs">
              <span className="text-green-600">Reshared by </span>
              <span className="font-medium text-green-700">
                @{String(resharedBy.acct || resharedBy.username).replace(/^@+|@+$/g, "")}
              </span>
            </div>
          )}
          <div className="flex items-start justify-between mb-2">
            {authorName && (
              <div className="flex items-center space-x-1">
                {profileId ? (
                  <Link
                    to={`/profile/${profileId}`}
                    className="hover:underline font-semibold text-gray-900"
                    onClick={(event) => event.stopPropagation()}
                  >
                    {authorName}
                  </Link>
                ) : (
                  <span className="font-semibold text-gray-900">
                    {authorName}
                  </span>
                )}
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
            <ReactMarkdown>{content}</ReactMarkdown>
          </div>

          <div className="flex justify-between mt-3 text-gray-500 text-sm max-w-md">
            {commentCount !== undefined && (
              <button
                type="button"
                className="flex items-center space-x-1 hover:text-blue-500 transition cursor-pointer"
                onClick={(event) => {
                  event.stopPropagation();
                  if (onCommentClick) {
                    onCommentClick(post);
                  }
                }}
              >
                <MessageCircle size={16} />{" "}
                <span>{commentCount}</span>
              </button>
            )}
            {shareCount !== undefined && (
              <button
                type="button"
                className={`flex items-center space-x-1 transition ${
                  isReblogged
                    ? "text-green-600"
                    : "text-gray-500 hover:text-green-500"
                }`}
                onClick={(event) => {
                  event.stopPropagation();
                  void shareAction();
                }}
              >
                <Repeat size={16} /> <span>{shareCount}</span>
              </button>
            )}
            {likeCount !== undefined && (
              <button
                type="button"
                className={`flex items-center space-x-1 transition ${
                  isFavourited
                    ? "text-pink-600"
                    : "text-gray-500 hover:text-pink-500"
                }`}
                onClick={(event) => {
                  event.stopPropagation();
                  void likeAction();
                }}
              >
                <Heart size={16} />
                <span>{likeCount}</span>
              </button>
            )}
            {shareCount !== undefined && (
              <button
                type="button"
                className="flex items-center space-x-1 hover:text-gray-700 transition cursor-pointer"
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
