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
  last_verified_at?: string;
  verify_in_progress: boolean;
}

export interface GetBioRequest {
  handle: string;
}

export interface Bio {
  handle: string;
  full_name: string;
  short_bio: string;
  long_bio: string;
  verified_mail_domains?: string[];
}

export interface UpdateBioRequest {
  handle?: string;
  full_name?: string;
  short_bio?: string;
  long_bio?: string;
}

export interface UploadProfilePictureRequest {
  image: Uint8Array;
}
