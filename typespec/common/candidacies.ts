import { CandidacyState } from "./interviews";

export interface GetCandidacyInfoRequest {
  candidacy_id: string;
}

export interface Candidacy {
  candidacy_id: string;
  opening_id: string;
  opening_title: string;
  opening_description: string;
  candidacy_state: CandidacyState;
  applicant_name: string;
  applicant_handle: string;
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
