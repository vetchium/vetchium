export interface AuthState {
  isAuthenticated: boolean;
  token: string | null;
  user: HubUser | null;
  loading: boolean;
  error: string | null;
}

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

export interface SignInCredentials {
  email: string;
  password: string;
}

export interface SignInResponse {
  token: string;
}

export interface SignUpRequest {
  email: string;
  password: string;
  name: string;
  phone?: string;
  linkedin_url?: string;
  github_url?: string;
  portfolio_url?: string;
}

export interface SignUpResponse {
  token: string;
}

export interface UpdateProfileRequest {
  name?: string;
  phone?: string;
  linkedin_url?: string;
  github_url?: string;
  portfolio_url?: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
} 