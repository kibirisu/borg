import {
  createContext,
  type Dispatch,
  type JSX,
  type RefObject,
  type SetStateAction,
} from "react";

export interface AppState {
  token: [string | null, Dispatch<SetStateAction<string | null>>];
  tokenRef: RefObject<string | null>;
  username: string | null;
}

const AppContext = createContext<AppState | undefined>(undefined);

export default AppContext;

interface Props {
  children: JSX.Element;
  token: [string | null, Dispatch<SetStateAction<string | null>>];
  tokenRef: RefObject<string | null>;
  username: string | null;
}

export const AppStateProvider = (props: Props) => {
  const { children, token, tokenRef, username } = props;
  return (
    <AppContext.Provider
      value={{ token: token, tokenRef: tokenRef, username: username }}
    >
      {children}
    </AppContext.Provider>
  );
};
