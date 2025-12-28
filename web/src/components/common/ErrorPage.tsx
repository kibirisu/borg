import { Link, isRouteErrorResponse, useRouteError } from "react-router";
import { useState } from "react";

const UnauthorizedBanner = () => {
  const [open, setOpen] = useState(true);
  if (!open) return null;

  return (
    <div
      id="marketing-banner"
      tabIndex={-1}
      className="fixed z-50 flex flex-col md:flex-row justify-between w-[calc(100%-2rem)] p-4 -translate-x-1/2 bg-neutral-100 border border-gray-200 rounded-lg shadow-sm lg:max-w-7xl left-1/2 top-6 text-gray-800"
    >
      <div className="flex flex-col items-start mb-3 me-4 md:items-center md:flex-row md:mb-0">
        <div className="flex items-center mb-2 md:pe-4 md:me-4 md:border-e md:border-gray-200 md:mb-0">
          <span className="text-lg font-semibold whitespace-nowrap">Borg</span>
        </div>
        <p className="flex items-center text-sm font-normal text-gray-700">
          You need to sign in to access this page.
        </p>
      </div>
      <div className="flex items-center shrink-0">
        <Link
          to="/signin"
          className="text-white bg-indigo-600 hover:bg-indigo-700 border border-transparent focus:ring-4 focus:ring-indigo-300 shadow-xs font-medium leading-5 rounded-lg text-xs px-3 py-1.5 focus:outline-none me-2"
        >
          Sign in
        </Link>
        <button
          type="button"
          onClick={() => setOpen(false)}
          className="hidden shrink-0 md:inline-flex justify-center text-sm w-7 h-7 items-center text-gray-500 hover:bg-gray-200 hover:text-gray-800 rounded-sm"
        >
          <svg
            className="w-4 h-4"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke="currentColor"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M6 18 17.94 6M18 18 6.06 6"
            />
          </svg>
          <span className="sr-only">Close banner</span>
        </button>
        <button
          type="button"
          onClick={() => setOpen(false)}
          className="md:hidden text-gray-600 bg-neutral-200 border border-gray-300 hover:bg-gray-300 hover:text-gray-800 focus:ring-4 focus:ring-gray-200 shadow-xs font-medium leading-5 rounded-lg text-xs px-3 py-1.5 focus:outline-none"
        >
          Close
        </button>
      </div>
    </div>
  );
};

const ErrorPage = () => {
  const error = useRouteError();

  const status =
    (isRouteErrorResponse(error) && error.status) ||
    (error instanceof Response && error.status) ||
    (typeof error === "object" && error !== null && "status" in error
      ? Number((error as any).status)
      : null) ||
    (typeof error === "object" &&
      error !== null &&
      "response" in error &&
      (error as any).response?.status);

  if (status === 401) {
    return <UnauthorizedBanner />;
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <div className="max-w-md text-center space-y-3">
        <p className="text-sm font-semibold text-gray-500">Something went wrong</p>
        <p className="text-gray-700">
          We couldn&apos;t load this page. Please try again or return to the home page.
        </p>
        <div className="flex items-center justify-center gap-3">
          <Link
            to="/signin"
            className="inline-flex items-center px-3 py-1.5 rounded-lg text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
          >
            Sign In
          </Link>
          <button
            type="button"
            onClick={() => window.location.reload()}
            className="inline-flex items-center px-3 py-1.5 rounded-lg text-sm font-medium text-gray-700 bg-white border border-gray-300 hover:bg-gray-100"
          >
            Reload
          </button>
        </div>
      </div>
    </div>
  );
};

export default ErrorPage;
