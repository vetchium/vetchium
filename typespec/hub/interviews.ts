import {
  InterviewState,
  InterviewType,
  RSVPStatus,
} from "../common/interviews";

export interface GetHubInterviewsByCandidacyRequest {
  candidacy_id: string;
  states?: InterviewState[];
}

export interface HubInterview {
  interview_id: string;
  interview_state: InterviewState;
  start_time: Date;
  end_time: Date;
  interview_type: InterviewType;
  description?: string;
  interviewers?: string[];
}

export interface HubRSVPInterviewRequest {
  interview_id: string;
  rsvp: RSVPStatus;
}
