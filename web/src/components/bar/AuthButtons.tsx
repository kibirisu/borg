import { useRef } from "react";

const AuthButtons = () => {
  const dialogRef = useRef<HTMLDialogElement>(null);

  return (
    <>
      <a className="btn">Sign Up</a>
      <button
        className="btn"
        onClick={() => {
          if (dialogRef.current) {
            dialogRef.current.showModal();
          }
        }}
      >
        Sign In
      </button>
      <dialog ref={dialogRef} className="modal">
        <div className="modal-box max-w-min">
          <fieldset className="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
            <legend className="fieldset-legend">Login</legend>

            <label className="label">Email</label>
            <input type="email" className="input" placeholder="Email" />

            <label className="label">Password</label>
            <input type="password" className="input" placeholder="Password" />

            <button className="btn btn-neutral mt-4">Login</button>
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
