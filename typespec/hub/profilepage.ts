import { EmailAddress } from "../common/common";

export interface AddProfessionalEmailRequest {
  email: EmailAddress;
}

export interface VerifyProfessionalEmailRequest {
  email: EmailAddress;
}

export interface TriggerVerificationRequest {
  email: EmailAddress;
}

export interface DeleteProfessionalEmailRequest {
  email: EmailAddress;
}

export interface ProfessionalEmail {
  email: EmailAddress;
  lastVerifiedAt?: string;
  verifyInProgress: boolean;
}
