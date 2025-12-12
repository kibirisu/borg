import { useContext } from "react";
import AppContext from "../../lib/state";

const ProfileButton = () => {
  const context = useContext(AppContext);

  // biome-ignore lint/style/noNonNullAssertion: AppContext will never be `underfined`
  const [_, setToken] = context!.token;

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
          <div className="justify-between">
            Profile
            <span className="badge">New</span>
          </div>
        </li>
        <li>
          <div>Settings</div>
        </li>
        <li>
          <button onClick={logoutAction} type="button">
            Logout
          </button>
        </li>
      </ul>
    </div>
  );
};

export default ProfileButton;
