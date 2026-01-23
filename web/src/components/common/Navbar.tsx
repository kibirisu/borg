import { useContext, useEffect } from "react";
import { Link } from "react-router";
import ClientContext from "../../lib/client";
import AppContext from "../../lib/state";

const Navbar = () => {
  const context = useContext(AppContext);
  const client = useContext(ClientContext);
  const username = context?.username;
  const setUsername = context?.setUsername;
  const userId = context?.userId ?? null;
  const [, setToken] = context?.token ?? [];

  const handleLogout = () => {
    localStorage.removeItem("jwt");
    setToken?.(null);
  };

  useEffect(() => {
    if (!client || !setUsername || username || !userId) {
      return;
    }
    const fetchAccount = async () => {
      const res = await client.fetchClient.GET("/api/accounts/{id}", {
        params: { path: { id: userId } },
      });
      if (res.data) {
        setUsername(res.data.display_name || res.data.username || null);
      }
    };
    fetchAccount().catch((err) => {
      console.warn("Failed to fetch account profile", err);
    });
  }, [client, setUsername, username, userId]);

  return (
    <header className="border-b border-indigo-700 bg-indigo-600">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
        {/* Logo */}
        <Link to="/" className="text-xl font-bold text-white">
          Borg
        </Link>

        {/* Actions */}
        <div className="flex items-center gap-4">
          {username ? (
            <>
              <span className="text-sm font-medium text-white">{username}</span>
              <button
                type="button"
                onClick={handleLogout}
                className="rounded-md border border-white/30 px-3 py-1.5 text-sm font-semibold text-white hover:bg-white/10"
              >
                Log out
              </button>
            </>
          ) : (
            <>
              <Link
                to="/signin"
                className="text-sm font-medium text-gray-300 hover:text-white"
              >
                Sign in
              </Link>

              <Link
                to="/signup"
                className="rounded-md bg-indigo-500 px-3 py-1.5 text-sm font-semibold text-white hover:bg-indigo-400"
              >
                Sign up
              </Link>
            </>
          )}
        </div>
      </div>
    </header>
  );
};

export default Navbar;
