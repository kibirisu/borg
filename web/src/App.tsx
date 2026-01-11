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
import "bootstrap-icons/font/bootstrap-icons.css";

const fetchClient = createFetchClient<paths>();
const $api = createClient(fetchClient);
const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 1000 * 10 } },
});
const client: AppClient = { $api, fetchClient, queryClient };

const normalizeToken = (value: string | null) => {
  if (!value) {
    return null;
  }
  const cleaned = value.replace(/^Bearer:\s*/i, "");
  if (cleaned !== value) {
    localStorage.setItem("jwt", cleaned);
  }
  return cleaned;
};

const App = () => {
  const [token, setToken] = useState(() =>
    normalizeToken(localStorage.getItem("jwt")),
  );
  const tokenRef = useRef(token);
  const decoded = useMemo(() => {
    return decodeToken(token);
  }, [token]);
  const username = decoded?.username ?? null;
  const userId = decoded?.userId ?? null;

  useEffect(() => {
    tokenRef.current = token;
  }, [token]);

  return (
    <AppStateProvider
      token={[token, setToken]}
      tokenRef={tokenRef}
      username={username}
      userId={userId}
    >
      <ClientProvider client={client}>
        {/* biome-ignore lint/complexity/noUselessFragments: ClientProvider takes a single JSX child element */}
        <>
          <RouterProvider />
          <ReactQueryDevtools />
        </>
      </ClientProvider>
    </AppStateProvider>
  );
};

export default App;
