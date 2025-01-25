import { CandidacyState } from "./interviews";

export interface GetCandidacyInfoRequest {
  candidacy_id: string;
}

export interface GetCandidacyCommentsRequest {
  candidacy_id: string;
}

export type CommenterType = "ORG_USER" | "HUB_USER";

export const CommenterTypes = {
  ORG_USER: "ORG_USER" as CommenterType,
  HUB_USER: "HUB_USER" as CommenterType,
};

export interface CandidacyComment {
  comment_id: string;
  commenter_name: string;
  commenter_type: CommenterType;
  content: string;
  created_at: Date;
}
