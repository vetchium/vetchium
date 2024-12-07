import { EmailAddress, Password } from "./common";

export type HubUserState =
  | "ACTIVE_HUB_USER"
  | "DELETED_HUB_USER"
  | "DISABLED_HUB_USER";

export interface LoginRequest {
  email: EmailAddress;
  password: Password;
}

export interface LoginResponse {
  token: string;
}

export interface HubTFARequest {
  tfa_token: string;
  tfa_code: string;
  remember_me: boolean;
}

export interface HubTFAResponse {
  session_token: string;
}

export interface InviteUserRequest {
  email: EmailAddress;
}

export interface GetMyHandleResponse {
  handle: string;
}

export interface ChangePasswordRequest {
  old_password: Password;
  new_password: Password;
}

export interface ForgotPasswordRequest {
  email: EmailAddress;
}

export interface ForgotPasswordResponse {
  token: string;
}

export interface HubUserResetPasswordRequest {
  token: string;
  password: Password;
}
