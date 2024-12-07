import { EmailAddress, Password } from "./common";

export type OrgUserRole =
  | "ADMIN"
  | "APPLICATIONS_CRUD"
  | "APPLICATIONS_VIEWER"
  | "COST_CENTERS_CRUD"
  | "COST_CENTERS_VIEWER"
  | "LOCATIONS_CRUD"
  | "LOCATIONS_VIEWER"
  | "OPENINGS_CRUD"
  | "OPENINGS_VIEWER"
  | "ORG_USERS_CRUD"
  | "ORG_USERS_VIEWER";

export type OrgUserState =
  | "ACTIVE_ORG_USER"
  | "ADDED_ORG_USER"
  | "DISABLED_ORG_USER"
  | "REPLICATED_ORG_USER";

export interface OrgUser {
  name: string;
  email: EmailAddress;
  state: OrgUserState;
  roles: OrgUserRole[];
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
  pagination_key?: EmailAddress;
  limit?: number;
  state?: OrgUserState[];
}

export interface OrgUserShort {
  name: string;
  email: EmailAddress;
  vetchi_handle?: string;
}

export interface SignupOrgUserRequest {
  name: string;
  password: Password;
  invite_token: string;
}

export interface UpdateOrgUserRequest {
  email: EmailAddress;
  name: string;
  roles: OrgUserRole[];
}
