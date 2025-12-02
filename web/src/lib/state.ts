import type { OpenapiQueryClient } from "openapi-react-query";
import { createContext } from "react";
import type { paths } from "./api/v1";

export default interface AppState {
  $api: OpenapiQueryClient<paths>;
  username: string | null;
}

export const AppContext = createContext<AppState | undefined>(undefined);

export const AppContextProvider = AppContext.Provider;
