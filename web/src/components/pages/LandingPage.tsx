import { Link } from "react-router";

const LandingPage = () => {
  return (
    <div className="min-h-screen bg-white">
      <div className="mx-auto flex min-h-screen max-w-3xl flex-col items-center justify-center px-6 text-center">
        <h1 className="text-4xl font-bold text-indigo-600 sm:text-6xl">
          Welcome to Borg!
        </h1>
        <p className="mt-4 text-base text-gray-600 sm:text-lg">
          If you want to use our service please sign up, or if you have a page,
          sign in to continue.
        </p>
        <div className="mt-8 flex flex-wrap items-center justify-center gap-4">
          <Link
            to="/signup"
            className="rounded-full bg-indigo-600 px-6 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-700"
          >
            Sign up
          </Link>
          <Link
            to="/signin"
            className="rounded-full border border-indigo-600 px-6 py-2 text-sm font-semibold text-indigo-600 hover:bg-indigo-50"
          >
            Sign in
          </Link>
        </div>
      </div>
    </div>
  );
};

export default LandingPage;
