import type { AppClient } from "../../lib/client";
import { redirect, type ActionFunctionArgs } from "react-router";
import type { components } from "../../lib/api/v1";

export const toggleLikeAction =
  (client: AppClient) =>
  async ({ request, params }: ActionFunctionArgs) => {
    if (!params.postId) {
      throw new Response("Post ID required", { status: 400 });
    }

    const postId = parseInt(params.postId);
    if (isNaN(postId)) {
        throw new Response("Invalid Post ID format", { status: 400 });
    }
    const newCommentOps: components["schemas"]["NewLike"] = {
      postID: postId,
      userID: 1, //TODO
    };

    const mutationOpts = client.$api.queryOptions(
        "post",
        "/api/posts/{id}/likes",
        {
          params: { path: { id: postId } },
          body: newCommentOps,
        },
    );

    client.queryClient.invalidateQueries({
        queryKey: ["get", "/api/posts/{id}", { id: postId }], 
    });
    client.queryClient.invalidateQueries({
        queryKey: ["get", "/api/posts/{id}/likes", { id: postId }],
    });
    client.queryClient.invalidateQueries({
        queryKey: ["get", "/api/posts"],
    });

    await client.queryClient.ensureQueryData(mutationOpts);

    return null;
  };
