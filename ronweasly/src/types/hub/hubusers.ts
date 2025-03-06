export type EmailAddress = string;

export interface HubUserInviteRequest {
  email: EmailAddress;
}

export interface LoginRequest {
  email: EmailAddress;
  password: string;
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
