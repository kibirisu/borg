import { useEffect, useRef } from "react";
import { type ActionFunctionArgs, Form, useNavigation } from "react-router";
import type { AppClient } from "../../lib/client";

export const action =
  (client: AppClient) =>
  async ({ request, params }: ActionFunctionArgs) => {
    // if (!params.postId) {
    //   throw new Error("No post ID provided");
    // }
    //
    // const postId = parseInt(params.postId);
    //
    // const formData = await request.formData();
    // const content = formData.get("content") as string;
    // if (!content) {
    //   throw new Error("Comment content cannot be empty");
    // }
    //
    // const newCommentOps: components["schemas"]["NewComment"] = {
    //   postID: postId,
    //   userID: 1, //TODO
    //   content,
    // };
    //
    // const mutationOpts = client.$api.queryOptions(
    //   "post",
    //   "/api/posts/{id}/comments",
    //   {
    //     params: { path: { id: postId } },
    //     body: newCommentOps,
    //   },
    // );
    // client.queryClient.invalidateQueries({
    //   queryKey: ["get", "/api/posts/{id}/comments", { id: postId }],
    // });
    //
    // await client.queryClient.ensureQueryData(mutationOpts);
    return { form: "Comments are not available yet." };
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
      className="p-4 border-t border-gray-200 bg-white flex flex-col gap-3 rounded-b-2xl"
    >
      <div className="rounded-lg bg-gray-50 border border-dashed border-gray-300 p-3 text-sm text-gray-600">
        Commenting is temporarily disabled while the API endpoint is being
        implemented.
      </div>
      <textarea
        ref={textareaRef}
        name="content"
        required
        placeholder="Write a comment..."
        className="border border-gray-300 p-3 rounded-xl w-full resize-none shadow-sm focus:border-indigo-500 focus:ring-indigo-500 bg-gray-50"
        rows={3}
        disabled
      />

      <button
        type="submit"
        disabled={isSubmitting}
        className="self-end rounded-full bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm disabled:opacity-50"
      >
        {isSubmitting ? "Posting…" : "Post"}
      </button>
    </Form>
  );
}
