import { QueryClient } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import createFetchClient from "openapi-fetch";
import createClient from "openapi-react-query";
import { useEffect, useMemo, useRef, useState } from "react";
import type { paths } from "./lib/api/v1";
import { type AppClient, ClientProvider } from "./lib/client.tsx";
import decodeToken from "./lib/decode.ts";
import { AppStateProvider } from "./lib/state.tsx";
import RouterProvider from "./routes/router.tsx";

const fetchClient = createFetchClient<paths>();
const $api = createClient(fetchClient);
const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 1000 * 10 } },
});
const client: AppClient = { $api, fetchClient, queryClient };

const App = () => {
  const [token, setToken] = useState(localStorage.getItem("jwt"));
  const tokenRef = useRef(token);
  const username = useMemo(() => {
    return decodeToken(token);
  }, [token]);

  useEffect(() => {
    tokenRef.current = token;
  }, [token]);

  return (
    <AppStateProvider
      token={[token, setToken]}
      tokenRef={tokenRef}
      username={username}
    >
      <ClientProvider client={client}>
        <>
          <RouterProvider />
          <ReactQueryDevtools />
        </>
      </ClientProvider>
    </AppStateProvider>
  );
};

export default App;
