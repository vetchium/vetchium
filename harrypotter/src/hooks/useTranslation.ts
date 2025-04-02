import { en } from "@/i18n/en";
import { useCallback } from "react";

interface TranslationObject {
  [key: string]: string | TranslationObject;
}

type TranslationParams = Record<string, string | number>;

export function useTranslation() {
  const t = useCallback((key: string, params?: TranslationParams): string => {
    // First try direct access with the full key
    if ((en as TranslationObject)[key] !== undefined) {
      const value = (en as TranslationObject)[key];
      if (typeof value === "string") {
        if (params) {
          return Object.entries(params).reduce(
            (str, [key, value]) => str.replace(`{${key}}`, String(value)),
            value
          );
        }
        return value;
      }
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

    if (params) {
      return Object.entries(params).reduce(
        (str, [key, value]) => str.replace(`{${key}}`, String(value)),
        value
      );
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
