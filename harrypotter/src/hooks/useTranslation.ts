import { en } from "@/i18n/en";

export function useTranslation() {
  const t = (key: string) => {
    const keys = key.split(".");
    let value: any = en;

    for (const k of keys) {
      if (value[k] === undefined) {
        console.warn(`Translation key not found: ${key}`);
        return key;
      }
      value = value[k];
    }

    return value;
  };

  return { t };
}
