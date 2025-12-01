import { useState } from "react";
import ProfileButton from "./ProfileButton";
import AuthButtons from "./AuthButtons";

export default function TopAppBar() {
  const [signedIn] = useState(false);

  return (
    <div className="navbar bg-base-100 shadow-sm">
      <div className="navbar-start">
        <a className="btn btn-ghost text-xl">borg</a>
      </div>
      <div className="navbar-center">
        <input
          type="text"
          placeholder="Search"
          className="input input-bordered w-24 md:w-auto"
        />
      </div>
      <div className="navbar-end gap-2">
        {signedIn ? <ProfileButton /> : <AuthButtons />}
      </div>
    </div>
  );
}
