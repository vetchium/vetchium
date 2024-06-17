import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import en from "./en/index";
import ta from "./ta/index";

i18n.use(initReactI18next).init({
  compatibilityJSON: "v3",
  resources: {
    en: {
      translation: en,
    },
    ta: {
      translation: ta,
    },
  },
  lng: "ta", // default language to use.
  interpolation: {
    escapeValue: false,
  },
});

export default i18n.t;
