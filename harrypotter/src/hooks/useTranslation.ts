import { en } from "@/i18n/en";

type TranslationObject = {
  [key: string]: string | TranslationObject;
};

export function useTranslation() {
  const t = (key: string) => {
    // First try direct access with the full key
    if ((en as TranslationObject)[key] !== undefined) {
      return (en as TranslationObject)[key] as string;
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
  };

  return { t };
}
