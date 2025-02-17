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

export interface GetBioRequest {
  handle: string;
}

export interface Bio {
  handle: string;
  fullName: string;
  shortBio: string;
  longBio: string;
}

export interface UpdateBioRequest {
  shortBio: string;
  longBio: string;
}

export interface UploadProfilePictureRequest {
  image: Uint8Array;
}
