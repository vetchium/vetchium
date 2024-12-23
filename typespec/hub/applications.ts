import { ApplicationState } from '../common/applications';

export interface MyApplicationsRequest {
    state?: ApplicationState;
    pagination_key?: string;
    limit: number;
}

export interface HubApplication {
    application_id: string;
    state: ApplicationState;
    opening_id: string;
    opening_title: string;
    employer_name: string;
    employer_domain: string;
    created_at: Date;
}

export interface WithdrawApplicationRequest {
    application_id: string;
} 