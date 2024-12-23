import { RSVPStatus } from '../common/interviews';

export interface HubRSVPInterviewRequest {
    interview_id: string;
    rsvp: RSVPStatus;
} 