export type OpeningState =
  | "DRAFT_OPENING_STATE"
  | "ACTIVE_OPENING_STATE"
  | "SUSPENDED_OPENING_STATE"
  | "CLOSED_OPENING_STATE";

export const OpeningStates = {
  DRAFT: "DRAFT_OPENING_STATE" as OpeningState,
  ACTIVE: "ACTIVE_OPENING_STATE" as OpeningState,
  SUSPENDED: "SUSPENDED_OPENING_STATE" as OpeningState,
  CLOSED: "CLOSED_OPENING_STATE" as OpeningState,
} as const;

export type OpeningType =
  | "FULL_TIME_OPENING"
  | "PART_TIME_OPENING"
  | "CONTRACT_OPENING"
  | "INTERNSHIP_OPENING"
  | "UNSPECIFIED_OPENING";

export const OpeningTypes = {
  FULL_TIME: "FULL_TIME_OPENING" as OpeningType,
  PART_TIME: "PART_TIME_OPENING" as OpeningType,
  CONTRACT: "CONTRACT_OPENING" as OpeningType,
  INTERNSHIP: "INTERNSHIP_OPENING" as OpeningType,
  UNSPECIFIED: "UNSPECIFIED_OPENING" as OpeningType,
} as const;

export type EducationLevel =
  | "BACHELOR_EDUCATION"
  | "MASTER_EDUCATION"
  | "DOCTORATE_EDUCATION"
  | "NOT_MATTERS_EDUCATION"
  | "UNSPECIFIED_EDUCATION";

export const EducationLevels = {
  BACHELOR: "BACHELOR_EDUCATION" as EducationLevel,
  MASTER: "MASTER_EDUCATION" as EducationLevel,
  DOCTORATE: "DOCTORATE_EDUCATION" as EducationLevel,
  NOT_MATTERS: "NOT_MATTERS_EDUCATION" as EducationLevel,
  UNSPECIFIED: "UNSPECIFIED_EDUCATION" as EducationLevel,
} as const;

export function isValidOpeningState(state: string): state is OpeningState {
  return Object.values(OpeningStates).includes(state as OpeningState);
}

export function isValidOpeningType(type: string): type is OpeningType {
  return Object.values(OpeningTypes).includes(type as OpeningType);
}

export function isValidEducationLevel(level: string): level is EducationLevel {
  return Object.values(EducationLevels).includes(level as EducationLevel);
}
