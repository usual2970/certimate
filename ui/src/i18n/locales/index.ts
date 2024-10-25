import { Resource } from "i18next";

import zh from "./zh";
import en from "./en";
import de from "./de";

const resources: Resource = {
  zh: {
    name: "简体中文",
    translation: zh,
  },
  en: {
    name: "English",
    translation: en,
  },
  de: {
    name: "Deutsch",
    translation: de,
  },
};

export default resources;
