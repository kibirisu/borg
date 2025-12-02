import { QueryClient } from "@tanstack/react-query";
import createFetchClient, { type Middleware } from "openapi-fetch";
import createClient, { type OpenapiQueryClient } from "openapi-react-query";
import type { paths } from "./api/v1.d.ts";
import decodeToken from "./decode.ts";

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
  return { queryClient, $api };
}

export function checkToken(): string | null {
  const token = localStorage.getItem("jwt");
  if (token) {
    const username = decodeToken(token);
    if (username) {
      accessToken = token;
      return username;
    }
  }
  return null;
}

let accessToken: string | null = null;

const authMiddleware: Middleware = {
  async onRequest({ request }) {
    if (accessToken) {
      request.headers.set("Authorization", accessToken);
    }
    return request;
  },

  async onResponse({ response, schemaPath }) {
    if (schemaPath === "/api/auth/login") {
      const token = response.headers.get("Authorization");
      if (token) {
        accessToken = token;
        localStorage.setItem("jwt", token);
      }
    }
    return response;
  },
};
