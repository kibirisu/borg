import { type QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type createFetchClient from "openapi-fetch";
import type { Middleware } from "openapi-fetch";
import type { OpenapiQueryClient } from "openapi-react-query";
import {
  createContext,
  type Dispatch,
  type JSX,
  type RefObject,
  type SetStateAction,
  useContext,
  useEffect,
} from "react";
import type { paths } from "./api/v1";
import AppContext from "./state.tsx";

export interface AppClient {
  queryClient: QueryClient;
  fetchClient: ReturnType<typeof createFetchClient<paths>>;
  $api: OpenapiQueryClient<paths>;
}

const createMiddleware = (
  token: RefObject<string | null>,
  setToken: Dispatch<SetStateAction<string | null>>,
): Middleware => {
  return {
    async onRequest({ request }) {
      if (token.current) {
        request.headers.set("Authorization", token.current);
      }
      return request;
    },

    async onResponse({ response, schemaPath }) {
      if (schemaPath === "/api/auth/login") {
        const token = response.headers.get("Authorization");
        if (token) {
          localStorage.setItem("jwt", token);
          setToken(token);
        }
      }
      return response;
    },
  };
};

const ClientContext = createContext<AppClient | undefined>(undefined);

export default ClientContext;

interface Props {
  client: AppClient;
  children: JSX.Element;
}

export const ClientProvider = ({ client, children }: Props) => {
  const context = useContext(AppContext);
  const { token, tokenRef } = context!;

  // biome-ignore lint/correctness/useExhaustiveDependencies: Probably linter is right but I dunno what i am doing
  useEffect(() => {
    const middleware = createMiddleware(tokenRef, token[1]);
    client.fetchClient.use(middleware);
    return () => client.fetchClient.eject(middleware);
  }, []);

  return (
    <ClientContext.Provider value={client}>
      <QueryClientProvider client={client.queryClient}>
        {children}
      </QueryClientProvider>
    </ClientContext.Provider>
  );
};
