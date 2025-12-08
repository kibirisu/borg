import { useRef } from "react";
import type { AppClient } from "../../lib/client";

const RegisterButton = (props: Props) => {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const { mutate: register } = props.client.$api.useMutation(
    "post",
    "/api/auth/register",
  );

  const openDialog = () => {
    if (dialogRef.current) {
      dialogRef.current.showModal();
    }
  };

  const registerAction = async (data: FormData) => {
    if (dialogRef.current) {
      const username = data.get("username")?.toString();
      const password = data.get("password")?.toString();
      if (username && password) {
        register({ body: { username: username, password: password } });
      }
    }
  };

  return (
    <>
      <button className="btn" onClick={openDialog}>
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
      </dialog>
    </>
  );
};

interface Props {
  client: AppClient;
}

export default RegisterButton;
