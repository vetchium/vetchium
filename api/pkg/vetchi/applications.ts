import { EmailAddress } from "./common";

export type ApplicationState =
  | "APPLIED"
  | "REJECTED"
  | "SHORTLISTED"
  | "WITHDRAWN"
  | "EXPIRED";

export type ApplicationColorTag = "GREEN" | "YELLOW" | "RED";

export interface ApplyForOpeningRequest {
  opening_id_within_company: string;
  company_domain: string;
  resume: string;
  filename: string;
  cover_letter?: string;
}

export interface GetApplicationsRequest {
  state: ApplicationState;
  search_query?: string;
  color_tag_filter?: ApplicationColorTag;
  opening_id: string;
  pagination_key?: string;
  limit: number;
}

export interface Application {
  id: string;
  cover_letter?: string;
  created_at: Date;
  filename: string;
  hub_user_handle: string;
  hub_user_last_employer_domain?: string;
  resume: string;
  state: ApplicationState;
}

export interface UpdateApplicationStateRequest {
  id: string;
  from_state: ApplicationState;
  to_state: ApplicationState;
}

export interface SetApplicationColorTagRequest {
  application_id: string;
  color_tag: ApplicationColorTag;
}

export interface RemoveApplicationColorTagRequest {
  application_id: string;
}

export interface ShortlistApplicationRequest {
  application_id: string;
}

export interface RejectApplicationRequest {
  application_id: string;
}

export interface MyApplicationsRequest {
  state?: ApplicationState;
  pagination_key?: string;
  limit: number;
}

export interface HubApplication {
  application_id: string;
  state: ApplicationState;
  opening_id: string;
  opening_title: string;
  employer_name: string;
  employer_domain: string;
  created_at: Date;
}

export interface WithdrawApplicationRequest {
  application_id: string;
}
