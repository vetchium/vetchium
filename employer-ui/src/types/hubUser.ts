export interface HubUser {
  id: string;
  email: string;
  name: string;
  state: HubUserState;
  resume_url?: string;
  phone?: string;
  linkedin_url?: string;
  github_url?: string;
  portfolio_url?: string;
  created_at: string;
  last_updated_at: string;
}

export enum HubUserState {
  ACTIVE_HUB_USER = 'ACTIVE_HUB_USER',
  SUSPENDED_HUB_USER = 'SUSPENDED_HUB_USER',
  DEFUNCT_HUB_USER = 'DEFUNCT_HUB_USER',
}

export interface GetHubUserRequest {
  id: string;
}

export interface GetHubUsersRequest {
  pagination_key?: string;
  limit?: number;
} 