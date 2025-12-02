import { Outlet } from "react-router";
import type { Client } from "../../lib/client";
import NewPostForm from "./NewPostForm";

export const loader = (client: Client) => async () => {
  try {
    const opts = client.$api.queryOptions("get", "/api/posts", {});
    if (!opts) {
      throw new Error("queryOptions returned undefined");
    }
    client.queryClient.prefetchQuery({ ...opts, staleTime: 0 });
    return { opts };
  } catch (error) {
    console.error("Loader error:", error);
    return { opts: undefined as any };
  }
};
export default function MainFeed() {
  return (
    <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
      <NewPostForm />
      <Outlet />
    </div>
  );
}
