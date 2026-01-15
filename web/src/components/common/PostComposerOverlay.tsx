import { type ReactNode, useEffect, useState } from "react";
import type { PostPresentable } from "./PostItem";

interface PostComposerOverlayProps {
  isOpen: boolean;
  onClose: () => void;
  replyTo?: PostPresentable | null;
  onSubmit?: (content: string) => Promise<void> | void;
  initialContent?: string;
  title?: string;
  submitLabel?: string;
}

const PostComposerOverlay = ({
  isOpen,
  onClose,
  replyTo,
  onSubmit,
  initialContent,
  title = "Share with others ❤️",
  submitLabel = "Post me",
}: PostComposerOverlayProps) => {
  const [message, setMessage] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (isOpen) {
      setMessage(initialContent ?? "");
    } else {
      setMessage("");
    }
  }, [initialContent, isOpen]);

  if (!isOpen) {
    return null;
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!message.trim() || isSubmitting) {
      return;
    }
    setIsSubmitting(true);
    try {
      if (onSubmit) {
        await onSubmit(message.trim());
      } else {
        console.info("[Composer] submitting post:", message);
      }
      onClose();
      setMessage("");
    } catch (err) {
      console.error("[Composer] failed to submit post", err);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-gray-900/40 px-4 backdrop-blur-sm"
      role="dialog"
      aria-modal="true"
      onClick={onClose}
      onKeyDown={(event) => {
        if (event.key === "Escape") {
          onClose();
        }
      }}
      tabIndex={-1}
    >
      <div
        className="relative w-full max-w-2xl rounded-3xl bg-white p-6 shadow-2xl"
        onClick={(event) => event.stopPropagation()}
        onKeyDown={(event) => event.stopPropagation()}
        role="document"
        tabIndex={-1}
      >
        <button
          type="button"
          className="absolute right-4 top-4 rounded-full bg-gray-100 p-2 text-gray-500 hover:text-gray-900"
          onClick={onClose}
          aria-label="Close composer"
        >
          <span aria-hidden="true">&times;</span>
        </button>
        <div className="space-y-4">
          <div>
            <p className="text-2xl font-semibold text-gray-900">{title}</p>
            {replyTo && "username" in replyTo.data && replyTo.data.username && (
              <p className="text-sm text-gray-500">
                Responding to{" "}
                <span className="font-medium">@{replyTo.data.username}</span>
              </p>
            )}
          </div>
          <form onSubmit={handleSubmit} className="space-y-3">
            <label htmlFor="post-composer" className="sr-only">
              Your post
            </label>
            <textarea
              id="post-composer"
              rows={4}
              className="w-full rounded-2xl border border-gray-200 bg-gray-50 p-4 text-sm text-gray-900 placeholder:text-gray-500 focus:border-indigo-500 focus:ring-indigo-500 focus:outline-none"
              placeholder="Write a comment..."
              value={message}
              onChange={(event) => setMessage(event.target.value)}
              required
              disabled={isSubmitting}
            />
            <div className="flex items-center justify-between border-t border-gray-200 pt-3">
              <div className="flex items-center gap-2">
                <IconButton label="Add emoji">
                  <svg
                    className="h-5 w-5"
                    aria-hidden="true"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <title>Add emoji</title>
                    <path
                      stroke="currentColor"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M15 9h.01M8.99 9H9m12 3a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM6.6 13a5.5 5.5 0 0 0 10.81 0H6.6Z"
                    />
                  </svg>
                </IconButton>
              </div>
              <button
                type="submit"
                className="rounded-full bg-indigo-600 px-6 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus:outline-none focus:ring-4 focus:ring-indigo-300"
                disabled={isSubmitting}
              >
                {isSubmitting ? "Posting…" : submitLabel}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

interface IconButtonProps {
  label: string;
  children: ReactNode;
}

const IconButton = ({ label, children }: IconButtonProps) => (
  <button
    type="button"
    className="rounded-md p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-900"
  >
    {children}
    <span className="sr-only">{label}</span>
  </button>
);

export default PostComposerOverlay;
