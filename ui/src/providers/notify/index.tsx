import { ReactNode, useContext, createContext, useEffect, useReducer, useCallback } from "react";

import { NotifyChannel, Setting } from "@/domain/settings";
import { getSetting } from "@/repository/settings";
import { notifyReducer } from "./reducer";

export type NotifyContext = {
  config: Setting;
  setChannel: (data: { channel: string; data: NotifyChannel }) => void;
  setChannels: (data: Setting) => void;
};

const Context = createContext({} as NotifyContext);

export const useNotifyContext = () => useContext(Context);

interface NotifyProviderProps {
  children: ReactNode;
}

export const NotifyProvider = ({ children }: NotifyProviderProps) => {
  const [notify, dispatchNotify] = useReducer(notifyReducer, {});

  useEffect(() => {
    const featchData = async () => {
      const chanels = await getSetting("notifyChannels");
      dispatchNotify({
        type: "SET_CHANNELS",
        payload: chanels,
      });
    };
    featchData();
  }, []);

  const setChannel = useCallback((data: { channel: string; data: NotifyChannel }) => {
    dispatchNotify({
      type: "SET_CHANNEL",
      payload: data,
    });
  }, []);

  const setChannels = useCallback((setting: Setting) => {
    dispatchNotify({
      type: "SET_CHANNELS",
      payload: setting,
    });
  }, []);

  return (
    <Context.Provider
      value={{
        config: notify,
        setChannel,
        setChannels,
      }}
    >
      {children}
    </Context.Provider>
  );
};
