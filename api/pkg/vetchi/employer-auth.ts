import { EmailAddress, Password } from "./common";

export type OnboardStatus =
  | "DOMAIN_NOT_VERIFIED"
  | "DOMAIN_VERIFIED_ONBOARD_PENDING"
  | "DOMAIN_ONBOARDED";

export interface GetOnboardStatusRequest {
  client_id: string;
}

export interface GetOnboardStatusResponse {
  status: OnboardStatus;
}

export interface SetOnboardPasswordRequest {
  client_id: string;
  password: Password;
  token: string;
}

export interface EmployerSignInRequest {
  client_id: string;
  email: EmailAddress;
  password: Password;
}

export interface EmployerSignInResponse {
  token: string;
}

export interface EmployerTFARequest {
  tfa_code: string;
  tfa_token: string;
  remember_me?: boolean;
}

export interface EmployerTFAResponse {
  session_token: string;
}
