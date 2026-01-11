import { useQuery } from "@tanstack/react-query";
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

  if (!opts || !client) {
    return null;
  }
  const { data, isPending } = useQuery(opts);

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
