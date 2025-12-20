import { useQuery } from "@tanstack/react-query";
import { useContext } from "react";
import { useLoaderData } from "react-router";
import type { AppClient } from "../../lib/client";
import ClientContext from "../../lib/client";
import { PostItem } from "../common/PostItem";
import Sidebar from "../layout/Sidebar";

export const loader = (client: AppClient) => async () => {
  const opts = client.$api.queryOptions("get", "/api/posts", {});
  await client.queryClient.ensureQueryData(opts);
  return { opts };
};

export default function ExplorePage() {
  const client = useContext(ClientContext);
  const { opts } = useLoaderData() as Awaited<
    ReturnType<ReturnType<typeof loader>>
  >;
  const { data, isPending } = useQuery(opts);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="grid grid-cols-[1fr_256px] gap-6">
        <main className="px-6 py-6 space-y-6">
          <section className="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
            <h1 className="text-2xl font-semibold text-gray-800">
              Explore trending posts
            </h1>
            <p className="text-gray-500">
              See what others are talking about right now.
            </p>
          </section>
          <section className="bg-white rounded-2xl border border-gray-200 p-4 space-y-4 min-h-[400px]">
            {isPending && (
              <p className="text-gray-500 text-center">Loading…</p>
            )}
            {!isPending &&
              data?.map((post) => (
                <PostItem key={post.id} post={{ data: post }} client={client!} />
              ))}
            {!isPending && !data?.length && (
              <p className="text-center text-gray-500">
                Nothing posted yet.
              </p>
            )}
          </section>
        </main>
        <Sidebar />
      </div>
    </div>
  );
}
