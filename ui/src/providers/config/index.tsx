import { createContext, useCallback, useContext, useEffect, useReducer, type ReactNode } from "react";

import { type Access as AccessType } from "@/domain/access";
import { list as listAccess } from "@/repository/access";
import { configReducer } from "./reducer";

export type ConfigData = {
  accesses: AccessType[];
};

export type ConfigContext = {
  config: ConfigData;
  addAccess: (access: AccessType) => void;
  updateAccess: (access: AccessType) => void;
  deleteAccess: (id: string) => void;
};

const Context = createContext({} as ConfigContext);

export const useConfigContext = () => useContext(Context);

export const ConfigProvider = ({ children }: { children: ReactNode }) => {
  const [config, dispatchConfig] = useReducer(configReducer, {
    accesses: [],
  });

  useEffect(() => {
    const featchData = async () => {
      const data = await listAccess();
      dispatchConfig({ type: "SET_ACCESSES", payload: data });
    };
    featchData();
  }, []);

  const deleteAccess = useCallback((id: string) => {
    dispatchConfig({ type: "DELETE_ACCESS", payload: id });
  }, []);

  const addAccess = useCallback((access: AccessType) => {
    dispatchConfig({ type: "ADD_ACCESS", payload: access });
  }, []);

  const updateAccess = useCallback((access: AccessType) => {
    dispatchConfig({ type: "UPDATE_ACCESS", payload: access });
  }, []);

  return (
    <Context.Provider
      value={{
        config: {
          accesses: config.accesses,
        },
        addAccess,
        updateAccess,
        deleteAccess,
      }}
    >
      {children}
    </Context.Provider>
  );
};
