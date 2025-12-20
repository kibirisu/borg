import { useQuery } from "@tanstack/react-query";
import { useContext } from "react";
import { type LoaderFunctionArgs, useLoaderData } from "react-router";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import { PostItem } from "./PostItem";

export const loader =
  (client: AppClient) =>
  async ({ params }: LoaderFunctionArgs) => {
    const userId = parseInt(String(params.handle));
    const queryParams = { params: { path: { id: userId } } };
    const opts = client.$api.queryOptions(
      "get",
      "/api/users/{id}/posts",
      queryParams,
    );
    await client.queryClient.ensureQueryData(opts);
    return { opts: opts };
  };

export default function Feed() {
  const client = useContext(ClientContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data, isPending } = useQuery(opts);
  if (isPending) {
    return <></>;
  }
  return (
    <div className="bg-white rounded-2xl border border-gray-200 shadow-sm divide-y">
      {data?.map((post) => (
        <PostItem key={post.id} post={{ data: post }} client={client!} />
      ))}
    </div>
  );
}
