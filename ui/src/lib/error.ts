export const getErrMessage = (error: unknown): string => {
  if (error instanceof Error) {
    return error.message;
  } else if (
    typeof error === "object" &&
    error !== null &&
    "message" in error
  ) {
    return String(error.message);
  } else if (typeof error === "string") {
    return error;
  } else {
    return "Something went wrong";
  }
};
