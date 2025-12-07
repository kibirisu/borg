import {
  createContext,
  type Dispatch,
  type JSX,
  type SetStateAction,
} from "react";

export interface AppState {
  token: [string | null, Dispatch<SetStateAction<string | null>>];
  username: string | null;
}

export interface UserState {
  username: string | null;
  token: string | null;
}

const AppContext = createContext<AppState | undefined>(undefined);

export default AppContext;

interface Props {
  children: JSX.Element;
  token: [string | null, Dispatch<SetStateAction<string | null>>];
  username: string | null;
}

export const AppStateProvider = (props: Props) => {
  const { children, token, username } = props;
  return (
    <AppContext.Provider value={{ token: token, username: username }}>
      {children}
    </AppContext.Provider>
  );
};
