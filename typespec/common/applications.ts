export type ApplicationState =
  | "APPLIED"
  | "REJECTED"
  | "SHORTLISTED"
  | "WITHDRAWN"
  | "EXPIRED";

export const ApplicationStates = {
  APPLIED: "APPLIED" as ApplicationState,
  REJECTED: "REJECTED" as ApplicationState,
  SHORTLISTED: "SHORTLISTED" as ApplicationState,
  WITHDRAWN: "WITHDRAWN" as ApplicationState,
  EXPIRED: "EXPIRED" as ApplicationState,
} as const;

export function isValidApplicationState(
  state: string
): state is ApplicationState {
  return Object.values(ApplicationStates).includes(state as ApplicationState);
}
