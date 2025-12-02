import type { Client } from "../../lib/client";

export const loader = (client: Client) => async () => {
  try {
    const opts = client.$api.queryOptions("get", "/api/posts", {});
    if (!opts) {
      throw new Error("queryOptions returned undefined");
    }
    return { opts };
  } catch (error) {
    console.error("Loader error:", error);
    return { opts: undefined as any };
  }
};
