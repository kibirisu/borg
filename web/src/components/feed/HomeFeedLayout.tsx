import { Outlet } from "react-router";
import type { AppClient } from "../../lib/client";
import NewPostForm from "./NewPostForm";

export const loader = (client: AppClient) => async () => {
  // Posts listing not implemented on backend yet; skip prefetch.
  return { opts: undefined, disabled: true };
};
export default function MainFeed() {
  return (
    <div className="max-w-2xl mx-auto border-x border-gray-300 min-h-screen bg-white">
      <NewPostForm />
      <Outlet />
    </div>
  );
}
