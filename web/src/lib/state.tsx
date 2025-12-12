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

export const AppStateProvider = ({
  children,
  token,
  tokenRef,
  username,
}: Props) => {
  return (
    <AppContext.Provider value={{ token, tokenRef, username }}>
      {children}
    </AppContext.Provider>
  );
};
