import { EmailAddress } from "../common/common";

export interface AddOfficialEmailRequest {
  email: EmailAddress;
}

export interface VerifyOfficialEmailRequest {
  email: EmailAddress;
  code: string;
}

export interface TriggerVerificationRequest {
  email: EmailAddress;
}

export interface DeleteOfficialEmailRequest {
  email: EmailAddress;
}

export interface OfficialEmail {
  email: EmailAddress;
  lastVerifiedAt?: string;
  verifyInProgress: boolean;
}
