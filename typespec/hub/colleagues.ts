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

export interface MyColleagueSeeksRequest {
  pagination_key?: string;
  limit?: number;
}

export interface ApproveColleagueRequest {
  handle: Handle;
}

export interface RejectColleagueRequest {
  handle: Handle;
}
