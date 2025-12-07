import { Outlet } from "react-router";
import TopAppBar from "./bar/TopAppBar";

const Root = () => {
  return (
    <>
      <TopAppBar />
      <Outlet />
    </>
  );
};

export default Root;
