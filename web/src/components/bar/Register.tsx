import { useRef, useState } from "react";
import type { AppClient } from "../../lib/client";

const RegisterButton = ({ client }: Props) => {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const [usernameError, setUsernameError] = useState<string | null>(null);
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const { mutateAsync: register } = client.$api.useMutation(
    "post",
    "/api/auth/register",
  );

  const openDialog = () => {
    setUsernameError(null); // Wyczyść błędy przy otwieraniu
    dialogRef?.current?.showModal();
  };

  const registerAction = async (data: FormData) => {
    if (dialogRef.current) {
      const username = data.get("username")?.toString();
      const password = data.get("password")?.toString();

      // Wyczyść poprzednie błędy
      setUsernameError(null);
      setPasswordError(null);

      // Walidacja pól
      let hasErrors = false;
      if (!username || username.trim() === "") {
        setUsernameError("Username is required");
        hasErrors = true;
      }
      if (!password || password.trim() === "") {
        setPasswordError("Password is required");
        hasErrors = true;
      }

      if (hasErrors) {
        return; // Nie kontynuuj jeśli są błędy walidacji
      }

      try {
        await register({ body: { username: username, password: password } });
        // Sukces - zamknij dialog
        if (dialogRef.current) {
          dialogRef.current.close();
        }
      } catch (error: any) {
        // Obsłuż błąd - openapi-fetch zwraca błędy w strukturze { status, data }
        const status = error?.status || error?.response?.status;
        if (status === 409) {
          // Spróbuj wyciągnąć komunikat z odpowiedzi
          const errorData = error?.data || error?.response?.data;
          if (
            errorData &&
            typeof errorData === "object" &&
            "error" in errorData
          ) {
            setUsernameError(errorData.error as string);
          } else {
            setUsernameError("Username already taken");
          }
        } else {
          setUsernameError("Registration failed. Please try again.");
        }
      }
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
                className={`input ${usernameError ? "input-error" : ""}`}
                placeholder="Username"
                name="username"
              />
              {usernameError && (
                <div className="label">
                  <span className="label-text-alt text-error">
                    {usernameError}
                  </span>
                </div>
              )}

              <label className="label">Password</label>
              <input
                type="password"
                className={`input ${passwordError ? "input-error" : ""}`}
                placeholder="Password"
                name="password"
              />
              {passwordError && (
                <div className="label">
                  <span className="label-text-alt text-error">
                    {passwordError}
                  </span>
                </div>
              )}

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
