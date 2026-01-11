import { useQuery } from "@tanstack/react-query";
import { useContext } from "react";
import { type LoaderFunctionArgs, useLoaderData } from "react-router";
import type { components } from "../../lib/api/v1";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import { PostItem } from "./PostItem";

export const loader =
  (_client: AppClient) =>
  async (_args: LoaderFunctionArgs) => {
    // const userId = parseInt(String(params.handle));
    // const queryParams = { params: { path: { id: userId } } };
    // const opts = client.$api.queryOptions(
    //   "get",
    //   "/api/users/{id}/posts",
    //   queryParams,
    // );
    // await client.queryClient.ensureQueryData(opts);
    // return { opts };
    return { opts: undefined };
  };

export default function Feed() {
  const client = useContext(ClientContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;

  const queryOptions =
    opts ??
    ({
      queryKey: ["posts-feed-disabled"],
      queryFn: async () => [] as components["schemas"]["Post"][],
      enabled: false,
    } satisfies Parameters<typeof useQuery>[0]);

  const { data, isPending } = useQuery<components["schemas"]["Post"][]>(
    queryOptions,
  );

  if (!opts || !client) {
    return null;
  }

  if (isPending) {
    return (
      <div className="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm">
        <p className="text-sm text-gray-500">Loading feed…</p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-2xl border border-gray-200 shadow-sm divide-y">
      {data?.map((post: components["schemas"]["Post"]) => (
        <PostItem key={post.id} post={{ data: post }} client={client} />
      ))}
    </div>
  );
}
