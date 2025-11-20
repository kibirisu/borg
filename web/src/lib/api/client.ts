import type { paths } from "./v1";
import { QueryClient } from "@tanstack/react-query";
import createFetchClient from "openapi-fetch";
import createClient, { type OpenapiQueryClient } from "openapi-react-query";

export interface Client {
  queryClient: QueryClient;
  $api: OpenapiQueryClient<paths>;
}

export function newClient(): Client {
  const fetchClient = createFetchClient<paths>();
  const $api = createClient(fetchClient);
  const queryClient = new QueryClient();
  return { queryClient, $api };
}
