import { useContext } from "react";
import { ClientContext } from "../../lib/client";
import AppContext from "../../lib/state";

const ProfileButton = () => {
  const client = useContext(ClientContext);
  const context = useContext(AppContext);
  if (!client || !context) {
    throw Error();
  }

  const [_, setToken] = context.token;

  const logoutAction = () => {
    localStorage.removeItem("jwt");
    setToken(null);
  };

  return (
    <div className="dropdown dropdown-end">
      <div
        tabIndex={0}
        role="button"
        className="btn btn-ghost btn-circle avatar"
      >
        <div className="w-10 rounded-full">
          <img
            alt="Tailwind CSS Navbar component"
            src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp"
          />
        </div>
      </div>
      <ul
        tabIndex={-1}
        className="menu menu-sm dropdown-content bg-base-100 rounded-box z-1 mt-3 w-52 p-2 shadow"
      >
        <li>
          <a className="justify-between">
            Profile
            <span className="badge">New</span>
          </a>
        </li>
        <li>
          <a>Settings</a>
        </li>
        <li>
          <button onClick={logoutAction}>Logout</button>
        </li>
      </ul>
    </div>
  );
};

export default ProfileButton;
