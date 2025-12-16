import { Outlet } from "react-router";
import Navbar from "./common/Navbar";

const Root = () => {
  return (
    <>
      <Navbar />
      <main>
        <Outlet />
      </main>
    </>
  );
};

export default Root;
