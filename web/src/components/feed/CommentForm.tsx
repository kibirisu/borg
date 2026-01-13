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
    if (!params.postId) {
      return { form: "No post ID provided" };
    }

    const postId = Number(params.postId);
    const formData = await request.formData();
    const contentRaw = formData.get("content")?.toString() ?? "";
    const userIdRaw = formData.get("userId")?.toString() ?? "";
    const userId = Number(userIdRaw);

    if (!contentRaw.trim()) {
      return { form: "Comment content cannot be empty" };
    }
    if (!userId || Number.isNaN(userId)) {
      return { form: "User not authenticated" };
    }

    const body = {
      postID: postId,
      userID: userId,
      content: contentRaw.trim(),
    };

    const res = await client.fetchClient.POST("/api/posts/{id}/comments", {
      params: { path: { id: postId } },
      body,
    });

    if (res.error) {
      return { form: "Failed to post comment" };
    }

    client.queryClient.invalidateQueries({ queryKey: ["comments", postId] });
    client.queryClient.invalidateQueries({ queryKey: ["user-posts"] });
    return null;
  };

export default function CommentForm() {
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
      className="p-4 border-t border-gray-200 bg-white flex flex-col gap-3"
    >
      <input type="hidden" name="userId" value={userId ?? ""} />
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
        {isSubmitting ? "Commenting…" : "Comment"}
      </button>
    </Form>
  );
}
