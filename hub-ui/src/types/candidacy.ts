import { Application } from './application';

export interface Candidacy {
  id: string;
  application: Application;
  state: CandidacyState;
  timeline: CandidacyEvent[];
  interviews: Interview[];
  created_at: string;
  last_updated_at: string;
}

export enum CandidacyState {
  SCREENING_CANDIDACY = 'SCREENING_CANDIDACY',
  INTERVIEWING_CANDIDACY = 'INTERVIEWING_CANDIDACY',
  OFFERED_CANDIDACY = 'OFFERED_CANDIDACY',
  ACCEPTED_CANDIDACY = 'ACCEPTED_CANDIDACY',
  DECLINED_CANDIDACY = 'DECLINED_CANDIDACY',
  REJECTED_CANDIDACY = 'REJECTED_CANDIDACY',
  WITHDRAWN_CANDIDACY = 'WITHDRAWN_CANDIDACY',
}

export interface CandidacyEvent {
  id: string;
  candidacy_id: string;
  event_type: CandidacyEventType;
  description: string;
  created_at: string;
}

export enum CandidacyEventType {
  STATE_CHANGE = 'STATE_CHANGE',
  INTERVIEW_SCHEDULED = 'INTERVIEW_SCHEDULED',
  INTERVIEW_COMPLETED = 'INTERVIEW_COMPLETED',
  INTERVIEW_FEEDBACK = 'INTERVIEW_FEEDBACK',
  OFFER_MADE = 'OFFER_MADE',
  OFFER_ACCEPTED = 'OFFER_ACCEPTED',
  OFFER_DECLINED = 'OFFER_DECLINED',
  CANDIDACY_WITHDRAWN = 'CANDIDACY_WITHDRAWN',
}

export interface Interview {
  id: string;
  candidacy_id: string;
  round: number;
  scheduled_at: string;
  duration_minutes: number;
  meeting_link?: string;
  interviewers: string[];
  feedback?: InterviewFeedback;
  state: InterviewState;
  created_at: string;
  last_updated_at: string;
}

export enum InterviewState {
  SCHEDULED_INTERVIEW = 'SCHEDULED_INTERVIEW',
  COMPLETED_INTERVIEW = 'COMPLETED_INTERVIEW',
  CANCELLED_INTERVIEW = 'CANCELLED_INTERVIEW',
}

export interface InterviewFeedback {
  id: string;
  interview_id: string;
  positives: string;
  negatives: string;
  overall: string;
  rating: InterviewRating;
  feedback_to_candidate?: string;
  created_at: string;
}

export enum InterviewRating {
  STRONG_YES = 'STRONG_YES',
  YES = 'YES',
  NO = 'NO',
  STRONG_NO = 'STRONG_NO',
}

export interface FilterCandidaciesRequest {
  state?: CandidacyState[];
  from_date?: string;
  to_date?: string;
  pagination_key?: string;
  limit?: number;
}

export interface GetCandidacyRequest {
  id: string;
}

export interface WithdrawCandidacyRequest {
  id: string;
}

export interface RespondToOfferRequest {
  candidacy_id: string;
  accept: boolean;
} 