import { ClientResponseError } from "pocketbase";

export const getErrMsg = (error: unknown): string => {
  if (error instanceof ClientResponseError) {
    return Object.keys(error.response ?? {}).length ? getErrMsg(error.response) : error.message;
  } else if (error instanceof Error) {
    return error.message;
  } else if (typeof error === "object" && error != null) {
    if ("message" in error) {
      return getErrMsg(error.message);
    } else if ("msg" in error) {
      return getErrMsg(error.msg);
    }
  } else if (typeof error === "string") {
    return error || "Unknown error";
  }

  return "Unknown error";
};
