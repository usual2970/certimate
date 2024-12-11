import { ReactNode, useContext, createContext, useEffect, useReducer, useCallback } from "react";

import { NotifyChannel, NotifyChannels, Settings } from "@/domain/settings";
import { get } from "@/repository/settings";
import { notifyReducer } from "./reducer";

export type NotifyContext = {
  config: Settings<NotifyChannels>;
  setChannel: (data: { channel: string; data: NotifyChannel }) => void;
  setChannels: (data: Settings<NotifyChannels>) => void;
  initChannels: () => void;
};

const Context = createContext({} as NotifyContext);

export const useNotifyContext = () => useContext(Context);

interface NotifyProviderProps {
  children: ReactNode;
}

export const NotifyProvider = ({ children }: NotifyProviderProps) => {
  const [notify, dispatchNotify] = useReducer(notifyReducer, {});

  useEffect(() => {
    featchData();
  }, []);

  const featchData = async () => {
    const chanels = await get<NotifyChannels>("notifyChannels");
    dispatchNotify({
      type: "SET_CHANNELS",
      payload: chanels,
    });
  };

  const initChannels = useCallback(() => {
    featchData();
  }, []);

  const setChannel = useCallback((data: { channel: string; data: NotifyChannel }) => {
    dispatchNotify({
      type: "SET_CHANNEL",
      payload: data,
    });
  }, []);

  const setChannels = useCallback((setting: Settings<NotifyChannels>) => {
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
        initChannels,
      }}
    >
      {children}
    </Context.Provider>
  );
};
