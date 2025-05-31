import { EmailAddress, OrgUserRole, Password } from "../common/common";

export type OrgUserState =
  | "ACTIVE_ORG_USER"
  | "ADDED_ORG_USER"
  | "DISABLED_ORG_USER"
  | "REPLICATED_ORG_USER";

export const OrgUserStates = {
  ACTIVE: "ACTIVE_ORG_USER" as OrgUserState,
  ADDED: "ADDED_ORG_USER" as OrgUserState,
  DISABLED: "DISABLED_ORG_USER" as OrgUserState,
  REPLICATED: "REPLICATED_ORG_USER" as OrgUserState,
} as const;

export interface OrgUser {
  name: string;
  email: EmailAddress;
  roles: OrgUserRole[];
  state: OrgUserState;
}

export interface AddOrgUserRequest {
  name: string;
  email: EmailAddress;
  roles: OrgUserRole[];
}

export interface DisableOrgUserRequest {
  email: EmailAddress;
}

export interface EnableOrgUserRequest {
  email: EmailAddress;
}

export interface FilterOrgUsersRequest {
  prefix?: string;
  state?: OrgUserState[];
  pagination_key?: EmailAddress;
  limit?: number;
}

export interface OrgUserTiny {
  name: string;
  email: EmailAddress;
}

export interface OrgUserShort {
  name: string;
  email: EmailAddress;
  vetchi_handle?: string;
}

export interface UpdateOrgUserRequest {
  email: EmailAddress;
  name: string;
  roles: OrgUserRole[];
}

export interface SignupOrgUserRequest {
  name: string;
  password: Password;
  invite_token: string;
}

export interface EmployerForgotPasswordRequest {
  email: EmailAddress;
}

export interface EmployerResetPasswordRequest {
  token: string;
  password: Password;
}

export function isValidOrgUserState(state: string): state is OrgUserState {
  return Object.values(OrgUserStates).includes(state as OrgUserState);
}
