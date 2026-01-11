import { Navigate, Outlet, useLocation } from "react-router";
import Navbar from "./common/Navbar";

const Root = () => {
  const location = useLocation();
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
