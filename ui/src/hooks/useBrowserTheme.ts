import { useTheme } from "ahooks";

export default function () {
  return useTheme({ localStorageKey: "certimate-ui-theme" });
}
