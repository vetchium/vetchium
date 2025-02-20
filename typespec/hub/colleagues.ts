import { Handle } from "../common/common";

export interface ConnectColleagueRequest {
  handle: Handle;
}

export interface UnlinkColleagueRequest {
  handle: Handle;
}

export interface MyColleagueApprovalsRequest {
  handle: Handle;
}

export interface MyColleagueSeeksRequest {
  handle: Handle;
}

export interface ApproveColleagueRequest {
  handle: Handle;
}

export interface RejectColleagueRequest {
  handle: Handle;
}
