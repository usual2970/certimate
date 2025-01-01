import { Resource } from "i18next";

import en from "./en";
import zh from "./zh";

export const LOCALE_ZH_NAME = "zh" as const;
export const LOCALE_EN_NAME = "en" as const;

const resources: Resource = {
  [LOCALE_ZH_NAME]: {
    name: "简体中文",
    translation: zh,
  },
  [LOCALE_EN_NAME]: {
    name: "English",
    translation: en,
  },
};

export default resources;
