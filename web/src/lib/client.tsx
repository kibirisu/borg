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
  const stripPrefix = (value: string) => value.replace(/^Bearer:\s*/i, "");
  return {
    async onRequest({ request }) {
      if (token.current) {
        request.headers.set("Authorization", `Bearer: ${token.current}`);
      }
      return request;
    },

    async onResponse({ response, schemaPath }) {
      if (schemaPath === "/auth/login") {
        const bearer = response.headers.get("Authorization");
        if (bearer) {
          const raw = stripPrefix(bearer);
          localStorage.setItem("jwt", raw);
          setToken(raw);
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
