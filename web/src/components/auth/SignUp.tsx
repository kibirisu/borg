import { Form, Link } from "react-router";

export const SignUp = () => {
  return (
    <div className="min-h-full bg-white">
      <div className="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
        <div className="sm:mx-auto sm:w-full sm:max-w-sm">
          <img
            src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
            alt="Your Company"
            className="mx-auto h-10 w-auto"
          />
          <h2 className="mt-10 text-center text-2xl/9 font-bold tracking-tight text-gray-900">
            Create a new account
          </h2>
        </div>

        <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
          <Form method="post" className="space-y-6">
            {/* USERNAME */}
            <div>
              <label
                htmlFor="username"
                className="block text-sm/6 font-medium text-gray-900"
              >
                Username
              </label>
              <div className="mt-2">
                <input
                  id="username"
                  name="username"
                  type="text"
                  required
                  autoComplete="username"
                  className="block w-full rounded-md bg-white px-3 py-1.5 text-base
                    text-gray-900 outline-1 -outline-offset-1 outline-gray-300
                    placeholder:text-gray-400
                    focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600
                    sm:text-sm/6"
                />
              </div>
            </div>

            {/* PASSWORD */}
            <div>
              <label
                htmlFor="password"
                className="block text-sm/6 font-medium text-gray-900"
              >
                Password
              </label>
              <div className="mt-2">
                <input
                  id="password"
                  name="password"
                  type="password"
                  required
                  autoComplete="new-password"
                  className="block w-full rounded-md bg-white px-3 py-1.5 text-base
                    text-gray-900 outline-1 -outline-offset-1 outline-gray-300
                    placeholder:text-gray-400
                    focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600
                    sm:text-sm/6"
                />
              </div>
            </div>

            {/* CONFIRM PASSWORD */}
            <div>
              <label
                htmlFor="confirmPassword"
                className="block text-sm/6 font-medium text-gray-900"
              >
                Confirm password
              </label>
              <div className="mt-2">
                <input
                  id="confirmPassword"
                  name="confirmPassword"
                  type="password"
                  required
                  autoComplete="new-password"
                  className="block w-full rounded-md bg-white px-3 py-1.5 text-base
                    text-gray-900 outline-1 -outline-offset-1 outline-gray-300
                    placeholder:text-gray-400
                    focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600
                    sm:text-sm/6"
                />
              </div>
            </div>

            {/* SUBMIT */}
            <div>
              <button
                type="submit"
                className="flex w-full justify-center rounded-md bg-indigo-600
                  px-3 py-1.5 text-sm/6 font-semibold text-white
                  hover:bg-indigo-500
                  focus-visible:outline-2 focus-visible:outline-offset-2
                  focus-visible:outline-indigo-600"
              >
                Create a new account
              </button>
            </div>
          </Form>

          {/* FOOTER */}
          <p className="mt-10 text-center text-sm/6 text-gray-500">
            Already have an account?{" "}
            <Link
              to="/signin"
              className="font-semibold text-indigo-600 hover:text-indigo-500"
            >
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};
