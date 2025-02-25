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
