import { useContext, useRef } from "react";
import { ClientContext } from "../../lib/client";

const AuthButtons = () => {
  const client = useContext(ClientContext);
  if (!client) {
    throw Error();
  }

  const dialogRef = useRef<HTMLDialogElement>(null);
  const { mutate } = client.$api.useMutation("post", "/api/auth/login");

  const openDialog = () => {
    if (dialogRef.current) {
      dialogRef.current.showModal();
    }
  };

  const loginAction = (data: FormData) => {
    if (dialogRef.current) {
      const username = data.get("username")?.toString();
      const password = data.get("password")?.toString();
      if (username && password) {
        mutate({ body: { username: username, password: password } });
      }
      dialogRef.current.close();
    }
  };

  return (
    <>
      <button className="btn">Sign Up</button>
      <button className="btn" onClick={openDialog}>
        Sign In
      </button>
      <dialog ref={dialogRef} className="modal">
        <div className="modal-box max-w-min">
          <form action={loginAction}>
            <fieldset className="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
              <legend className="fieldset-legend">Login</legend>

              <label className="label">Email</label>
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
                Login
              </button>
            </fieldset>
          </form>
        </div>

        <form method="dialog" className="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>
    </>
  );
};

export default AuthButtons;
