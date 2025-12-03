import type { OpenapiQueryClient } from "openapi-react-query";
import { createContext, type JSX } from "react";
import type { paths } from "./api/v1";

export default interface AppState {
  $api: OpenapiQueryClient<paths>;
  username: string | null;
}

export const AppContext = createContext<AppState | undefined>(undefined);

// export const AppContextProvider = AppContext.Provider;

interface Children {
  children: JSX.Element;
  state: AppState;
}

export const AppStateProvider = (props: Children) => {
  const { children, state } = props;
  return <AppContext.Provider value={state}>{children}</AppContext.Provider>;
};
