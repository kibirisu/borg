import { Outlet } from "react-router";
import "./App.css";
import TopAppBar from "./components/TopAppBar";

const App = () => {
  return (
    <>
      <TopAppBar />
      <Outlet />
    </>
  );
};

export default App;
