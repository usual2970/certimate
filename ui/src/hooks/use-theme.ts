import { useTheme } from "ahooks";

export default function () {
  const { theme, themeMode, setThemeMode } = useTheme({ localStorageKey: "certimate-ui-theme" });
  return { theme, themeMode, setThemeMode };
}
