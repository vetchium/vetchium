import { en } from "@/i18n/en";
import { useCallback } from "react";

type TranslationParams = Record<string, string | number>;

// For now, we'll just use English. In the future, this can be expanded to support multiple languages
export function useTranslation() {
  const t = useCallback((key: string, params?: TranslationParams) => {
    const keys = key.split(".");
    let value: any = en;

    for (const k of keys) {
      if (value[k] === undefined) {
        console.warn(`Translation key not found: ${key}`);
        return key;
      }
      value = value[k];
    }

    if (params) {
      return Object.entries(params).reduce(
        (str, [key, value]) => str.replace(`{${key}}`, String(value)),
        value
      );
    }

    return value;
  }, []); // Empty dependencies since en is static

  return { t };
}
