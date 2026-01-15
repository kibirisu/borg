import {
  type QueryKey,
  type UseQueryOptions,
  useQuery,
} from "@tanstack/react-query";
import { useContext } from "react";
import { type LoaderFunctionArgs, useLoaderData } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import { PostItem } from "./PostItem";

export const loader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    const handle = params.handle;
    const userId = handle ? Number(handle) : NaN;
    if (!handle || Number.isNaN(userId)) {
      return { opts: undefined };
    }

    const opts = client.$api.queryOptions("get", "/api/users/{id}/posts", {
      params: { path: { id: userId } },
    });
    await client.queryClient.ensureQueryData(opts);
    return { opts };
  };

export default function Feed() {
  const client = useContext(ClientContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;

  const queryOptions = (opts ??
    ({
      queryKey: ["posts-feed-disabled"],
      queryFn: async () => [] as components["schemas"]["Post"][],
      enabled: false,
    } satisfies UseQueryOptions<
      components["schemas"]["Post"][],
      Error,
      components["schemas"]["Post"][],
      QueryKey
    >)) as UseQueryOptions<
    components["schemas"]["Post"][],
    Error,
    components["schemas"]["Post"][],
    QueryKey
  >;

  const { data, isPending } = useQuery(queryOptions);

  if (!opts || !client) {
    return null;
  }

  const posts = Array.isArray(data) ? data : [];

  if (isPending) {
    return (
      <div className="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm">
        <p className="text-sm text-gray-500">Loading feed…</p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      {posts.map((post: components["schemas"]["Post"]) => (
        <PostItem key={post.id} post={{ data: post }} client={client} />
      ))}
    </div>
  );
}
