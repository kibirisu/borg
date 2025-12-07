import type React from "react";
import { useEffect, useRef } from "react";
import { type ActionFunctionArgs, Form, useNavigation } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";

export const action =
  (client: AppClient) =>
  async ({ request, params }: ActionFunctionArgs) => {
    const formData = await request.formData();
    const content = formData.get("content") as string;

    if (!content || !content.trim()) {
      throw new Error("Post content cannot be empty");
    }

    const newPostData: components["schemas"]["NewPost"] = {
      userID: 1,
      content,
    };

    const mutationOpts = client.$api.queryOptions("post", "/api/posts", {
      body: newPostData,
    });

    await client.queryClient.ensureQueryData(mutationOpts);
    const listOpts = client.$api.queryOptions("get", "/api/posts", {});
    client.queryClient.invalidateQueries({
      queryKey: listOpts.queryKey,
    });

    return null;
  };

export default function NewPostForm() {
  const navigation = useNavigation();
  const isSubmitting = navigation.state === "submitting";
  const textareaRef = useRef<HTMLTextAreaElement | null>(null);

  useEffect(() => {
    if (!isSubmitting && textareaRef.current) {
      textareaRef.current.value = "";
    }
  }, [isSubmitting]);

  return (
    <Form
      method="post"
      className="border-b border-gray-300 p-3 flex items-start space-x-3 bg-white"
    >
      <div className="flex-1 overflow-hidden">
        <textarea
          name="content"
          ref={textareaRef}
          placeholder="What's happening?"
          className="w-full resize-none outline-none text-gray-800 placeholder-gray-400 bg-transparent min-h-[60px] overflow-hidden"
          required
        ></textarea>

        <div className="flex justify-end mt-2">
          <button
            type="submit"
            className="bg-indigo-600 text-white px-4 py-1.5 rounded-full text-sm font-semibold hover:bg-indigo-700 disabled:opacity-50"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Posting…" : "Post"}
          </button>
        </div>
      </div>
    </Form>
  );
}
