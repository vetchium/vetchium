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

export type ColleagueConnectionState = {
  CAN_SEND_REQUEST: {};
  CANNOT_SEND_REQUEST: {};
  REQUEST_SENT_PENDING: {};
  REQUEST_RECEIVED_PENDING: {};
  CONNECTED: {};
  REJECTED_BY_ME: {};
  REJECTED_BY_THEM: {};
  UNLINKED_BY_ME: {};
  UNLINKED_BY_THEM: {};
}[keyof {
  CAN_SEND_REQUEST: {};
  CANNOT_SEND_REQUEST: {};
  REQUEST_SENT_PENDING: {};
  REQUEST_RECEIVED_PENDING: {};
  CONNECTED: {};
  REJECTED_BY_ME: {};
  REJECTED_BY_THEM: {};
  UNLINKED_BY_ME: {};
  UNLINKED_BY_THEM: {};
}];

export interface Bio {
  handle: string;
  full_name: string;
  short_bio: string;
  long_bio: string;
  verified_mail_domains?: string[];
  colleague_connection_state: ColleagueConnectionState;
}

export interface UpdateBioRequest {
  full_name?: string;
  short_bio?: string;
  long_bio?: string;
}

export interface UploadProfilePictureRequest {
  image: Uint8Array;
}
