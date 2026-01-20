import { useContext, useEffect, useRef } from "react";
import {
  type ActionFunctionArgs,
  Form,
  useActionData,
  useNavigation,
} from "react-router";
import type { AppClient } from "../../lib/client";
import AppContext from "../../lib/state";

export const action =
  (client: AppClient) =>
  async ({ request, params }: ActionFunctionArgs) => {
    const fallbackPostId = params.postId ? String(params.postId) : "";
    const formData = await request.formData();
    const replyToId = formData.get("replyToId")?.toString() ?? "";
    const postId = replyToId || fallbackPostId;

    if (!postId) {
      return { form: "No post ID provided" };
    }

    const contentRaw = formData.get("content")?.toString() ?? "";
    const userId = formData.get("userId")?.toString() ?? "";

    if (!contentRaw.trim()) {
      return { form: "Comment content cannot be empty" };
    }
    if (!userId) {
      return { form: "User not authenticated" };
    }

    const res = await client.fetchClient.POST("/api/statuses", {
      body: { status: contentRaw.trim(), in_reply_to_id: postId },
    });

    if (res.error) {
      return { form: "Failed to post comment" };
    }

    client.queryClient.invalidateQueries({
      queryKey: [
        "get",
        "/api/statuses/{id}/replies",
        { params: { path: { id: postId } } },
      ],
    });
    client.queryClient.invalidateQueries({
      queryKey: [
        "get",
        "/api/statuses/{id}",
        { params: { path: { id: postId } } },
      ],
    });
    client.queryClient.invalidateQueries({
      queryKey: ["account-statuses"],
      exact: false,
    });
    client.queryClient.invalidateQueries({
      queryKey: ["get", "/api/timelines/home", {}],
    });
    return null;
  };

export default function CommentForm({ postId }: { postId?: string }) {
  const appState = useContext(AppContext);
  const errors = useActionData() as { form?: string } | undefined;
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";
  const userId = appState?.userId ?? null;
  const isAuthenticated = userId !== null;

  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    if (!isSubmitting && textareaRef.current) {
      textareaRef.current.value = "";
    }
  }, [isSubmitting]);

  return (
    <Form
      method="post"
      className="fixed bottom-0 left-0 right-0 z-50 w-full border-t border-gray-200 bg-white flex flex-col gap-3 p-4 shadow-[0_-8px_24px_-12px_rgba(0,0,0,0.3)]"
    >
      <input type="hidden" name="userId" value={userId ?? ""} />
      <input type="hidden" name="replyToId" value={postId ?? ""} />
      {!isAuthenticated && (
        <div className="rounded-lg bg-yellow-50 border border-dashed border-yellow-200 p-3 text-sm text-gray-700">
          Sign in to comment.
        </div>
      )}
      {errors?.form && (
        <div className="rounded-lg bg-red-50 border border-red-200 p-3 text-sm text-red-700">
          {errors.form}
        </div>
      )}
      <textarea
        ref={textareaRef}
        name="content"
        required
        placeholder="Write a comment..."
        className="border-none p-3 rounded-xl w-full resize-none shadow-sm focus:outline-none bg-gray-50 text-black"
        rows={3}
        disabled={!isAuthenticated}
      />

      <button
        type="submit"
        disabled={isSubmitting || !isAuthenticated}
        className="self-end rounded-full bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm disabled:opacity-50"
      >
        {isSubmitting ? "Posting…" : "Post"}
      </button>
    </Form>
  );
}
