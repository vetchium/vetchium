import { en } from "@/i18n/en";
import { useCallback } from "react";

type TranslationObject = {
  [key: string]: string | TranslationObject;
};

export function useTranslation() {
  const t = useCallback((key: string): string => {
    // First try direct access with the full key
    if ((en as TranslationObject)[key] !== undefined) {
      const value = (en as TranslationObject)[key];
      if (typeof value === "string") return value;
      console.warn(`Translation key ${key} points to an object, not a string`);
      return key;
    }

    // If not found, try nested access
    const keys = key.split(".");
    let value: TranslationObject | string = en;

    for (const k of keys) {
      if (typeof value === "string" || value[k] === undefined) {
        console.warn(`Translation key not found: ${key}`);
        return key;
      }
      value = value[k];
    }

    if (typeof value !== "string") {
      console.warn(`Translation key ${key} points to an object, not a string`);
      return key;
    }

    return value;
  }, []); // Empty dependencies since en is static

  const tObject = useCallback(
    (key: string): Record<string, string> => {
      const value = t(key);
      if (value === key) return {}; // Key not found or points to a string

      const obj = (en as TranslationObject)[key];
      if (typeof obj === "object") {
        return Object.entries(obj).reduce((acc, [k, v]) => {
          if (typeof v === "string") acc[k] = v;
          return acc;
        }, {} as Record<string, string>);
      }

      return {};
    },
    [t]
  ); // Add t as a dependency since we use it in the callback

  return { t, tObject };
}
