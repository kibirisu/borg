import { useContext } from "react";
import { Navigate, Outlet, useLocation } from "react-router";
import Navbar from "./common/Navbar";
import AppContext from "../lib/state";

const Root = () => {
  const location = useLocation();
  const appState = useContext(AppContext);
  const hideNavbar = ["/signin", "/signup"].includes(location.pathname);

  if (location.pathname === "/") {
    return <Navigate to="/explore" replace />;
  }

  return (
    <>
      {!hideNavbar && <Navbar />}
      <main>
        <Outlet />
      </main>
    </>
  );
};

export default Root;
