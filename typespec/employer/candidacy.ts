import { CandidacyState } from "../common/interviews";
import { InterviewState, InterviewersDecision } from "../common/interviews";
import { OrgUserTiny } from "./orgusers";

export type InterviewType =
  | "IN_PERSON"
  | "VIDEO_CALL"
  | "TAKE_HOME"
  | "UNSPECIFIED";

export const InterviewTypes = {
  IN_PERSON: "IN_PERSON" as InterviewType,
  VIDEO_CALL: "VIDEO_CALL" as InterviewType,
  TAKE_HOME: "TAKE_HOME" as InterviewType,
  UNSPECIFIED: "UNSPECIFIED" as InterviewType,
} as const;

export interface FilterCandidacyInfosRequest {
  opening_id?: string;
  recruiter_email?: string;
  state?: CandidacyState;

  pagination_key?: string;
  limit?: number;
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

export interface AddEmployerCandidacyCommentRequest {
  candidacy_id: string;
  comment: string;
}

export interface AddInterviewRequest {
  candidacy_id: string;
  start_time: Date;
  end_time: Date;
  interview_type: InterviewType;
  description?: string;
}

export interface AddInterviewResponse {
  interview_id: string;
}

export interface EmployerInterview {
  interview_id: string;
  interview_state: InterviewState;
  start_time: Date;
  end_time: Date;
  interview_type: InterviewType;
  description?: string;
  interviewers?: OrgUserTiny[];
  interviewers_decision?: InterviewersDecision;
  positives?: string;
  negatives?: string;
  overall_assessment?: string;
  feedback_to_candidate?: string;
  feedback_submitted_by?: OrgUserTiny;
  feedback_submitted_at?: Date;
  created_at: Date;
}

export function isValidInterviewType(type: string): type is InterviewType {
  return Object.values(InterviewTypes).includes(type as InterviewType);
}
