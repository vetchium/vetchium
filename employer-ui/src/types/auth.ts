export interface AuthState {
  isAuthenticated: boolean;
  token: string | null;
  user: OrgUser | null;
  loading: boolean;
  error: string | null;
}

export interface OrgUser {
  email: string;
  name: string;
  roles: string[];
  state: string;
}

export interface OrgUserShort {
  id: string;
  name: string;
  email: string;
}

export interface SignInCredentials {
  client_id: string;
  email: string;
  password: string;
}

export interface SignInResponse {
  token: string;
}

export interface TFARequest {
  tfa_code: string;
  tfa_token: string;
  remember_me: boolean;
}

export interface TFAResponse {
  session_token: string;
} 