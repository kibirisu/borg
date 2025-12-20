import { Form, useActionData, useNavigate } from "react-router";

type SignInErrors = {
  username?: string;
  password?: string;
  form?: string;
};

export const SignIn = () => {
  const navigate = useNavigate();
  const errors = useActionData() as SignInErrors | undefined;

  return (
    <div className="min-h-full bg-white">
      <div className="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
        <div>
          <button
            type="button"
            onClick={() => navigate("/")}
            aria-label="Go back to feed"
            className="inline-flex items-center justify-center border border-black text-black rounded-[7px] text-sm p-2.5"
          >
            <i className="bi bi-arrow-left" />
          </button>
        </div>
        <div className="sm:mx-auto sm:w-full sm:max-w-sm">
          <img
            src="https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=600"
            alt="Your Company"
            className="mx-auto h-10 w-auto"
          />

          <h2 className="mt-10 text-center text-2xl/9 font-bold tracking-tight text-gray-900">
            Sign in to your account
          </h2>
        </div>

        <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
          {/* Form з react-router */}
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
                  aria-invalid={errors?.username ? true : undefined}
                  aria-describedby={errors?.username ? "username-error" : undefined}
                  className="block w-full rounded-md bg-white px-3 py-1.5 text-base
                    text-gray-900 outline-1 -outline-offset-1 outline-gray-300
                    placeholder:text-gray-400
                    focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600
                    sm:text-sm/6"
                />
              </div>
              {errors?.username ? (
                <p id="username-error" className="mt-2 text-sm text-red-600">
                  {errors.username}
                </p>
              ) : null}
            </div>

            {/* PASSWORD */}
            <div>
              <div className="flex items-center justify-between">
                <label
                  htmlFor="password"
                  className="block text-sm/6 font-medium text-gray-900"
                >
                  Password
                </label>
              </div>

              <div className="mt-2">
                <input
                  id="password"
                  name="password"
                  type="password"
                  required
                  autoComplete="current-password"
                  aria-invalid={errors?.password ? true : undefined}
                  aria-describedby={errors?.password ? "password-error" : undefined}
                  className="block w-full rounded-md bg-white px-3 py-1.5 text-base
                    text-gray-900 outline-1 -outline-offset-1 outline-gray-300
                    placeholder:text-gray-400
                    focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-600
                    sm:text-sm/6"
                />
              </div>
              {errors?.password ? (
                <p id="password-error" className="mt-2 text-sm text-red-600">
                  {errors.password}
                </p>
              ) : null}
            </div>

            {/* SUBMIT */}
            <div>
              <button
                type="submit"
                className="flex w-full justify-center rounded-md bg-indigo-600
                  px-3 py-1.5 text-sm/6 font-semibold text-white shadow-xs
                  hover:bg-indigo-500
                  focus-visible:outline-2 focus-visible:outline-offset-2
                  focus-visible:outline-indigo-600"
              >
                Sign in
              </button>
            </div>
          </Form>
        </div>
      </div>
    </div>
  );
};
