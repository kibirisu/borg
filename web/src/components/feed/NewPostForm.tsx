import { useContext, useEffect, useMemo, useRef } from "react";
import {
  type ActionFunctionArgs,
  Form,
  useActionData,
  useNavigation,
} from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";
import AppContext from "../../lib/state";
import decodeToken from "../../lib/decode";

export const action =
  (client: AppClient) =>
  async ({ request }: ActionFunctionArgs) => {
    const formData = await request.formData();
    const content = formData.get("content") as string;
    const userIdRaw = formData.get("userId");
    const userId = typeof userIdRaw === "string" ? Number(userIdRaw) : null;

    if (!content || !content.trim()) {
      throw new Error("Post content cannot be empty");
    }
    if (!userId || Number.isNaN(userId)) {
      return { form: "You must be logged in to post." };
    }

    const newPostData: components["schemas"]["NewPost"] = {
      userID: userId,
      content,
    };

    const mutationOpts = client.$api.queryOptions("post", "/api/posts", {
      body: newPostData,
    });
    await client.queryClient.ensureQueryData(mutationOpts);
    client.queryClient.invalidateQueries({
      queryKey: ["get", "/api/posts", {}],
    });

    return null;
  };

export default function NewPostForm() {
  const navigation = useNavigation();
  const errors = useActionData() as { form?: string } | undefined;
  const isSubmitting = navigation.state === "submitting";
  const textareaRef = useRef<HTMLTextAreaElement | null>(null);
  const appState = useContext(AppContext);
  const decoded = useMemo(
    () => decodeToken(appState?.tokenRef?.current ?? null),
    [appState?.tokenRef?.current],
  );
  const userId = decoded?.userId ?? null;
  const isAuthenticated = userId !== null;

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
      <input type="hidden" name="userId" value={userId ?? ""} />
      <div className="flex-1 overflow-hidden">
        {!isAuthenticated && (
          <p className="mb-2 text-sm text-gray-500">Sign in to create a post.</p>
        )}
        {errors?.form && (
          <p className="mb-2 text-sm text-red-600">{errors.form}</p>
        )}
        <textarea
          name="content"
          ref={textareaRef}
          placeholder="What's happening?"
          className="w-full resize-none outline-none text-gray-800 placeholder-gray-400 bg-transparent min-h-[60px] overflow-hidden"
          required
          disabled={!isAuthenticated}
        ></textarea>

        <div className="flex justify-end mt-2">
          <button
            type="submit"
            className="bg-indigo-600 text-white px-4 py-1.5 rounded-full text-sm font-semibold hover:bg-indigo-700 disabled:opacity-50"
            disabled={isSubmitting || !isAuthenticated}
          >
            {isSubmitting ? "Posting…" : "Post"}
          </button>
        </div>
      </div>
    </Form>
  );
}
