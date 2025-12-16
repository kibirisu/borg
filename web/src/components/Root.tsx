import { Outlet, useLocation } from "react-router";
import Navbar from "./common/Navbar";

const Root = () => {
  const location = useLocation();
  const hideNavbar = ["/signin", "/signup"].includes(location.pathname);

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
