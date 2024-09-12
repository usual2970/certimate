import { Access } from "@/domain/access";
import { list } from "@/repository/access";

import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useReducer,
} from "react";
import { configReducer } from "./reducer";
import { getEmails } from "@/repository/settings";
import { Setting } from "@/domain/settings";

export type ConfigData = {
  accesses: Access[];
  emails: Setting;
};

export type ConfigContext = {
  config: ConfigData;
  deleteAccess: (id: string) => void;
  addAccess: (access: Access) => void;
  updateAccess: (access: Access) => void;
  setEmails: (email: Setting) => void;
};

const Context = createContext({} as ConfigContext);

export const useConfig = () => useContext(Context);

interface ContainerProps {
  children: ReactNode;
}

export const ConfigProvider = ({ children }: ContainerProps) => {
  const [config, dispatchConfig] = useReducer(configReducer, {
    accesses: [],
    emails: { content: { emails: [] } },
  });

  useEffect(() => {
    const featchData = async () => {
      const data = await list();
      dispatchConfig({ type: "SET_ACCESSES", payload: data });
    };
    featchData();
  }, []);

  useEffect(() => {
    const featchEmails = async () => {
      const emails = await getEmails();
      dispatchConfig({ type: "SET_EMAILS", payload: emails });
    };
    featchEmails();
  }, []);

  const setEmails = useCallback((emails: Setting) => {
    dispatchConfig({ type: "SET_EMAILS", payload: emails });
  }, []);

  const deleteAccess = useCallback((id: string) => {
    dispatchConfig({ type: "DELETE_ACCESS", payload: id });
  }, []);

  const addAccess = useCallback((access: Access) => {
    dispatchConfig({ type: "ADD_ACCESS", payload: access });
  }, []);

  const updateAccess = useCallback((access: Access) => {
    dispatchConfig({ type: "UPDATE_ACCESS", payload: access });
  }, []);

  return (
    <Context.Provider
      value={{
        config: {
          accesses: config.accesses,
          emails: config.emails,
        },
        deleteAccess,
        addAccess,
        setEmails,
        updateAccess,
      }}
    >
      {children && children}
    </Context.Provider>
  );
};
