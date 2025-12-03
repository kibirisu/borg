import { Outlet } from "react-router";
import "./App.css";
import TopAppBar from "./components/bar/TopAppBar";

const App = () => {
  return (
    <>
      <TopAppBar />
      <Outlet />
    </>
  );
};

export default App;
