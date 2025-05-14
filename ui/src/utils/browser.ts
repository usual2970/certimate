export const isBrowserHappy = () => {
  return typeof Promise.withResolvers === "function";
};
