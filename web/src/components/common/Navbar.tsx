import { Link } from "react-router";

const Navbar = () => {
  return (
    <header className="border-b border-white/10 bg-gray-900">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
        {/* Logo */}
        <Link to="/" className="text-xl font-bold text-white">
          Borg
        </Link>

        {/* Actions */}
        <div className="flex items-center gap-4">
          <Link
            to="/signin"
            className="text-sm font-medium text-gray-300 hover:text-white"
          >
            Sign in
          </Link>

          <Link
            to="/signup"
            className="rounded-md bg-indigo-500 px-3 py-1.5 text-sm font-semibold text-white hover:bg-indigo-400"
          >
            Sign up
          </Link>
        </div>
      </div>
    </header>
  );
};

export default Navbar;