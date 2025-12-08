import { useContext } from "react";
import AppContext from "../../lib/state";
import AuthButtons from "./AuthButtons";
import ProfileButton from "./ProfileButton";

export default function TopAppBar() {
  const state = useContext(AppContext);
  if (!state) {
    throw Error();
  }

  return (
    <div className="navbar bg-base-100 shadow-sm sticky top-0">
      <div className="navbar-start">
        <div className="btn btn-ghost text-xl">borg</div>
      </div>
      <div className="navbar-center">
        <input
          type="text"
          placeholder="Search"
          className="input input-bordered w-24 md:w-auto"
        />
      </div>
      <div className="navbar-end gap-2">
        {state.username ? <ProfileButton /> : <AuthButtons />}
      </div>
    </div>
  );
}
