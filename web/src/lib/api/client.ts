import type { paths } from "./v1";
import { QueryClient } from "@tanstack/react-query";
import createFetchClient, { type Middleware } from "openapi-fetch";
import createClient, { type OpenapiQueryClient } from "openapi-react-query";

export interface Client {
  queryClient: QueryClient;
  $api: OpenapiQueryClient<paths>;
}

export default function newClient(): Client {
  const fetchClient = createFetchClient<paths>();
  fetchClient.use(authMiddleware);
  const $api = createClient(fetchClient);
  const queryClient = new QueryClient({
    defaultOptions: { queries: { staleTime: 1000 * 10 } },
  });
  accessToken = localStorage.getItem("jwt");
  return { queryClient, $api };
}

let accessToken: string | null = null;

const authMiddleware: Middleware = {
  async onRequest({ request }) {
    if (accessToken) {
      request.headers.set("Authorization", `Bearer ${accessToken}`);
    }
    return request;
  },
};
