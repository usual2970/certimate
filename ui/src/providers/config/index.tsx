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

export type ConfigData = {
  accesses: Access[];
};

export type ConfigContext = {
  config: ConfigData;
  deleteAccess: (id: string) => void;
  addAccess: (access: Access) => void;
  updateAccess: (access: Access) => void;
};

const Context = createContext({} as ConfigContext);

export const useConfig = () => useContext(Context);

interface ContainerProps {
  children: ReactNode;
}

export const ConfigProvider = ({ children }: ContainerProps) => {
  const [config, dispatchConfig] = useReducer(configReducer, { accesses: [] });

  useEffect(() => {
    const featchData = async () => {
      const data = await list();
      dispatchConfig({ type: "SET_ACCESSES", payload: data });
    };
    featchData();
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
        },
        deleteAccess,
        addAccess,
        updateAccess,
      }}
    >
      {children && children}
    </Context.Provider>
  );
};
