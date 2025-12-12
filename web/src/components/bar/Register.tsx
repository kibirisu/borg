import { useRef } from "react";
import type { AppClient } from "../../lib/client";

const RegisterButton = ({ client }: Props) => {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const { mutateAsync: register } = client.$api.useMutation(
    "post",
    "/api/auth/register",
  );

  const openDialog = () => dialogRef?.current?.showModal();

  const registerAction = async (data: FormData) => {
    const username = data.get("username")?.toString();
    const password = data.get("password")?.toString();
    if (username && password) {
      await register({ body: { username: username, password: password } });
    }
  };

  return (
    <>
      <button className="btn" onClick={openDialog} type="button">
        Sign Up
      </button>
      <dialog ref={dialogRef} className="modal">
        <div className="modal-box max-w-min">
          <form action={registerAction}>
            <fieldset className="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
              <legend className="fieldset-legend">Register</legend>

              <label className="label">Username</label>
              <input
                type="text"
                className="input"
                placeholder="Username"
                name="username"
              />

              <label className="label">Password</label>
              <input
                type="password"
                className="input"
                placeholder="Password"
                name="password"
              />

              <button type="submit" className="btn btn-neutral mt-4">
                Register
              </button>
            </fieldset>
          </form>
        </div>

        <form method="dialog" className="modal-backdrop">
          <button type="submit">close</button>
        </form>
      </dialog>
    </>
  );
};

interface Props {
  client: AppClient;
}

export default RegisterButton;
