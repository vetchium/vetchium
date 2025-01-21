import { CandidacyState } from "./interviews";

export interface GetCandidacyInfoRequest {
  candidacy_id: string;
}

export interface GetCandidacyCommentsRequest {
  candidacy_id: string;
}

export interface CandidacyComment {
  comment_id: string;
  commenter_name: string;
  commenter_type: string;
  content: string;
  created_at: Date;
}
