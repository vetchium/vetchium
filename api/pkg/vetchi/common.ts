export type CountryCode = string; // ISO 3166-1 alpha-3 code
export type Currency = string; // ISO 4217 currency code
export type EmailAddress = string;
export type Password = string;

export interface ValidationErrors {
  errors: string[];
}
