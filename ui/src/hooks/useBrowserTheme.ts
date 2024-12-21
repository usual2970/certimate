import { useTheme } from "ahooks";

export default () => {
  return useTheme({ localStorageKey: "certimate-ui-theme" });
};
