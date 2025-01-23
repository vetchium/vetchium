import { CandidacyState, InterviewType } from "../common/interviews";
import { InterviewState, InterviewersDecision } from "../common/interviews";
import { OrgUserTiny } from "./orgusers";

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
  interviewer_emails?: string[];
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
