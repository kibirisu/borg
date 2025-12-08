import { type QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type createFetchClient from "openapi-fetch";
import type { Middleware } from "openapi-fetch";
import type { OpenapiQueryClient } from "openapi-react-query";
import {
  createContext,
  type Dispatch,
  type JSX,
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
  tokenState: [string | null, Dispatch<SetStateAction<string | null>>],
): Middleware => {
  const [token, setToken] = tokenState;
  return {
    async onRequest({ request }) {
      if (token) {
        request.headers.set("Authorization", token);
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

export const ClientContext = createContext<AppClient | undefined>(undefined);

interface Props {
  client: AppClient;
  children: JSX.Element;
}

export const ClientProvider = (props: Props) => {
  const client = props.client;
  const context = useContext(AppContext);
  if (!context) {
    throw Error();
  }

  useEffect(() => {
    const middleware = createMiddleware(context.token);
    client.fetchClient.use(middleware);
    // FIXME: bellow leads to bug
    // return () => {
    //   client.fetchClient.eject(middleware);
    // };
  }, []);

  return (
    <ClientContext.Provider value={client}>
      <QueryClientProvider client={client.queryClient}>
        {props.children}
      </QueryClientProvider>
    </ClientContext.Provider>
  );
};
