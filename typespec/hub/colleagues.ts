import { Handle } from "../common/common";

export interface HubUserShort {
  handle: Handle;
  name: string;
  short_bio: string;
}

export interface ConnectColleagueRequest {
  handle: Handle;
}

export interface UnlinkColleagueRequest {
  handle: Handle;
}

export interface MyColleagueApprovalsRequest {
  pagination_key?: string;
  limit?: number;
}

export interface MyColleagueApprovals {
  approvals: HubUserShort[];
  pagination_key?: string;
}

export interface MyColleagueSeeksRequest {
  pagination_key?: string;
  limit?: number;
}

export interface MyColleagueSeeks {
  seeks: HubUserShort[];
  pagination_key?: string;
}

export interface ApproveColleagueRequest {
  handle: Handle;
}

export interface RejectColleagueRequest {
  handle: Handle;
}

export interface FilterColleaguesRequest {
  filter: string;
  limit: number;
}

export enum EndorsementState {
  SoughtEndorsement = "SOUGHT_ENDORSEMENT",
  Endorsed = "ENDORSED",
  DeclinedEndorsement = "DECLINED_ENDORSEMENT",
}

export interface MyEndorseApprovalsRequest {
  pagination_key?: string;
  state: EndorsementState[];
  limit?: number;
}

export interface MyEndorseApproval {
  application_id: string;
  applicant_handle: Handle;
  applicant_name: string;
  applicant_short_bio: string;
  employer_name: string;
  employer_domain: string;
  opening_title: string;
  opening_url: string;
  application_status: string;
  application_created_at: string;
  endorsement_status: EndorsementState;
}

export interface MyEndorseApprovalsResponse {
  endorsements: MyEndorseApproval[];
  pagination_key?: string;
}
