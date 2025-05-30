import { ApplicationState } from "../common/applications";

export type ApplicationColorTag = "GREEN" | "YELLOW" | "RED";

export const ApplicationColorTags = {
  GREEN: "GREEN" as ApplicationColorTag,
  YELLOW: "YELLOW" as ApplicationColorTag,
  RED: "RED" as ApplicationColorTag,
} as const;

export interface GetApplicationsRequest {
  state: ApplicationState;
  search_query?: string;
  color_tag_filter?: ApplicationColorTag;
  opening_id: string;
  pagination_key?: string;
  limit: number;
}

export interface Endorser {
  full_name: string;
  short_bio: string;
  handle: string;
  current_company_domains?: string[];
}

export interface ModelScore {
  model_name: string;
  score: number;
}

export interface Application {
  id: string;
  cover_letter?: string;
  created_at: Date;
  hub_user_handle: string;
  hub_user_name: string;
  hub_user_short_bio: string;
  hub_user_last_employer_domains?: string[];
  state: ApplicationState;
  color_tag?: ApplicationColorTag;
  endorsers: Endorser[];
  scores: ModelScore[];
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

export interface GetResumeRequest {
  application_id: string;
}

export function isValidApplicationColorTag(
  tag: string
): tag is ApplicationColorTag {
  return Object.values(ApplicationColorTags).includes(
    tag as ApplicationColorTag
  );
}
