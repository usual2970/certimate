import { NotifyChannel, Setting } from "@/domain/settings";

type Action =
  | {
      type: "SET_CHANNEL";
      payload: {
        channel: string;
        data: NotifyChannel;
      };
    }
  | {
      type: "SET_CHANNELS";
      payload: Setting;
    };

export const notifyReducer = (state: Setting, action: Action) => {
  switch (action.type) {
    case "SET_CHANNEL": {
      const channel = action.payload.channel;
      return {
        ...state,
        content: {
          ...state.content,
          [channel]: action.payload.data,
        },
      };
    }
    case "SET_CHANNELS": {
      return { ...action.payload };
    }

    default:
      return state;
  }
};
