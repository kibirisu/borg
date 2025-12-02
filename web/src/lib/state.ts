import type { OpenapiQueryClient } from "openapi-react-query";
import { type Context, createContext } from "react";
import type { paths } from "./api/v1";

export default interface AppState {
  $api: OpenapiQueryClient<paths>;
}

export function createAppContext(
  $api: OpenapiQueryClient<paths>,
): Context<AppState> {
  return createContext({ $api });
}
