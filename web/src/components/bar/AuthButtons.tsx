import { useContext, useRef } from "react";
import { AppContext } from "../../lib/state";

const AuthButtons = () => {
  const context = useContext(AppContext);
  if (!context) {
    return <></>;
  }

  const dialogRef = useRef<HTMLDialogElement>(null);
  const usernameInputRef = useRef<HTMLInputElement>(null);
  const passwordInputRef = useRef<HTMLInputElement>(null);

  const { mutate } = context.$api.useMutation("post", "/api/auth/login");

  const openDialog = () => {
    if (dialogRef.current) {
      dialogRef.current.showModal();
    }
  };

  const submitForm = () => {
    if (context && usernameInputRef.current && passwordInputRef.current) {
      const username = usernameInputRef.current.value;
      const password = passwordInputRef.current.value;
      mutate({ body: { username: username, password: password } });
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
          <fieldset className="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
            <legend className="fieldset-legend">Login</legend>

            <label className="label">Email</label>
            <input
              ref={usernameInputRef}
              type="text"
              className="input"
              placeholder="Username"
            />

            <label className="label">Password</label>
            <input
              ref={passwordInputRef}
              type="password"
              className="input"
              placeholder="Password"
            />

            <button className="btn btn-neutral mt-4" onClick={submitForm}>
              Login
            </button>
          </fieldset>
        </div>

        <form method="dialog" className="modal-backdrop">
          <button>close</button>
        </form>
      </dialog>
    </>
  );
};

export default AuthButtons;
