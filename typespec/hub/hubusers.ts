import { EmailAddress, Password } from "../common/common";

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
  remember_me?: boolean;
}

export interface HubTFAResponse {
  session_token: string;
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

export interface ResetPasswordRequest {
  token: string;
  password: Password;
}

export interface GetMyHandleResponse {
  handle: string;
}

export type HubUserState = "ACTIVE_HUB_USER";

export const HubUserStates = {
  ACTIVE: "ACTIVE_HUB_USER" as HubUserState,
} as const;

export function isValidHubUserState(state: string): state is HubUserState {
  return Object.values(HubUserStates).includes(state as HubUserState);
}

export interface HubUserInviteRequest {
  email: EmailAddress;
}
