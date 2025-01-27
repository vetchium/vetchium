import { InterviewState, InterviewersDecision } from "../common/interviews";

export interface AddInterviewerRequest {
  interview_id: string;
  org_user_email: string;
}

export interface RemoveInterviewerRequest {
  interview_id: string;
  org_user_email: string;
}

export interface GetEmployerInterviewsByOpeningRequest {
  opening_id: string;
  states?: InterviewState[];
  pagination_key?: string;
  limit?: number;
}

export interface GetEmployerInterviewsByCandidacyRequest {
  candidacy_id: string;
  states?: InterviewState[];
}

export interface GetInterviewDetailsRequest {
  interview_id: string;
}

export interface PutAssessmentRequest {
  interview_id: string;
  decision?: InterviewersDecision;
  positives?: string;
  negatives?: string;
  overall_assessment?: string;
  feedback_to_candidate?: string;
}

export interface Assessment {
  interview_id: string;
  decision?: InterviewersDecision;
  positives?: string;
  negatives?: string;
  overall_assessment?: string;
  feedback_to_candidate?: string;
  feedback_submitted_by?: string;
  feedback_submitted_at?: Date;
}

export interface GetAssessmentRequest {
  interview_id: string;
}
