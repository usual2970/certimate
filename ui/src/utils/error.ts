import { ClientResponseError } from "pocketbase";

export const getErrMsg = (error: unknown): string => {
  if (error instanceof ClientResponseError) {
    return error.response != null ? getErrMsg(error.response) : error.message;
  } else if (error instanceof Error) {
    return error.message;
  } else if (typeof error === "object" && error != null) {
    if ("message" in error) {
      return String(error.message);
    } else if ("msg" in error) {
      return String(error.msg);
    }
  } else if (typeof error === "string") {
    return error;
  }

  return String(error ?? "Unknown error");
};
