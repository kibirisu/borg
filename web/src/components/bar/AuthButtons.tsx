import { useContext } from "react";
import ClientContext from "../../lib/client";
import LoginButton from "./Login";
import RegisterButton from "./Register";

const AuthButtons = () => {
  const client = useContext(ClientContext);

  return (
    <>
      <LoginButton client={client!} />
      <RegisterButton client={client!} />
    </>
  );
};

export default AuthButtons;
