import { useTheme } from "ahooks";

export type UseBrowserThemeReturns = ReturnType<typeof useTheme>;

/**
 * 获取并设置当前浏览器系统主题。
 * @returns {UseBrowserThemeReturns}
 */
const useBrowserTheme = (): UseBrowserThemeReturns => {
  return useTheme({ localStorageKey: "certimate-ui-theme" });
};

export default useBrowserTheme;
