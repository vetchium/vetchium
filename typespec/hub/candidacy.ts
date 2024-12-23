import { CandidacyState } from '../common/interviews';
import { GetCandidacyCommentsRequest, CandidacyComment } from '../common/candidacies';

export interface AddHubCandidacyCommentRequest {
    candidacy_id: string;
    comment: string;
}

export interface MyCandidaciesRequest {
    candidacy_states?: CandidacyState[];
    pagination_key?: string;
    limit?: number;
}

export interface MyCandidacy {
    candidacy_id: string;
    company_name: string;
    company_domain: string;
    opening_id: string;
    opening_title: string;
    opening_description: string;
    candidacy_state: CandidacyState;
}

export type { GetCandidacyCommentsRequest, CandidacyComment }; 