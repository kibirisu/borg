import { useEffect, useRef } from "react";
import { type ActionFunctionArgs, Form, useNavigation } from "react-router";
import type { components } from "../../lib/api/v1";
import type { Client } from "../../lib/client";

export const action =
  (client: Client) =>
  async ({ request, params }: ActionFunctionArgs) => {
    if (!params.postId) {
      throw new Error("No post ID provided");
    }

    const postId = parseInt(params.postId);

    const formData = await request.formData();
    const content = formData.get("content") as string;
    if (!content) {
      throw new Error("Comment content cannot be empty");
    }

    const newCommentOps: components["schemas"]["NewComment"] = {
      postID: postId,
      userID: 1, //TODO
      content,
    };

    const mutationOpts = client.$api.queryOptions(
      "post",
      "/api/posts/{id}/comments",
      {
        params: { path: { id: postId } },
        body: newCommentOps,
      },
    );
    client.queryClient.invalidateQueries({
      queryKey: ["get", "/api/posts/{id}/comments", { id: postId }],
    });

    await client.queryClient.ensureQueryData(mutationOpts);
    return null;
  };

export default function CommentForm() {
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";

  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    if (!isSubmitting && textareaRef.current) {
      textareaRef.current.value = "";
    }
  }, [isSubmitting]);

  return (
    <Form
      method="post"
      className="p-4 border-b border-gray-200 bg-white flex flex-col gap-2"
    >
      <textarea
        ref={textareaRef}
        name="content"
        required
        placeholder="Write a comment..."
        className="border p-2 rounded-md w-full resize-none"
        rows={3}
      />

      <button
        type="submit"
        disabled={isSubmitting}
        className="btn btn-primary self-end"
      >
        {isSubmitting ? "Posting…" : "Post"}
      </button>
    </Form>
  );
}
